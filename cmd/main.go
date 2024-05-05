package main

import (
	"flag"
	"fmt"
	"llm-check/checker"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	openaiApiKey := flag.String("key", "", "openai api key")
	flag.Parse()

	if *openaiApiKey == "" {
		panic("must set openai api key")
	}

	checker := checker.NewOpenAI("gpt-4", *openaiApiKey)
	fails, err := checker.Check(*dir)
	if err != nil {
		panic(err)
	}

	for _, fail := range fails {
		fmt.Println(fail.String())
	}
}
