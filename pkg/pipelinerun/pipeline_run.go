package pipelinerun

type PipelineRun []struct {
	Class string `json:"_class"`
	Links struct {
		Self struct {
			Class string `json:"_class"`
			Href  string `json:"href"`
		} `json:"self"`
		Actions struct {
			Class string `json:"_class"`
			Href  string `json:"href"`
		} `json:"actions"`
		Steps struct {
			Class string `json:"_class"`
			Href  string `json:"href"`
		} `json:"steps"`
	} `json:"_links"`
	Actions            []interface{} `json:"actions"`
	DisplayDescription interface{}   `json:"displayDescription"`
	DisplayName        string        `json:"displayName"`
	DurationInMillis   int           `json:"durationInMillis"`
	ID                 string        `json:"id"`
	Input              interface{}   `json:"input"`
	Result             string        `json:"result"`
	StartTime          string        `json:"startTime"`
	State              string        `json:"state"`
	Type               string        `json:"type"`
	CauseOfBlockage    interface{}   `json:"causeOfBlockage"`
	Edges              []struct {
		Class string `json:"_class"`
		ID    string `json:"id"`
		Type  string `json:"type"`
	} `json:"edges"`
	FirstParent interface{} `json:"firstParent"`
	Restartable bool        `json:"restartable"`
}
