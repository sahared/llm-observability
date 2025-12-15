package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Aditya-Pimpalkar/clarity/internal/models"
	"github.com/Aditya-Pimpalkar/clarity/internal/repository"
	"github.com/Aditya-Pimpalkar/clarity/internal/services"
)

func main() {
	// This is an example - not meant to be run directly
	// Shows how to use the trace service

	// Initialize repository
	repo, err := repository.NewClickHouseRepository("localhost:9000")
	if err != nil {
		log.Fatal(err)
	}

	// Create trace service
	traceService := services.NewTraceService(repo)

	// Create a sample trace request
	req := &models.TraceRequest{
		OrganizationID: "org-example",
		ProjectID:      "proj-example",
		Model:          "gpt-4",
		Provider:       "openai",
		TraceType:      "single_call",
		Spans: []models.SpanRequest{
			{
				Name:             "chat_completion",
				Model:            "gpt-4",
				Provider:         "openai",
				Input:            "What is the meaning of life?",
				Output:           "The meaning of life is a philosophical question...",
				PromptTokens:     10,
				CompletionTokens: 50,
				DurationMs:       1500,
				Status:           "success",
				Tags:             map[string]string{"environment": "production"},
			},
		},
	}

	// Create the trace
	ctx := context.Background()
	response, err := traceService.CreateTrace(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Trace created: %s\n", response.TraceID)
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Timestamp: %s\n", response.CreatedAt)
}

func exampleBatchTraces() {
	// Example of creating multiple traces
	repo, _ := repository.NewClickHouseRepository("localhost:9000")
	traceService := services.NewTraceService(repo)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		req := &models.TraceRequest{
			OrganizationID: "org-example",
			ProjectID:      "proj-example",
			Model:          "gpt-4",
			Provider:       "openai",
			TraceType:      "single_call",
			Spans: []models.SpanRequest{
				{
					Name:             fmt.Sprintf("request_%d", i),
					Model:            "gpt-4",
					Provider:         "openai",
					Input:            fmt.Sprintf("Question %d", i),
					Output:           fmt.Sprintf("Answer %d", i),
					PromptTokens:     int(50 + i*10),
					CompletionTokens: int(30 + i*5),
					DurationMs:       int64(150 + i*50),
					Status:           "success",
					Tags:             map[string]string{"batch": "example"},
				},
			},
		}

		_, err := traceService.CreateTrace(ctx, req)
		if err != nil {
			log.Printf("Error creating trace %d: %v", i, err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
