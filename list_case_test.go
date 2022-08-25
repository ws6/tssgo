package tssgo

import (
	"context"
	"encoding/json"

	"os"
	"testing"
)

func getConfig() map[string]string {
	ret := make(map[string]string)
	ret[AUTH_TOKEN] = os.Getenv("TSS_TEST_AUTH_TOKEN")
	ret[ILMN_DOMAIN] = os.Getenv("TSS_TEST_ILMN_DOMAIN")
	ret[ILMN_WORKGROUP] = os.Getenv("TSS_TEST_ILMN_WORKGROUP")
	ret[BASE_URL] = os.Getenv("TSS_TEST_BASE_URL")

	return ret
}

func getNewClient() *Client {
	ret, err := NewClient(getConfig())
	if err != nil {
		panic(err)
	}
	return ret
}

type AnalysisResp struct {
	FormatVersion string `json:"formatVersion"`
	//...more here
}

func TestGetReport(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	caseId := `5d17ee9a-d514-4e03-bfd9-de30c08fdbf6`
	res, err := client.GetRerpotJsonContentByCaseId(ctx, caseId)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("%s\n", res)

}

func _TestGetAnalysisMetrics(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	_url := `/cfs/api/v1/analysis/fwb.252d7815a24d4fb08809xxxfe/qc`
	body, err := client.GetBytes(ctx, _url)
	if err != nil {
		t.Fatal(err.Error())
	}
	ret := new(AnalysisResp)
	if err := json.Unmarshal(body, ret); err != nil {
		t.Fatal(err.Error())
	}

	t.Log(string(body))
	t.Logf("%+v\n", ret)

}

func _TestCompleteCaseStatus(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	caseId := `eb0370fd-2dd8-4f20-968d-xxx`
	err := client.CloseCaseWithCompletedTimeNow(ctx,
		caseId,
	)
	if err != nil {
		t.Fatal(err.Error())
	}

}

//GetCaseByIdWithPHI
func _TestGetCaseByIdWithPHI(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	caseId := `c7e68e6c-aae3-4a3e-a48f-xxx`
	ret, err := client.GetCaseByIdWithPHI(ctx, caseId)
	if err != nil {
		t.Fatal(err.Error())
	}
	body, _ := json.MarshalIndent(ret, "", "  ")
	t.Log(string(body))
	return
	for _, csj := range ret.CaseSubjects {
		t.Logf("%+v\n", csj) //patient information
	}

}
func _TestSearchCaseBySearchTerm(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	st := `internal-sample-id123`

	retCh := client.SearchCaseByExternalSampleId(ctx, st)
	numFound := 0
	for found := range retCh {
		numFound++
		t.Logf(`%+v`, found)
	}

	if numFound == 0 {
		t.Fatal(`not found `)
	}
	t.Log(numFound)
}

func _TestSearchCase(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	params := map[string]string{
		`externalSampleId`: `internal-sample-id123`,
	}
	ret, err := client.SearchCase(ctx, params)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, ret)
}
func _TestListCase(t *testing.T) {
	client := getNewClient()
	queryParam := new(CaseRequestParams)
	queryParam.ShowOverDueCases = true
	queryParam.ShowUnspecifiedCases = true
	queryParam.PageSize = 25
	queryParam.SortByColumns = []string{`+createdDate`}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ret, err := client.ListCase(ctx, queryParam)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, ret)

	for _, found := range ret.Content {

		t.Logf(`%+v`, found)
	}

}

func _TestSearchAuditLog(t *testing.T) {
	client := getNewClient()
	queryParam := make(map[string]string)
	queryParam[`fromDate`] = `2021-10-25T14:12:06+0000`
	queryParam[`orderBy`] = `+createdDate`

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ret, err := client.SearchAuditLog(ctx, queryParam)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, ret)

	for _, found := range ret.Content {
		t.Logf(`%+v`, found)
	}

}
