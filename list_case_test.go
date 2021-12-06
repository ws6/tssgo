package tssgo

import (
	"context"

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

//GetCaseByIdWithPHI
func TestGetCaseByIdWithPHI(t *testing.T) {
	client := getNewClient()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	caseId := `a0c26f99-8bfb-42f7-97c6-9348babc26b9`
	ret, err := client.GetCaseByIdWithPHI(ctx, caseId)
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, csj := range ret.CaseSubjects {
		t.Logf("%+v\n", csj.Subject) //patient information
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
	queryParam[`orderBy`] = `-createdDate`

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
