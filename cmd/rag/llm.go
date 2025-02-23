package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/philippgille/chromem-go"
	"github.com/sashabaranov/go-openai"
)

const (
	ollamaBaseURL = "http://localhost:11434/v1"
	llmModel      = "gemma2:2b"
)

func rag(md []chromem.Result, q string) {

	aiClient := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    ollamaBaseURL,
		HTTPClient: http.DefaultClient,
	})

	systemPrompt := fmt.Sprintf(`You are a helpful research assistant tasked with answering questions
	relating to knowlege bases you have access to.
	
	Answer the question in a factual manner strictly based only on your research findings below:
%s
	Do not extrapolate. If the data is not available, respond with "I don't know".
	`, relevantDocs(md))

	question := "Question: " + strings.TrimPrefix(q, "search_query: ")
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: question,
		},
	}

	fmt.Println(systemPrompt)
	fmt.Println("Asking gemma2:2b on a local machine without a GPU. This will take a minute...")
	fmt.Println(question)

	ctx := context.Background()
	res, err := aiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    llmModel,
		Messages: msgs,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.TrimSpace(res.Choices[0].Message.Content))
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
