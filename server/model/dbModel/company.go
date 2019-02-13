package dbModel

type Company struct {
	BaseModel

	// 多个talks
	CarrerTalks []CareerTalk `gorm:"ForeignKey:Id;AssociationForeignKey:Id" json:"carrer_talks"`
}
