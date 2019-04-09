package leancloud

import (
	"bytes"
	"errors"
	"goframework/gLog"
	"goframework/utils/appJson"
	"net/http"
	"time"
)

const (
	AppKey = "iXpuFsScQm6YzIjK0fXGnob9"
	AppId = "Wg3eXD1ftMSGqoDJhsFgy5xk-gzGzoHsz"
)


type registerUserRep struct {
	SessionToken string `json:"sessionToken"`
	CreatedAt  time.Time `json:"createdAt"`
	ObjectId string `json:"objectId"`
}

func CreateAccout(phone, name, password string) (string, error){

	url := "https://wg3exd1f.api.lncld.net/1.1/users"


	body := map[string]interface{}{
		"username": name,
		"phone": phone,
		"password": password,
	}
	b, _ := appJson.JsonMarsh(body)
 	req, err := http.NewRequest(http.MethodPost, url,bytes.NewReader(b))
 	if err != nil{
 		gLog.LOG_ERROR(err)
		return "", err
	}
	req.Header.Set("X-LC-Id", AppId)
 	req.Header.Set("X-LC-Key", AppKey)
 	req.Header.Set("Content-Type", "application/json")

 	client := &http.Client{}

 	res ,err := client.Do(req)

 	if err != nil{
 		gLog.LOG_ERROR(err)
		return "", err
	}

 	defer res.Body.Close()
	var result registerUserRep
 	gLog.LOG_INFO(res.StatusCode)
 	if res.StatusCode == http.StatusCreated{
		err = appJson.JsonDecode(res.Body, &result)
		if err != nil{
			gLog.LOG_ERROR(err)
			return  "", err
		}
	}else{
		return "", errors.New("failed")
	}

	return result.ObjectId, nil


}
