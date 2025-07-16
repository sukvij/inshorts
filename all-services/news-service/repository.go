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

func (repo *NewsRepository) GetNewsArticlesByCategory(category string) (*[]NewsArticle, int64, string, error) {
	marshaledCategory, err1 := json.Marshal(category)
	if err1 != nil {
		return nil, 0, "", err1
	}
	var result []NewsArticle

	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&NewsArticle{}).Where("JSON_CONTAINS(category, ?)", string(marshaledCategory)).Order("publication_date desc").Limit(1).Find(&NewsArticle{})
	})
	fmt.Println(queryDetails)
	var totalRecords int64
	// find total records
	repo.DB.Model(&NewsArticle{}).Where("JSON_CONTAINS(category, ?)", string(marshaledCategory)).Count(&totalRecords)
	err := repo.DB.Where("JSON_CONTAINS(category, ?)", string(marshaledCategory)).Order("publication_date desc").Limit(1).Find(&result).Error

	if err != nil {
		return nil, 0, "", err
	}
	return &result, totalRecords, queryDetails, nil
}

func (repo *NewsRepository) GetNewsArticlesByScore(score float64) (*[]NewsArticle, int64, string, error) {
	var result []NewsArticle

	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&NewsArticle{}).Where("relevance_score > ?", score).Order("relevance_score desc").Limit(1).Find(&NewsArticle{})
	})

	var totalRecords int64
	// find total records
	repo.DB.Model(&NewsArticle{}).Where("relevance_score > ?", score).Count(&totalRecords)
	err := repo.DB.Where("relevance_score > ?", score).Order("relevance_score desc").Limit(1).Find(&result).Error
	if err != nil {
		return nil, 0, "", err
	}
	return &result, totalRecords, queryDetails, nil
}

func (repo *NewsRepository) GetNewsArticlesBySource(source string) (*[]NewsArticle, int64, string, error) {
	var result []NewsArticle

	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&NewsArticle{}).Where("source_name = ?", source).Order("publication_date desc").Limit(1).Find(&NewsArticle{})
	})

	var totalRecords int64
	// find total records
	repo.DB.Model(&NewsArticle{}).Where("source_name = ?", source).Order("publication_date desc").Count(&totalRecords)

	err := repo.DB.Where("source_name = ?", source).Order("publication_date desc").Limit(1).Find(&result).Error
	if err != nil {
		return nil, 0, "", err
	}
	return &result, totalRecords, queryDetails, nil
}

func (repo *NewsRepository) GetNearByNewsArticle(lat, lon, radius float64) (*[]NewsArticle, int64, string, error) {
	var result []NewsArticle

	query := `
								SELECT
								*,
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS( ? )) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS( ? )) +
										SIN(RADIANS( ? )) * SIN(RADIANS(latitude))
									)
								) AS distance 
							FROM
								news_articles
							WHERE
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS( ? )) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS( ? )) +
										SIN(RADIANS( ? )) * SIN(RADIANS(latitude))
									)
								) <= ?   -- here radius is in km
							
								`
	original_query := fmt.Sprintf(`%v  ORDER BY distance limit 4`, query)
	temp := repo.DB.Raw(original_query, lat, lon, lat, lat, lon, lat, radius).Scan(&result)
	if temp.Error != nil {
		return nil, 0, "", fmt.Errorf("failed to fetch points of interest within radius: %w", temp.Error)
	}

	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(original_query, lat, lon, lat, lat, lon, lat, radius).Find(&NewsArticle{})
	})

	// var totalRecords int64
	// find total records
	totalRecords := repo.DB.Model(&NewsArticle{}).Raw(query, lat, lon, lat, lat, lon, lat, radius).Scan(&NewsArticle{}).RowsAffected
	fmt.Println(totalRecords)
	return &result, totalRecords, queryDetails, nil
}

func (repo *NewsRepository) GetNewsArticleBySearch(whereClause string, arg []interface{}) (*[]NewsArticle, int64, string, error) {

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
    `

	// originalQuery := sqlQuery + "limit 4"
	originalQuery := fmt.Sprintf(`%v limit 4`, sqlQuery)
	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(originalQuery, variable, variable, variable).Find(&NewsArticle{})
	})
	var result []NewsArticle
	err := repo.DB.Raw(originalQuery, variable, variable, variable).Scan(&result).Error

	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to fetch articles from DB: %v", err)
	}

	var totalRecords int64
	// find total records
	repo.DB.Raw(sqlQuery, variable, variable, variable).Count(&totalRecords)
	fmt.Println("totalrec", totalRecords)

	return &result, totalRecords, queryDetails, nil
}
