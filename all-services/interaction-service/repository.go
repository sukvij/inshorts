package interactionservice

import (
	"fmt"

	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"gorm.io/gorm"
)

type InteractionRepository struct {
	DB           *gorm.DB
	Interactions *[]UserInteraction
}

func _NewRepository(db *gorm.DB, interatcions *[]UserInteraction) *InteractionRepository {
	return &InteractionRepository{DB: db, Interactions: interatcions}
}

func (repo *InteractionRepository) CreateUserInteraction() error {

	for _, interaction := range *repo.Interactions {
		err := repo.DB.Create(&interaction).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *InteractionRepository) TrendingNewsArticles(lat, lon float64, limit, radiusMeter int) (*[]newsservice.NewsArticle, error, int64, string) {
	query := `
				SELECT
				na.id, na.title, na.description, na.url, na.publication_date, na.source_name, na.category, na.relevance_score, na.latitude, na.longitude,
				SUM(
					CASE ui.event_type
						WHEN 'view' THEN 1.0
						WHEN 'click' THEN 2.0
						WHEN 'like' THEN 3.0
						WHEN 'share' THEN 5.0
						ELSE 0.0
					END
				) AS weighted_interaction_score,
				COUNT(DISTINCT ui.user_id) AS unique_users_count,
				(
					SUM(
						CASE ui.event_type
							WHEN 'view' THEN 1.0
							WHEN 'click' THEN 2.0
							WHEN 'like' THEN 3.0
							WHEN 'share' THEN 5.0
							ELSE 0.0
						END
					) * 0.6
				) +
				(
					COUNT(DISTINCT ui.user_id) * 0.4
				) AS combined_trending_score
			FROM
				news_articles AS na
			JOIN
				user_interactions AS ui ON na.id = ui.article_id
			WHERE
				ui.event_time_stamp >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
				AND ST_Distance_Sphere(POINT(na.longitude, na.latitude), POINT(?, ?)) <= ?
			GROUP BY
				na.id, na.title, na.description, na.url, na.publication_date, na.source_name, na.category, na.relevance_score, na.latitude, na.longitude
			HAVING
				combined_trending_score > 0
			ORDER BY
				combined_trending_score DESC, 
				na.relevance_score DESC,      	
				na.id ASC
	
	`
	var radius int = radiusMeter // in meters
	var result []newsservice.NewsArticle
	temp := repo.DB.Raw(query+` limit ?`, lon, lat, radius, limit).Scan(&result)
	queryDetails := repo.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		// Build the query without executing Scan here
		return tx.Raw(query+` limit ?`, lon, lat, radius, limit)
	})
	fmt.Println(queryDetails)
	err := temp.Error
	fmt.Println(result)
	if err != nil {
		return nil, err, 0, queryDetails
	}
	// var totalRecords int64
	totalRecords := repo.DB.Raw(query, lon, lat, radius).Scan(&newsservice.NewsArticle{}).RowsAffected
	fmt.Println("tot", totalRecords)
	return &result, err, totalRecords, queryDetails
}
