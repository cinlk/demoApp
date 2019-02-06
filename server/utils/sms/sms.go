package sms

func GetSmsService(name string) SmsServer {
	switch name {
	case "yupian":
		return nil
	}
	return newYunpian()
}
