package goprimepayer

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/shopspring/decimal"
)

var ErrBadSign = errors.New("bad sign")
var ErrBadPaymentID = errors.New("unable to get payment id")

type Config struct {
	ShopID int
	Secret string
}

type Shop struct {
	Config *Config
}

func New(config *Config) *Shop {
	return &Shop{
		Config: config,
	}
}

func (s *Shop) NewPayment(id *big.Int, currencyID int, amount decimal.Decimal, description string) *Payment {
	return &Payment{
		ID:            id,
		CurrencyID:    currencyID,
		Amount:        amount,
		Description:   description,
		userVariables: make(map[string]string),
		uvMu:          &sync.RWMutex{},
		s:             s,
	}
}

func (s *Shop) VerifyNotification(variables map[string]string) (*Notification, error) {
	sign, ok := variables["sign"]
	if !ok {
		return nil, ErrBadSign
	}

	delete(variables, "sign")
	vv := make([]string, 0, len(variables))
	for _, key := range sortByKeys(variables) {
		vv = append(vv, variables[key])
	}

	mySign := hash(fmt.Sprintf("%s:%s", strings.Join(vv, ":"), s.Config.Secret))
	if mySign != sign {
		return nil, ErrBadSign
	}

	notification := &Notification{
		userVariables: map[string]string{},
	}

	paymentID, ok := new(big.Int).SetString(variables["payment"], 10)
	if !ok {
		return nil, ErrBadPaymentID
	}
	systemPaymentID, ok := new(big.Int).SetString(variables["systemPayment"], 10)
	if !ok {
		return nil, ErrBadPaymentID
	}
	amount, err := decimal.NewFromString(variables["amount"])
	if err != nil {
		return nil, err
	}

	notification.PaymentID = paymentID
	notification.SystemPaymentID = systemPaymentID
	notification.Amount = amount

	for key, value := range variables {
		if strings.HasPrefix(key, "uv_") {
			notification.userVariables[key[3:]] = value
		}
	}
	return notification, nil
}

type Payment struct {
	// ID платежа в вашей системе, не должен повторяться (payment)
	ID *big.Int
	// Валюта (currency)
	CurrencyID int
	// Сумма (amount)
	Amount decimal.Decimal
	// Описание (description)
	Description string

	// Шлюз платежа
	Via string

	// Адрес перенаправления при успешной оплате (если разрешено в настройках)
	SuccessURL string
	// Адрес перенаправления при ошибке оплаты (если разрешено в настройках)
	FailURL string

	userVariables map[string]string
	uvMu          *sync.RWMutex

	s *Shop
}

// Изменить пользовательский параметр
func (p *Payment) Set(key string, value interface{}) {
	p.uvMu.Lock()
	p.userVariables[key] = fmt.Sprintf("%v", value)
	p.uvMu.Unlock()
}

// Получить пользовательский параметр
func (p *Payment) Get(key string) (string, bool) {
	p.uvMu.RLock()
	value, ok := p.userVariables[key]
	p.uvMu.RUnlock()
	return value, ok
}

// Подпись платежа
func (p *Payment) Sign() string {
	vv := make(map[string]string)

	vv["payment"] = fmt.Sprintf("%v", p.ID)
	vv["currency"] = fmt.Sprintf("%v", p.CurrencyID)
	vv["amount"] = fmt.Sprintf("%v", p.Amount)
	if len(p.Description) > 0 {
		vv["description"] = fmt.Sprintf("%v", p.Description)
	}
	if len(p.Via) > 0 {
		vv["via"] = p.Via
	}
	if len(p.SuccessURL) > 0 {
		vv["success"] = p.SuccessURL
	}
	if len(p.FailURL) > 0 {
		vv["fail"] = p.FailURL
	}
	vv["shop"] = fmt.Sprintf("%v", p.s.Config.ShopID)

	p.uvMu.RLock()
	for k, v := range p.userVariables {
		vv["uv_"+k] = v
	}
	p.uvMu.RUnlock()

	svv := make([]string, 0, len(vv))
	for _, key := range sortByKeys(vv) {
		svv = append(svv, vv[key])
	}

	return hash(fmt.Sprintf("%s:%s", strings.Join(svv, ":"), p.s.Config.Secret))
}

type Notification struct {
	// ID платежа в вашей системе, не должен повторяться (payment)
	PaymentID *big.Int
	// ID платежа в PrimePayer (systemPayment)
	SystemPaymentID *big.Int
	// Валюта (currency, 3 - rub)
	CurrencyID int
	// Сумма (amount)
	Amount decimal.Decimal

	userVariables map[string]string

	// ID магазина
	ShopID int
}

func (n *Notification) Get(key string) (string, bool) {
	value, ok := n.userVariables[key]
	return value, ok
}
