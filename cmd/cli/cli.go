package main

import (
	"flag"
	"fmt"
	"llm-check/checker"
)

func main() {
	dir := flag.String("dir", ".", "directory to scan")
	openAiKey := flag.String("key", "", "openai api key")
	flag.Parse()

	if *openAiKey == "" {
		panic("must set openai api key")
	}

	checker := checker.NewOpenAI("gpt-4", *openAiKey)
	fails, err := checker.Check(*dir)
	if err != nil {
		panic(err)
	}

	for _, fail := range fails {
		fmt.Println(fail.String())
	}
}
