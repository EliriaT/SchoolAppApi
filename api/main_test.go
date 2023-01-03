package api

import (
	"github.com/EliriaT/SchoolAppApi/config"
	"github.com/EliriaT/SchoolAppApi/db/seed"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	configOb := config.Config{
		TokenSymmetricKey:   seed.RandomString(32),
		AccessTokenDuration: time.Hour,
	}

	serverService, err := service.NewServerService(store)
	if err != nil {
		log.Fatal("cannot create create service", err)
	}

	server, err := NewServer(serverService, configOb)

	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
