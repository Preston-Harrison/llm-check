package comments

import (
	"fmt"
	"llm-check/common"
	"strings"
)

const FUNCTION_SCOPE string = "function"
const CALLER_SCOPE string = "caller"

func ParseCheckComment(comment string) (common.Check, error) {
	if !strings.HasPrefix(comment, "// llm-check") {
		return common.Check{}, fmt.Errorf("no leading '// llm-check'")
	}
	trimmed := strings.TrimPrefix(comment, "// llm-check")
	colonIx := strings.Index(trimmed, ":")
	if colonIx == -1 {
		return common.Check{}, fmt.Errorf("no colon in llm-check comment")
	}

	scopeStr := trimmed[0:colonIx]
	var scope string
	switch scopeStr {
	case "(caller)":
		scope = CALLER_SCOPE
	case "":
		scope = FUNCTION_SCOPE
	case "(function)":
		scope = FUNCTION_SCOPE
	default:
		return common.Check{}, fmt.Errorf("invalid scope, expected 'caller', got '%s'", scope)
	}

	return common.Check{
		Content: strings.TrimSpace(trimmed[colonIx+1:]),
		Scope:   scope,
	}, nil

}

func ParseChecksFromComments(comments []string) []common.Check {
	var checks []common.Check
	for _, comment := range comments {
		check, err := ParseCheckComment(comment)
		if err != nil {
			continue
		}
		checks = append(checks, check)
	}
	return checks
}
