package interactionservice

import "time"

type UserInteraction struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	UserID         string    `json:"user_id"`
	ArticleID      string    `json:"article_id"`
	EventType      string    `json:"event_type"`
	EventTimeStamp time.Time `json:"event_time_stamp" gorm:"autoCreateTime"`
	UserLat        float64   `json:"user_lat"`
	UserLon        float64   `json:"user_lon"`
}
