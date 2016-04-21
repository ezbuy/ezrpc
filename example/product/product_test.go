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
	r := productUrl + purchaseSource
	result.ProductUrl = &r
	return result, nil
}

func (s *productServiceImplementation) Ping() error {
	return nil
}

func (s *productServiceImplementation) OnExchangeUpdate() error {
	return nil
}

func (s *productServiceImplementation) OnCacheEvict(arg string) error {
	return nil
}

func (s *productServiceImplementation) DirectGetProductDetail(productUrl string, purchaseSource string) (*TProduct, error) {
	result := new(TProduct)
	r := productUrl + purchaseSource + "!"
	result.ProductUrl = &r
	return result, nil
}

func (s *productServiceImplementation) DirectOnCacheEvict() error {
	return nil
}

func TestMain(t *testing.T) {
	var nc *nats.Conn
	nc, _ = nats.Connect(nats.DefaultURL)
	server := new(productServiceImplementation)
	s := NewProductServer(server, nc)
	directKey := "testServer1"
	s.SetDirectKey(directKey)
	time.Sleep(10 * time.Millisecond)

	client := ezrpc.NewClient("Product", nc)
	client.DirectKey = directKey
	scr := ProductClient{Client: client}

	err := scr.Ping()
	if err != nil {
		t.Error(err)
	}

	product, err := scr.GetProductDetail("productUrl", "surf")
	if err != nil {
		t.Error(err)
	}
	if *product.ProductUrl != "productUrlsurf" {
		t.Error("server response error")
	}

	product, err = scr.DirectGetProductDetail("productUrl", "direct")
	if err != nil {
		t.Error(err)
	}
	if *product.ProductUrl != "productUrldirect!" {
		t.Error("server response error")
	}
}
