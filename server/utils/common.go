package utils

import (
	"bytes"
	"goframework/utils"
	"math/rand"
	"net/http"
	"time"
)

//

const (
	SEED_STR = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SUCCESS  = "SUCCESS"
	FAILED   = "FAILED"
)

// 加密密码

func RandStrings(length int) string {

	b := bytes.Buffer{}
	rand.Seed(time.Now().UnixNano())
	lens := len(SEED_STR)
	for i := 0; i < length; i++ {
		b.WriteByte(SEED_STR[rand.Intn(lens)])
	}
	return b.String()

}

func JsonResponse(w http.ResponseWriter, data interface{}, code int) {
	utils.JsonResponse(w, data, code)
}
