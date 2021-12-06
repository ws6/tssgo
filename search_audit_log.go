package tssgo

import (
	"context"
	"encoding/json"
	"fmt"

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

func (self *Client) SearchAuditLog(ctx context.Context, qp map[string]string) (*CaseListResp, error) {
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

	body, err := self.GetBytes(ctx, base.String())
	if err != nil {
		return nil, err
	}

	ret := new(CaseListResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil

	return nil, nil
}
