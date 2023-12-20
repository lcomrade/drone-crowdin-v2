// Copyright (C) 2022-2023 Leonid Maslakov.

// This file is part of drone-crowdin-v2.

// drone-crowdin-v2 is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// drone-crowdin-v2 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with drone-crowdin-v2.
// If not, see <https://www.gnu.org/licenses/>.

package crowdin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	baseAddr        = "https://api.crowdin.com"
	paginationLimit = 500
	userAgent       = "drone-crowdin-v2"
	tmpFilePattern  = "drone-crowdin-v2-"
)

var BadSymbols = []rune{'\\', '/', ':', '*', '?', '"', '<', '>', '|'}

type Client struct {
	key    string
	client *http.Client
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func readBody(r io.Reader) string {
	if r == nil {
		return "<nil>"
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

func (client *Client) get(s string, goodCode int) (*http.Response, error) {
	req, err := http.NewRequest("GET", baseAddr+s, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+client.key)

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, errors.New("crowdin api: " + err.Error())
	}

	if resp.StatusCode != goodCode {
		return nil, errors.New("crowdin api: GET " + s + ": " + resp.Status + ": " + readBody(resp.Body))
	}

	return resp, nil
}

func (client *Client) dlToTmpFile(u string) (string, error) {
	// Request
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	//req.Header.Set("Authorization", "Bearer "+client.key)

	resp, err := client.client.Do(req)
	if err != nil {
		return "", errors.New("crowdin api: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return "", errors.New("crowdin api: GET " + u + ": " + resp.Status + ": " + readBody(resp.Body))
	}

	// Save to temp file
	tmpFile, err := os.CreateTemp("", tmpFilePattern)
	if err != nil {
		return "", errors.New("crowdin api: failed create temp file: " + err.Error())
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", errors.New("crowdin api: failed create temp file: " + err.Error())
	}

	return tmpFile.Name(), nil
}

func (client *Client) sendJSON(method string, s string, goodCode int, body interface{}) (*http.Response, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("crowdin api: " + method + " " + s + ": " + err.Error())
	}

	req, err := http.NewRequest(method, baseAddr+s, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+client.key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, errors.New("crowdin api: " + err.Error())
	}

	if resp.StatusCode != goodCode {
		return nil, errors.New("crowdin api: " + method + " " + s + ": " + resp.Status + ": " + readBody(resp.Body))
	}

	return resp, nil
}

func (client *Client) uploadFileExtra(s string, goodCode int, body io.Reader, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", baseAddr+s, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+client.key)
	req.Header.Set("Content-Type", "application/octet-stream")
	if header != nil {
		for key, val := range header {
			req.Header.Set(key, val)
		}
	}

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, errors.New("crowdin api: " + err.Error())
	}

	if resp.StatusCode != goodCode {
		return nil, errors.New("crowdin api: POST " + s + ": " + resp.Status + ": " + readBody(resp.Body))
	}

	return resp, nil
}
