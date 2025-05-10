package models

type NotificationType string

const (
	NotificationPayment NotificationType = "payment"
	NotificationCredit  NotificationType = "credit"
)

type Notification struct {
	Type    NotificationType
	Email   string
	Subject string
	Content string
}
