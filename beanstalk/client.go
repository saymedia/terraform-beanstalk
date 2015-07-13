package beanstalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ClientConfig struct {
	AccountName string
	Username    string
	AccessToken string
}

type Client struct {
	httpClient  *http.Client
	apiURL      *url.URL
	username    string
	accessToken string
}

func NewClient(config *ClientConfig) (*Client, error) {
	httpClient := &http.Client{}

	apiURL, err := url.Parse(fmt.Sprintf("https://%s.beanstalkapp.com/api/", config.AccountName))
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient:  httpClient,
		apiURL:      apiURL,
		username:    config.Username,
		accessToken: config.AccessToken,
	}, nil
}

type request struct {
	Method string
	PathParts []string
	QueryArgs map[string]string
	Headers map[string]string
	BodyBytes []byte
}

func (c *Client) rawRequest(req *request) ([]byte, error) {
	httpReq := req.MakeHTTPRequest(c)
	log.Printf("Beanstalk %v request to %v", httpReq.Method, httpReq.URL)
	log.Printf("Request body is %v", string(req.BodyBytes))
	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	log.Printf("Response body is %v", string(resBodyBytes))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, &NotFoundError{}
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP Error %v", res.StatusCode)
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, nil
	}

	return resBodyBytes, nil
}

func (c *Client) jsonRequest(method string, pathParts []string, queryArgs map[string]string, reqBody interface{}, result interface{}) error {

	var err error
	var reqBodyBytes []byte
	reqBodyBytes = nil
	if reqBody != nil {
		reqBodyBytes, err = json.Marshal(reqBody)
	}

	req := &request{
		Method: method,
		PathParts: pathParts,
		QueryArgs: queryArgs,
		BodyBytes: reqBodyBytes,
		Headers: map[string]string{},
	}

	if reqBody != nil {
		req.Headers["Content-Type"] = "application/json"
	}

	resBodyBytes, err := c.rawRequest(req)
	if err != nil {
		return err
	}

	if result != nil {
		if resBodyBytes == nil {
			return fmt.Errorf("server did not return a JSON payload")
		}
		err = json.Unmarshal(resBodyBytes, result)
		if err != nil {
			return fmt.Errorf("error decoding response JSON payload: %s", err.Error())
		}
		log.Printf("Response structure is %v\n", result)
	}

	return nil
}

func (c *Client) Get(pathParts []string, queryArgs map[string]string, result interface{}) error {
	return c.jsonRequest("GET", pathParts, queryArgs, nil, result)
}

func (c *Client) Post(pathParts []string, reqBody interface{}, result interface{}) error {
	return c.jsonRequest("POST", pathParts, nil, reqBody, result)
}

func (c *Client) Put(pathParts []string, reqBody interface{}, result interface{}) error {
	return c.jsonRequest("PUT", pathParts, nil, reqBody, result)
}

func (c *Client) Delete(pathParts []string) error {
	return c.jsonRequest("DELETE", pathParts, nil, nil, nil)
}

func (r *request) MakeHTTPRequest(client *Client) *http.Request {
	req := &http.Request{
		Method: r.Method,
		Header: http.Header{},
	}

	req.Header.Add("User-Agent", "Terraform-Beanstalk")
	req.SetBasicAuth(client.username, client.accessToken)

	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	urlPath := &url.URL{
		Path: strings.Join(r.PathParts, "/") + ".json",
	}
	reqURL := client.apiURL.ResolveReference(urlPath)
	req.URL = reqURL

	if len(r.QueryArgs) > 0 {
		urlQuery := url.Values{}
		for k, v := range r.QueryArgs {
			urlQuery.Add(k, v)
		}
		reqURL.RawQuery = urlQuery.Encode()
	}

	if r.BodyBytes != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(r.BodyBytes))
		req.ContentLength = int64(len(r.BodyBytes))
	}

	return req
}

type NotFoundError struct {}

func (err NotFoundError) Error() string {
	return "not found"
}
