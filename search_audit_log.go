package tssgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

//searchAuditLog
//supported query parameters
//userName
//userId
//toDate
//pageSize int
//pageNumber int
//orderBy
//fromDate 2021-10-25T14:12:06+0000
//eventName
//entityType repeatable header
//entityId
//caseId

func (self *Client) searchAuditLog(ctx context.Context, qp map[string]string) (*CaseResp, error) {
	_url := `/als/api/v1/auditlogs/search`
	base, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	for k, v := range qp {
		q.Add(k, v)
	}
	base.RawQuery = q.Encode()

	resp, err := self.NewRequestWithContext(ctx, `GET`, base.String(), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer fmt.Println(string(body))
	if resp.StatusCode != 200 {

		return nil, fmt.Errorf(`bad status code-%d:%s`, resp.StatusCode, string(body))
	}

	ret := new(CaseResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil

	return nil, nil
}
