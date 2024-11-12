package doc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// TipoHuella is the SHA-256 fingerprint type for Verifactu - L12
const TipoHuella = "01"

// FormatField returns a formatted field as key=value or key= if the value is empty.
func FormatField(key, value string) string {
	value = strings.TrimSpace(value) // Remove whitespace
	if value == "" {
		return fmt.Sprintf("%s=", key)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// Concatenatef builds the concatenated string based on Verifactu requirements.
func (d *VeriFactu) fingerprintAlta(inv *RegistroAlta) error {
	f := []string{
		FormatField("IDEmisorFactura", inv.IDFactura.IDEmisorFactura),
		FormatField("NumSerieFactura", inv.IDFactura.NumSerieFactura),
		FormatField("FechaExpedicionFactura", inv.IDFactura.FechaExpedicionFactura),
		FormatField("TipoFactura", inv.TipoFactura),
		FormatField("CuotaTotal", inv.CuotaTotal),
		FormatField("ImporteTotal", inv.ImporteTotal),
		FormatField("Huella", inv.Encadenamiento.RegistroAnterior.Huella),
		FormatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	st := strings.Join(f, "&")
	hash := sha256.New()
	hash.Write([]byte(st))

	d.RegistroFactura.RegistroAlta.Huella = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	return nil
}

func (d *VeriFactu) fingerprintAnulacion(inv *RegistroAnulacion) error {
	f := []string{
		FormatField("IDEmisorFactura", inv.IDFactura.IDEmisorFactura),
		FormatField("NumSerieFactura", inv.IDFactura.NumSerieFactura),
		FormatField("FechaExpedicionFactura", inv.IDFactura.FechaExpedicionFactura),
		FormatField("Huella", inv.Encadenamiento.RegistroAnterior.Huella),
		FormatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	st := strings.Join(f, "&")
	hash := sha256.New()
	hash.Write([]byte(st))

	d.RegistroFactura.RegistroAnulacion.Huella = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	return nil
}

// GenerateHash generates the SHA-256 hash for the invoice data.
func (d *VeriFactu) GenerateHash(prev *Encadenamiento) error {
	if prev == nil {
		return fmt.Errorf("previous document is required")
	}
	// Concatenate f according to Verifactu specifications
	if d.RegistroFactura.RegistroAlta != nil {
		d.RegistroFactura.RegistroAlta.Encadenamiento = prev
		if err := d.fingerprintAlta(d.RegistroFactura.RegistroAlta); err != nil {
			return err
		}
	} else if d.RegistroFactura.RegistroAnulacion != nil {
		d.RegistroFactura.RegistroAnulacion.Encadenamiento = prev
		if err := d.fingerprintAnulacion(d.RegistroFactura.RegistroAnulacion); err != nil {
			return err
		}
	}

	return nil
}
