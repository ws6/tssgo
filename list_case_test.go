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

func TestSearchCase(t *testing.T) {
	client := getNewClient()
	queryParam := new(CaseRequestParams)
	queryParam.SortByColumns = []string{`+createdDate`}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	ret, err := client.SearchCase(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, ret)
}
func TestListCase(t *testing.T) {
	client := getNewClient()
	queryParam := new(CaseRequestParams)
	queryParam.ShowOverDueCases = true
	queryParam.ShowUnspecifiedCases = true
	// queryParam.PageSize = 5
	queryParam.SortByColumns = []string{`+createdDate`}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ret, err := client.ListCase(ctx, queryParam)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, ret)
}
