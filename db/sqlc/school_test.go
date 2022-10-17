package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

//func CreateRandomSchool(t *testing.T) School {
//	arg := "I.P.L.T. Mircea Eliade"
//
//	account, err := testQueries.CreateSchool(context.Background(), arg)
//
//	require.NoError(t, err)
//	require.NotEmpty(t, account)
//	require.Equal(t, arg, account.Name)
//	require.NotZero(t, account.ID)
//	require.NotZero(t, account.CreatedAt)
//	return School
//}

func TestCreateSchool(t *testing.T) {
	arg := "I.P.L.T. Mircea Eliade"

	account, err := testQueries.CreateSchool(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg, account.Name)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

//
//func TestGetSchool()
