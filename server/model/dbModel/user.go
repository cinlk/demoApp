package dbModel

import (
	utils2 "demoApp/server/utils"
	"demoApp/server/utils/leancloud"
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
	Phone      string    `gorm:"unique" json:"phone"`
	Email      string    `gorm:"unique" json:"email"`
	Password   string    `json:"password"`
	Uuid       string    `gorm:"primary_key;unique" json:"uuid"`
	User       User      `gorm:"ForeignKey:uuid;AssociationForeignKey:uuid" json:"user"`
	Recruiter  Recruiter `gorm:"ForeignKey:uuid;AssociationForeignKey:uuid" json:"recruiter"`
	// 当前角色 和 终端类型
	Role       string `gorm:"type:role" json:"role"`
	DeviceType string `json:"device_type"`

	// 可以有多个不同类型的第三方账号
	RelatedAccount []SocialAccount `gorm:"ForeignKey:uuid;AssociationForeignKey:uuid" json:"-"`
	// 关联一个leancloud 账号
	LeanCloud LeanCloudAccount `gorm:"ForeignKey:uuid;AssociationForeignKey:uuid"`
}

// LeanCloud 账号
type LeanCloudAccount struct {
	gorm.Model `json:"-"`
	//Phone string  `gorm:"ForeignKey:Phone"`
	Uuid string `gorm:"unique" json:"uuid"`
	// leancloud 自身账号
	UserId   string `gorm:"unique" json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// 第三方账号(相同type的phone 唯一)
type SocialAccount struct {
	gorm.Model `json:"-"`
	//Phone      string     `gorm:"ForeignKey:Phone" json:"phone"`
	Uuid      string     `json:"uuid"`
	RelatedID string     `gorm:"primary_key;unique" json:"related_id"`
	Type      SocialType `gorm:"type:text" json:"type"`
	Account   *Account   `gorm:"ForeignKey:Phone;AssociationForeignKey:Phone"`
}

type User struct {
	gorm.Model `json:"-"`
	//Phone      string     `gorm:"ForeignKey:phone" json:"phone"`
	Uuid      string     `gorm:"unique;type:text" json:"uuid"`
	Name      string     `gorm:"type:text" json:"name"`
	Online    bool       `gorm:"default:false" json:"online"`
	LastLogin *time.Time `gorm:"default:now()" json:"-"`
	//Account    Account   `gorm:"ForeignKey:Phone;AssociationForeignKey:Phone"`
	UserIcon string `json:"user_icon"`
	//CarrerTalks []UserCarrerTalk `gorm:"ForeignKey:UserId;AssociationForeignKey:UserId" json:"carrer_talks"`
	CheckVisitorTime *time.Time `json:"check_visitor_time"`
	// 关联的我的访问者
	MyVisitors []RecruiterVisitorUser `gorm:"ForeignKey:UserId,AssociationForeignKey:UserId" json:"my_visitors"`
}

type Recruiter struct {
	gorm.Model `json:"-"`
	//Phone      string        `gorm:"ForeignKey:phone" json:"phone"`
	Uuid      string `gorm:"unique; type:text" json:"uuid"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	CompanyId string `json:"company_id"`
	UserIcon  string `json:"user_icon"`
	Title     string `json:"title"`
	// 审核状态 TODO
	State     int        `json:"state"`
	Online    bool       `gorm:"default:false" json:"online"`
	LastLogin *time.Time `gorm:"default:now()" json:"-"`
	//Account    Account      `gorm:"ForeignKey:Phone"`
	CompusJobs []CompuseJobs `gorm:"ForeignKey:RecruiterUUID" json:"compus_jobs,omitempty"`
	InternJobs []InternJobs  `gorm:"ForeignKey:RecruiterUUID" json:"intern_jobs,omitempty"`
}

type RecruiterVisitorUser struct {
	gorm.Model  `json:"-"`
	RecruiterId string
	UserId      string
	Checked     bool `gorm:"default:false" json:"checked"`
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

// 普通用户注册, recruoter 注册 TODO
func CreateNewAccount(orm *gorm.DB, phone, password string) (uuid, lid string, err error) {

	if password == "" {

		password, err = utils.DefaultCryptPassword(utils2.RandStrings(30))
		if err != nil {
			return "", "", err
		}
	} else {
		password, err = utils.DefaultCryptPassword(password)
		if err != nil {
			return "", "", err
		}
	}

	session := orm.Begin()
	var uid = utils.GetUUID()
	ob := Account{
		Phone:    phone,
		Password: password,
		Uuid:     uid,
		Email:    utils.GetUUID(),
		//User: User{
		//	Uuid:  uid,
		//},
	}
	err = session.Create(&ob).Error
	if err != nil {
		session.Rollback()
		return "", "", err
	}
	// 创建普通用户
	err = session.Create(&User{
		Uuid: uid,
	}).Error
	if err != nil {
		session.Rollback()
		return "", "", err
	}
	// 创建recruiter  公司关联延迟到审核决定  TODO
	err = session.Create(&Recruiter{
		Uuid: uid,
	}).Error
	if err != nil {
		session.Rollback()
		return "", "", err
	}

	// TODO 创建leancloud 账号
	//var password = "password"
	id, err := leancloud.CreateAccout(phone, phone, phone)
	if err != nil {
		session.Rollback()
		return "", "", err
	}

	err = session.Create(&LeanCloudAccount{
		Uuid:     uid,
		UserId:   id,
		Name:     phone,
		Password: phone,
	}).Error
	if err != nil {
		session.Rollback()
		return "", "", err
	}

	session.Commit()
	return ob.Uuid, id, nil

}
