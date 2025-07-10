package newsservice

type NewsArticleUserInuut struct {
	Id              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	URL             string   `json:"url"`
	PublicationDate string   `json:"publication_date"`
	SourceName      string   `json:"source_name"`
	Category        []string `json:"category"`
	RelevanceScore  float64  `json:"relevance_score"`
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
}

type NewsArticle struct {
	Id              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	URL             string  `json:"url"`
	PublicationDate string  `json:"publication_date"`
	SourceName      string  `json:"source_name"`
	Category        []byte  `gorm:"type:json" json:"category"`
	RelevanceScore  float64 `json:"relevance_score"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}
