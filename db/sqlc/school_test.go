package db

import (
	"context"
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/dbSeed"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomSchool(t *testing.T) School {
	arg := dbSeed.RandomSchool()

	school, err := testQueries.CreateSchool(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, school)
	require.Equal(t, arg, school.Name)
	require.NotZero(t, school.ID)
	require.NotZero(t, school.CreatedAt)
	return school
}

func TestCreateSchool(t *testing.T) {
	CreateRandomSchool(t)
}

func TestGetSchool(t *testing.T) {
	school1 := CreateRandomSchool(t)
	school2, err := testQueries.GetSchoolbyId(context.Background(), school1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, school2)
	require.Equal(t, school2.Name, school1.Name)
	require.Equal(t, school2.ID, school1.ID)
	require.Equal(t, school2.CreatedAt, school1.CreatedAt)

}

func TestUpdateSchool(t *testing.T) {
	school1 := CreateRandomSchool(t)

	arg := UpdateSchoolParams{
		ID:   school1.ID,
		Name: dbSeed.RandomSchool(),
	}

	school2, err := testQueries.UpdateSchool(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, school2)
	require.Equal(t, arg.Name, school2.Name)
	require.Equal(t, school2.ID, school1.ID)
	require.Equal(t, school2.CreatedAt, school1.CreatedAt)
	require.NotEmpty(t, school2.UpdatedAt)

}

func TestDeleteSchool(t *testing.T) {
	school1 := CreateRandomSchool(t)
	err := testQueries.DeleteSchool(context.Background(), school1.ID)
	require.NoError(t, err)

	school2, err := testQueries.GetSchoolbyId(context.Background(), school1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, school2)
}

func TestListSchools(t *testing.T) {
	for i := 1; i < 10; i++ {
		CreateRandomSchool(t)
	}
	arg := ListSchoolsParams{
		Limit:  5,
		Offset: 6,
	}
	schools, err := testQueries.ListSchools(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, schools, 5)

	for _, school := range schools {
		require.NotEmpty(t, school)
	}
}
