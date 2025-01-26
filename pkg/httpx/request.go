package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	netUrl "net/url"
	"os"

	"github.com/qf0129/gox/pkg/logx"
	"golang.org/x/net/proxy"
)

func Get[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodGet, Url: url, Headers: headers, Params: params})
}

func Post[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodPost, Url: url, Headers: headers, Body: body})
}
func Put[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodPut, Url: url, Headers: headers, Body: body})
}

func Delete[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodDelete, Url: url, Headers: headers, Body: body})
}

func Head[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodHead, Url: url, Headers: headers, Params: params})
}

func Options[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](&RequestOption{Method: http.MethodOptions, Url: url, Headers: headers, Params: params})
}

type RequestOption struct {
	Method  string
	Url     string
	Headers map[string]string
	Params  map[string]string
	Body    any
	Socks5  string
}

func Request[T any](opt *RequestOption) (*T, error) {
	logx.Debug("HttpRequestStart", "method", opt.Method, "url", opt.Url, "params", opt.Params, "headers", opt.Headers, "body", opt.Body)
	if len(opt.Params) > 0 {
		urlValues := netUrl.Values{}
		for k, v := range opt.Params {
			urlValues.Add(k, v)
		}
		opt.Url = opt.Url + "?" + urlValues.Encode()
	}

	var reqBody io.Reader
	if opt.Body != nil {
		jsonBytes, err := json.Marshal(opt.Body)
		if err != nil {
			return nil, errors.New("JsonMarshalError, " + err.Error())
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = nil
	}

	req, err := http.NewRequest(opt.Method, opt.Url, reqBody)
	if err != nil {
		return nil, errors.New("CreateRequestError, " + err.Error())
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(opt.Headers) > 0 {
		for k, v := range opt.Headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	if opt.Socks5 != "" {
		dialer, err := proxy.SOCKS5("tcp", opt.Socks5, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			return nil, errors.New("RequestDeeplError, " + err.Error())
		}

		httpTransport := &http.Transport{}
		httpTransport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		client.Transport = httpTransport
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("SendRequestError, " + err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ReadResponseError, " + err.Error())
	}
	logx.Info("HttpReponse", "status", resp.Status)
	logx.Info("HttpReponse", "response", string(respBody))
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HttpRequestError, resp=" + string(respBody))
	}

	var result T
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.New("JsonUnmarshalErr, " + err.Error())
	}
	logx.Info("HttpRequestSuccess", "response", result)
	return &result, nil
}
