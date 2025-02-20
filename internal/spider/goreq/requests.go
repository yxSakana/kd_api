package goreq

import (
	"crypto/tls"
	pcookejar "github.com/juju/persistent-cookiejar"
	"io"
	"net/http"
	urllib "net/url"
	"os"
	"path/filepath"
)

type Session struct {
	http.Client
}

func NewSession(verify bool) *Session {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: verify,
		},
	}
	jar, _ := pcookejar.New(&pcookejar.Options{
		Filename: "cookies_file",
	})
	//jar, _ := cookiejar.New(nil)
	return &Session{http.Client{Transport: tr, Jar: jar}}
}

func (s *Session) Get(url string, headers map[string]string, query map[string]interface{}) (*http.Response, error) {
	return s.Fetch("GET", url, headers, query, nil)
}

func (s *Session) Post(url string, headers map[string]string, query map[string]interface{}, data io.Reader) (*http.Response, error) {
	return s.Fetch("POST", url, headers, query, data)
}

func (s *Session) Download(filename string, mkdir bool, url string, headers map[string]string, query map[string]interface{}) error {
	if mkdir {
		err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		if err != nil {
			return err
		}
	}
	res, err := s.Get(url, headers, query)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) GetCookie(u *urllib.URL, key string) string {
	for _, cookie := range s.Jar.Cookies(u) {
		if cookie.Name == key {
			return cookie.Value
		}
	}
	return ""
}

func (s *Session) Fetch(
	method, url string, headers map[string]string, query map[string]interface{}, data io.Reader,
) (*http.Response, error) {
	queryStr := QueryMapToString(query)
	if queryStr != "" {
		url += "?" + queryStr
	}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return s.Do(req)
}
