package service

import (
	"bank-api/pkg/mail"
	"fmt"
)

type NotificationService struct {
	mailer *mail.Mailer
}

func NewNotificationService(mailer *mail.Mailer) *NotificationService {
	return &NotificationService{mailer: mailer}
}

func (s *NotificationService) SendPaymentNotification(email string, amount float64) error {
	subject := "Платеж успешно проведен"
	content := fmt.Sprintf(`
		<h1>Спасибо за оплату!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount)

	return s.mailer.Send(email, subject, content)
}

func (s *NotificationService) SendCreditNotification(email string, amount float64, term int) error {
	subject := "Кредит успешно оформлен"
	content := fmt.Sprintf(`
		<h1>Ваш кредит оформлен!</h1>
		<p>Сумма: <strong>%.2f RUB</strong></p>
		<p>Срок: <strong>%d месяцев</strong></p>
		<small>Это автоматическое уведомление</small>
	`, amount, term)

	return s.mailer.Send(email, subject, content)
}
