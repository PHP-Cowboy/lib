package myhttp

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

func DoHttpPost(path string, data string) (resp string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 0, Transport: tr}
	req, _ := http.NewRequest("POST", path, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	resp = string(body)
	return
}

func DoHttpPostJson(path string, data string) (resp string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 0, Transport: tr}
	req, _ := http.NewRequest("POST", path, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	resp = string(body)
	return
}

func DoHttpPostHeader(path string, data string, header map[string]string) (resp string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 0, Transport: tr}
	req, _ := http.NewRequest("POST", path, strings.NewReader(data))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	resp = string(body)
	return
}

func DoHttpGet(path string, header map[string]string) (resp string, err error) {
	client := &http.Client{Timeout: 0}
	req, _ := http.NewRequest("GET", path, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := io.ReadAll(httpResp.Body)
	defer httpResp.Body.Close()
	resp = string(body)
	return
}
