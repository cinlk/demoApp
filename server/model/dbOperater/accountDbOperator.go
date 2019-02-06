package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/utils/errorStatus"
	"demoApp/server/utils/jwt_auth"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"goframework/orm"
	"goframework/utils"
	"net/http"
	"time"
)

const (
	DEFAULT_ROLE    = "seeker"
	ANONYMOUSE_ROLE = "anonymous"
)

type AccountDbOperator struct {
	orm   *gorm.DB
	claim jwt_auth.CreateClaim
}

func (a *AccountDbOperator) ExistAccount(phone string) error {

	_, err := dbModel.FindAccount(a.orm, phone)
	if err != nil {
		return &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.Wrap(err, "not found account"),
		}
	}
	return nil
}

func (a *AccountDbOperator) LoginAnonymous() (string, error) {

	return a.claim.CreateDefaultToken("", ANONYMOUSE_ROLE)
}

func (a *AccountDbOperator) LoginByPwd(phone, password string) (token string, err error) {

	account, err := dbModel.FindAccount(a.orm, phone)
	if err != nil {
		return "", &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.Wrap(err, "not found account"),
		}
	}

	if utils.CompareHashPassword(account.Password, password) == false {
		return "", &errorStatus.AppError{
			Code: http.StatusForbidden,
			Err:  errors.New("password not match"),
		}
	}
	uuid, i := account.FindAssociateUserByColume(a.orm, "Uuid")
	if uuid == "" {
		return "", &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("not found user with account"),
		}
	}
	a.resetLoginState(i, uuid)

	return a.claim.CreateDefaultToken(uuid, account.Role)

}

func (a *AccountDbOperator) LoginByCode(phone, code string) (token string, err error) {

	// 从数据库查找该用户，有返回该用户数据
	account, err := dbModel.FindAccount(a.orm, phone)

	if err == gorm.ErrRecordNotFound {
		// 创建新用户 密码随机
		uuid, err := dbModel.CreateNewAccount(a.orm, phone, "")
		if err != nil {
			return "", err
		}

		return a.claim.CreateDefaultToken(uuid, DEFAULT_ROLE)

	} else if err == nil {
		uuid, i := account.FindAssociateUserByColume(a.orm, "Uuid")
		if uuid == "" {
			return "", errors.New("not found associate user")
		}
		a.resetLoginState(i, uuid)

		return a.claim.CreateDefaultToken(uuid, account.Role)

	} else {
		return "", err
	}

}

func (a *AccountDbOperator) RegistryAccount(phone, password string) (token string, err error) {

	account, err := dbModel.FindAccount(a.orm, phone)
	if account != nil {
		return "", &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("account is exist"),
		}
	}
	if err == gorm.ErrRecordNotFound {
		uuid, err := dbModel.CreateNewAccount(a.orm, phone, password)
		if err != nil {
			return "", err
		}
		// 生成token
		return a.claim.CreateDefaultToken(uuid, DEFAULT_ROLE)

	}

	return "", err

}

func (a *AccountDbOperator) ResetPassword(phone, password string) (token string, err error) {

	hashPwd, err := utils.DefaultCryptPassword(password)
	if err != nil {
		return "", err
	}
	account, err := dbModel.FindAccount(a.orm, phone)
	if err == gorm.ErrRecordNotFound {

		return "", &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("account not exist"),
		}
	} else if err == nil {
		err = a.orm.Model(account).Update("password", hashPwd).Error
		if err != nil {
			return "", err
		}
		uuid, ins := account.FindAssociateUserByColume(a.orm, "Uuid")
		if uuid == "" {
			return "", errors.New("not found user")
		}
		switch ins.(type) {
		case dbModel.User:
			_ = a.orm.Model(&dbModel.User{}).Where("uuid = ?", uuid).
				Updates(map[string]interface{}{"online": true, "last_login": time.Now()}).Error

		case dbModel.Recruiter:
			_ = a.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", uuid).
				Updates(map[string]interface{}{"online": true, "last_login": time.Now()}).Error
		}

		// 生成token
		return a.claim.CreateDefaultToken(uuid, account.Role)
	} else {

		return "", err
	}

}

func (a *AccountDbOperator) LoginByRelatedAccount(relatedID, Type string) (token string, err error) {

	socialAccount, err := dbModel.FindSocitalAccount(a.orm, relatedID, Type)
	if err != nil {
		code := http.StatusUnprocessableEntity
		if err == gorm.ErrRecordNotFound {
			code = http.StatusNotFound
		}
		return "", &errorStatus.AppError{
			Code: code,
			Err:  errors.Wrap(err, "can't find social account with"),
		}
	}

	// 根据 account 获取用户信息
	socialAccount.Account = &dbModel.Account{}
	err = a.orm.Model(socialAccount).Related(&socialAccount.Account, "Account").Error
	if err != nil {
		return "", errors.New("not found related Account")
	}

	uuid, i := socialAccount.Account.FindAssociateUserByColume(a.orm, "Uuid")
	if uuid == "" {
		return "", errors.New("associate user not found")
	}
	a.resetLoginState(i, uuid)
	return a.claim.CreateDefaultToken(uuid, socialAccount.Account.Role)

}

func (a *AccountDbOperator) BindRelatedAccount(phone, relatedID, Type string) (token string, err error) {

	account, err := dbModel.FindAccount(a.orm, phone)

	if err == gorm.ErrRecordNotFound {
		return "", &errorStatus.AppError{
			Code: http.StatusNotFound,
			Err:  errors.New("account not found"),
		}
	} else if err == nil {
		t := dbModel.SocialType(Type)
		if t.IsValide() == false {
			return "", errors.New("social type is invalide")
		}

		err = a.orm.Model(account).Association("RelatedAccount").Append(&dbModel.SocialAccount{
			Phone:     phone,
			RelatedID: relatedID,
			Type:      t,
		}).Error
		if err != nil {
			return "", err
		}
		uuid, i := account.FindAssociateUserByColume(a.orm, "Uuid")
		if uuid == "" {
			return "", errors.New("not found user with account")
		}
		a.resetLoginState(i, uuid)

		return a.claim.CreateDefaultToken(uuid, account.Role)
	} else {

		return "", err
	}

}

func (a *AccountDbOperator) resetLoginState(i interface{}, uuid string) {

	switch i.(type) {
	case dbModel.User:
		_ = a.orm.Model(&dbModel.User{}).Where("uuid = ?", uuid).
			Updates(map[string]interface{}{"online": true, "last_login": time.Now()}).Error

	case dbModel.Recruiter:
		_ = a.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", uuid).
			Updates(map[string]interface{}{"online": true, "last_login": time.Now()}).Error
	}

}

func (a *AccountDbOperator) SetLogOut(role, uuid string) {
	switch role {
	case "seeker":
		_ = a.orm.Model(&dbModel.User{}).Where("uuid = ?", uuid).
			Updates(map[string]interface{}{"online": false}).Error
	case "recruiter":
		_ = a.orm.Model(&dbModel.Recruiter{}).Where("uuid = ?", uuid).
			Updates(map[string]interface{}{"online": false}).Error
	}

}

func NewAccountDbOperator() *AccountDbOperator {

	return &AccountDbOperator{
		orm:   orm.DB,
		claim: jwt_auth.DefaultUserClaim,
	}
}
