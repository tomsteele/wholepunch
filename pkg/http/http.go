package http

import (
	"crypto/tls"
	"errors"
	"net/http"
)

// GetSpoofHostHeader sends a GET request to a URL with a modified Host header.
func GetSpoofHostHeader(url, hostHeader, userAgent string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Host = hostHeader
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, errors.New("connection was successful but an invalid status code was returned")
	}
	return true, nil
}

// TLSGetSpoofHostHeader sends a GET request to a URL with a modified Host header.
func TLSGetSpoofHostHeader(url, hostHeader, userAgent string) (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Host = hostHeader
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, errors.New("connection was successful but an invalid status code was returned")
	}
	return true, nil
}

// Get send an HTTP GET request to a URL.
func Get(url, userAgent string) (bool, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, errors.New("connection was successful but an invalid status code was returned")
	}
	return true, nil
}

// TLSGet send an HTTP GET request to a URL.
func TLSGet(url, userAgent string) (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, errors.New("connection was successful but an invalid status code was returned")
	}
	return true, nil
}

// TLSGetSpoofSNI sends an HTTPS GET request with a spoofed SNI in the TLS ClientHello.
func TLSGetSpoofSNI(url, serverName, userAgent string) (bool, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName:         serverName,
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, errors.New("connection was successful but an invalid status code was returned")
	}
	return true, nil
}
