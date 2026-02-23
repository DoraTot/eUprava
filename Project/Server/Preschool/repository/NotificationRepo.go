package repository

import (
	"database/sql"
	"log"
	"time"

	"main.go/model"
)

type NotificationRepository struct {
	DB *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) EnsureTableExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id INT AUTO_INCREMENT PRIMARY KEY,
		message VARCHAR(500) NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		is_read BOOLEAN DEFAULT FALSE,
		created_at DATETIME NOT NULL
	) 
	`
	_, err := r.DB.Exec(query)
	if err != nil {
		log.Println("Error creating notifications table:", err)
	}
	return err
}

func (r *NotificationRepository) Save(notification *model.Notification) error {
	query := `
	INSERT INTO notifications (message, user_id, is_read, created_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := r.DB.Exec(query, notification.Message, notification.UserId, false, time.Now())
	return err
}

func (r *NotificationRepository) GetByUser(userID string) ([]model.Notification, error) {
	query := `
	SELECT id, message, user_id, is_read, created_at
	FROM notifications
	WHERE user_id = ? AND is_read = true
	ORDER BY created_at DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.Message, &n.UserId, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(id string) error {
	_, err := r.DB.Exec("UPDATE notifications SET is_read = TRUE WHERE user_id = ? AND is_read = false", id)
	return err
}
