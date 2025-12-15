package main

import (
	"bytes"
	"os"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	apiURL = "http://localhost:8080/api/v1/traces"
	apiKey = "demo-key-456"
)

// Models and providers with realistic pricing
var models = []Model{
	{"gpt-4", "openai", 0.03, 0.06},
	{"gpt-4-turbo", "openai", 0.01, 0.03},
	{"gpt-3.5-turbo", "openai", 0.0005, 0.0015},
	{"claude-3-opus", "anthropic", 0.015, 0.075},
	{"claude-3-sonnet", "anthropic", 0.003, 0.015},
	{"claude-3-haiku", "anthropic", 0.00025, 0.00125},
}

var prompts = []string{
	"Summarize this document for me",
	"Write a product description",
	"Generate a blog post about AI",
	"Translate this text to Spanish",
	"Debug this code snippet",
	"Explain quantum computing simply",
	"Create a marketing email",
	"Analyze customer feedback",
	"Generate SQL query from description",
	"Write unit tests for this function",
}

var organizations = []string{"org-acme", "org-techcorp", "org-startup", "org-enterprise"}
var projects = []string{"proj-prod", "proj-staging", "proj-dev", "proj-test"}

type Model struct {
	Name     string
	Provider string
	PromptPrice float64
	CompletionPrice float64
}

type TraceRequest struct {
	OrganizationID string        `json:"organization_id"`
	ProjectID      string        `json:"project_id"`
	Model          string        `json:"model"`
	Provider       string        `json:"provider"`
	TraceType      string        `json:"trace_type"`
	Spans          []SpanRequest `json:"spans"`
}

type SpanRequest struct {
	Name             string            `json:"name"`
	Model            string            `json:"model"`
	Provider         string            `json:"provider"`
	Input            string            `json:"input"`
	Output           string            `json:"output"`
	PromptTokens     int               `json:"prompt_tokens"`
	CompletionTokens int               `json:"completion_tokens"`
	DurationMs       int               `json:"duration_ms"`
	Status           string            `json:"status"`
	Tags             map[string]string `json:"tags,omitempty"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	count := 100 // Default: 100 traces
	
	// Check command line args
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &count)
	}

	log.Printf("🎲 Generating %d sample traces...\n", count)
	log.Println("─────────────────────────────────────────")

	successCount := 0
	errorCount := 0

	startTime := time.Now()

	for i := 0; i < count; i++ {
		trace := generateRandomTrace()
		
		if err := sendTrace(trace); err != nil {
			errorCount++
			if errorCount <= 5 { // Only show first 5 errors
				log.Printf("❌ Error creating trace %d: %v", i+1, err)
			}
		} else {
			successCount++
			if (i+1) % 10 == 0 {
				log.Printf("✅ Created %d/%d traces...", i+1, count)
			}
		}

		// Small delay to avoid overwhelming the API
		time.Sleep(50 * time.Millisecond)
	}

	duration := time.Since(startTime)

	log.Println("─────────────────────────────────────────")
	log.Printf("🎉 Generation complete!")
	log.Printf("✅ Success: %d", successCount)
	log.Printf("❌ Failed: %d", errorCount)
	log.Printf("⏱️  Duration: %s", duration)
	log.Printf("📊 Rate: %.1f traces/sec", float64(successCount)/duration.Seconds())
}

func generateRandomTrace() *TraceRequest {
	model := models[rand.Intn(len(models))]
	org := organizations[rand.Intn(len(organizations))]
	project := projects[rand.Intn(len(projects))]
	prompt := prompts[rand.Intn(len(prompts))]

	// Realistic token counts
	promptTokens := rand.Intn(500) + 50      // 50-550 tokens
	completionTokens := rand.Intn(1000) + 100 // 100-1100 tokens

	// Realistic duration based on tokens
	baseDuration := 200
	tokenDuration := (promptTokens + completionTokens) / 2
	duration := baseDuration + tokenDuration + rand.Intn(500)

	// 95% success rate
	status := "success"
	if rand.Float64() > 0.95 {
		if rand.Float64() > 0.5 {
			status = "error"
		} else {
			status = "timeout"
		}
	}

	return &TraceRequest{
		OrganizationID: org,
		ProjectID:      project,
		Model:          model.Name,
		Provider:       model.Provider,
		TraceType:      "single_call",
		Spans: []SpanRequest{
			{
				Name:             "llm_completion",
				Model:            model.Name,
				Provider:         model.Provider,
				Input:            prompt,
				Output:           generateOutput(completionTokens),
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				DurationMs:       duration,
				Status:           status,
				Tags: map[string]string{
					"environment": randomEnv(),
					"user_type":   randomUserType(),
				},
			},
		},
	}
}

func generateOutput(tokens int) string {
	words := tokens / 4 // Rough estimate: 4 tokens per word
	output := "Generated response with "
	for i := 0; i < words && i < 20; i++ { // Limit to 20 words for brevity
		output += "word "
	}
	return output + "..."
}

func randomEnv() string {
	envs := []string{"production", "staging", "development"}
	return envs[rand.Intn(len(envs))]
}

func randomUserType() string {
	types := []string{"free", "pro", "enterprise"}
	return types[rand.Intn(len(types))]
}

func sendTrace(trace *TraceRequest) error {
	jsonData, err := json.Marshal(trace)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}
