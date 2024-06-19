package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ResponseStatusType string

const (
	ResponseOK  ResponseStatusType = "ok"
	ResponseErr ResponseStatusType = "error"
)

type Result struct {
	Status       ResponseStatusType `json:"status"`
	ErrorMessage string             `json:"error"`
	Response     json.RawMessage    `json:"response"`
}

func (r *Result) Error() string {
	return r.ErrorMessage
}

const UserAgent string = "fprotimaru"

type Request struct {
	Method   string
	Endpoint string
	Query    url.Values
	Header   http.Header
	FullUrl  string
	Body     io.Reader
}

type Client struct {
	baseUrl    string
	token      string
	httpClient *http.Client
}

func NewClient(baseUrl, token string, httpClient *http.Client) *Client {
	return &Client{
		baseUrl:    baseUrl,
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) CallAPI(ctx context.Context, r *Request) (*Result, error) {
	if err := c.parseRequest(r); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.Method, r.FullUrl, r.Body)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)

	req.URL.RawQuery = r.Query.Encode()

	fmt.Println("---->", req.URL.RawQuery)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	p, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := new(Result)
	if err := json.Unmarshal(p, result); err != nil {
		return nil, err
	}
	defer func() {
		closeErr := res.Body.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	switch result.Status {
	case ResponseOK:
		return result, nil
	case ResponseErr:
		return nil, result
	default:
		return result, nil
	}
}

func (c *Client) parseRequest(r *Request) error {
	if len(r.FullUrl) == 0 {
		r.FullUrl = fmt.Sprintf("%s%s", c.baseUrl, r.Endpoint)
	}

	header := http.Header{}
	if r.Header != nil {
		header = r.Header.Clone()
	}
	header.Set("User-Agent", UserAgent)

	r.Header = header

	r.Query.Set("token", c.token)

	return nil
}
