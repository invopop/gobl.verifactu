package verifactu

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/xmldsig"
	"github.com/invopop/xmldsig/profiles/verifactu"
)

// SignDocument signs any XML-marshalable struct using the VeriFactu XAdES
// configuration. It can be used to sign registration, cancellation, and event
// records (or any other document type) with the provided certificate.
func SignDocument(doc any, cert *xmldsig.Certificate, opts ...xmldsig.Option) (*xmldsig.Signature, error) {
	data, err := xml.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("marshaling document: %w", err)
	}

	signOpts := []xmldsig.Option{
		xmldsig.WithXMLDSigConfig(verifactu.XMLDSigConfig()),
		xmldsig.WithXAdESConfig(verifactu.XAdESConfig()),
		xmldsig.WithCertificate(cert),
	}
	signOpts = append(signOpts, opts...)

	sig, err := xmldsig.Sign(data, signOpts...)
	if err != nil {
		return nil, fmt.Errorf("signing document: %w", err)
	}

	return sig, nil
}
