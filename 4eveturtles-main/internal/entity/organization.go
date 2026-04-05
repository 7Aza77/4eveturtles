package entity

import "time"

type Organization struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	UniversityID string    `json:"university_id" db:"university_id"`
	GroupChatLink string   `json:"group_chat_link" db:"group_chat_link"`
	OwnerID      int64     `json:"owner_id" db:"owner_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
