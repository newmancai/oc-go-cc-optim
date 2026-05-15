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

// Package token provides token counting utilities using tiktoken encoding.
package token

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// Counter handles token counting for text and message arrays.
type Counter struct {
	tiktoken *tiktoken.Tiktoken
}

// NewCounter creates a new token counter with cl100k_base encoding.
func NewCounter() (*Counter, error) {
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding: %w", err)
	}
	return &Counter{tiktoken: enc}, nil
}

// CountTokens counts tokens in a string.
func (c *Counter) CountTokens(text string) (int, error) {
	tokens := c.tiktoken.Encode(text, nil, nil)
	return len(tokens), nil
}

// MessageContent represents a single message in a conversation.
type MessageContent struct {
	Role    string
	Content string
}

// CountMessages counts tokens in a message array.
// Estimates tokens for system prompt + messages with formatting overhead.
func (c *Counter) CountMessages(system string, messages []MessageContent) (int, error) {
	// Base tokens for message formatting
	total := 3 // Start token

	if system != "" {
		sysTokens, err := c.CountTokens(system)
		if err != nil {
			return 0, err
		}
		total += sysTokens + 5 // System prompt overhead
	}

	for _, msg := range messages {
		msgTokens, err := c.CountTokens(msg.Content)
		if err != nil {
			return 0, err
		}
		total += msgTokens + 5 // Per-message overhead
	}

	return total, nil
}
