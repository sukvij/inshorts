package llmservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const LLMInstruction = `
				i will provide title and description and you need to provide short summery for this which explain title and description.

				Return the output as a JSON object with the following structure:
				{
					"summery": "string"
				}

				you can take the below examples
					{
						"title": "Article Title 1",
						"description": "Article Description 1"
					}
					for this  --> "llm_summary": "This article discusses the latest developments in...",

				Now, analyze the following query:
				Input Query: "{{.Query}}"
				and fetch title and description from .Query and return 
				Output:
		`

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type GenerateContentRequest struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type GenerateContentResponse struct {
	Candidates []Candidate `json:"candidates"`
}

func GenerateSummeryLLM(title, description string) string {

	// Construct the prompt
	prompt := fmt.Sprintf(`Summarize the following news article title and description concisely in one or two sentences:

					Title: %s
					Description: %s

					Summary:`, title, description)
	payload := GenerateContentRequest{
		Contents: []Content{
			{Parts: []Part{
				{Text: prompt},
			}},
		},
	}
	jsonPayload, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		log.Fatalf("Error marshaling payload: %v", jsonErr)
	}

	// Create and send HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", geminiAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("key", geminiAPIKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to Gemini API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Gemini API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}
	var llmResponse GenerateContentResponse
	var summary string
	err = json.Unmarshal(body, &llmResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %v, body: %s", err, string(body))
	}
	fmt.Println(llmResponse)
	if len(llmResponse.Candidates) > 0 && len(llmResponse.Candidates[0].Content.Parts) > 0 {
		summary = llmResponse.Candidates[0].Content.Parts[0].Text
		fmt.Println("Generated Summary:")
		fmt.Println(summary)
	} else {
		fmt.Println("No summary generated by the LLM.")
	}
	return summary
}
