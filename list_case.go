package tssgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type CaseResp struct {
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
	Size   int  `json:"size"`
	Number int  `json:"number"`
	First  bool `json:"first"`
	Empty  bool `json:"empty"`
}

func (self *Client) SearchCase(ctx context.Context) (*CaseResp, error) {
	url := `/crs/api/v2/cases/search?subState=READY_FOR_INTERPRETATION`

	resp, err := self.NewRequestWithContext(ctx, `GET`, url, nil)
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
}

func (self *Client) ListCase(ctx context.Context, opt *CaseRequestParams) (*CaseResp, error) {

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

	ret := new(CaseResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil
}
