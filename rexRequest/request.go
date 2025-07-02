package rexRequest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type (
	RequestClient interface {
		GetClient() *http.Client
		Head(url string) (resp *http.Response, err error)
		HeadSync(url string) func() ([]byte, error)
		Get(url string) (resp *http.Response, err error)
		GetSync(url string) func() ([]byte, error)
		Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
		PostSyncJson(url string, sendBody interface{}) func() ([]byte, error)
		PostSyncJsonWithHeaders(url string, sendBody interface{}, headers map[string]string) func() ([]byte, error)
		PostSyncJsonBodyWithHeaders(url string, jsonBody string, headers map[string]string) func() ([]byte, error)
		PostSyncJsonWithFile(url string, fileFieldName, fileName string, file io.Reader, fields map[string]string) func() ([]byte, error)
		Put(url, contentType string, body io.Reader) (resp *http.Response, err error)
		PutSyncJson(url string, sendBody interface{}) func() ([]byte, error)
		Delete(url string) (resp *http.Response, err error)
		DeleteSync(url string) func() ([]byte, error)
	}
	requestClient struct {
		client *http.Client
	}
)

func NewRequestClient() RequestClient {
	return &requestClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c *requestClient) GetClient() *http.Client {
	return c.client
}

func (c *requestClient) Head(url string) (resp *http.Response, err error) {
	return c.client.Head(url)
}

func (c *requestClient) HeadSync(url string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		var res *http.Response
		res, err = c.Head(url)

		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) Get(url string) (resp *http.Response, err error) {
	return c.client.Get(url)
}

func (c *requestClient) GetSync(url string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		var res *http.Response
		res, err = c.Get(url)

		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	postRequest, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	postRequest.Header.Set("Content-Type", contentType)
	return c.client.Do(postRequest)
}

func (c *requestClient) PostSyncJson(url string, sendBody interface{}) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		var res *http.Response

		sendBodyBt, err := json.Marshal(sendBody)
		if err != nil {
			return
		}

		res, err = c.Post(url, "application/json", bytes.NewBuffer(sendBodyBt))

		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) PostSyncJsonWithHeaders(url string, sendBody interface{}, headers map[string]string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		sendBodyBt, err := json.Marshal(sendBody)
		if err != nil {
			return
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(sendBodyBt))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		for key, val := range headers {
			req.Header.Set(key, val)
		}

		var res *http.Response
		res, err = c.client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) PostSyncJsonBodyWithHeaders(url string, jsonBody string, headers map[string]string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(jsonBody)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		for key, val := range headers {
			req.Header.Set(key, val)
		}

		var res *http.Response
		res, err = c.client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) PostSyncJsonWithFile(url string, fileFieldName, fileName string, file io.Reader, fields map[string]string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		sendBody := new(bytes.Buffer)
		writer := multipart.NewWriter(sendBody)
		formFile, err := writer.CreateFormFile(fileFieldName, fileName)
		if err != nil {
			return
		}
		_, err = io.Copy(formFile, file)
		if err != nil {
			return
		}

		for key, val := range fields {
			_ = writer.WriteField(key, val)
		}
		if err = writer.Close(); err != nil {
			return
		}
		req, err := http.NewRequest(http.MethodPost, url, sendBody)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())

		var res *http.Response
		res, err = c.client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) Put(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	putRequest, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	putRequest.Header.Set("Content-Type", contentType)
	return c.client.Do(putRequest)
}

func (c *requestClient) PutSyncJson(url string, sendBody interface{}) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		var res *http.Response

		sendBodyBt, err := json.Marshal(sendBody)
		if err != nil {
			return
		}

		res, err = c.Put(url, "application/json", bytes.NewBuffer(sendBodyBt))

		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}

func (c *requestClient) Delete(url string) (resp *http.Response, err error) {
	delRequest, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(delRequest)
}

func (c *requestClient) DeleteSync(url string) func() ([]byte, error) {
	var body []byte
	var err error

	ch := make(chan struct{}, 1)
	go func() {
		defer close(ch)

		var res *http.Response

		res, err = c.Delete(url)

		if err != nil {
			return
		}

		defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		return
	}()
	return func() ([]byte, error) {
		_, ok := <-ch
		if !ok {
			//fmt.Println("channel closed!")
			return body, err
		}
		return body, err
	}
}
