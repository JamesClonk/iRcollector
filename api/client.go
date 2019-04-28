package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/JamesClonk/iRcollector/env"
	"github.com/JamesClonk/iRcollector/log"
)

type Client struct {
	CookieJar *cookiejar.Jar
}

func New() *Client {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return &Client{
		CookieJar: cookieJar,
	}
}

func (c *Client) Login() error {
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

	if strings.Contains(strings.ToLower(string(data)), "email address or password was invalid") ||
		strings.Contains(strings.ToLower(string(data)), "invalid email address or password") ||
		resp.Header.Get("Location") == "https://members.iracing.com/membersite/failedlogin.jsp" ||
		resp.Header.Get("Location") == "http://members.iracing.com/membersite/failedlogin.jsp" {
		return fmt.Errorf("login failed")
	}
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
