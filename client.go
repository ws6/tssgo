package tssgo

import (
	"context"
	"fmt"
	"io"
	"strings"

	"net/http"
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
	ret.httpclient.Timeout = time.Second * time.Duration(DEFAULT_TIMEOUT_SECOND)
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

func (self *Client) AttachHeaders(req *http.Request) {
	for _, rk := range requiredKeys() {
		req.Header.Set(rk, self.cfg[rk])
	}

	req.Header.Set(CONTENT_TYPE, `application/json`)

	v := req.Header.Get(AUTH_TOKEN)
	if strings.HasPrefix(v, `apiKkey`) {
		req.Header.Set(AUTH_TOKEN, `apiKkey`+" "+v)
	}
	//for compatible with als domain api
	// https://tss-test.trusight.illumina.com/als/swagger-ui/index.html#/audit-log-controller/getAuditLogsUsingGET
	req.Header.Set(`Authorization`, req.Header.Get(AUTH_TOKEN))

}

//NewRequestWithContext over write http.NewRequestWithContext wih auth headers
func (self *Client) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	absUrl := self.GetBaseUrl() + url

	req, err := http.NewRequestWithContext(ctx, method, absUrl, body)
	if err != nil {
		return nil, fmt.Errorf(`NewRequest:%s`, err.Error())
	}

	self.AttachHeaders(req)

	return self.httpclient.Do(req)
}
