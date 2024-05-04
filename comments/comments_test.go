package comments_test

import (
	"llm-check/comments"
	"testing"
)

func TestParseComments(t *testing.T) {
    t.Run("valid comment for caller", func(t *testing.T) {
        comment := "// llm-check(caller): some content"
        check, err := comments.ParseCheckComment(comment)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if check.Scope != comments.CALLER_SCOPE {
            t.Errorf("unexpected scope: %s", check.Scope)
        }
        if check.Content != "some content" {
            t.Errorf("unexpected content: %s", check.Content)
        }
    })

    t.Run("valid comment for caller", func(t *testing.T) {
        comment := "// llm-check(function): some other content   "
        check, err := comments.ParseCheckComment(comment)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if check.Scope != comments.FUNCTION_SCOPE {
            t.Errorf("unexpected scope: %s", check.Scope)
        }
        if check.Content != "some other content" {
            t.Errorf("unexpected content: %s", check.Content)
        }
    })

    t.Run("valid comment with default scope", func(t *testing.T) {
        comment := "// llm-check: some other content   "
        check, err := comments.ParseCheckComment(comment)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if check.Scope != comments.FUNCTION_SCOPE {
            t.Errorf("unexpected scope: %s", check.Scope)
        }
        if check.Content != "some other content" {
            t.Errorf("unexpected content: %s", check.Content)
        }
    })

    t.Run("invalid comment", func(t *testing.T) {
        comment := "// llm-check(something weird): some content"
        _, err := comments.ParseCheckComment(comment)
        if err == nil {
            t.Fatal("expected error, got nil")
        }
    })
}
