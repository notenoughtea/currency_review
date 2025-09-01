package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/notenoughtea/currency_review/currency/internal/config"
	"github.com/notenoughtea/currency_review/currency/internal/handler"
	"github.com/notenoughtea/currency_review/currency/internal/service"
	"github.com/notenoughtea/currency_review/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	config.Load()
	conf := config.GetGrpcConfig()
	port := fmt.Sprintf(":%d", conf.Port)
	lis, err := net.Listen(conf.Protocol, port)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	s := grpc.NewServer()
	svc := service.GetRatesService()
	pkg.RegisterRatesServiceServer(s, handler.NewRatesHandler(svc))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func(port string) {
		log.Printf("Сервер Currency(gPRC) запущен, tcp, %s", port)
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			grpclog.Errorf("serve error: %v", err)
		}
	}(port)

	<-ctx.Done()
	log.Println("Сигнал к отключению, идет остановка сервера")

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Currency(gPRC) остановлен")
	case <-time.After(10 * time.Second):
		log.Println("Выключение по graceful timeout, принудительная остановка")
		s.Stop()
	}

	_ = lis.Close()
	log.Println("Сервер остановлен")
}
