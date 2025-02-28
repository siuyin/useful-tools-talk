package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/siuyin/dflt"
)

// const llmModel = "gemma2:2b"
// const llmModel = "deepseek-r1:1.5b"
var llmModel string

func init() {
	llmModel = dflt.EnvString("MODEL", "gemma2:2b")
}

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for {
		// GenerateRequest streams via respFunc.
		req := &api.GenerateRequest{
			Model:  llmModel,
			Prompt: getQuery(),
		}

		if err := client.Generate(ctx, req, respFunc); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
}

func getQuery() string {
	fmt.Printf("\n\n\nEnter your question below for %s :\n", llmModel)
	r := bufio.NewReader(os.Stdin)
	q, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return q
}
func respFunc(res api.GenerateResponse) error {
	fmt.Print(res.Response)
	return nil
}
