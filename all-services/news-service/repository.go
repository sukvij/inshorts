package newsservice

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type NewsRepository struct {
	DB          *gorm.DB
	NewsArticle *[]NewsArticle
}

func _NewRepository(db *gorm.DB, news *[]NewsArticle) *NewsRepository {
	return &NewsRepository{DB: db, NewsArticle: news}
}

func (repo *NewsRepository) CreateNewsArticle() error {

	for _, article := range *repo.NewsArticle {
		err := repo.DB.Create(&article).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *NewsRepository) GetNewsArticlesByCategory(category string) (*[]NewsArticle, error) {
	marshaledCategory, err1 := json.Marshal(category)
	if err1 != nil {
		return nil, err1
	}
	var result []NewsArticle

	// Use GORM's Where clause with MySQL's JSON_CONTAINS function.
	// JSON_CONTAINS(json_doc, val) returns 1 if val is found in json_doc.
	err := repo.DB.Where("JSON_CONTAINS(category, ?)", string(marshaledCategory)).Order("publication_date desc").Limit(1).Find(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *NewsRepository) GetNewsArticlesByScore(score float64) (*[]NewsArticle, error) {
	var result []NewsArticle
	err := repo.DB.Where("relevance_score > ?", score).Order("relevance_score desc").Limit(1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *NewsRepository) GetNewsArticlesBySource(source string) (*[]NewsArticle, error) {
	var result []NewsArticle
	err := repo.DB.Where("source_name = ?", source).Order("publication_date desc").Limit(1).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *NewsRepository) GetNearByNewsArticle(lat, lon, radius float64) (*[]NewsArticle, error) {
	var result []NewsArticle

	query := fmt.Sprintf(`
								SELECT
								*,
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS(%v)) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS(%v)) +
										SIN(RADIANS(%v)) * SIN(RADIANS(latitude))
									)
								) AS distance 
							FROM
								news_articles
							WHERE
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS(%v)) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS(%v)) +
										SIN(RADIANS(%v)) * SIN(RADIANS(latitude))
									)
								) <= %v 
							ORDER BY
								distance 
							LIMIT 1;
								
								`, lat, lon, lat, lat, lon, lat, radius)

	err := repo.DB.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch points of interest within radius: %w", err)
	}

	return &result, nil
}

func (repo *NewsRepository) GetNewsArticleBySearch(whereClause string, arg []interface{}) (*[]NewsArticle, error) {

	variable := whereClause
	sqlQuery := `
        SELECT
            *,
            -- Calculate the text matching score using MATCH AGAINST in NATURAL LANGUAGE MODE
            MATCH (p.title, p.description) AGAINST (? IN NATURAL LANGUAGE MODE) AS text_match_score,
            -- Calculate the combined score using weights
            (COALESCE(p.relevance_score, 0.0) * 0.7) +
            (MATCH (p.title, p.description) AGAINST (? IN NATURAL LANGUAGE MODE) * 0.3) AS combined_score
        FROM
            news_articles AS p
        WHERE
            -- Ensure at least one match for the query terms
            MATCH (p.title, p.description) AGAINST (? IN NATURAL LANGUAGE MODE)
        ORDER BY
            combined_score DESC,    -- Primary sort by the combined ranking score (highest first)
            p.relevance_score DESC
		LIMIT 1
    `

	var result []NewsArticle
	err := repo.DB.Raw(sqlQuery, variable, variable, variable).Limit(1).Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles from DB: %v", err)
	}
	// fmt.Println(result)

	return &result, nil
}
