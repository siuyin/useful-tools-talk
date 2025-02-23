package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		return nil
	}

	for {
		// By default, GenerateRequest is streaming.
		req := &api.GenerateRequest{
			Model:  "gemma2:2b",
			Prompt: getQuery(),
		}

		if err := client.Generate(ctx, req, respFunc); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
}
func getQuery() string {
	fmt.Println("\n\n\nEnter your question:")
	r := bufio.NewReader(os.Stdin)
	q, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return q
}
