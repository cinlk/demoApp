package dbModel

import "github.com/lib/pq"

type Company struct {
	BaseModel
	Type          string         `json:"type"`
	Latitude      float64        `gorm:"type:numeric" json:"latitude"`
	Longitude     float64        `gorm:"type:numeric" json:"longitude"`
	Citys         pq.StringArray `gorm:"type:text[]" json:"citys"`
	BusinessField pq.StringArray `gorm:"type:text[]" json:"business_field"`
	// 被关注次数
	ReviewCounts  int64          `json:"review_counts"`
	// 多个talks
	CarrerTalks []CareerTalk `gorm:"ForeignKey:Id;AssociationForeignKey:Id" json:"carrer_talks"`
}
