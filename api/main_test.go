package api

import (
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/dbSeed"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	configOb := config.Config{
		TokenSymmetricKey:   dbSeed.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, configOb)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
