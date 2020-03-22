package main

import (
	"auth-service/cmd/auth/app"
	"auth-service/pkg/managers"
	"auth-service/pkg/token"
	"context"
	"flag"
	"github.com/ParvizBoymurodov/mux/pkg/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("host", "", "server host")
	port = flag.String("port", "", "server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)

const (
	envHost = "HOST"
	envPort = "PORT"
	envDSN  = "DATABASE_URL"
)

func fromFLagOrEnv(flag *string, envName string) (server string, ok bool) {
	if *flag != "" {
		return *flag, true
	}
	return os.LookupEnv(envName)
}

func main() {
	flag.Parse()
	hostf, _ := fromFLagOrEnv(host, envHost)
	portf, _ := fromFLagOrEnv(port, envPort)
	dsnf, _ := fromFLagOrEnv(dsn, envDSN)

	addr := net.JoinHostPort(hostf, portf)
	start(addr, dsnf)
}

func start(addr string, dsn string) {
	var secret = []byte("secret")
	router := mux.NewExactMux()

	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	managerSvc := managers.NewService(pool)
	managerSvc.Start()
	tokenSvc:=token.NewService(secret,pool)
	server := app.NewServer(
		router,
		pool,
		secret,
		tokenSvc,
		managerSvc,
	)
	server.Start()
	panic(http.ListenAndServe(addr, server))
}
