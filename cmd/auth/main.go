package main

import (
	"context"
	"flag"
	"github.com/ParvizBoymurodov/auth-service/cmd/auth/app"
	"github.com/ParvizBoymurodov/auth-service/pkg/managers"
	"github.com/ParvizBoymurodov/auth-service/pkg/token"
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
	hostf, ok := fromFLagOrEnv(host, envHost)
	if !ok {
		hostf = *host
	}
	portf, ok := fromFLagOrEnv(port, envPort)
	if !ok {
		portf = *port
	}
	dsnf, ok := fromFLagOrEnv(dsn, envDSN)
	if !ok {
		dsnf = *dsn
	}

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
