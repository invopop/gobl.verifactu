package doc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// TipoHuella is the SHA-256 fingerprint type for Verifactu - L12
// Might include support for other encryption types in the future.
const TipoHuella = "01"

// ChainData contains the fields of this invoice that will be
// required for fingerprinting the next invoice. JSON tags are
// provided to help with serialization.
type ChainData struct {
	IDIssuer    string `json:"issuer"`
	NumSeries   string `json:"num_series"`
	IssueDate   string `json:"issue_date"`
	Fingerprint string `json:"fingerprint"`
}

func formatField(key, value string) string {
	value = strings.TrimSpace(value) // Remove whitespace
	if value == "" {
		return fmt.Sprintf("%s=", key)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

func (d *Envelope) fingerprintAlta(inv *RegistroAlta) error {
	var h string
	if inv.Encadenamiento.PrimerRegistro == "S" {
		h = ""
	} else {
		h = inv.Encadenamiento.RegistroAnterior.Huella
	}
	f := []string{
		formatField("IDEmisorFactura", inv.IDFactura.IDEmisorFactura),
		formatField("NumSerieFactura", inv.IDFactura.NumSerieFactura),
		formatField("FechaExpedicionFactura", inv.IDFactura.FechaExpedicionFactura),
		formatField("TipoFactura", inv.TipoFactura),
		formatField("CuotaTotal", inv.CuotaTotal),
		formatField("ImporteTotal", inv.ImporteTotal),
		formatField("Huella", h),
		formatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	st := strings.Join(f, "&")
	hash := sha256.New()
	hash.Write([]byte(st))

	d.Body.VeriFactu.RegistroFactura.RegistroAlta.Huella = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	return nil
}

func (d *Envelope) fingerprintAnulacion(inv *RegistroAnulacion) error {
	var h string
	if inv.Encadenamiento.PrimerRegistro == "S" {
		h = ""
	} else {
		h = inv.Encadenamiento.RegistroAnterior.Huella
	}
	f := []string{
		formatField("IDEmisorFacturaAnulada", inv.IDFactura.IDEmisorFactura),
		formatField("NumSerieFacturaAnulada", inv.IDFactura.NumSerieFactura),
		formatField("FechaExpedicionFacturaAnulada", inv.IDFactura.FechaExpedicionFactura),
		formatField("Huella", h),
		formatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	st := strings.Join(f, "&")
	hash := sha256.New()
	hash.Write([]byte(st))

	d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.Huella = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	return nil
}

func (d *Envelope) generateHashAlta(prev *ChainData) error {
	if prev == nil {
		d.Body.VeriFactu.RegistroFactura.RegistroAlta.Encadenamiento = &Encadenamiento{
			PrimerRegistro: "S",
		}
		if err := d.fingerprintAlta(d.Body.VeriFactu.RegistroFactura.RegistroAlta); err != nil {
			return err
		}
		return nil
	}
	// Concatenate f according to Verifactu specifications
	d.Body.VeriFactu.RegistroFactura.RegistroAlta.Encadenamiento = &Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        prev.IDIssuer,
			NumSerieFactura:        prev.NumSeries,
			FechaExpedicionFactura: prev.IssueDate,
			Huella:                 prev.Fingerprint,
		},
	}
	if err := d.fingerprintAlta(d.Body.VeriFactu.RegistroFactura.RegistroAlta); err != nil {
		return err
	}
	return nil
}

func (d *Envelope) generateHashAnulacion(prev *ChainData) error {
	if prev == nil {
		d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.Encadenamiento = &Encadenamiento{
			PrimerRegistro: "S",
		}
		if err := d.fingerprintAnulacion(d.Body.VeriFactu.RegistroFactura.RegistroAnulacion); err != nil {
			return err
		}
		return nil
	}
	d.Body.VeriFactu.RegistroFactura.RegistroAnulacion.Encadenamiento = &Encadenamiento{
		RegistroAnterior: &RegistroAnterior{
			IDEmisorFactura:        prev.IDIssuer,
			NumSerieFactura:        prev.NumSeries,
			FechaExpedicionFactura: prev.IssueDate,
			Huella:                 prev.Fingerprint,
		},
	}
	if err := d.fingerprintAnulacion(d.Body.VeriFactu.RegistroFactura.RegistroAnulacion); err != nil {
		return err
	}
	return nil
}
