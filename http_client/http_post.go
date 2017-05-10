package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Jar struct {
	m      sync.Mutex
	perURL map[string][]*http.Cookie
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.m.Lock()
	defer j.m.Unlock()
	if j.perURL == nil {
		j.perURL = make(map[string][]*http.Cookie)
	}
	j.perURL[u.Host] = cookies
}

func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	j.m.Lock()
	defer j.m.Unlock()
	return j.perURL[u.Host]
}
func main() {

	login_url := "http://52.34.131.226:10086/login"
	var client = &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	client.Jar = &Jar{perURL: make(map[string][]*http.Cookie)}
	resp, err := client.Get(login_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	_ = string(body)
	cookies := resp.Cookies()
	fmt.Println(cookies)
	fmt.Println(resp.StatusCode)
	//fmt.Println(client.Jar.Cookies())

	// Add form data
	v := url.Values{}
	v.Set("username", "public")
	v.Set("password", "public@derby")

	// Values.Encode() encodes the values into "URL encoded" form sorted by key.
	s := v.Encode()
	//	fmt.Printf("v.Encode(): %v\n", s)
	//	resp, err = http.PostForm(login_url, v)
	//	if err != nil {
	//		panic(err)
	//	}

	req, err := http.NewRequest("POST", login_url, strings.NewReader(s))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//	c := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}
	_ = string(data)

	//fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
	cookies_login := resp.Header
	fmt.Println(cookies_login)
	fmt.Println(resp.StatusCode)
}
