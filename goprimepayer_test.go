package goprimepayer

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
)

var s = New(&Config{
	ShopID: 1,
	Secret: "f6482bd9a166bf2s43ssc9fe60eb4774",
})

func TestPayments(t *testing.T) {
	p := s.NewPayment(big.NewInt(int64(1)), 3, decimal.New(int64(5), 0), "Оплата товара")
	if p.Sign() != "8ce441e3438556e4ffdce8fec789af6f0ea3a4f221f648ea2b41c4f4d799620e" {
		t.Errorf("bad sign (%v)", p)
	}
	p.Set("user_id", 123456)
	if p.Sign() != "771fd70dfa04179cfde68293d3737ed41854efe060a7bd620bdb28f85bd7a18e" {
		t.Errorf("bad sign (%v)", p)
	}
}

func TestNotifications(t *testing.T) {
	vv := map[string]string{
		"shop":          "1",
		"payment":       "42",
		"systemPayment": "42",
		"currency":      "3",
		"amount":        "150.3",
		"uv_user_id":    "123456",
		"sign":          "d52fa12ab53cd3f701efdcaefd2db5dac1ec58871bcfdf0dba49b76fb4bbd6ff",
	}
	if notification, err := s.VerifyNotification(vv); err != nil {
		t.Error(err)
	} else {
		if fmt.Sprintf("%d", notification.PaymentID) != vv["payment"] {
			t.Errorf("wrong payment id (%v)", notification)
		}
	}
}
