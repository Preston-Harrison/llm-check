package checker

import (
	"fmt"
	"llm-check/common"
	"llm-check/go_checks"
	"llm-check/lang"
)

type CheckFail struct {
	Origin common.FunctionWithCheck
	FailFn common.Function
	Reason string
}

func (cf CheckFail) String() string {
    return fmt.Sprintf("Function %s failed at %s because %s", cf.Origin.Name(), cf.FailFn.Name(), cf.Reason)
}

type Checker struct {
	llm llm
}

func NewOpenAI(model string, apiKey string) *Checker {
    return &Checker{
        llm: NewClient(apiKey, "https://api.openai.com/v1/chat/completions", model),
    }
}

func (c *Checker) Check(dir string) ([]CheckFail, error) {
	language, err := lang.Detect(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect language: %w", err)
	}
	switch language {
	case lang.GO:
		proj, err := go_checks.NewProject(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to create project: %w", err)
		}

		functions := go_checks.GetFnsWithCheck(proj)
		fails, err := check(c.llm, functions)
		if err != nil {
			return nil, fmt.Errorf("failed to check: %w", err)
		}

		return fails, nil
	default:
		return nil, fmt.Errorf("unable to detect language")
	}
}
