package sms

import "context"

type SMSService interface {
	SendCode(ctx context.Context, biz string, phone string) error
}
