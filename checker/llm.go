package checker

import (
	"fmt"
	"llm-check/comments"
	"llm-check/common"
	"strings"
	"sync"
)

type llm interface {
	call([]Message) ([]ResponseChoice, error)
}

func check(llm llm, funcs []common.FunctionWithCheck) ([]CheckFail, error) {
	var wg sync.WaitGroup
	errs := make(chan error, len(funcs))
	fails := make(chan CheckFail, len(funcs))

	for _, fn := range funcs {
		for _, callSite := range fn.CallSites() {
			wg.Add(1)
			go func(fn common.FunctionWithCheck, callSite common.Function) {
				defer wg.Done()
				if fn.Check().Scope != comments.CALLER_SCOPE {
					errs <- fmt.Errorf("function %s does not have a caller scope", fn.Check().Content)
					return
				}
				messages := []Message{
					{Role: "system", Content: sysPrompt},
					{Role: "user", Content: getCallSitePrompt(fn, callSite, fn.Check().Content)},
				}
				res, err := llm.call(messages)
				if err != nil {
					errs <- err
					return
				}
				response := res[0].Message.Content
				fail, ok := parseResponse(response)
				if ok {
					fails <- CheckFail{Origin: fn, FailFn: callSite, Reason: fail.Reason}
				}

			}(fn, callSite)
		}
	}

	wg.Wait()
	close(errs)
	close(fails)

	for err := range errs {
		return nil, err
	}

	var result []CheckFail
	for fail := range fails {
		result = append(result, fail)
	}

	return result, nil
}

const sysPrompt = `
You are an AI assistant that checks whether a function is implemented correctly
according to the provided specification. You will be given a function implementation
and call site and must determine if it satisfies a set of requirements.

You response must fit exactly one of the following formats:

If the function satisfies the specification:
"The function is correct."

If the function does not satisfy the specification:
"This function is incorrect.", followed by a reason why the function is does not
satisfy the specification.

If you do not have enough information to determine if the function is correct:
"I cannot determine if the function is correct.", followed by a reason why you
cannot determine if the function is correct.

The user will now provide you with a set of requirements, a function implementation, 
and a call site.
`

const callSitePrompt = `
Here is the function that is being called:
{{ %s }}

The caller of this function must satisfy the following requirements:
- {{ %s }}

Here is the function call site:
{{ %s }}

Please determine if the function call site satisfies the requirements, using the words
correct or incorrect in your response.
`

func getCallSitePrompt(fn, callSite common.Function, check string) string {
	return fmt.Sprintf(callSitePrompt, fn.String(), check, callSite.String())
}

func parseResponse(response string) (CheckFail, bool) {
	if strings.Contains(response, "incorrect") {
		reason := strings.TrimPrefix(response, "This function is incorrect.")
		reason = strings.TrimSpace(reason)
		return CheckFail{Reason: reason}, true
	}
	if !strings.Contains(response, "correct") {
		fmt.Printf("Invalid response: %s\n", response)
	}
	return CheckFail{}, false
}
