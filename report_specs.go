package tssgo

//start definition of specs from report API
type Report struct {
	ExternalCaseId                      interface{}            `json:"externalCaseId"`
	Client                              map[string]interface{} `json:"client"` //hide it regardless
	Tags                                []string               `json:"tags"`
	CaseId                              string                 `json:"id"` //this is actually caseId
	DisplayId                           string                 `json:"displayId"`
	UpdatedDate                         string                 `json:"updatedDate"`
	AdditionalInstructionsForLaboratory string                 `json:"additionalInstructionsForLaboratory"`
	TestDefinition                      map[string]interface{} `json:"testDefinition"`
	Tat                                 string                 `json:"tat"`
	StartDate                           string                 `json:"startDate"`
	DueDate                             string                 `json:"dueDate"`
	CustomMetadata                      map[string]interface{} `json:"customMetadata"`
	CaseOwners                          []string               `json:"caseOwners"`
	CreatedDate                         string                 `json:"createdDate"`
	ClientId                            string                 `json:"clientId"`
	ClientAddressId                     string                 `json:"clientAddressId"`
	CreatedBy                           string                 `json:"createdBy"`
	UpdatedBy                           string                 `json:"updatedBy"`
	CaseSubjects                        []CaseSubject          `json:"caseSubjects"`
	TruSightSoftwareSuiteVersion        string                 `json:"truSightSoftwareSuiteVersion"`
	DragenVersion                       string                 `json:"dragenVersion"`
	DataSourceVersions                  string                 `json:"dataSourceVersions"`
	Status                              string                 `json:"status"`
	SubState                            string                 `json:"subState"`
	VariantIndexId                      interface{}            `json:"variantIndexId"`
	Phi                                 map[string]interface{} `json:"phi"`
	ActivationState                     string                 `json:"activationState"`
	SelectedReportId                    interface{}            `json:"selectedReportId"`
	ReportDataVersion                   interface{}            `json:"reportDataVersion"`
	PedigreeSize                        interface{}            `json:"pedigreeSize"`
	TestType                            string                 `json:"testType"`
	CompletedDate                       interface{}            `json:"completedDate"`
}

type CaseSubject struct {
	Id                    string `json:"id"`
	RelationshipToProband string `json:"relationshipToProband"`

	Samples []struct {
		ExternalSampleId   string `json:"externalSampleId"`
		SampleName         string `json:"sampleName"`
		ExternalSampleName string `json:"externalSampleName"`
		ReportTypes        []struct {
			Id            string `json:"id"`
			ReportType    string `json:"reportType"`
			ReportTypeId  string `json:"reportTypeId"`
			ReportDetails struct {
				Id               string `json:"id"`
				CreationDate     string `json:"creationDate"`
				UpdateDate       string `json:"updateDate"`
				ReportId         string `json:"reportId"`
				Status           string `json:"status"`
				SubState         string `json:"subState"`
				CaseId           string `json:"caseId"`
				CaseDisplayId    string `json:"caseDisplayId"`
				TestDefinitionId string `json:"testDefinitionId"`
				Name             string `json:"name"`       //actual report definition - "Secondary Findings" or "RUGD"
				EditStatus       string `json:"editStatus"` //actual report type, Corrected, Amended, null
				Variants         []Variant
			} `json:"reportDetails"`
		} `json:"reportTypes"`
	} `json:"samples"`
}

type Variant struct {
	VariantId  string `json:"variantId"`
	Category   string `json:"category"`
	Gene       string `json:"gene"`
	Chromosome string `json:"chromosome"`
	Start      int64  `json:"start"`
	Stop       int64  `json:"stop"`
	RefAllele  string `json:"refAllele"`
	AltAllele  string `json:"altAllele"`
	Transcript struct {
		Name         string `json:"name"` //NM name. nucleotides change; most important unique naming
		Hgvsc        string `json:"hgvsc"`
		Hgvsp        string `json:"hgvsp"`
		Exons        string `json:"exons"`
		ProteinId    string `json:"proteinId"`
		Consequences []struct {
			Label       string `json:"label"`
			Consequence string `json:"consequence"`
		} `json:"consequences"`
	} `json:"transcript"`

	Zygosity string `json:"variantId"`
}
