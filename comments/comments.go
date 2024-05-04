package comments

import (
	"fmt"
	"strings"
)

const FUNCTION_SCOPE scope = "function"
const CALLER_SCOPE scope = "caller"

type scope string

type Check struct {
	Content string
	Scope   scope
}

func ParseCheckComment(comment string) (Check, error) {
	if !strings.HasPrefix(comment, "// llm-check") {
		return Check{}, fmt.Errorf("no leading '// llm-check'")
	}
	trimmed := strings.TrimPrefix(comment, "// llm-check")
	colonIx := strings.Index(trimmed, ":")
	if colonIx == -1 {
		return Check{}, fmt.Errorf("no colon in llm-check comment")
	}

	scopeStr := trimmed[0:colonIx]
	var scope scope
	switch scopeStr {
	case "(caller)":
		scope = CALLER_SCOPE
	case "":
		scope = FUNCTION_SCOPE
	case "(function)":
		scope = FUNCTION_SCOPE
	default:
		return Check{}, fmt.Errorf("invalid scope, expected 'caller', got '%s'", scope)
	}

	return Check{
		Content: strings.TrimSpace(trimmed[colonIx+1:]),
		Scope:   scope,
	}, nil

}

func ParseChecksFromComments(comments []string) []Check {
	var checks []Check
	for _, comment := range comments {
		check, err := ParseCheckComment(comment)
		if err != nil {
			continue
		}
		checks = append(checks, check)
	}
	return checks
}
