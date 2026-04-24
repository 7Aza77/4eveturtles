package entity

import "time"

type Event struct {
	ID              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	Date            time.Time `json:"date" db:"date"`
	Location        string    `json:"location" db:"location"`
	MaxParticipants int       `json:"max_participants" db:"max_participants"`
	CreatorID       int64     `json:"creator_id" db:"creator_id"`
	OrganizationID  *int64    `json:"organization_id,omitempty" db:"organization_id"`
	GroupChatLink   string    `json:"group_chat_link" db:"group_chat_link"`
	Tags            []Tag     `json:"tags,omitempty" db:"-"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
