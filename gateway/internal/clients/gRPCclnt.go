package clients

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cc   pkg.RatesServiceClient
	conn *grpc.ClientConn
}

func dialAddr() string {
	if v := os.Getenv("CURRENCY_GRPC_ADDR"); v != "" {
		return v
	}
	conf := config.GetGrpcConfig()
	server := config.GetServerConfig()
	if server.Host != "" {
		return fmt.Sprintf("%v:%v", server.Host, conf.Port)
	}
	return fmt.Sprintf("127.0.0.1:%v", conf.Port)
}

func New(addr string) (*Client, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		cc:   pkg.NewRatesServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetAll(ctx context.Context) (*pkg.CurrencyRates, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.cc.GetAllRates(ctx, &pkg.Empty{})
}

func (c *Client) Get(ctx context.Context, code string) (*pkg.GetRateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return c.cc.GetRate(ctx, &pkg.GetRateRequest{Code: code})
}

func GetAllRatesHandler() *pkg.CurrencyRates {
	cl, err := New(dialAddr())
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()
	all, err := cl.GetAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return all
}

func GetRateHandler(code string) (float64, error) {
	cl, err := New(dialAddr())
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()
	r, err := cl.Get(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}
	if r.Found {
		return r.Value, nil
	}
	log.Println("not found")
	return 0, nil
}
