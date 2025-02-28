package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/philippgille/chromem-go"
	"github.com/siuyin/dflt"
)

const (
	ollamaBaseURL = "http://localhost:11434/v1"
	// llmModel      = "gemma2:2b"
)

var llmModel string

func init() {
	llmModel = dflt.EnvString("MODEL", "gemma2:2b")
}

func rag(md []chromem.Result, q string) {

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	question := "Question: " + strings.TrimPrefix(q, "search_query: ")
	prompt := fmt.Sprintf(`You are a helpful research assistant tasked with answering questions
	relating to knowlege bases you have access to.
	
	Answer the question in a factual manner strictly based only on your research findings below:
%s
	Do not extrapolate. If the data is not available, respond with "I don't know".
	
%s
	`, relevantDocs(md), question)

	req := &api.GenerateRequest{
		Model:  llmModel,
		Prompt: prompt,
	}
	fmt.Println("------------------------")
	fmt.Println(prompt)
	fmt.Println("Asking " + llmModel + " on a local machine without a GPU. This will take a minute...")

	ctx := context.Background()
	if err := client.Generate(ctx, req, respFunc); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

func relevantDocs(md []chromem.Result) string {
	docs := ""
	for _, d := range md {
		if d.Similarity > 0.6 {
			docs += d.Content + "\n"
		}
	}
	return docs
}

func respFunc(resp api.GenerateResponse) error {
	fmt.Print(resp.Response)
	return nil
}
