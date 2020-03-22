package token

import (
	"context"
	"errors"
	"github.com/ParvizBoymurodov/jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	secret []byte
	pool *pgxpool.Pool
}

func NewService(secret []byte, pool *pgxpool.Pool) *Service {
	return &Service{secret: secret, pool: pool}
}

type Payload struct {
	Id    int64    `json:"id"`
	Exp   int64    `json:"exp"`
}

type RequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDTO struct {
	Token string `json:"token"`
}

var ErrInvalidLogin = errors.New("invalid password")
var ErrInvalidPassword = errors.New("invalid password")



func (s *Service) Generate(context context.Context, request *RequestDTO) (response ResponseDTO, err error) {
	var pass string
	var id int64


	err = s.pool.QueryRow(context,
		`SELECT password, id 
		FROM managers
		WHERE removed = FALSE AND login = $1`,
		request.Username,
	).Scan(
		&pass,
		&id,
	)

	if err != nil {
		err = ErrInvalidLogin
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		err = ErrInvalidPassword
		return
	}

	response.Token, err = jwt.Encode(Payload{
		Id:    id,
		Exp:   time.Now().Add(time.Hour).Unix(),
	}, s.secret)
	return
}
