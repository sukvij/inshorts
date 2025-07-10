package newsservice

import "gorm.io/gorm"

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

func (repo *NewsRepository) GetNewsArticles() (*[]NewsArticle, error) {
	var conditions NewsArticle = (*repo.NewsArticle)[0]
	var result []NewsArticle
	err := repo.DB.Where(&conditions).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
