package infrastructure

import (
	"fmt"
	"homeserver/common"
	"io"
	"net/http"
	"net/url"
	"time"
)

const DefaultTimeout = time.Duration(10) * time.Second

func HttpGet(targetURL string, params map[string]any, proxy string, timeout time.Duration) ([]byte, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	client := &http.Client{
		Timeout: timeout,
	}
	if len(proxy) > 0 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			common.Log.Error().
				Err(err).
				Str("proxy", proxy).
				Msg("Parse proxy url error.")
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	reqURL, err := buildURL(targetURL, params)
	if err != nil {
		common.Log.Error().
			Err(err).
			Str("targerURL", targetURL).
			Any("param", params).
			Msg("Build request URL error.")
		return nil, err
	}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		common.Log.Error().
			Err(err).
			Str("targerURL", targetURL).
			Any("param", params).
			Msg("Request error.")
		return nil, err
	}
	rawBody, err := parseResponse(resp)
	if err != nil {
		common.Log.Error().
			Err(err).
			Str("targerURL", targetURL).
			Any("resp", resp).
			Msg("Parse resposne error")
		return nil, err
	}
	return rawBody, nil
}

func parseResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			common.Log.Error().Err(err).Msg("Close http body error")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code not ok, code=%d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func buildURL(targetURL string, params map[string]any) (string, error) {
	reqURL, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	if len(params) > 0 {
		param := url.Values{}
		for key, value := range params {
			param.Add(key, fmt.Sprintf("%v", value))
		}
		reqURL.RawQuery = param.Encode()
	}
	return reqURL.String(), nil
}
