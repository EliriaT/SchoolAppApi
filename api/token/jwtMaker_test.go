package token

import (
	"github.com/EliriaT/SchoolAppApi/db/seed"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(seed.RandomString(32))
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
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(seed.RandomString(32))
	require.NoError(t, err)

	token, _, err := maker.CreateToken(seed.RandomEmail(), []int64{seed.RandomInt(1, 3), seed.RandomInt(1, 3)}, seed.RandomInt(1, 100), seed.RandomInt(1, 100), seed.RandomInt(1, 100), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// TestInvalidJWTTokenAlgNone tests if vulnerable to token forgery by chaning alg in header
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(seed.RandomEmail(), []int64{seed.RandomInt(1, 3), seed.RandomInt(1, 3)}, seed.RandomInt(1, 100), seed.RandomInt(1, 100), seed.RandomInt(1, 100), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(seed.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
