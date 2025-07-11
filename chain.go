package verifactu

import (
	"fmt"
	"strings"
)

// TipoHuella is the SHA-256 fingerprint type for Verifactu - L12
// Might include support for other encryption types in the future.
const TipoHuella = "01"

// ChainData contains the fields of this invoice that will be
// required for fingerprinting the _next_ invoice. JSON tags are
// provided to help with serialization.
type ChainData struct {
	IDIssuer    string `json:"issuer"`
	NumSeries   string `json:"num_series"`
	IssueDate   string `json:"issue_date"`
	Fingerprint string `json:"fingerprint"`
}

// Encadenamiento contains chaining information between documents
type Encadenamiento struct {
	PrimerRegistro   string            `xml:"sum1:PrimerRegistro,omitempty"`
	RegistroAnterior *RegistroAnterior `xml:"sum1:RegistroAnterior,omitempty"`
}

// RegistroAnterior contains information about the previous registration
type RegistroAnterior struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
	Huella                 string `xml:"sum1:Huella"`
}

// Software contains the details about the software that is using this library to
// generate VeriFactu documents. These details are included in the final
// document.
type Software struct {
	NombreRazon                 string `xml:"sum1:NombreRazon"`
	NIF                         string `xml:"sum1:NIF"`
	NombreSistemaInformatico    string `xml:"sum1:NombreSistemaInformatico"`
	IdSistemaInformatico        string `xml:"sum1:IdSistemaInformatico"` //nolint:revive,staticcheck
	Version                     string `xml:"sum1:Version"`
	NumeroInstalacion           string `xml:"sum1:NumeroInstalacion"`
	TipoUsoPosibleSoloVerifactu string `xml:"sum1:TipoUsoPosibleSoloVerifactu,omitempty"`
	TipoUsoPosibleMultiOT       string `xml:"sum1:TipoUsoPosibleMultiOT,omitempty"`
	IndicadorMultiplesOT        string `xml:"sum1:IndicadorMultiplesOT,omitempty"`
}

// formatChainField is a helper method to help prepare an entry in the chain's
// string used for hashing.
func formatChainField(key, value string) string {
	value = strings.TrimSpace(value) // Remove whitespace
	if value == "" {
		return fmt.Sprintf("%s=", key)
	}
	return fmt.Sprintf("%s=%s", key, value)
}
