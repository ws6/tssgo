package tssgo

import (
	"context"
	"encoding/json"
	"fmt"
)

func (self *Client) GetRerpotByCaseId(ctx context.Context, caseId string) (map[string]interface{}, error) {

	url := fmt.Sprintf(`/crs/api/v1/cases/%s/reports/json`, caseId)
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

//GetRerpotsJsonContentByCaseId only return the key - jsonContent
func (self *Client) GetRerpotJsonContentByCaseId(ctx context.Context, caseId string) (map[string]interface{}, error) {
	res, err := self.GetRerpotByCaseId(ctx, caseId)
	if err != nil {
		return nil, fmt.Errorf(`GetRerpotsByCaseId:%s`, err.Error())
	}

	if _str, ok := res[`jsonContent`]; ok {
		str, ok2 := _str.(string)
		if !ok2 {
			return nil, fmt.Errorf(`jsonContent is not string`)
		}
		ret := make(map[string]interface{})
		if err := json.Unmarshal([]byte(str), &ret); err != nil {
			return nil, fmt.Errorf(`jsonContent is not msi:%s`, err.Error())
		}
		return ret, nil
	}

	return nil, fmt.Errorf(`no jsonContent`)

}
