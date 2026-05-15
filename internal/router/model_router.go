// oc-go-cc — Anthropic-to-OpenAI proxy for OpenCode Go + Claude Code
//
// Copyright (C) 2026  Samuel Tuyizere
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package router defines HTTP route registration and middleware chaining,
// as well as model selection based on request scenarios.
package router

import (
	"fmt"

	"oc-go-cc/internal/config"
)

// ModelRouter handles model selection based on scenarios.
type ModelRouter struct {
	atomic *config.AtomicConfig
}

// NewModelRouter creates a new model router.
func NewModelRouter(atomic *config.AtomicConfig) *ModelRouter {
	return &ModelRouter{atomic: atomic}
}

// RouteResult contains the selected model and fallback chain.
type RouteResult struct {
	Primary   config.ModelConfig
	Fallbacks []config.ModelConfig
	Scenario  Scenario
}

// RouteByModel checks if the given model name has a direct entry in the config.
// Returns the RouteResult and true if found, or zero value and false if not.
func (r *ModelRouter) RouteByModel(modelName string) (RouteResult, bool) {
	cfg := r.atomic.Get()
	primary, ok := cfg.Models[modelName]
	if !ok {
		return RouteResult{}, false
	}

	// Look up fallbacks by model name first, then fall back to "default"
	fallbacks := cfg.Fallbacks[modelName]
	if len(fallbacks) == 0 {
		fallbacks = cfg.Fallbacks["default"]
	}

	return RouteResult{
		Primary:   primary,
		Fallbacks: fallbacks,
		Scenario:  ScenarioModelBased,
	}, true
}

// Route determines which model to use for a request.
// If modelName matches a key in the config, it is used directly (model-based routing).
// Otherwise, content-based scenario detection is used.
func (r *ModelRouter) Route(messages []MessageContent, tokenCount int, modelName string) (RouteResult, error) {
	// 1. Check for model-based routing first (exact match on Claude model name)
	if result, ok := r.RouteByModel(modelName); ok {
		return result, nil
	}

	// 2. Fall back to content-based scenario detection
	cfg := r.atomic.Get()
	result := DetectScenario(messages, tokenCount, cfg)

	// Get primary model for scenario
	primary, ok := cfg.Models[string(result.Scenario)]
	if !ok {
		// Fall back to default if scenario model not configured
		primary, ok = cfg.Models["default"]
		if !ok {
			return RouteResult{}, fmt.Errorf("no default model configured")
		}
	}

	// Get fallbacks for scenario
	fallbacks := cfg.Fallbacks[string(result.Scenario)]
	if len(fallbacks) == 0 {
		// Fall back to default fallbacks
		fallbacks = cfg.Fallbacks["default"]
	}

	return RouteResult{
		Primary:   primary,
		Fallbacks: fallbacks,
		Scenario:  result.Scenario,
	}, nil
}

// IsStreamingScenarioRoutingEnabled returns whether streaming requests should use
// scenario-based routing instead of always routing to the fast model.
func (r *ModelRouter) IsStreamingScenarioRoutingEnabled() bool {
	return r.atomic.Get().EnableStreamingScenarioRouting
}

// GetModelChain returns the full chain of models to try (primary + fallbacks).
func (rr *RouteResult) GetModelChain() []config.ModelConfig {
	chain := []config.ModelConfig{rr.Primary}
	chain = append(chain, rr.Fallbacks...)
	return chain
}

// RouteForStreaming determines which model to use for streaming requests.
// Prioritizes fast TTFT (time-to-first-token) over capability.
// If modelName matches a key in the config, it is used directly (model-based routing).
func (r *ModelRouter) RouteForStreaming(messages []MessageContent, tokenCount int, modelName string) RouteResult {
	// 1. Check for model-based routing first
	if result, ok := r.RouteByModel(modelName); ok {
		return result
	}

	// 2. Fall back to content-based streaming scenario detection
	cfg := r.atomic.Get()
	result := RouteForStreaming(messages, tokenCount, cfg)

	// Get primary model for scenario
	primary, ok := cfg.Models[string(result.Scenario)]
	if !ok {
		// Fall back to fast scenario if not configured
		primary, ok = cfg.Models["fast"]
		if !ok {
			// Fall back to default
			primary = cfg.Models["default"]
		}
	}

	// Get fallbacks for scenario
	fallbacks := cfg.Fallbacks[string(result.Scenario)]
	if len(fallbacks) == 0 {
		// Fall back to fast fallbacks
		fallbacks = cfg.Fallbacks["fast"]
	}

	return RouteResult{
		Primary:   primary,
		Fallbacks: fallbacks,
		Scenario:  result.Scenario,
	}
}
