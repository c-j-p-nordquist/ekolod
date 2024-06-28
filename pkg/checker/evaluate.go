package checker

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func EvaluateCheck(check Check, response Response) CheckResult {
	if check.HTTPStatus != nil {
		result := evaluateCondition("HTTP Status", *check.HTTPStatus, response.StatusCode)
		if !result.Success {
			return result
		}
	}

	if check.ResponseTime != nil {
		result := evaluateCondition("Response Time", *check.ResponseTime, response.Duration.Seconds())
		if !result.Success {
			return result
		}
	}

	if check.ResponseBody != nil {
		result := evaluateCondition("Response Body", *check.ResponseBody, response.Body)
		if !result.Success {
			return result
		}
	}

	return CheckResult{Success: true, Message: "All checks passed"}
}

func evaluateCondition(checkType string, condition Condition, value interface{}) CheckResult {
	switch condition.Type {
	case "eq":
		return checkEquality(checkType, condition.Value, value)
	case "in":
		return checkInclusion(checkType, condition.Values, value)
	case "contains":
		return checkContains(checkType, condition.Value, value)
	case "regex":
		return checkRegex(checkType, condition.Value, value)
	case "below", "above":
		return evaluateThreshold(checkType, condition, value)
	default:
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Unknown condition type: %s", checkType, condition.Type)}
	}
}

func checkEquality(checkType string, expected, actual interface{}) CheckResult {
	if expected == actual {
		return CheckResult{Success: true, Message: fmt.Sprintf("%s: Equality check passed", checkType)}
	}
	return CheckResult{Success: false, Message: fmt.Sprintf("%s: Expected %v, got %v", checkType, expected, actual)}
}

func checkInclusion(checkType string, expectedValues []interface{}, actual interface{}) CheckResult {
	for _, v := range expectedValues {
		if v == actual {
			return CheckResult{Success: true, Message: fmt.Sprintf("%s: Inclusion check passed", checkType)}
		}
	}
	return CheckResult{Success: false, Message: fmt.Sprintf("%s: Value %v not in allowed set %v", checkType, actual, expectedValues)}
}

func checkContains(checkType string, expected, actual interface{}) CheckResult {
	strExpected, okExpected := expected.(string)
	strActual, okActual := actual.(string)
	if !okExpected || !okActual {
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Contains check requires string values", checkType)}
	}
	if strings.Contains(strActual, strExpected) {
		return CheckResult{Success: true, Message: fmt.Sprintf("%s: Contains check passed", checkType)}
	}
	return CheckResult{Success: false, Message: fmt.Sprintf("%s: String does not contain '%s'", checkType, strExpected)}
}

func checkRegex(checkType string, pattern, actual interface{}) CheckResult {
	strPattern, okPattern := pattern.(string)
	strActual, okActual := actual.(string)
	if !okPattern || !okActual {
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Regex check requires string values", checkType)}
	}
	match, err := regexp.MatchString(strPattern, strActual)
	if err != nil {
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Regex error: %v", checkType, err)}
	}
	if match {
		return CheckResult{Success: true, Message: fmt.Sprintf("%s: Regex check passed", checkType)}
	}
	return CheckResult{Success: false, Message: fmt.Sprintf("%s: String does not match pattern '%s'", checkType, strPattern)}
}

func evaluateThreshold(checkType string, condition Condition, value interface{}) CheckResult {
	thresholdValue, ok := toFloat64(condition.Value)
	if !ok {
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Invalid threshold value: %v", checkType, condition.Value)}
	}

	actualValue, ok := toFloat64(value)
	if !ok {
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Invalid actual value: %v", checkType, value)}
	}

	switch condition.Type {
	case "below":
		if actualValue < thresholdValue {
			return CheckResult{Success: true, Message: fmt.Sprintf("%s: Below threshold check passed", checkType)}
		}
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Value %v not below threshold %v", checkType, actualValue, thresholdValue)}
	case "above":
		if actualValue > thresholdValue {
			return CheckResult{Success: true, Message: fmt.Sprintf("%s: Above threshold check passed", checkType)}
		}
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Value %v not above threshold %v", checkType, actualValue, thresholdValue)}
	default:
		return CheckResult{Success: false, Message: fmt.Sprintf("%s: Unknown threshold type: %s", checkType, condition.Type)}
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch value := v.(type) {
	case float64:
		return value, true
	case float32:
		return float64(value), true
	case int:
		return float64(value), true
	case int64:
		return float64(value), true
	case time.Duration:
		return value.Seconds(), true
	default:
		return 0, false
	}
}
