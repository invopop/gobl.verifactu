package verifactu

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// FingerprintType is the SHA-256 fingerprint type for Verifactu - L12. Might include
// support for other encryption types in the future.
const FingerprintType = "01"

// ChainData contains the fields of this invoice that will be required for fingerprinting
// the _next_ invoice. JSON tags are provided to help with serialization.
type ChainData struct {
	IDIssuer    string `json:"issuer"`
	NumSeries   string `json:"num_series"`
	IssueDate   string `json:"issue_date"`
	Fingerprint string `json:"fingerprint"`
}

// Encadenamiento contains chaining information between invoice documents
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

// EventChainData contains the fields of this event that will be required for
// fingerprinting the _next_ event. JSON tags are provided to help with serialization.
type EventChainData struct {
	EventType           string `json:"event_type"`
	GenerationTimestamp string `json:"generation_timestamp"`
	Fingerprint         string `json:"fingerprint"`
}

// EventChaining contains chaining information between event registrations
type EventChaining struct {
	FirstEvent    string         `xml:"sf:PrimerEvento,omitempty"`
	PreviousEvent *PreviousEvent `xml:"sf:EventoAnterior,omitempty"`
}

// PreviousEvent contains information about the previous event registration
// used for chaining.
type PreviousEvent struct {
	EventType           string `xml:"sf:TipoEvento"`
	GenerationTimestamp string `xml:"sf:FechaHoraHusoGenEvento"`
	Fingerprint         string `xml:"sf:HuellaEvento"`
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
	NumeroInstalacion           string `xml:"sum1:NumeroInstalacion"` // may need to be overridden at run time
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

// computeFingerprint joins the provided chain fields with "&", computes the SHA-256 hash,
// and returns the result as an uppercase hex string.
func computeFingerprint(fields []string) string {
	st := strings.Join(fields, "&")
	hash := sha256.New()
	hash.Write([]byte(st))
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}
