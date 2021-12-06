package tssgo

//case_client.go represents a client model of the case

type ClientResp struct {
	Id          string             `json:"id"`
	Domain      string             `json:"domain"`
	WorkGroupId string             `json:"workGroupId"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	Email       string             `json:"email"`
	Status      string             `json:"status"`
	Addresses   []*AddressItemResp `json:"addresses"`
}

type AddressItemResp struct {
	Id              string `json:"id"`
	InstitutionName string `json:"institutionName"`
	Phone           string `json:"phone"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	City            string `json:"city"`
	Region          string `json:"region"`
}
