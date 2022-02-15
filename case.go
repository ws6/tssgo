package tssgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	CASE_STATUS_IN_PROGRESS             = `In Progress`
	CASE_STATUS_COMPLETE                = `Complete`
	CASE_SUB_STATE_CLOSED               = `CLOSED`
	CASE_SUB_STATE_DELIVERED            = `DELIVERED`
	CASE_SUB_STATE_DELIVERY_IN_PROGRESS = `DELIVERY_IN_PROGRESS`
	CASE_SUB_STATE_READY_FOR_DELIVERY   = `READY_FOR_DELIVERY `
	CASE_SUB_STATE_READY_HAS_ISSUE      = `HAS_ISSUE `
)

var ERR_NOT_FOUND = fmt.Errorf(`not found`)

//case.go any case related API calls

//CaseResp mapping the model of API resposne
///crs/api/v1/cases/$caseId
type CaseResp struct {
	Id                 string      `json:"id"`
	DisplayId          string      `json:"displayId"`
	CreatedDate        string      `json:"createdDate"`
	UpdatedDate        string      `json:"updatedDate"`
	ClientId           string      `json:"clientId"`
	ClientAddressId    string      `json:"clientAddressId"`
	ClientRecipientIds string      `json:"clientRecipientIds"`
	Client             *ClientResp `json:"client"`
	StartDate          string      `json:"startDate"`
	DueDate            string      `json:"dueDate"`
	CompletedDate      string      `json:"completedDate"`
	Status             string      `json:"status"`
	SubState           string      `json:"subState"`
	CaseSubjects       []struct {
		Phenotypes  []map[string]string `json:"phenotypes"`
		ReportTypes []map[string]string `json:"reportTypes"`
		Subject     map[string]string   `json:"subject"`
		Samples     []struct {
			CreatedDate        string `json:"createdDate"`
			UpdatedDate        string `json:"updatedDate"`
			Id                 string `json:"id"`
			ExternalSampleId   string `json:"externalSampleId"`
			SampleName         string `json:"sampleName"`
			ExternalSampleName string `json:"externalSampleName"`
			SampleType         string `json:"sampleType"`
			Status             string `json:"status"`
			SubState           string `json:"subState"`
			DateReceived       string `json:"dateReceived"`
			DateCollected      string `json:"dateCollected"`
			SampleSourceType   string `json:"sampleSourceType"`
		} `json:"samples"`
		ActiveSample struct {
			CreatedDate        string `json:"createdDate"`
			UpdatedDate        string `json:"updatedDate"`
			Id                 string `json:"id"`
			ExternalSampleId   string `json:"externalSampleId"`
			SampleName         string `json:"sampleName"`
			ExternalSampleName string `json:"externalSampleName"`
			SampleType         string `json:"sampleType"`
			Status             string `json:"status"`
			SubState           string `json:"subState"`
			DateReceived       string `json:"dateReceived"`
			DateCollected      string `json:"dateCollected"`
			SampleSourceType   string `json:"sampleSourceType"`
		} `json:"activeSample"`
	} `json:"caseSubjects"`

	ActiveSample struct {
		ReatedDate         string `json:"reatedDate"`
		UpdatedDate        string `json:"updatedDate"`
		Id                 string `json:"id"`
		ExternalSampleId   string `json:"externalSampleId"`
		SampleName         string `json:"sampleName"`
		ExternalSampleName string `json:"externalSampleName"`
		SampleType         string `json:"sampleType"`
		Status             string `json:"status"`
		SubState           string `json:"subState"`
		DateReceived       string `json:"dateReceived"`
		DateCollected      string `json:"dateCollected"`
		SampleSourceType   string `json:"sampleSourceType"`
	} `json:"activeSample"`
	TestDefinition struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Order   int    `json:"order"`
		Subject string `json:"subject"`
	} `json:"testDefinition"`
}

//GetCaseById get almost full resposne specs for a case
//directIdentifiers=true will return PHI
// https://support-docs.illumina.com/SW/TSSS/TruSight_SW_API/Content/SW/TSSS/API/GetCase_fTSSS.htm
func (self *Client) GetCaseById(ctx context.Context, caseId string, params map[string]string) (*CaseResp, error) {
	_url := fmt.Sprintf(`/crs/api/v1/cases/%s`, caseId)

	base, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}
	q := url.Values{}
	if params != nil {
		for k, v := range params {
			q.Add(k, v)
		}
	}

	base.RawQuery = q.Encode()

	body, err := self.GetBytes(ctx, base.String())
	if err != nil {
		return nil, err
	}
	ret := new(CaseResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *Client) GetCaseByIdWithPHI(ctx context.Context, caseId string) (*CaseResp, error) {
	params := map[string]string{
		`directIdentifiers`: `true`,
	}
	return self.GetCaseById(ctx, caseId, params)
}

func msiToCaseResp(m map[string]interface{}) (*CaseResp, error) {
	body, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	ret := new(CaseResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

///api/v1/cases/list/search?searchTerm=internal-sample-id123&pageNumber=0&orderBy=ASC&pageSize=2
func (self *Client) GetCasesSearchChan(ctx context.Context, params map[string]string) chan *CaseResp {
	pageSize := 30
	if ns, ok := params[`pageSize`]; ok {
		if n, err := strconv.Atoi(ns); err == nil && n > 0 {
			pageSize = n
		}
	}
	pageNumber := 0
	if ns, ok := params[`pageNumber`]; ok {
		if n, err := strconv.Atoi(ns); err == nil && n > 0 {
			pageNumber = n
		}
	}
	orderBy := `ASC`

	if _orderBy, ok := params[`orderBy`]; ok {
		orderBy = _orderBy
	}

	ret := make(chan *CaseResp, 2*pageSize)
	makeParams := func() {
		params[`pageSize`] = fmt.Sprintf(`%d`, pageSize)
		params[`pageNumber`] = fmt.Sprintf(`%d`, pageNumber)
		params[`orderBy`] = orderBy
	}

	go func() {
		defer close(ret)

		for {
			makeParams()
			pageNumber++
			resp, err := self.SearchCaseByListAPI(ctx, params)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, item := range resp.Content {
				topush, err := msiToCaseResp(item)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ret <- topush
			}

			if resp.NumberOfElements < pageSize || len(resp.Content) < pageSize {
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}

		}

	}()

	return ret
}

type CaseSearchFilterFn func(*CaseResp) bool

func (self *Client) GetCasesSearchBySearchTerm(ctx context.Context, searchTerm string, filterFn CaseSearchFilterFn) chan *CaseResp {

	pageSize := 30
	ret := make(chan *CaseResp, 2*pageSize)
	go func() {

		defer close(ret)

		params := map[string]string{
			`searchTerm`: searchTerm,
			`pageSize`:   fmt.Sprintf(`%d`, pageSize),
		}
		ctx0, cancelFn := context.WithCancel(ctx)
		defer cancelFn()

		ch := self.GetCasesSearchChan(ctx0, params)

		for item := range ch {

			select {
			case <-ctx.Done():
				return
			default:
				if filterFn(item) {
					ret <- item
				}
				continue
			}
		}

	}()

	return ret
}

func (self *Client) SearchCaseByExternalSampleId(ctx context.Context, externalSampleId string) chan *CaseResp {

	stFn := func(c *CaseResp) bool {
		for _, csj := range c.CaseSubjects {
			for _, sample := range csj.Samples {
				if sample.ExternalSampleId == externalSampleId {
					return true
				}
			}

		}

		return false
	}

	return self.GetCasesSearchBySearchTerm(ctx, externalSampleId, stFn)
}

func (self *Client) UpdateCaseStatusAndSubStatus(ctx context.Context, caseId, status, subState string) error {
	url := fmt.Sprintf(`/crs/api/v1/cases/%s/status`, caseId)
	body := fmt.Sprintf(`{
		"status":"%s",
		"subState":"%s"
	}`,
		status,
		subState,
	)
	resp, err := self.NewRequestWithContext(ctx, `POST`, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	ebody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(ebody))
	if resp.StatusCode != 200 {

		return fmt.Errorf(`wrong status code:%d`, resp.StatusCode)
	}
	return nil
}

//CloseCase close case from new
// https://jira.illumina.com/browse/OLYM-16130
//!!! in 2.6 TSS only. there is a limitation only appy to case status is New

func (self *Client) CloseCase(ctx context.Context, caseId string, body []byte) error {
	url := fmt.Sprintf(`/crs/api/v1/cases/%s/complete`, caseId)
	resp, err := self.NewRequestWithContext(ctx, `POST`, url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	ebody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(ebody))
	if resp.StatusCode != 200 {

		return fmt.Errorf(`wrong status code:%d`, resp.StatusCode)
	}
	return nil
}

func (self *Client) CloseCaseWithCompletedTime(ctx context.Context, caseId string, t time.Time) error {
	javaTimeFormatter := `2006-01-02T15:04:05.000Z`

	body := fmt.Sprintf(`
	{
		"completedDate":"%s"
	}
	`, t.Format(javaTimeFormatter))
	return self.CloseCase(ctx, caseId, []byte(body))
}

func (self *Client) CloseCaseWithCompletedTimeNow(ctx context.Context, caseId string) error {
	//TSS might be using UTC
	return self.CloseCaseWithCompletedTime(ctx, caseId, time.Now().UTC())
}
