package doc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// FormatField returns a formatted field as key=value or key= if the value is empty.
func FormatField(key, value string) string {
	value = strings.TrimSpace(value) // Remove whitespace
	if value == "" {
		return fmt.Sprintf("%s=", key)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

// ConcatenateFields builds the concatenated string based on Verifactu requirements.
func makeRegistroAltaFields(inv *RegistroAlta) string {
	fields := []string{
		FormatField("IDEmisorFactura", inv.IDFactura.IDEmisorFactura),
		FormatField("NumSerieFactura", inv.IDFactura.NumSerieFactura),
		FormatField("FechaExpedicionFactura", inv.IDFactura.FechaExpedicionFactura),
		FormatField("TipoFactura", inv.TipoFactura),
		FormatField("CuotaTotal", inv.CuotaTotal),
		FormatField("ImporteTotal", inv.ImporteTotal),
		FormatField("Huella", inv.Encadenamiento.RegistroAnterior.Huella),
		FormatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	return strings.Join(fields, "&")
}

func makeRegistroAnulacionFields(inv *RegistroAnulacion) string {
	fields := []string{
		FormatField("IDEmisorFactura", inv.IDFactura.IDEmisorFactura),
		FormatField("NumSerieFactura", inv.IDFactura.NumSerieFactura),
		FormatField("FechaExpedicionFactura", inv.IDFactura.FechaExpedicionFactura),
		FormatField("Huella", inv.Encadenamiento.RegistroAnterior.Huella),
		FormatField("FechaHoraHusoGenRegistro", inv.FechaHoraHusoGenRegistro),
	}
	return strings.Join(fields, "&")
}

// GenerateHash generates the SHA-256 hash for the invoice data.
func GenerateHash(inv *RegistroFactura) string {
	// Concatenate fields according to Verifactu specifications
	var concatenatedString string
	if inv.RegistroAlta != nil {
		concatenatedString = makeRegistroAltaFields(inv.RegistroAlta)
	} else if inv.RegistroAnulacion != nil {
		concatenatedString = makeRegistroAnulacionFields(inv.RegistroAnulacion)
	}

	// Convert to UTF-8 byte array and hash it with SHA-256
	hash := sha256.New()
	hash.Write([]byte(concatenatedString))

	// Convert the hash to hexadecimal and make it uppercase
	return strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}
