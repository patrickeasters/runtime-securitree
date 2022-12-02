package sds

type ActionResultsType struct {
	Type            string `json:"type"`
	Successful      bool   `json:"successful"`
	Token           string `json:"token"`
	SysdigCaptureID *int   `json:"sysdigCaptureId,omitempty"`
	BeforeEventNs   *int64 `json:"beforeEventNs,omitempty"`
	AfterEventNs    *int64 `json:"afterEventNs,omitempty"`
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Severity int

func (s Severity) String() string {
	switch {
	case s >= 0 && s <= 3:
		return "High"
	case s >= 4 && s <= 5:
		return "Medium"
	case s == 6:
		return "Low"
	case s == 7:
		return "Info"
	}
	return ""
}

type PolicyEvent struct {
	ID                   string               `json:"id"`
	ContainerID          string               `json:"containerId"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	Severity             Severity             `json:"severity"`
	PolicyID             int                  `json:"policyId"`
	ActionResults        []*ActionResultsType `json:"actionResults"`
	Output               string               `json:"output"`
	RuleType             string               `json:"ruleType"`
	RuleSubtype          string               `json:"ruleSubtype,omitempty"`
	MatchedOnDefault     bool                 `json:"matchedOnDefault"`
	Fields               []*KeyValuePair      `json:"fields"`
	EventLabels          []*KeyValuePair      `json:"eventLabels"`
	FalsePositive        bool                 `json:"falsePositive"`
	BaselineID           string               `json:"baselineId"`
	PolicyVersion        int                  `json:"policyVersion"`
	Origin               string               `json:"origin"`
	Timestamp            int64                `json:"timestamp"`
	TimestampNano        int64                `json:"timestampNs,omitempty"`
	TimestampRFC3339Nano string               `json:"timestampRFC3339Nano,omitempty"`
	HostMac              string               `json:"hostMac"`
	IsAggregated         bool                 `json:"isAggregated"`
	URL                  string               `json:"url,omitempty"`
}
