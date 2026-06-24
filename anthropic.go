// Package anthropic is an Anthropic (Claude) driver for togo ai using the Messages API.
// Blank-import it and set AI_DRIVER=anthropic + ANTHROPIC_API_KEY.
package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/ai"
	"github.com/togo-framework/togo"
)

const (
	endpoint     = "https://api.anthropic.com/v1/messages"
	apiVersion   = "2023-06-01"
	defaultModel = "claude-3-5-sonnet-latest"
)

func init() {
	ai.RegisterDriver("anthropic", func(k *togo.Kernel) (ai.Provider, error) {
		key := os.Getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, errors.New("ai-anthropic: ANTHROPIC_API_KEY not set")
		}
		return &provider{key: key, model: defaultModel, client: &http.Client{Timeout: 120 * time.Second}}, nil
	})
}

type provider struct {
	key, model string
	client     *http.Client
}

func (p *provider) Chat(ctx context.Context, req ai.ChatRequest) (ai.ChatResponse, error) {
	model := req.Model
	if model == "" {
		model = p.model
	}
	maxTok := req.MaxTokens
	if maxTok == 0 {
		maxTok = 1024
	}
	var system string
	var msgs []map[string]string
	for _, m := range req.Messages {
		if m.Role == ai.RoleSystem || m.Role == "system" {
			system += m.Content + "\n"
			continue
		}
		msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
	}
	body := map[string]any{"model": model, "max_tokens": maxTok, "messages": msgs}
	if strings.TrimSpace(system) != "" {
		body["system"] = strings.TrimSpace(system)
	}
	if req.Temperature != 0 {
		body["temperature"] = req.Temperature
	}
	buf, _ := json.Marshal(body)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(buf))
	if err != nil {
		return ai.ChatResponse{}, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("x-api-key", p.key)
	r.Header.Set("anthropic-version", apiVersion)
	resp, err := p.client.Do(r)
	if err != nil {
		return ai.ChatResponse{}, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return ai.ChatResponse{}, fmt.Errorf("ai-anthropic: %s: %s", resp.Status, string(data))
	}
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model string `json:"model"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return ai.ChatResponse{}, err
	}
	var sb strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	return ai.ChatResponse{
		Content: sb.String(),
		Model:   out.Model,
		Usage:   ai.Usage{PromptTokens: out.Usage.InputTokens, CompletionTokens: out.Usage.OutputTokens, TotalTokens: out.Usage.InputTokens + out.Usage.OutputTokens},
	}, nil
}

func (p *provider) Embed(ctx context.Context, req ai.EmbedRequest) (ai.EmbedResponse, error) {
	return ai.EmbedResponse{}, errors.New("ai-anthropic: Anthropic has no embeddings API — use ai-openai (or ai-voyage) for embeddings")
}
