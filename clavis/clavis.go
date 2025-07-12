package clavis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Clavis struct {
	client  *http.Client
	baseURL string
	token   string
}

func New() (*Clavis, error) {
	clavisURL := os.Getenv("CLAVIS_URL")
	if clavisURL == "" {
		return nil, errors.New("CLAVIS_URL is empty")
	}
	token := os.Getenv("CLAVIS_TOKEN")
	if token == "" {
		return nil, errors.New("CLAVIS_TOKEN is empty")
	}
	return &Clavis{
		client:  &http.Client{},
		baseURL: strings.TrimRight(clavisURL, "/"),
		token:   token,
	}, nil
}

func (c *Clavis) Get(key string) (string, error) {
	url := fmt.Sprintf(`%s/get?key=%s`, c.baseURL, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	rsp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	type Response struct {
		Code uint32 `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}

	rspBody := &Response{}
	err = json.Unmarshal(body, rspBody)
	if err != nil {
		return "", err
	}
	if rspBody.Code != 0 {
		return "", errors.New(rspBody.Msg)
	}

	return rspBody.Data, nil
}

func (c *Clavis) Set(key string, value string) error {
	data := map[string]string{
		"key":   key,
		"value": value,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/set", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	type Response struct {
		Code uint32 `json:"code"`
		Msg  string `json:"msg"`
	}

	rspBody := &Response{}
	err = json.Unmarshal(body, rspBody)
	if err != nil {
		return err
	}
	if rspBody.Code != 0 {
		return errors.New(rspBody.Msg)
	}

	return nil
}
