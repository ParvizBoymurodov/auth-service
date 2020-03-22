package app

import (
	"auth-service/pkg/token"
	"github.com/ParvizBoymurodov/mux/pkg/middleware/auth"
	"github.com/ParvizBoymurodov/mux/pkg/middleware/jwt"
	"github.com/ParvizBoymurodov/mux/pkg/middleware/logger"
	"reflect"
)

func (s *server) InitRoutes() {
	s.router.POST(
		"/api/tokens",
		s.handleCreateToken(),
		logger.Logger("TOKEN"),
	)

	s.router.POST("/api/managers",
		s.handleAddManager(),
		logger.Logger("MANAGERS"))


	s.router.GET(
		"/api/managers/{id}",
		s.handleProfile(),
		auth.Auth(),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("USERS"),
	)

	s.router.DELETE(
		"/api/managers/1",
		s.handleProfile(),
		auth.Auth(),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("USERS"),
	)
}
