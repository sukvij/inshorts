package newsservice

import (
	"encoding/json"
)

type NewsArticleResponse struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	URL             string   `json:"url"`
	PublicationDate string   `json:"publication_date"`
	SourceName      string   `json:"source_name"`
	Category        []string `json:"category"`
	RelevanceScore  float64  `json:"relevance_score"`
	LLMSummery      string   `json:"llm_summery"`
}

func ConvertUserInputToNewsArticle(userInput *[]NewsArticleUserInuut) *[]NewsArticle {
	var result []NewsArticle
	for _, input := range *userInput {
		var temp NewsArticle
		res, _ := json.Marshal(input)
		json.Unmarshal(res, &temp)
		x, _ := json.Marshal((input.Category))
		temp.Category = x
		result = append(result, temp)
	}
	return &result
}

func Convert_NewsArticle_To_NewsArticleResponse(userInput *[]NewsArticle) *[]NewsArticleResponse {
	var result []NewsArticleResponse
	for _, input := range *userInput {
		var temp NewsArticleResponse
		res, _ := json.Marshal(input)
		json.Unmarshal(res, &temp)
		var haha []string
		json.Unmarshal(input.Category, &haha)
		temp.Category = haha
		result = append(result, temp)
	}
	return &result
}
