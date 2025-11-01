package main

import (
	"net/http"
	"net/url"
	"strconv"
)

func (t *TorrentFile) CreateURL(peerID []byte, port int) (string, error) {
	base, err := url.Parse(t.announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.info_hash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(port)},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.info.length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func sendAnnounce(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
