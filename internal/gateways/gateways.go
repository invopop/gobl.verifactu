// Package gateways provides the VeriFactu gateway
package gateways

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/xmldsig"
)

// Environment defines the environment to use for connections
type Environment string

// Environment to use for connections
const (
	EnvironmentProduction Environment = "production"
	EnvironmentSandbox    Environment = "sandbox"

	// Production environment not published yet
	ProductionBaseURL = "xxxxxxxx"
	TestingBaseURL    = "https://prewww1.aeat.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"

	correctStatus = "Correcto"
)

// Connection defines what is expected from a connection to a gateway.
type Connection struct {
	client *resty.Client
}

// New instantiates and configures a new connection to the VeriFactu gateway.
func New(env Environment, cert *xmldsig.Certificate) (*Connection, error) {
	tlsConf, err := cert.TLSAuthConfig()
	if err != nil {
		return nil, doc.ErrValidation.WithMessage(fmt.Errorf("preparing TLS config: %v", err).Error())
	}
	c := new(Connection)
	c.client = resty.New()

	switch env {
	case EnvironmentProduction:
		return nil, doc.ErrValidation.WithMessage("production environment not available yet")
	default:
		c.client.SetBaseURL(TestingBaseURL)
	}
	tlsConf.InsecureSkipVerify = true
	c.client.SetTLSClientConfig(tlsConf)
	c.client.SetDebug(os.Getenv("DEBUG") == "true")
	return c, nil
}

// Post sends the VeriFactu document to the gateway
func (c *Connection) Post(ctx context.Context, doc doc.VeriFactu) error {
	pyl, err := doc.Envelop()
	if err != nil {
		return fmt.Errorf("generating payload: %w", err)
	}
	return c.post(ctx, TestingBaseURL, pyl)
}

func (c *Connection) post(ctx context.Context, path string, payload []byte) error {
	out := new(Envelope)
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
		return doc.ErrValidation.WithCode(strconv.Itoa(res.StatusCode()))
	}
	if out.Body.Respuesta.EstadoEnvio != correctStatus {
		err := doc.ErrValidation.WithCode(strconv.Itoa(res.StatusCode()))
		if len(out.Body.Respuesta.RespuestaLinea) > 0 {
			e1 := out.Body.Respuesta.RespuestaLinea[0]
			err = err.WithMessage(e1.DescripcionErrorRegistro).WithCode(e1.CodigoErrorRegistro)
		}
		return err
	}

	return nil
}
