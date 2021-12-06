package tssgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

const (
	DEFAULT_LIST_PAGE_SIZE = 25
)

// /crs/api/v1/cases/list

// example request body
// {
// 	"filters": [ ],
// 	"onlyUnassigned": "false",
// 	"onlyOverDueCases": "false",
// 	"showOverDueCases": true,
// 	"showUnspecifiedCases": true,
// 	"pageNumber": 0,
// 	"pageSize": 1,
// 	"sortByColumns": ["+createdDate"],
// 	"tags": [],
// 	"phiTypes": ["summary"],
// 	"clientTimeZoneId": "America/Los_Angeles"
// }

type CaseRequestParams struct {
	Filters              []string `json:"filters"`
	OnlyUnassigned       string   `json:"onlyUnassigned"`
	OnlyOverDueCases     string   `json:"onlyOverDueCases"`
	ShowOverDueCases     bool     `json:"showOverDueCases"`
	ShowUnspecifiedCases bool     `json:"showUnspecifiedCases"`
	PageNumber           int      `json:"pageNumber"`
	PageSize             int      `json:"pageSize"`
	SortByColumns        []string `json:"sortByColumns,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	PhiTypes             []string `json:"phiTypes,omitempty"`
}

type CaseListResp struct {
	Content []map[string]interface{} `json:"content"`
	// Pageable struct {
	// 	Sort struct {
	// 		Sorted   bool `json:"sorted"`
	// 		Unsorted bool `json:"unsorted"`
	// 		Empty    bool `json:"empty"`
	// 	} `json:"sort"`
	// 	Offset     int  `json:"offset"`
	// 	PageNumber int  `json:"pageNumber"`
	// 	PageSize   int  `json:"pageSize"`
	// 	Paged      bool `json:"paged"`
	// 	Unpaged    bool `json:"unpaged"`
	// } `json:"pageable"`
	Last          bool `json:"last"`
	TotalPages    int  `json:"totalPages"`
	TotalElements int  `json:"totalElements"`
	Sort          struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	Size             int  `json:"size"`
	Number           int  `json:"number"`
	NumberOfElements int  `json:"numberOfElements"`
	First            bool `json:"first"`
	Empty            bool `json:"empty"`
}

//SearchCaseByListAPI an internal API found in TSS UI only. it is not documented
//https://icsl-test.trusight.illumina.com/crs/api/v1/cases/list/search?searchTerm=internal-sample-id123&pageNumber=17&orderBy=ASC&pageSize=1
func (self *Client) SearchCaseByListAPI(ctx context.Context, params map[string]string) (*CaseListResp, error) {
	_url := `/crs/api/v1/cases/list/search`

	base, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}
	q := url.Values{}

	for k, v := range params {
		q.Add(k, v)
	}

	base.RawQuery = q.Encode()

	resp, err := self.NewRequestWithContext(ctx, `POST`, base.String(), nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := new(CaseListResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil
}

//SearchCase supported v2 search parameters
// https://icsl-test.trusight.illumina.com/crs/swagger-ui/index.html#/Cases/searchCasesUsingGET
func (self *Client) SearchCase(ctx context.Context, params map[string]string) (*CaseListResp, error) {
	_url := `/crs/api/v2/cases/search`

	base, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}
	q := url.Values{}

	for k, v := range params {
		q.Add(k, v)
	}

	base.RawQuery = q.Encode()

	resp, err := self.GetBytes(ctx, base.String())
	if err != nil {
		return nil, err
	}

	ret := new(CaseListResp)
	if err := json.Unmarshal(resp, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil
}

func (self *Client) ListCase(ctx context.Context, opt *CaseRequestParams) (*CaseListResp, error) {

	if opt.PageSize <= 0 {
		opt.PageSize = DEFAULT_LIST_PAGE_SIZE
	}

	rb, err := json.Marshal(opt)
	if err != nil {
		return nil, err
	}

	url := `/crs/api/v1/cases/list`
	resp, err := self.NewRequestWithContext(ctx, `POST`, url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(`bad status code-%d:%s`, resp.StatusCode, string(body))
	}

	ret := new(CaseListResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil
}
