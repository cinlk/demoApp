package dbModel

import (
	"github.com/lib/pq"
	"time"
)

type CompuseJobs struct {
	BaseModel
	// 职位类型
	Type string `json:"type"`
	// 职位附加描述
	Benefits      string         `json:"benefits"`
	NeedSkills    string         `json:"need_skills"`
	WorkContent   string         `json:"work_content"`
	BussinesField pq.StringArray `gorm:"type:text[]" json:"bussines_field"`
	Major         pq.StringArray `gorm:"type:text[]" json:"major"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	LocationCity  pq.StringArray `gorm:"type:text[]" json:"location_city"`
	Salary        string         `json:"salary"`
	Education     string         `json:"education"`
	ApplyEndTime  *time.Time     `json:"apply_end_time"`
	CompanyID     string         `json:"company_id"`
	Company       Company        `gorm:"ForeignKey:CompanyID;AssociationForeignKey:CompanyID" json:"company"`
	Publisher     Recruiter      `gorm:"ForeignKey:Uuid;AssociationForeignKey:Uuid" json:"-"`
}
