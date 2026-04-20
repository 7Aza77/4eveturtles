package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresDB_Error(t *testing.T) {
	db, err := NewPostgresDB("localhost", "5431", "invalid", "invalid", "invalid")
	
	assert.Error(t, err)
	assert.Nil(t, db)
}
