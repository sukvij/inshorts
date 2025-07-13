package llmservice

import "time"

// --- Article Struct (matches your table schema) ---
type Article struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	PublicationDate time.Time `json:"publication_date"`
	SourceName      string    `json:"source_name"`
	Category        []string  `json:"category"`
	RelevanceScore  float64   `json:"relevance_score"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
}

// --- LLM Request/Response Structs ---

// QueryRequest represents the incoming JSON payload from the user.
type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

// LLMRequestPayload for the Gemini API.
type LLMRequestPayload struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		ResponseMimeType string                 `json:"responseMimeType"`
		ResponseSchema   map[string]interface{} `json:"responseSchema"`
	} `json:"generationConfig"`
}

// LLMResponse from the Gemini API.
type LLMResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// ParsedLLMOutput is the structured output we expect from the LLM.
type ParsedLLMOutput struct {
	Entities string   `json:"entities"` // Comma-separated string of entities
	Concepts []string `json:"concepts"` // Array of concepts
	Intent   string   `json:"intent"`   // The identified intent
}

// --- Constants for LLM API ---
const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"
	geminiAPIKey = "AIzaSyAEQBsC13woRJDddZwNY_7ISfkPMxX39CI" // In Canvas, this will be automatically provided at runtime.
	modelName    = "gemini-2.0-flash"
)
