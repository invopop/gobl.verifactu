package gateways

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/gobl.verifactu/internal/doc"
)

// Environment defines the environment to use for connections
type Environment string

// Environment to use for connections
const (
	EnvironmentProduction Environment = "production"
	EnvironmentTesting    Environment = "testing"

	// Production environment not published yet
	ProductionBaseURL = "xxxxxxxx"
	TestingBaseURL    = "https://prewww1.aeat.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"
)

// Standard gateway error responses
var (
	ErrConnection     = errors.New("connection")
	ErrInvalidRequest = errors.New("invalid request")
)

// Connection defines what is expected from a connection to a gateway.
type VerifactuConn struct {
	client *resty.Client
}

// New instantiates a new connection using the provided config.
func NewVerifactu(env Environment) *VerifactuConn {
	c := new(VerifactuConn)
	c.client = resty.New()

	switch env {
	case EnvironmentProduction:
		c.client.SetBaseURL(ProductionBaseURL)
	default:
		c.client.SetBaseURL(TestingBaseURL)
	}
	c.client.SetDebug(os.Getenv("DEBUG") == "true")
	return c
}

func (v *VerifactuConn) Post(ctx context.Context, doc doc.VeriFactu) error {
	payload, err := doc.Bytes()
	if err != nil {
		return fmt.Errorf("generating payload: %w", err)
	}

	res, err := v.client.R().
		SetContext(ctx).
		SetBody(payload).
		Post(v.client.BaseURL)

	if err != nil {
		return fmt.Errorf("%w: verifactu: %s", ErrConnection, err.Error())
	}

	if res.StatusCode() != 200 {
		return fmt.Errorf("%w: verifactu: status %d", ErrInvalidRequest, res.StatusCode())
	}

	return nil
}
