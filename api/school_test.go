package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/EliriaT/SchoolAppApi/db/mock"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/dbSeed"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetSchoolApi(t *testing.T) {
	school := randomSchool()

	testCases := []struct {
		name          string
		schoolId      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			schoolId: school.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Eq(school.ID)).Times(1).Return(school, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSchool(t, recorder.Body, school)
			},
		},

		{
			name:     "NotFound",
			schoolId: school.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs, returning empty school
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Eq(school.ID)).Times(1).Return(db.School{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name:     "InternalError",
			schoolId: school.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs, returning empty school
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Eq(school.ID)).Times(1).Return(db.School{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:     "InvalidId",
			schoolId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			//start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/schools/%d", tc.schoolId)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			require.NoError(t, err)
			//we get the response from the server in the recorder
			server.router.ServeHTTP(recorder, request)

			//check the response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomSchool() db.School {
	return db.School{
		ID:        dbSeed.RandomInt(1, 1000),
		Name:      dbSeed.RandomSchool(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

}

func requireBodyMatchSchool(t *testing.T, body *bytes.Buffer, school db.School) {
	data, err := io.ReadAll(body)

	var respSchool db.School
	err = json.Unmarshal(data, &respSchool)
	require.NoError(t, err)
	require.EqualValues(t, school.ID, respSchool.ID)
	//require.
}
