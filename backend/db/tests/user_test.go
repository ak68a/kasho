package db_test

import (
	"context"
	db "github/kasho/backend/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	arg := db.CreateUserParams{
		Email: "test1@example.com",
		HashedPassword: "root",
	}

	user, err :=testQuery.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user) 

	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 2*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 2*time.Second)
}