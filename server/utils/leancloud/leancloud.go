package leancloud

import (
	"bytes"
	"github.com/pkg/errors"

	"fmt"
	"goframework/gLog"
	"goframework/utils/appJson"
	"io"
	"net/http"
	"time"
)

const (
	appKey              = "iXpuFsScQm6YzIjK0fXGnob9"
	appId               = "Wg3eXD1ftMSGqoDJhsFgy5xk-gzGzoHsz"
	masterKey           = "PVrKTVDpgMveiCFBrjJGqKmc,master"
	baseUrl             = "https://wg3exd1f.api.lncld.net/1.1/"
	expireTimeDuration  = time.Minute
	systemNotifyChannel = "systemNotify"
)

type alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type leanCloudNotifyMessage struct {
	Data  map[string]interface{} `json:"data"`
	Where map[string]interface{} `json:"where,omitempty"`
	Cql   string                 `json:"cql"`
	Prod  string                 `json:"prod"`
	// 最长 16 个字符且只能由英文字母和数字组成
	ReqId string `json:"req_id,omitempty"`
	// 最长 16 个字符且只能由英文字母和数字组成
	NotificationId string `json:"notification_id,omitempty"`

	ExpirationInterval float64 `json:"expiration_interval,omitempty"`
}

type registerUserRep struct {
	SessionToken string    `json:"sessionToken"`
	CreatedAt    time.Time `json:"createdAt"`
	ObjectId     string    `json:"objectId"`
}

func CreateAccout(phone, name, password string) (string, error) {

	url := "https://wg3exd1f.api.lncld.net/1.1/users"

	body := map[string]interface{}{
		"username": name,
		"phone":    phone,
		"password": password,
	}
	b, _ := appJson.JsonMarsh(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		gLog.LOG_ERROR(err)
		return "", err
	}
	req.Header.Set("X-LC-Id", appId)
	req.Header.Set("X-LC-Key", appKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		gLog.LOG_ERROR(err)
		return "", err
	}

	defer res.Body.Close()
	var result registerUserRep
	gLog.LOG_INFO(res.StatusCode)
	if res.StatusCode == http.StatusCreated {
		err = appJson.JsonDecode(res.Body, &result)
		if err != nil {
			gLog.LOG_ERROR(err)
			return "", err
		}
	} else {
		return "", errors.New("failed")
	}

	return result.ObjectId, nil

}

var LeanCloudUtil *leanCloudUtil

type leanCloudUtil struct {
	baseUrl string
	target  map[string]interface{}
	headers map[string]string
}

func (l *leanCloudUtil) GetTargetObject(name string) interface{} {

	return l.target[name]
}

// 没有证书
func (l *leanCloudUtil) httpRequest(url, method string, headers map[string]string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		l.headers[k] = v
	}
	for k, v := range l.headers {
		req.Header.Set(k, v)
	}
	fmt.Println(req.Header)

	client := &http.Client{}
	client.Timeout = time.Second * 30
	return client.Do(req)
}

// 有证书的请求

type leanCloudInstallation struct {
	leanCloudUtil
	url       string
	masterKey string
	pushUrl   string
}

func (l *leanCloudInstallation) InstallationDetail(objectID string) (*http.Response, error) {
	var url = l.url + "/" + objectID
	return l.httpRequest(url, http.MethodGet, nil, nil)

}

// 给某些channel 推送消息
func (l *leanCloudInstallation) PushDataToChannel(content interface{}) error {

	// prod 只对ios 设置有效
	data, err := appJson.JsonMarsh(content)
	if err != nil {

		return err
	}
	fmt.Println(string(data))

	res, err := l.httpRequest(l.pushUrl, http.MethodPost, map[string]string{
		"X-LC-Key": l.masterKey,
	}, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("send notify failed"), res.Status)
	}
	var m map[string]interface{}
	_ = appJson.JsonDecode(res.Body, &m)
	fmt.Println(m)
	defer res.Body.Close()
	return nil
}

func newLeanCloudInstallation(lc *leanCloudUtil) *leanCloudInstallation {
	return &leanCloudInstallation{
		*lc,
		lc.baseUrl + "installations",
		masterKey,
		"https://wg3exd1f.push.lncld.net/1.1/push",
	}
}

// 对外暴露的接口

func LeanCloudSendSystemNotify(channel, title, notificationId, content string) error {

	if l, ok := LeanCloudUtil.GetTargetObject("installation").(*leanCloudInstallation); ok {

		return l.PushDataToChannel(leanCloudNotifyMessage{
			Data: map[string]interface{}{
				"alert": content,
				"badge": "Increment",
				"sound": "default",
				"title": title,
			},
			Where: map[string]interface{}{
				"channels": channel,
			},
			Prod:               "dev",
			ExpirationInterval: expireTimeDuration.Seconds(),
			NotificationId:     notificationId,
		})
	}
	return errors.New("not found leancloud installation utility")
}

func LeanCloudSendUserNotify(userId, title, notificationId, content string) error {
	// 更加userid 查询最新使用的设备(一个) 发送消息
	// 与userid 关联的所有设置 ？？？
	if l, ok := LeanCloudUtil.GetTargetObject("installation").(*leanCloudInstallation); ok {

		return l.PushDataToChannel(leanCloudNotifyMessage{
			Data: map[string]interface{}{
				"alert": content,
				"badge": "Increment",
				"sound": "default",
				"title": title,
			},
			Cql:                "select * from _Installation where userId=\"" + userId + "\" limit 1 order by createdAt desc  ",
			Prod:               "dev",
			ExpirationInterval: expireTimeDuration.Seconds(),
			NotificationId:     notificationId,
		})
	}

	return errors.New("not found leancloud installation utility")
}

func init() {
	println("initial leancloud")
	LeanCloudUtil = &leanCloudUtil{
		baseUrl: baseUrl,
		target:  map[string]interface{}{},
		headers: map[string]string{
			"X-LC-Id":      appId,
			"X-LC-Key":     appKey,
			"Content-Type": "application/json",
		},
	}
	//LeanCkoudUtil.baseUrl = baseUrl
	LeanCloudUtil.target["installation"] = newLeanCloudInstallation(LeanCloudUtil)

}
