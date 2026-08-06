package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aclgo/grpc-orders/config"
	"github.com/aclgo/grpc-orders/internal/orders/delivery/grpc/service"
	"github.com/aclgo/grpc-orders/internal/orders/repository"
	"github.com/aclgo/grpc-orders/internal/orders/usecase"
	"github.com/aclgo/grpc-orders/migrations"
	"github.com/aclgo/grpc-orders/proto"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// user
// admin
// product
// orders
// balance
// gateway
func main() {

	cfg := config.NewConfig(".")
	if err := cfg.Load(); err != nil {
		log.Fatalf("cfg.Load: %v", err)
	}

	db, err := sqlx.Open(cfg.DBDriver, cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed sqlx.Open %v", err)
	}

	db.SetMaxIdleConns(15)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(time.Minute * 5)

	srcDriver, err := iofs.New(migrations.MigrateUsingGolangMigrateFs, ".")
	if err != nil {
		log.Fatalf("iofs.New: %v\n",err)
	}

	dbDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("postgres.NewIstance: %v\n",err)
	}

	m,err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"postgres",
		dbDriver,
	)

	if err != nil {
		log.Fatalf("migrate.NewWithInstance: %v\n",err)
	}

	if err := m.Up();err != nil && err != migrate.ErrNoChange{
		log.Fatalf("up migrate: %v\n",err)
	}

	// if cfg.MigrationRun {
	// 	migrations.SetAppMigrations(db, nil)

	// 	if err := migrations.Run(); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	repo := repository.NewRepository(db)
	orderUC := usecase.NewOrderUseCase(repo)
	serviceGRPC := service.NewServiceGprc(orderUC)

	// auth := grpc_auth.NewGrpcAuth(cfg)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServerPort))
	if err != nil {
		log.Fatalf("net.Listener: %v", err)
	}

	opts := []grpc.ServerOption{
		// grpc.ChainUnaryInterceptor(auth.AuthInterceptor),
	}

	server := grpc.NewServer(opts...)

	proto.RegisterServiceOrderServer(server, serviceGRPC)

	ec := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server grpc running port 50055")
		ec <- server.Serve(l)
	}()

	select {
	case err = <-ec:
		log.Fatalln(err)
	case <-ctx.Done():
		server.Stop()
		stop()
		err = <-ec
		if err != nil {
			log.Fatalf("application terminated by error: %v", err)
		}

		fmt.Println()
		log.Println("server stopped")
	}
}
