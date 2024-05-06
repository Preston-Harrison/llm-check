package go_checks_test

import (
	"llm-check/comments"
	"llm-check/go_checks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoFunction(t *testing.T) {
	src := `package baz

// llm-check: Make sure this always prints hello world
func bar() {
    println("hello world")
}

func bar2() {
    bar()
}

func bar3() {
    println("hello mars")
    bar()
}

func bar4() {
    println("hello mars")
}
`
	proj, err := go_checks.NewProjectFromSource(src)
	if err != nil {
		t.Fatalf("failed to create go project: %v", err)
	}

	checks := go_checks.GetFnsWithCheck(proj)

	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}

	check := checks[0]
	if check.Name() != "bar" {
		t.Fatalf("Expected check name to be 'bar', was '%s'", check.Name())
	}
	if check.Module() != "baz" {
		t.Fatalf("expected module to be 'baz', was '%s'", check.Module())
	}

    start, end := check.Location()
    if start.Line != 4 || end.Line != 6 {
        t.Fatalf("Expected start line to be 4 and end line to be 6, got %d and %d", start.Line, end.Line)
    }

	comments := comments.ParseChecksFromComments(check.Docs())
	if len(comments) != 1 {
		t.Fatalf("Expected check comment on function 'bar'")
	}
	comment := comments[0]
	if comment.Content != "Make sure this always prints hello world" {
		t.Fatalf("Did not receive correct comment on 'bar', received '%s'", comment)
	}

	callSites := check.CallSites()
	if len(callSites) != 2 {
		t.Fatalf("Expected call sites len to be 1, was %d", len(callSites))
	}

	assert.ElementsMatch(
        t, 
        []string{"bar2", "bar3"}, 
        []string{callSites[0].Name(), callSites[1].Name()},
    )
}

func TestRealModule(t *testing.T) {
    t.Skip("This changes too much for now")
	proj, err := go_checks.NewProject("../sample")
	if err != nil {
		t.Fatalf("failed to create go project: %v", err)
	}

	checks := go_checks.GetFnsWithCheck(proj)
	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}

	assert.ElementsMatch(
		t,
		[]string{"openFileForWriting", "MyExportedFunc"},
		[]string{checks[0].Name(), checks[1].Name()},
	)

    for _, check := range checks {
        if check.Name() == "openFileForWriting" {
            start, end := check.Location()
            if start.Line != 10 || end.Line != 13 {
                t.Fatalf("Expected start line to be 10 and end line to be 13, got %d and %d", start.Line, end.Line)
            }
        }
    }
}
