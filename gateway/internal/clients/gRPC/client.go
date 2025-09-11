package clients

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notenoughtea/currency_review/gateway/internal/config"
	"github.com/notenoughtea/currency_review/gateway/internal/dto"
	"github.com/notenoughtea/currency_review/gateway/internal/logger"
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
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.cc.GetAllRates(ctx, &pkg.Empty{})
}

func (c *Client) GetByDates(ctx context.Context, req *dto.ParsedCurrencyRequest) (*pkg.CurrencyRatesList, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.cc.GetRatesByDates(ctx, &pkg.CurrencyRequest{
		DateFrom: req.DateFrom.Format("2006-01-02"),
		DateTo:   req.DateTo.Format("2006-01-02"),
	})
}

func GetAllRatesHandler() *pkg.CurrencyRates {
	cl, err := New(dialAddr())
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer cl.Close()
	all, err := cl.GetAll(context.Background())
	if err != nil {
		logger.Log.Fatal(err)
	}
	return all
}

type CurrencyRate struct {
	Date time.Time
	Rate float32
}

func parseDate(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02", s)
}

func GetRatesByDates(req *dto.ParsedCurrencyRequest) ([]CurrencyRate, error) {
	cl, err := New(dialAddr())
	if err != nil {
		return nil, err
	}
	defer cl.Close()

	resp, err := cl.GetByDates(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		logger.Log.Info("not found")
		return nil, nil
	}

	var result []CurrencyRate
	for _, item := range resp.Items {
		for _, r := range item.Rates {
			t, err := parseDate(r.Date)
			if err != nil {
				return nil, err
			}
			result = append(result, CurrencyRate{
				Date: t,
				Rate: r.Rate,
			})
		}
	}

	return result, nil
}
