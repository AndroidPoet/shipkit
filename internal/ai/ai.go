package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AndroidPoet/shipkit/internal/doctor"
	"github.com/AndroidPoet/shipkit/internal/runner"
)

type openAIRequest struct {
	Model string        `json:"model"`
	Input []interface{} `json:"input"`
}

type openAIResponse struct {
	OutputText string `json:"output_text"`
}

func Run(ctx context.Context, r runner.Runner, stdout io.Writer) error {
	status := summarize(ctx, r)
	fmt.Fprintln(stdout, "Shipkit AI")
	fmt.Fprintln(stdout, status)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "AI guidance is optional. Set OPENAI_API_KEY to get a tailored release plan.")
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Local recommendation:")
		fmt.Fprintln(stdout, "- Run `shipkit install` if any provider tool is missing.")
		fmt.Fprintln(stdout, "- Run provider auth commands after install: `gpc setup --auto`, `rc login`, `asc auth login`.")
		fmt.Fprintln(stdout, "- Run `shipkit ci github` once the local setup is green.")
		return nil
	}

	advice, err := askOpenAI(ctx, apiKey, status)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, advice)
	return nil
}

func summarize(ctx context.Context, r runner.Runner) string {
	var builder strings.Builder
	builder.WriteString("Local release status:\n")
	for _, result := range doctor.Check(ctx, r) {
		if result.Ready {
			fmt.Fprintf(&builder, "- %s: installed at %s\n", result.Tool.Name, result.Message)
		} else {
			fmt.Fprintf(&builder, "- %s: missing executable %s\n", result.Tool.Name, result.Tool.Executable)
		}
	}
	if _, err := os.Stat(".shipkit.yaml"); err == nil {
		builder.WriteString("- Shipkit config: present\n")
	} else {
		builder.WriteString("- Shipkit config: missing\n")
	}
	return strings.TrimSpace(builder.String())
}

func askOpenAI(ctx context.Context, apiKey, status string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	body := openAIRequest{
		Model: "gpt-4.1-mini",
		Input: []interface{}{
			map[string]string{
				"role":    "system",
				"content": "You are Shipkit's concise mobile release copilot. Give exact next commands. Keep it under 10 lines.",
			},
			map[string]string{
				"role":    "user",
				"content": status,
			},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("OpenAI request failed: %s", strings.TrimSpace(string(responseBody)))
	}

	var decoded openAIResponse
	if err := json.Unmarshal(responseBody, &decoded); err != nil {
		return "", err
	}
	if strings.TrimSpace(decoded.OutputText) == "" {
		return "", fmt.Errorf("OpenAI response did not include guidance")
	}
	return strings.TrimSpace(decoded.OutputText), nil
}
