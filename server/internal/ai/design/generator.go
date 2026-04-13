// Package design provides the end-to-end pipeline for converting a
// natural-language prompt into a validated Meshery Design (PatternFile).
package design

import (
	"context"
	"fmt"

	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"

	"github.com/meshery/meshery/server/internal/ai"
	aicontext "github.com/meshery/meshery/server/internal/ai/context"
	"github.com/meshery/meshery/server/internal/ai/validation"
	patternV1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
)

// GenerateRequest carries the parameters for a design generation request.
type GenerateRequest struct {
	// Prompt is the user's natural-language infrastructure description.
	Prompt string

	// Provider is the AI provider name to use (e.g. "openai", "ollama").
	Provider string

	// Model optionally overrides the provider's default model.
	Model string
}

// GenerateResult is returned by Generator.Generate.
type GenerateResult struct {
	// Design is the generated PatternFile. May be non-nil even when Warnings
	// is non-empty (best-effort result).
	Design *patternV1beta1.PatternFile

	// Warnings contains validation errors from the final LLM attempt.
	// A non-empty slice means the design may be incomplete or partially invalid.
	Warnings []string
}

// Generator orchestrates the full NL → LLM → validate → re-prompt pipeline.
type Generator struct {
	aiService     *ai.Service
	schemaBuilder *aicontext.SchemaBuilder
	reprompt      *validation.RepromptLoop
}

// NewGenerator creates a Generator with the given dependencies.
// registryManager may be nil; in that case schema injection is disabled.
func NewGenerator(
	svc *ai.Service,
	registryManager *meshmodel.RegistryManager,
) *Generator {
	var sb *aicontext.SchemaBuilder
	if registryManager != nil {
		// Reserve ~1 500 tokens for schema context (≈ 6 000 characters).
		sb = aicontext.NewSchemaBuilder(registryManager, 1500)
	}

	validator := validation.NewValidator()
	repromptLoop := validation.NewRepromptLoop(svc, validator)

	return &Generator{
		aiService:     svc,
		schemaBuilder: sb,
		reprompt:      repromptLoop,
	}
}

// Generate converts a natural-language prompt into a Meshery PatternFile.
// It injects relevant component schemas from the registry, calls the LLM,
// validates the output, and re-prompts on validation failures (up to
// validation.MaxRepromptAttempts times).
//
// An error is returned only for hard failures (LLM unreachable, context
// cancelled, etc.). Validation warnings are returned in GenerateResult.Warnings.
func (g *Generator) Generate(ctx context.Context, req GenerateRequest) (GenerateResult, error) {
	if req.Prompt == "" {
		return GenerateResult{}, fmt.Errorf("prompt is required")
	}
	if req.Provider == "" {
		return GenerateResult{}, fmt.Errorf("provider is required")
	}

	// Build the system prompt, optionally augmented with schema context.
	systemPrompt := aicontext.MesheryDesignSystemPrompt
	if g.schemaBuilder != nil {
		schemaCtx := g.schemaBuilder.BuildContext(req.Prompt)
		systemPrompt = aicontext.BuildSystemPrompt(systemPrompt, schemaCtx)
	}

	design, warnings, err := g.reprompt.Run(ctx, req.Provider, systemPrompt, req.Prompt, req.Model)
	if err != nil {
		return GenerateResult{}, err
	}

	return GenerateResult{
		Design:   design,
		Warnings: warnings,
	}, nil
}
