package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
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

func (c *Client) Login() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	location, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		log.Fatalf("%v", err)
	}
	_, utcoffset := time.Now().In(location).Zone()

	values := url.Values{}
	values.Set("username", env.MustGet("IR_USERNAME"))
	values.Set("password", env.MustGet("IR_PASSWORD"))
	values.Set("utcoffset", fmt.Sprintf("%d", utcoffset/60))
	values.Set("todaysdate", "")

	req, err := http.NewRequest("POST", "https://members.iracing.com/membersite/Login", strings.NewReader(values.Encode()))
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//log.Debugf("%v", string(data))
	if strings.Contains(strings.ToLower(string(data)), "email address or password was invalid") ||
		strings.Contains(strings.ToLower(string(data)), "invalid email address or password") ||
		strings.Contains(strings.ToLower(resp.Header.Get("Location")), "failedlogin") {
		return fmt.Errorf("login failed")
	}
	c.lastLogin = time.Now()
	return nil
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) Post(url string, values url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// relogin if needed
	if c.lastLogin.Before(time.Now().Add(-5 * time.Minute)) {
		if err := c.Login(); err != nil {
			return nil, err
		}
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	time.Sleep(2345 * time.Millisecond)

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
		return nil, fmt.Errorf("failed request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}
	return data, nil
}
