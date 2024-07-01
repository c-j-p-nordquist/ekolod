package checkconverter

import (
	"github.com/c-j-p-nordquist/ekolod/pkg/checker"
	"github.com/c-j-p-nordquist/ekolod/pkg/config"
)

func ConvertConfigCheckToCheckerCheck(configCheck config.Check) checker.Check {
	return checker.Check{
		Path:         configCheck.Path,
		HTTPStatus:   convertCondition(configCheck.HTTPStatus),
		ResponseTime: convertCondition(configCheck.ResponseTime),
		ResponseBody: convertCondition(configCheck.ResponseBody),
	}
}

func convertCondition(configCondition *config.Condition) *checker.Condition {
	if configCondition == nil {
		return nil
	}
	return &checker.Condition{
		Type:   configCondition.Type,
		Value:  configCondition.Value,
		Values: configCondition.Values,
	}
}
