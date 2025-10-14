package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/models/entity"
	"google.golang.org/genai"
)

var GemniClient *genai.Client

func Init(ctx context.Context, cfg *config.Config) error {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GeminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}

	GemniClient = client
	return nil
}

func buildRubricPrompt(r entity.EvaluationRubrics, cv *dto.CVResponse) string {
	var sb strings.Builder

	sb.WriteString("You are a senior technical recruiter.\n")
	sb.WriteString("Evaluate the candidate CV based on the following rubrics and return JSON ONLY in this format:\n")
	sb.WriteString(`{
  "cv_scores": {...},
  "project_scores": {...},
  "weighted_cv_score": float between 0 and 1,
  "weighted_project_score": float between 0 and 5,
  "overall_summary": string,
  "feedback": string
}

--- CV RUBRICS ---
`)

	for _, item := range r.CV {
		sb.WriteString(fmt.Sprintf("- %s (Weight: %.0f%%): %s\n  Scale: %s\n",
			item.Name, item.Weight*100, item.Description, item.Scale))
	}

	sb.WriteString("\n--- PROJECT RUBRICS ---\n")
	for _, item := range r.Project {
		sb.WriteString(fmt.Sprintf("- %s (Weight: %.0f%%): %s\n  Scale: %s\n",
			item.Name, item.Weight*100, item.Description, item.Scale))
	}

	sb.WriteString("\nReturn only valid JSON — no markdown or explanations outside the JSON.\n")

	sb.WriteString(fmt.Sprintf("\nCV Title: %s\n", cv.Title))
	sb.WriteString(fmt.Sprintf("CV Summary: %s\n", cv.Summary))
	sb.WriteString(fmt.Sprintf("File Path (reference only): %s\n", cv.FilePath))

	return sb.String()
}

func GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	fmt.Printf("\ntext: %v\n", text)
	if GemniClient == nil {
		return nil, fmt.Errorf("gemini client not initialized")
	}

	model := "text-embedding-004"
	outD := int32(768)

	result, err := GemniClient.Models.EmbedContent(ctx,
		model,
		genai.Text(text),
		&genai.EmbedContentConfig{OutputDimensionality: &outD},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if result == nil || len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("empty embedding response or values")
	}

	return result.Embeddings[0].Values, nil
}

func EvaluateCV(ctx context.Context, cv *dto.CVResponse) (*dto.CVEvaluationResponse, error) {
	if GemniClient == nil {
		return nil, fmt.Errorf("gemini client not initialized")
	}
	rubrics := entity.NewDefaultRubrics()
	prompt := buildRubricPrompt(rubrics, cv)

	model := "gemini-2.0-flash"

	resp, err := GemniClient.Models.GenerateContent(ctx, model, genai.Text(prompt),
		&genai.GenerateContentConfig{},
	)
	if err != nil {
		return nil, fmt.Errorf("gemini content generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Extract response text
	content := resp.Candidates[0].Content
	if len(content.Parts) == 0 {
		return nil, fmt.Errorf("empty content from Gemini")
	}

	text := content.Parts[0].Text
	// clean texts
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var eval dto.CVEvaluationResponse
	if err := json.Unmarshal([]byte(text), &eval); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini JSON: %w\nraw response: %s", err, text)
	}

	return &eval, nil
}
