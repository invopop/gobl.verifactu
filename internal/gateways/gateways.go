// Package gateways provides the VeriFactu gateway
package gateways

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/gobl.verifactu/doc"
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

	correctStatus = "Correcto"
)

// Response defines the response fields from the VeriFactu gateway.
type Response struct {
	XMLName        xml.Name `xml:"RespuestaSuministro"`
	CSV            string   `xml:"CSV"`
	EstadoEnvio    string   `xml:"EstadoEnvio"`
	RespuestaLinea []struct {
		EstadoRegistro           string `xml:"EstadoRegistro"`
		DescripcionErrorRegistro string `xml:"DescripcionErrorRegistro,omitempty"`
	} `xml:"RespuestaLinea"`
}

// Connection defines what is expected from a connection to a gateway.
type Connection struct {
	client *resty.Client
}

// New instantiates a new connection using the provided config.
func New(env Environment) (*Connection, error) {
	c := new(Connection)
	c.client = resty.New()

	switch env {
	case EnvironmentProduction:
		c.client.SetBaseURL(ProductionBaseURL)
	default:
		c.client.SetBaseURL(TestingBaseURL)
	}
	c.client.SetDebug(os.Getenv("DEBUG") == "true")
	return c, nil
}

// Post sends the VeriFactu document to the gateway
func (c *Connection) Post(ctx context.Context, doc doc.VeriFactu) error {
	payload, err := doc.Bytes()
	if err != nil {
		return fmt.Errorf("generating payload: %w", err)
	}
	return c.post(ctx, TestingBaseURL, payload)
}

func (c *Connection) post(ctx context.Context, path string, payload []byte) error {
	out := new(Response)
	req := c.client.R().
		SetContext(ctx).
		SetDebug(true).
		SetHeader("Content-Type", "application/xml").
		SetContentLength(true).
		SetBody(payload).
		SetResult(out)

	res, err := req.Post(path)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return ErrInvalid.withCode(strconv.Itoa(res.StatusCode()))
	}

	if out.EstadoEnvio != correctStatus {
		err := ErrInvalid
		if len(out.RespuestaLinea) > 0 {
			e1 := out.RespuestaLinea[0]
			err = err.withMessage(e1.DescripcionErrorRegistro).withCode(e1.EstadoRegistro)
		}
		return err
	}

	return nil
}
