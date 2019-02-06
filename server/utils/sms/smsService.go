package sms

type SmsServer interface {
	Send(phone string) (string, error)
	SendGroup(phones []string) error
}

type yunPianSms struct {
	AcccessKey string
}

func (yp *yunPianSms) Send(phone string) (string, error) {

	return "789422", nil
}

func (yp *yunPianSms) SendGroup(phones []string) error {

	return nil
}

func newYunpian() *yunPianSms {

	return &yunPianSms{
		AcccessKey: "dqwdwq",
	}
}
