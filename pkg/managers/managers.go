package managers

import (
	"auth-service/pkg/token"
	"context"
	"errors"
	"github.com/ParvizBoymurodov/mux/pkg/middleware/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}
type ResponseDTO struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}

type Manager struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

func (s *Service) Profile(ctx context.Context) (response ResponseDTO, err error) {
	auth, ok := jwt.FromContext(ctx).(*token.Payload)
	if !ok {
		return ResponseDTO{}, errors.New("...")
	}

	return ResponseDTO{
		Id: auth.Id,
	}, nil
}

func (s *Service) Start() {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(errors.New("can't create database"))
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `
CREATE TABLE if not exists managers (
  id BIGSERIAL PRIMARY KEY,
  login    TEXT NOT NULL UNIQUE,
  password    TEXT NOT NULL ,
  removed BOOLEAN DEFAULT FALSE
);
`)
	if err != nil {
		panic(errors.New("can't create database"))
	}
}

func (s *Service) AddManager(m Manager)(err error)  {

	save, err := s.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't acuire: %d",err)
		return err
	}
	defer save.Release()
	password, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	_, err = save.Exec(
		context.Background(),
		`INSERT INTO managers( login, password) 
		VALUES ($1, $2);`,
		m.Login,
		password,
	)

	if err != nil {
		log.Printf("can't exec add user: %d", err)
		return err
	}
	return nil
}