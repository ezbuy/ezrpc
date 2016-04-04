package ezrpc

import (
	"erproduct"
	"testing"
	"time"

	"github.com/nats-io/nats"
)

type productServiceImplementation int

func (s *productServiceImplementation) GetProductDetail(productUrl string, purchaseSource string) (*erproduct.TProduct, error) {
	result := new(erproduct.TProduct)
	result.ProductUrl = productUrl + purchaseSource
	return result, nil
}

func TestMain(t *testing.T) {
	var nc *nats.Conn
	nc, _ = nats.Connect(nats.DefaultURL)
	server := new(productServiceImplementation)
	NewServer(server, nc)
	time.Sleep(10 * time.Millisecond)

	client := NewClient("Product", nc)
	scr := erproduct.ProductClient{Client: client}
	product, err := scr.GetProductDetail("productUrl", "surf")
	if err != nil {
		t.Error(err)
	}

	if product.ProductUrl != "productUrlsurf" {
		t.Error("server response error")
	}

}
