package llmservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// --- Main Application Logic ---

func FindLLMEntity(request1 string) (*ParsedLLMOutput, error) {
	responseSchemaForSearch := map[string]interface{}{
		"type": "OBJECT",
		"properties": map[string]interface{}{
			"entities": map[string]string{"type": "STRING"},
			"concepts": map[string]interface{}{
				"type":  "ARRAY",
				"items": map[string]string{"type": "STRING"},
			},
			"intent": map[string]string{"type": "STRING"},
		},
		"propertyOrdering": []string{"entities", "concepts", "intent"}, // Preserve order
	}

	var req QueryRequest
	req.Query = request1

	var parsedOutput ParsedLLMOutput
	var llmError error

	// --- Simulate LLM Call if API Key is not set ---
	// --- Real LLM Call ---
	// 1. Construct the LLM Prompt
	tmpl, err := template.New("prompt").Parse(llmPromptTemplateForSearch)
	if err != nil {
		log.Printf("Error parsing prompt template: %v", err)
		return nil, fmt.Errorf("failed to prepare LLM prompt. err is %s", err)
	}

	var promptBuffer bytes.Buffer
	if err := tmpl.Execute(&promptBuffer, gin.H{"Query": req.Query}); err != nil {
		log.Printf("Error executing prompt template: %v", err)
		return nil, fmt.Errorf("error executing prompt template: %v", err)
	}
	finalPrompt := promptBuffer.String()

	// 2. Prepare Payload for Gemini API
	payload := LLMRequestPayload{
		// Contents: This is a slice (array) of conversation turns.
		Contents: []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				// Role: "user" indicates this message is from our application (the user of the LLM).
				Role: "user",
				// Parts: This is a slice of message components. For text-only, it contains one Text part.
				Parts: []struct {
					Text string `json:"text"`
				}{
					{
						// Text: This is the actual prompt or instruction for the LLM.
						Text: finalPrompt},
				},
			},
		},
		// GenerationConfig: This object holds settings for how the LLM should generate its response.
		GenerationConfig: struct {
			ResponseMimeType string                 `json:"responseMimeType"`
			ResponseSchema   map[string]interface{} `json:"responseSchema"`
		}{
			// ResponseMimeType: Specifies the desired MIME type for the model's output.
			// "application/json" tells the LLM to output a pure JSON string.
			ResponseMimeType: "application/json",
			// ResponseSchema: Provides a formal JSON Schema that the LLM's output must adhere to.
			// This enforces a precise structure on the generated JSON, making it predictable.
			ResponseSchema: responseSchemaForSearch,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling LLM request payload: %v", err)
		return nil, fmt.Errorf("error marshaling LLM request payload: %v", err)
	}

	// 3. Call Gemini API
	// if llmError == nil { // Only proceed if no error so far
	client := &http.Client{}
	geminiReq, err := http.NewRequest("POST", geminiAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error creating Gemini API request: %v", err)
		llmError = fmt.Errorf("failed to create LLM request: %w", err)
		return nil, llmError
	}

	geminiReq.Header.Set("Content-Type", "application/json")
	q := geminiReq.URL.Query()
	q.Add("key", geminiAPIKey)
	geminiReq.URL.RawQuery = q.Encode()

	resp, err := client.Do(geminiReq)
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		llmError = fmt.Errorf("LLM service unavailable or error connecting: %w", err)
		return nil, llmError
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading Gemini API response body: %v", err)
			llmError = fmt.Errorf("failed to read LLM response: %w", err)
			return nil, llmError
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("Gemini API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
			llmError = fmt.Errorf("LLM API error: %s", string(body))
			return nil, llmError
		} else {
			var llmResponse LLMResponse
			if err := json.Unmarshal(body, &llmResponse); err != nil {
				log.Printf("Error unmarshaling Gemini API response: %v, body: %s", err, string(body))
				llmError = fmt.Errorf("failed to parse LLM response: %w", err)
				return nil, llmError
			} else if len(llmResponse.Candidates) == 0 || len(llmResponse.Candidates[0].Content.Parts) == 0 {
				log.Println("LLM response contains no candidates or content parts.")
				llmError = fmt.Errorf("LLM returned empty or malformed content")
				return nil, llmError
			} else {
				llmOutputText := llmResponse.Candidates[0].Content.Parts[0].Text
				if err := json.Unmarshal([]byte(llmOutputText), &parsedOutput); err != nil {
					log.Printf("Error unmarshaling LLM's JSON output text: %v, text: %s", err, llmOutputText)
					llmError = fmt.Errorf("failed to parse LLM's structured output: %w", err)
					return nil, llmError
				}

			}
		}
	}
	// }
	return &parsedOutput, nil
}
