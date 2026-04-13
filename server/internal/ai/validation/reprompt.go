package validation

import (
	"context"
	"fmt"

	"github.com/meshery/meshery/server/internal/ai"
	aicontext "github.com/meshery/meshery/server/internal/ai/context"
	patternV1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
)

const (
	// MaxRepromptAttempts is the maximum number of times the re-prompt loop
	// will ask the LLM to fix validation errors. Capped at 2 to limit latency
	// and token usage.
	MaxRepromptAttempts = 2
)

// RepromptLoop orchestrates the generate → validate → re-prompt cycle.
// It calls the LLM, validates the output, and if validation fails sends the
// errors back to the LLM to request a fix. It returns the best PatternFile
// produced across all attempts (even if imperfect) along with any final errors.
type RepromptLoop struct {
	service   *ai.Service
	validator *Validator
}

// NewRepromptLoop creates a RepromptLoop that uses the given AI service and
// validator.
func NewRepromptLoop(svc *ai.Service, v *Validator) *RepromptLoop {
	return &RepromptLoop{service: svc, validator: v}
}

// Run executes the generate → validate → re-prompt cycle.
//
// providerName is the AI provider to call (e.g. "openai").
// systemPrompt is the full system prompt (base + schema context).
// userPrompt is the user's natural-language request.
//
// Returns the best PatternFile found and any validation errors from the final
// attempt. The caller should check whether errs is empty before treating the
// result as authoritative.
func (l *RepromptLoop) Run(
	ctx context.Context,
	providerName string,
	systemPrompt string,
	userPrompt string,
	model string,
) (*patternV1beta1.PatternFile, []string, error) {

	messages := []ai.Message{
		{Role: ai.RoleUser, Content: userPrompt},
	}
	reqCtx := ai.RequestContext{
		Provider:     providerName,
		Model:        model,
		SystemPrompt: systemPrompt,
	}

	var bestDesign *patternV1beta1.PatternFile
	var bestErrors []string

	for attempt := 0; attempt <= MaxRepromptAttempts; attempt++ {
		resp, err := l.service.Complete(ctx, ai.CompleteRequest{
			Context:  reqCtx,
			Messages: messages,
		})
		if err != nil {
			return bestDesign, bestErrors, fmt.Errorf("llm call failed on attempt %d: %w", attempt+1, err)
		}

		result := l.validator.Validate(resp.Content)
		if result.Design != nil {
			bestDesign = result.Design
			bestErrors = result.Errors
		}

		if result.Valid {
			return bestDesign, nil, nil
		}

		// If this was the last attempt, stop and return the best we have.
		if attempt == MaxRepromptAttempts {
			break
		}

		// Append the LLM response and a repair request to the conversation.
		messages = append(messages,
			ai.Message{Role: ai.RoleAssistant, Content: resp.Content},
			ai.Message{
				Role:    ai.RoleUser,
				Content: aicontext.DesignRepairPrompt(resp.Content, result.Errors),
			},
		)
	}

	// Return whatever we have — the caller decides whether to surface it.
	return bestDesign, bestErrors, nil
}
