package product

import (
	"testing"
	"time"

	"github.com/ezbuy/ezrpc/ezrpc"
	"github.com/nats-io/nats"
)

type productServiceImplementation int

func (s *productServiceImplementation) GetProductDetail(productUrl string, purchaseSource string) (*TProduct, error) {
	result := new(TProduct)
	result.ProductUrl = productUrl + purchaseSource
	return result, nil
}

func (s *productServiceImplementation) Ping() error {
	return nil
}

func TestMain(t *testing.T) {
	var nc *nats.Conn
	nc, _ = nats.Connect(nats.DefaultURL)
	server := new(productServiceImplementation)
	NewServer(server, nc)
	time.Sleep(10 * time.Millisecond)

	client := ezrpc.NewClient("Product", nc)
	scr := ProductClient{Client: client}

	err := scr.Ping()
	if err != nil {
		t.Error(err)
	}

	product, err := scr.GetProductDetail("productUrl", "surf")
	if err != nil {
		t.Error(err)
	}

	if product.ProductUrl != "productUrlsurf" {
		t.Error("server response error")
	}
}
