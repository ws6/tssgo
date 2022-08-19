package tssgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//illumina tss golang client
//http://support-docs.illumina.com/SW/TSSS/TruSight_SW_API/Content/SW/FrontPages/TruSightSoftware_API.htm
type Client struct {
	cfg        map[string]string
	httpclient *http.Client
}

var (
	DEFAULT_TIMEOUT_SECOND = 60
)

const (
	AUTH_TOKEN     = `X-Auth-Token`
	ILMN_DOMAIN    = `X-ILMN-Domain`
	ILMN_WORKGROUP = `X-ILMN-Workgroup`
	CONTENT_TYPE   = `Content-Type`
	BASE_URL       = `BASE_URL`
)

func requiredKeys() []string {
	return []string{
		AUTH_TOKEN, ILMN_DOMAIN, ILMN_WORKGROUP,
	}
}

func getTimeoutSeconds(cfg map[string]string) int {
	//from config
	if s, ok := cfg[`TSSGO_TIMEOUT_SECOND`]; ok {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}
	//from env
	if s := os.Getenv(`TSSGO_TIMEOUT_SECOND`); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return n
		}
	}

	//from default
	return DEFAULT_TIMEOUT_SECOND
}

func NewClient(cfg map[string]string) (*Client, error) {
	ret := new(Client)
	ret.cfg = cfg
	if err := ret.ConfigCheck(); err != nil {
		return nil, fmt.Errorf(`ConfigCheck:%s`, err.Error())
	}

	baseUrl := ret.GetBaseUrl()
	if baseUrl == "" {
		return nil, fmt.Errorf(`no base url`)
	}

	ret.httpclient = new(http.Client)
	//you can use GetHttpClient then set the client paramters
	ret.httpclient.Timeout = time.Second * time.Duration(getTimeoutSeconds(cfg))
	return ret, nil
}

func (self *Client) GetHttpClient() *http.Client {
	return self.httpclient
}

func (self *Client) ConfigCheck() error {
	for _, rk := range requiredKeys() {
		if _, ok := self.cfg[rk]; !ok {
			return fmt.Errorf(`missing required key:%s`, rk)
		}
	}
	return nil
}

func (self *Client) GetBaseUrl() string {
	//expect https://tss-test.trusight.illumina.com
	return self.cfg[`BASE_URL`]
}

func (self *Client) AttachHeaders(req *http.Request, fromUrl string) {
	for _, rk := range requiredKeys() {
		req.Header.Set(rk, self.cfg[rk])
	}

	req.Header.Set(CONTENT_TYPE, `application/json`)

	v := req.Header.Get(AUTH_TOKEN)
	if !strings.HasPrefix(v, `apiKey`) {
		req.Header.Set(AUTH_TOKEN, `apiKey`+" "+v)
	}
	//for compatible with als domain api
	// https://tss-test.trusight.illumina.com/als/swagger-ui/index.html#/audit-log-controller/getAuditLogsUsingGET
	req.Header.Set(`Authorization`, req.Header.Get(AUTH_TOKEN))

}

type ModifyRequest func(*http.Request)

//NewRequestWithContext over write http.NewRequestWithContext wih auth headers
func (self *Client) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader, modifiers ...ModifyRequest) (*http.Response, error) {
	absUrl := self.GetBaseUrl() + url

	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		absUrl = url
	}

	req, err := http.NewRequestWithContext(ctx, method, absUrl, body)
	if err != nil {
		return nil, fmt.Errorf(`NewRequest:%s`, err.Error())
	}
	if len(modifiers) == 0 {
		self.AttachHeaders(req, url)
	}

	if len(modifiers) > 0 {
		for _, modFn := range modifiers {
			modFn(req)
		}
	}

	return self.httpclient.Do(req)
}

//GetBytes a GET method with a return type []byte
func (self *Client) GetBodyReader(ctx context.Context, url string, modFns ...ModifyRequest) (io.ReadCloser, error) {
	resp, err := self.NewRequestWithContext(ctx, `GET`, url, nil, modFns...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {

		return nil, fmt.Errorf(`bad status code:%d`, resp.StatusCode)
	}
	return resp.Body, nil

}

//GetBytes a GET method with a return type []byte
func (self *Client) GetBytes(ctx context.Context, url string, modFns ...ModifyRequest) ([]byte, error) {
	reader, err := self.GetBodyReader(ctx, url, modFns...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

//GetMsi a GET method with a return type map[string]interface{}
func (self *Client) GetMsi(ctx context.Context, url string) (map[string]interface{}, error) {
	//!!! no modFns passed in because set Content-Type:applicaon/json
	body, err := self.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
