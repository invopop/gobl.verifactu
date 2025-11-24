package verifactu

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/xmldsig"
)

// Environment defines the environment to use for connections
type Environment string

// Environment to use for connections
const (
	EnvironmentProduction Environment = "production"
	EnvironmentSandbox    Environment = "sandbox"

	BaseURLProduction         = "https://www1.agenciatributaria.gob.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"
	BaseURLProductionWithSeal = "https://www10.agenciatributaria.gob.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"
	BaseURLTesting            = "https://prewww1.aeat.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"
	BaseURLTestingWithSeal    = "https://prewww10.aeat.es/wlpl/TIKE-CONT/ws/SistemaFacturacion/VerifactuSOAP"
)

// connection defines what is expected from a connection to a gateway.
type connection struct {
	client *resty.Client
}

// newConnection instantiates and configures a new connection to the VeriFactu gateway.
func newConnection(env Environment, cert *xmldsig.Certificate, withSeal bool) (*connection, error) {
	// Prepare the tls configuration
	tlsConf, err := cert.TLSAuthConfig()
	if err != nil {
		return nil, ErrValidation.WithMessage(fmt.Errorf("preparing TLS config: %v", err).Error())
	}
	certs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("preparing cert pool: %w", err)
	}
	tlsConf.RootCAs = certs
	tlsConf.Renegotiation = tls.RenegotiateOnceAsClient

	c := new(connection)
	c.client = resty.New()

	switch env {
	case EnvironmentProduction:
		if withSeal {
			c.client.SetBaseURL(BaseURLProductionWithSeal)
		} else {
			c.client.SetBaseURL(BaseURLProduction)
		}
	default:
		if withSeal {
			c.client.SetBaseURL(BaseURLTestingWithSeal)
		} else {
			c.client.SetBaseURL(BaseURLTesting)
		}
	}
	// tlsConf.InsecureSkipVerify = true
	c.client.SetTLSClientConfig(tlsConf)
	c.client.SetDebug(os.Getenv("DEBUG") == "true")
	return c, nil
}

func (c *connection) post(ctx context.Context, payload []byte) (*EnvelopeResponse, error) {
	out := new(EnvelopeResponse)
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/xml").
		SetContentLength(true).
		SetBody(payload).
		SetResult(out)

	res, err := req.Post("")
	if err != nil {
		return nil, err
	}

	if out.Body.Fault != nil {
		if out.Body.Fault.Code == "env:Server" {
			return nil, ErrServer.WithMessage(out.Body.Fault.Message)
		}
		return nil, ErrValidation.WithMessage(out.Body.Fault.Message).WithCode(out.Body.Fault.Code)
	}
	if res.StatusCode() != http.StatusOK {
		return nil, ErrValidation.WithCode(strconv.Itoa(res.StatusCode())).WithMessage(res.String())
	}

	return out, nil
}
