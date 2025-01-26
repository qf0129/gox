package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	netUrl "net/url"

	"github.com/qf0129/gox/pkg/logx"
)

func Get[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](http.MethodGet, url, headers, params, nil)
}

func Post[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](http.MethodPost, url, headers, nil, body)
}

func Put[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](http.MethodPut, url, headers, nil, body)
}

func Delete[T any](url string, headers map[string]string, body any) (*T, error) {
	return Request[T](http.MethodDelete, url, headers, nil, body)
}

func Head[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](http.MethodHead, url, headers, params, nil)
}

func Options[T any](url string, headers, params map[string]string) (*T, error) {
	return Request[T](http.MethodOptions, url, headers, params, nil)
}

func Request[T any](method, url string, headers, params map[string]string, body any) (*T, error) {
	logx.Info("HttpRequestStart", "method", method, "url", url, "params", params, "headers", headers, "body", body)
	if len(params) > 0 {
		urlValues := netUrl.Values{}
		for k, v := range params {
			urlValues.Add(k, v)
		}
		url = url + "?" + urlValues.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New("JsonMarshalError, " + err.Error())
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = nil
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.New("CreateRequestError, " + err.Error())
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
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
