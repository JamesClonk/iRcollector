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
	CookieJar   *cookiejar.Jar
	Token       Token
	mutex       *sync.Mutex
	lastLogin   time.Time
	lastRefresh time.Time
}

type Token struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

func New() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &Client{
		CookieJar:   cookieJar,
		mutex:       &sync.Mutex{},
		lastLogin:   time.Now().Add(-24 * time.Hour),
		lastRefresh: time.Now().Add(-24 * time.Hour),
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

func (c *Client) LoginToken() error {
	log.Debugf("login via oauth.iracing.com ...")

	// https://oauth.iracing.com/oauth2/book/password_limited_flow.html
	hash := sha256.Sum256([]byte(env.MustGet("IR_PASSWORD") + strings.ToLower(env.MustGet("IR_USERNAME"))))
	hashedPassword := base64.StdEncoding.EncodeToString(hash[:])
	hash = sha256.Sum256([]byte(env.MustGet("IR_CLIENT_SECRET") + strings.ToLower(env.MustGet("IR_CLIENT_ID"))))
	hashedSecret := base64.StdEncoding.EncodeToString(hash[:])
	data := []byte(fmt.Sprintf(`grant_type=password_limited&client_id=%s&client_secret=%s&username=%s&password=%s&scope=iracing.auth`,
		url.QueryEscape(env.MustGet("IR_CLIENT_ID")),
		url.QueryEscape(hashedSecret),
		url.QueryEscape(env.MustGet("IR_USERNAME")),
		url.QueryEscape(hashedPassword),
	))

	req, err := http.NewRequest("POST", "https://oauth.iracing.com/oauth2/token", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	// read oauth token
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.Token = Token{}
	if err := json.Unmarshal(data, &c.Token); err != nil {
		clientLoginError.Inc()
		log.Errorf("could not unmarshal oauth token: %s", data)
		return err
	}

	// default values
	if c.Token.ExpiresIn == 0 {
		c.Token.ExpiresIn = 555
	}
	if c.Token.RefreshTokenExpiresIn == 0 {
		c.Token.ExpiresIn = 3456
	}

	// if we have no refresh-token, then set its expiry time to same as normal token, to force relogin before normal token expires
	if len(c.Token.RefreshToken) == 0 {
		c.Token.RefreshTokenExpiresIn = c.Token.ExpiresIn
	}

	return nil
}

func (c *Client) RefreshToken() error {
	log.Debugf("refreshing token via oauth.iracing.com ...")

	// https://oauth.iracing.com/oauth2/book/token_endpoint.html#refresh-token-grant
	hash := sha256.Sum256([]byte(env.MustGet("IR_CLIENT_SECRET") + strings.ToLower(env.MustGet("IR_CLIENT_ID"))))
	hashedSecret := base64.StdEncoding.EncodeToString(hash[:])
	data := []byte(fmt.Sprintf(`grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s`,
		url.QueryEscape(env.MustGet("IR_CLIENT_ID")),
		url.QueryEscape(hashedSecret),
		url.QueryEscape(c.Token.RefreshToken),
	))

	req, err := http.NewRequest("POST", "https://oauth.iracing.com/oauth2/token", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	// read oauth token
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.Token = Token{}
	if err := json.Unmarshal(data, &c.Token); err != nil {
		clientLoginError.Inc()
		log.Errorf("could not unmarshal oauth token: %s", data)
		return err
	}

	// default values
	if c.Token.ExpiresIn == 0 {
		c.Token.ExpiresIn = 555
	}
	if c.Token.RefreshTokenExpiresIn == 0 {
		c.Token.ExpiresIn = 3456
	}

	// if we have no refresh-token, then set its expiry time to same as normal token, to force relogin before normal token expires
	if len(c.Token.RefreshToken) == 0 {
		c.Token.RefreshTokenExpiresIn = c.Token.ExpiresIn
	}

	if len(c.Token.AccessToken) == 0 {
		return fmt.Errorf("refreshing token failed, no new access-token in response")
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

	// relogin after a long time, or if refresh token is about to expire
	if c.lastLogin.Before(time.Now().Add(-2*time.Hour)) ||
		c.lastLogin.Before(time.Now().Add(-1*time.Duration(c.Token.RefreshTokenExpiresIn)*time.Second).Add(30*time.Second)) {
		if err := c.LoginToken(); err != nil {
			clientLoginError.Inc()
			time.Sleep(3 * time.Second) // safety sleep
			return nil, err
		}
		c.lastLogin = time.Now()
		c.lastRefresh = c.lastLogin
	}
	// refresh token if needed
	if c.lastRefresh.Before(time.Now().Add(-1 * time.Duration(c.Token.ExpiresIn) * time.Second).Add(30 * time.Second)) {
		if err := c.RefreshToken(); err != nil {
			clientLoginError.Inc()
			time.Sleep(3 * time.Second) // safety sleep
			return nil, err
		}
		c.lastRefresh = time.Now()
	}

	req.Header.Add("User-Agent", "iRcollector")
	//req.Header.Add("Referer", "https://members.iracing.com/membersite/login.jsp")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Charset", "UTF-8,utf-8;q=0.7,*;q=0.3")
	req.Header.Add("Cache-Control", "max-age=0")
	//req.Header.Add("Origin", "members.iracing.com")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")

	client := &http.Client{
		Jar: c.CookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		clientRequestError.Inc()
		time.Sleep(2 * time.Second) // safety sleep
		return nil, fmt.Errorf("failed request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		clientRequestError.Inc()
		time.Sleep(2 * time.Second) // safety sleep
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	/*
		X-Ratelimit-Limit:[240]
		X-Ratelimit-Remaining:[239]
		X-Ratelimit-Reset:[1641553935]
	*/
	// check ratelimiting values
	ratelimitRemaining := resp.Header.Get("Ratelimit-Remaining")
	if len(ratelimitRemaining) == 0 {
		ratelimitRemaining = resp.Header.Get("X-Ratelimit-Remaining")
	}
	ratelimitReset := resp.Header.Get("Ratelimit-Reset")
	if len(ratelimitReset) == 0 {
		ratelimitReset = resp.Header.Get("X-Ratelimit-Reset")
	}
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
		time.Sleep(2 * time.Second)
	} else {
		//log.Debugln("could not determine ratelimit, will do a safety sleep ...")
		time.Sleep(444 * time.Millisecond) // safety sleep
	}
	time.Sleep(111 * time.Millisecond) // safety sleep

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clientRequestError.Inc()
		return nil, fmt.Errorf("read body: %v", err)
	}
	clientRequestTotal.Inc()
	return data, nil
}
