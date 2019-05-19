package qiNiuStorage

import (
	"bytes"
	"context"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"time"
)

var (
	bucket             = "syzxbucket"
	bucketDomain       = "http://pic.yihu.bingfengtech.com"
	publicBucketDomain = "http://pic-public.yihu.bingfengtech.com"
)

var mac *qbox.Mac

func StreamFileUpload(bs []byte, filename string) error {

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "logo",
		},
	}

	dataLen := int64(len(bs))
	err := formUploader.Put(context.Background(), &ret, upToken, filename,
		bytes.NewReader(bs), dataLen, &putExtra)

	if err != nil {

		return err
	}

	return nil
}

func DeleteFile(filename string) error {

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuadong
	manager := storage.NewBucketManager(mac, &cfg)
	err := manager.Delete(bucket, filename)
	if err != nil {
		return err
	}

	return nil

}

// private 的 bucket 生成有效外部链接 (bucket-domain 固定)

func LinkPrivateBucketFile(imageName string, timeout time.Duration) string {

	deadLine := time.Now().Add(timeout).Unix()
	privateURL := storage.MakePrivateURL(mac, bucketDomain, imageName, deadLine)
	return privateURL

}

// public 文件外链接
func PubluicBucketFile(imageName string) string {

	pURL := storage.MakePublicURL(publicBucketDomain, imageName)
	return pURL

}

func InitialQiNiuStorage() {
	//mac = qbox.NewMac(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"))
	//bucket = os.Getenv("BUCKET")
	//bucketDomain = os.Getenv("PRIVATE_BUCKET_DOMAIN")
	//publicBucketDomain = os.Getenv("PUBLIC_BUCKET_DOMAIN")

	mac = qbox.NewMac("Six34sWKrmERbPsTrY5kg3xJEET1-64Nm4Uk9dSu",
		"P5IZdaQ2OW0JqP1HkNVP1WGFB7duWlULhbFxB5G3")

	bucket = "syzxbucket"
	bucketDomain = "http://pic.yihu.bingfengtech.com"
	publicBucketDomain = "http://pic-public.yihu.bingfengtech.com"
}
