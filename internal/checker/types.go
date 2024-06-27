package checker

import "time"

type Check struct {
	Path         string
	HTTPStatus   *Condition
	ResponseTime *Condition
	ResponseBody *Condition
}

type Condition struct {
	Type      string
	Value     interface{}
	Values    []interface{}
	Threshold *Threshold
}

type Threshold struct {
	Type  string
	Value interface{}
}

type Response struct {
	StatusCode int
	Body       string
	Duration   time.Duration
}

type CheckResult struct {
	Success bool
	Message string
}
