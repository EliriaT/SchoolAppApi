package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/EliriaT/SchoolAppApi/api/token"
	mockdb "github.com/EliriaT/SchoolAppApi/db/mock"
	"github.com/EliriaT/SchoolAppApi/db/seed"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO
func TestGetSchoolApi(t *testing.T) {
	// i should create an admin user
	//user := seed.
	school := randomSchool()

	testCases := []struct {
		name          string
		schoolId      int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.TokenMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			schoolId: school.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, seed.RandomEmail(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetRoles(gomock.Any()).Times(1)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, seed.RandomEmail(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs, returning empty school
				store.EXPECT().GetRoles(gomock.Any()).Times(1)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, seed.RandomEmail(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs, returning empty school
				store.EXPECT().GetRoles(gomock.Any()).Times(1)
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, seed.RandomEmail(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetRoles(gomock.Any()).Times(1)
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "NoAuthorization",
			schoolId: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.TokenMaker) {

			},
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetRoles(gomock.Any()).Times(1)
				store.EXPECT().GetSchoolbyId(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//check the response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			tc.setupAuth(t, request, server.tokenMaker)

			//we get the response from the server in the recorder
			server.router.ServeHTTP(recorder, request)

			//check the response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomSchool() db.School {
	return db.School{
		ID:        seed.RandomInt(1, 1000),
		Name:      seed.RandomSchool(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

}

func requireBodyMatchSchool(t *testing.T, body *bytes.Buffer, school db.School) {
	data, err := io.ReadAll(body)

	var respSchool dto.SchoolResponse
	err = json.Unmarshal(data, &respSchool)
	require.NoError(t, err)
	require.EqualValues(t, school.ID, respSchool.ID)
	//require.
}
