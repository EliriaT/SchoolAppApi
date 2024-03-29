package token

import (
	"github.com/EliriaT/SchoolAppApi/db/seed"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(seed.RandomString(32))
	require.NoError(t, err)

	email := seed.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, _, err := maker.CreateToken(email, []int64{seed.RandomInt(1, 3), seed.RandomInt(1, 3)}, seed.RandomInt(1, 100), seed.RandomInt(1, 100), seed.RandomInt(1, 100), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(seed.RandomString(32))
	require.NoError(t, err)

	token, _, err := maker.CreateToken(seed.RandomEmail(), []int64{seed.RandomInt(1, 3), seed.RandomInt(1, 3)}, seed.RandomInt(1, 100), seed.RandomInt(1, 100), seed.RandomInt(1, 100), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
