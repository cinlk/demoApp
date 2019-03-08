package dbModel

import (
	utils2 "demoApp/server/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/gLog"
	"goframework/utils"
	"reflect"
	"time"
)

var roleTableMap = map[string]string{

	"seeker":    "User",
	"recruiter": "Recruiter",
}

const (
	WECHAT = "wechat"
	QQ     = "qq"
	WEIBO  = "weibo"
)

type SocialType string

func (s SocialType) IsValide() bool {
	switch s {
	case WECHAT, QQ, WEIBO:
		return true
	}
	return false
}

type Account struct {
	gorm.Model `json:"-"`
	Phone      string    `gorm:"primary_key;unique" json:"phone"`
	Password   string    `json:"password"`
	User       User      `gorm:"ForeignKey:Phone" json:"user"`
	Recruiter  Recruiter `gorm:"ForeignKey:Phone" json:"recruiter"`
	// 当前角色 和 终端类型
	Role       string `gorm:"type:text;default:'seeker'" json:"role"`
	DeviceType string `json:"device_type"`

	// 可以有多个不同类型的第三方账号
	RelatedAccount []SocialAccount `gorm:"ForeignKey:Phone" json:"-"`
}

// 第三方账号(相同type的phone 唯一)
type SocialAccount struct {
	gorm.Model `json:"-"`
	Phone      string     `gorm:"ForeignKey:Phone" json:"phone"`
	RelatedID  string     `gorm:"primary_key;unique" json:"related_id"`
	Type       SocialType `gorm:"type:text" json:"type"`
	Account    *Account   `gorm:"ForeignKey:Phone;AssociationForeignKey:Phone"`
}

type User struct {
	gorm.Model `json:"-"`
	Phone      string     `gorm:"ForeignKey:phone" json:"phone"`
	Uuid       string     `gorm:"primary_key;unique; type:text" json:"uuid"`
	Name       string     `gorm:"type:text" json:"name"`
	Online     bool       `gorm:"default:false" json:"online"`
	LastLogin  *time.Time `gorm:"default:now()" json:"-"`
	Account    *Account   `gorm:"ForeignKey:Phone;AssociationForeignKey:Phone"`

	//CarrerTalks []UserCarrerTalk `gorm:"ForeignKey:UserId;AssociationForeignKey:UserId" json:"carrer_talks"`
}

type Recruiter struct {
	gorm.Model `json:"-"`
	Phone      string        `gorm:"ForeignKey:phone" json:"phone"`
	Uuid       string        `gorm:"primary_key;unique; type:text" json:"uuid"`
	Name       string        `json:"name"`
	Company    string        `json:"company"`
	CompanyId  string        `json:"company_id"`
	UserIcon   string        `json:"user_icon"`
	Title      string        `json:"title"`
	State      int           `json:"state"`
	Online     bool          `gorm:"default:false" json:"online"`
	LastLogin  *time.Time    `gorm:"default:now()" json:"-"`
	Account    *Account      `gorm:"ForeignKey:Phone"`
	CompusJobs []CompuseJobs `gorm:"ForeignKey:RecruiterUUID" json:"compus_jobs,omitempty"`
	InternJobs []InternJobs  `gorm:"ForeignKey:RecruiterUUID" json:"intern_jobs,omitempty"`
}

func (a *Account) FindAssociateUserByColume(orm *gorm.DB, name string) (uuid string, t interface{}) {

	colume := roleTableMap[a.Role]
	v := reflect.Indirect(reflect.ValueOf(a))
	if v.FieldByName(colume).IsValid() == false {
		gLog.LOG_ERROR(errors.New("not found colume by name"))
		return "", nil
	}
	tr := reflect.New(v.FieldByName(colume).Type()).Interface()
	err := orm.Model(a).Association(colume).Find(tr).Error
	if err != nil {
		gLog.LOG_ERROR(err)
		return "", nil
	}

	// 获取name 的值
	u := reflect.Indirect(reflect.ValueOf(tr))
	if u.Kind() != reflect.Struct {
		return "", nil
	}

	if u.FieldByName(name).IsValid() {
		return u.FieldByName(name).String(), u.Interface()
	}

	return "", nil

}

func FindAccount(orm *gorm.DB, phone string) (*Account, error) {

	var account Account
	err := orm.Model(&account).Where("phone = ?", phone).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func FindSocitalAccount(orm *gorm.DB, relatedID, Type string) (*SocialAccount, error) {
	var account SocialAccount
	err := orm.Model(&account).Where("related_id = ? and type = ?", relatedID, Type).First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil

}

func CreateNewAccount(orm *gorm.DB, phone, password string) (uuid string, err error) {

	if password == "" {

		password, err = utils.DefaultCryptPassword(utils2.RandStrings(30))
		if err != nil {
			return "", err
		}
	} else {
		password, err = utils.DefaultCryptPassword(password)
		if err != nil {
			return "", err
		}
	}

	ob := Account{
		Phone:    phone,
		Password: password,
		User: User{
			Phone: phone,
			Uuid:  utils.GetUUID(),
		},
	}
	err = orm.Create(&ob).Error
	if err != nil {
		return "", err
	}
	return ob.User.Uuid, nil

}
