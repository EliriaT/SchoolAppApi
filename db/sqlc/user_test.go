package db

import (
	"context"
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/db/seed"
	"github.com/EliriaT/SchoolAppApi/db/service"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func CreateRandomUser(t *testing.T) User {
	hashedPassword, err := service.HashPassword(seed.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:       seed.RandomEmail(),
		Password:    hashedPassword,
		LastName:    seed.RandomString(6),
		FirstName:   seed.RandomString(6),
		Gender:      seed.RandomGender(),
		PhoneNumber: sql.NullString{Valid: true, String: seed.RandomPhoneNumber()},
		Domicile:    sql.NullString{Valid: true, String: seed.RandomResidence()},
		BirthDate:   seed.RandomBirthDate(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.Gender, user.Gender)
	require.Equal(t, arg.Domicile, user.Domicile)
	require.WithinDuration(t, arg.BirthDate, user.BirthDate, time.Hour*24)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserbyId(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.Gender, user2.Gender)
	require.Equal(t, user1.Domicile, user2.Domicile)
	require.WithinDuration(t, user1.BirthDate, user2.BirthDate, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)

}
