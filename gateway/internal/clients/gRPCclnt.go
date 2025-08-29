package clients

import (
	"context"
	"fmt"
	"log"
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
	conf := config.GetConfGRPC()
	hostString := fmt.Sprintf("localhost:%v", conf.Port)
	cl, err := New(hostString)
	if err != nil {
		log.Fatal(err)
	}
	all, err := cl.GetAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(all, len(all.ConversionRates))
	return all
}

func GetRateHandler(code string) (float64, error) {
	conf := config.GetConfGRPC()
	hostString := fmt.Sprintf("localhost:%v", conf.Port)
	cl, err := New(hostString)
	if err != nil {
		log.Fatal(err)
	}
	r, err := cl.Get(context.Background(), code)
	if err != nil {
		log.Fatal(err)
	}
	if r.Found {
		return r.Value, nil
	} else {
		log.Println("not found")
	}
	return 0, err
}
