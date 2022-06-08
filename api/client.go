package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	clientLoginError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ircollector_api_client_login_error_total",
		Help: "Total number of iRcollector API client login errors.",
	})
	clientRequestError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ircollector_api_client_request_error_total",
		Help: "Total number of iRcollector API client request errors.",
	})
	clientRequestTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ircollector_api_client_request_total",
		Help: "Total number of iRcollector API client request.",
	})
)

type Client struct {
	CookieJar *cookiejar.Jar
	mutex     *sync.Mutex
	lastLogin time.Time
}

func New() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &Client{
		CookieJar: cookieJar,
		mutex:     &sync.Mutex{},
		lastLogin: time.Now().Add(-24 * time.Hour),
	}
}

func (c *Client) LoginNG() error {
	log.Debugf("login to members-ng ...")

	// https://forums.iracing.com/discussion/22109/login-form-changes
	hash := sha256.Sum256([]byte(env.MustGet("IR_PASSWORD") + strings.ToLower(env.MustGet("IR_USERNAME"))))
	password := base64.StdEncoding.EncodeToString(hash[:])
	data := []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, env.MustGet("IR_USERNAME"), password))

	req, err := http.NewRequest("POST", "https://members-ng.iracing.com/auth", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Jar: c.CookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(1 * time.Minute)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with HTTP [%d]", resp.StatusCode)
	}
	return nil
}

func (c *Client) FollowLink(url string) ([]byte, error) {
	// get target link for caching first
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clientRequestError.Inc()
		return nil, err
	}
	data, err := c.doRequest(req)
	if err != nil {
		clientRequestError.Inc()
		return nil, err
	}

	link := Link{}
	if err := json.Unmarshal(data, &link); err != nil {
		clientRequestError.Inc()
		log.Errorf("could not unmarshal cache link: %s", data)
		return nil, err
	}

	// now get the actual data
	req, err = http.NewRequest("GET", link.Target, nil)
	if err != nil {
		clientRequestError.Inc()
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clientRequestError.Inc()
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) Post(url string, values url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		clientRequestError.Inc()
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// relogin if needed
	if c.lastLogin.Before(time.Now().Add(-5 * time.Minute)) {
		if err := c.LoginNG(); err != nil {
			clientLoginError.Inc()
			time.Sleep(2222 * time.Millisecond) // safety sleep
			return nil, err
		}
		c.lastLogin = time.Now()
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")
	req.Header.Add("Referer", "https://members.iracing.com/membersite/login.jsp")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Charset", "UTF-8,utf-8;q=0.7,*;q=0.3")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Origin", "members.iracing.com")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")

	client := &http.Client{
		Jar: c.CookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		clientRequestError.Inc()
		time.Sleep(2222 * time.Millisecond) // safety sleep
		return nil, fmt.Errorf("failed request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		clientRequestError.Inc()
		time.Sleep(2222 * time.Millisecond) // safety sleep
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	/*
		X-Ratelimit-Limit:[240]
		X-Ratelimit-Remaining:[239]
		X-Ratelimit-Reset:[1641553935]
	*/
	// check ratelimiting values
	ratelimitRemaining := resp.Header.Get("X-Ratelimit-Remaining")
	ratelimitReset := resp.Header.Get("X-Ratelimit-Reset")
	// do we have the necessary headers? (is it members-ng?)
	if len(ratelimitRemaining) > 0 && len(ratelimitReset) > 0 {
		remaining, err := strconv.Atoi(ratelimitRemaining)
		if err != nil {
			remaining = 0
		}
		if remaining < 10 {
			sleepEpoch, err := strconv.ParseInt(ratelimitReset, 10, 64)
			if err != nil {
				sleepEpoch = time.Now().Add(1 * time.Minute).Unix()
			}
			log.Debugf("sleeping for ratelimit, until: %v", time.Unix(sleepEpoch, 0))
			time.Sleep(time.Until(time.Unix(sleepEpoch, 0)))
		}
	} else if req.URL.Host == "members.iracing.com" {
		// old API, lets sleep a fixed amount
		log.Debugf("sleeping for 2s because of old API call to: [%s, %s]", req.URL.Host, req.URL.RequestURI())
		time.Sleep(2222 * time.Millisecond)
	} else {
		//log.Debugln("could not determine ratelimit, will do a safety sleep ...")
		time.Sleep(444 * time.Millisecond) // safety sleep
	}
	time.Sleep(222 * time.Millisecond) // safety sleep

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clientRequestError.Inc()
		return nil, fmt.Errorf("read body: %v", err)
	}
	clientRequestTotal.Inc()
	return data, nil
}
