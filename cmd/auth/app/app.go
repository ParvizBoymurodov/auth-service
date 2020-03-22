package app

import (
	"auth-service/pkg/managers"
	"auth-service/pkg/token"
	"errors"
	"github.com/ParvizBoymurodov/jwt/jwt"
	"github.com/ParvizBoymurodov/mux/pkg/mux"
	"github.com/ParvizBoymurodov/rest/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

type server struct {
	router     *mux.ExactMux
	pool       *pgxpool.Pool
	secret     jwt.Secret
	tokenSvc   *token.Service
	managerSvc *managers.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, tokenSvc *token.Service, managerSvc *managers.Service) *server {
	if router == nil {
		panic(errors.New("router can't be nil"))
	}
	if pool == nil {
		panic(errors.New("pool can't be nil"))
	}
	if secret == nil {
		panic(errors.New("secret can't be nil"))
	}
	if tokenSvc == nil {
		panic(errors.New("tokenSvc can't be nil"))
	}
	if managerSvc == nil {
		panic(errors.New("managerSvc can't be nil"))
	}
	return &server{
		router:     router,
		pool:       pool,
		secret:     secret,
		tokenSvc:   tokenSvc,
		managerSvc: managerSvc}
}

func (s *server) Start() {
	s.InitRoutes()
}

func (s *server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *server) handleCreateToken() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body token.RequestDTO
		err := rest.ReadJSONBody(request, &body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.json_invalid"},
			})
			log.Print(err)
			return
		}

		response, err := s.tokenSvc.Generate(request.Context(), &body)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrInvalidLogin):
				writer.WriteHeader(http.StatusBadRequest)
				err := rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.login_mismatch"},
				})
				log.Print(err)
			case errors.Is(err, token.ErrInvalidPassword):
				writer.WriteHeader(http.StatusBadRequest)
				err := rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.password_mismatch"},
				})
				log.Print(err)
			default:
				writer.WriteHeader(http.StatusBadRequest)
				err := rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.unknown"},
				})
				log.Print(err)
			}
			return
		}

		err = rest.WriteJSONBody(writer, &response)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *server) handleProfile() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response, err := s.managerSvc.Profile(request.Context())
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err := rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &response)
		if err != nil {
			log.Print(err)
		}

	}
}

func (s *server) handleAddManager() func(http.ResponseWriter, *http.Request)  {
	return func(writer http.ResponseWriter, request *http.Request) {
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			log.Println("can't")
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		manager := managers.Manager{}
		err := rest.ReadJSONBody(request, &manager)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = s.managerSvc.AddManager(manager)
		if err != nil {
			log.Printf("can't handle post add user: %d", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = rest.WriteJSONBody(writer, &manager)
		if err != nil {
			http.Error(writer,http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}