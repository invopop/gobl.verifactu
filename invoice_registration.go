package verifactu

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/nbio/xml"
)

var correctiveCodes = []cbc.Code{ // Credit or Debit notes
	"R1", "R2", "R3", "R4", "R5",
}

// InvoiceRegistration contains the details of an invoice registration
type InvoiceRegistration struct {
	XMLName                             xml.Name              `xml:"sum1:RegistroAlta"`
	NS                                  string                `xml:"xmlns:sum1,attr,omitempty"`
	IDVersion                           string                `xml:"sum1:IDVersion"`
	IDFactura                           *IDFactura            `xml:"sum1:IDFactura"`
	RefExterna                          string                `xml:"sum1:RefExterna,omitempty"`
	NombreRazonEmisor                   string                `xml:"sum1:NombreRazonEmisor"`
	Subsanacion                         string                `xml:"sum1:Subsanacion,omitempty"`
	RechazoPrevio                       string                `xml:"sum1:RechazoPrevio,omitempty"`
	TipoFactura                         string                `xml:"sum1:TipoFactura"`
	TipoRectificativa                   string                `xml:"sum1:TipoRectificativa,omitempty"`
	FacturasRectificadas                []*FacturaRectificada `xml:"sum1:FacturasRectificadas,omitempty"`
	FacturasSustituidas                 []*FacturaSustituida  `xml:"sum1:FacturasSustituidas,omitempty"`
	ImporteRectificacion                *ImporteRectificacion `xml:"sum1:ImporteRectificacion,omitempty"`
	FechaOperacion                      string                `xml:"sum1:FechaOperacion,omitempty"`
	DescripcionOperacion                string                `xml:"sum1:DescripcionOperacion"`
	FacturaSimplificadaArt7273          string                `xml:"sum1:FacturaSimplificadaArt7273,omitempty"`
	FacturaSinIdentifDestinatarioArt61d string                `xml:"sum1:FacturaSinIdentifDestinatarioArt61d,omitempty"`
	Macrodato                           string                `xml:"sum1:Macrodato,omitempty"`
	EmitidaPorTerceroODestinatario      string                `xml:"sum1:EmitidaPorTerceroODestinatario,omitempty"`
	Tercero                             *Party                `xml:"sum1:Tercero,omitempty"`
	Destinatarios                       []*Destinatario       `xml:"sum1:Destinatarios,omitempty"`
	Cupon                               string                `xml:"sum1:Cupon,omitempty"`
	Desglose                            *Desglose             `xml:"sum1:Desglose"`
	CuotaTotal                          string                `xml:"sum1:CuotaTotal"`
	ImporteTotal                        string                `xml:"sum1:ImporteTotal"`
	Encadenamiento                      *Encadenamiento       `xml:"sum1:Encadenamiento"`
	SistemaInformatico                  *Software             `xml:"sum1:SistemaInformatico"`
	FechaHoraHusoGenRegistro            string                `xml:"sum1:FechaHoraHusoGenRegistro"`
	NumRegistroAcuerdoFacturacion       string                `xml:"sum1:NumRegistroAcuerdoFacturacion,omitempty"`
	IdAcuerdoSistemaInformatico         string                `xml:"sum1:IdAcuerdoSistemaInformatico,omitempty"` //nolint:revive
	TipoHuella                          string                `xml:"sum1:TipoHuella"`
	Huella                              string                `xml:"sum1:Huella"`
	// Signature                           *xmldsig.Signature   `xml:"sum1:Signature,omitempty"`
}

// IDFactura contains the identifying information for an invoice
type IDFactura struct {
	IDEmisorFactura        string `xml:"sum1:IDEmisorFactura"`
	NumSerieFactura        string `xml:"sum1:NumSerieFactura"`
	FechaExpedicionFactura string `xml:"sum1:FechaExpedicionFactura"`
}

// FacturaRectificada represents a rectified invoice
type FacturaRectificada struct {
	IDFactura IDFactura `xml:"sum1:IDFacturaRectificada"`
}

// FacturaSustituida represents a substituted invoice
type FacturaSustituida struct {
	IDFactura IDFactura `xml:"sum1:IDFacturaSustituida"`
}

// ImporteRectificacion contains rectification amounts
type ImporteRectificacion struct {
	BaseRectificada         num.Amount `xml:"sum1:BaseRectificada"`
	CuotaRectificada        num.Amount `xml:"sum1:CuotaRectificada"`
	CuotaRecargoRectificado num.Amount `xml:"sum1:CuotaRecargoRectificado,omitempty"`
}

// Party represents a in the document, covering fields Generador, Tercero and IDDestinatario
type Party struct {
	NombreRazon string  `xml:"sum1:NombreRazon"`
	NIF         string  `xml:"sum1:NIF,omitempty"`
	IDOtro      *IDOtro `xml:"sum1:IDOtro,omitempty"`
}

// Destinatario represents a recipient in the document
type Destinatario struct {
	IDDestinatario *Party `xml:"sum1:IDDestinatario"`
}

// IDOtro contains alternative identifying information
type IDOtro struct {
	CodigoPais string `xml:"sum1:CodigoPais"`
	IDType     string `xml:"sum1:IDType"`
	ID         string `xml:"sum1:ID"`
}

// Desglose contains the breakdown details
type Desglose struct {
	DetalleDesglose []*DetalleDesglose `xml:"sum1:DetalleDesglose"`
}

// DetalleDesglose contains detailed breakdown information
type DetalleDesglose struct {
	Impuesto                      string `xml:"sum1:Impuesto,omitempty"`
	ClaveRegimen                  string `xml:"sum1:ClaveRegimen,omitempty"`
	CalificacionOperacion         string `xml:"sum1:CalificacionOperacion,omitempty"`
	OperacionExenta               string `xml:"sum1:OperacionExenta,omitempty"`
	TipoImpositivo                string `xml:"sum1:TipoImpositivo,omitempty"`
	BaseImponibleOImporteNoSujeto string `xml:"sum1:BaseImponibleOimporteNoSujeto"`
	BaseImponibleACoste           string `xml:"sum1:BaseImponibleACoste,omitempty"`
	CuotaRepercutida              string `xml:"sum1:CuotaRepercutida,omitempty"`
	TipoRecargoEquivalencia       string `xml:"sum1:TipoRecargoEquivalencia,omitempty"`
	CuotaRecargoEquivalencia      string `xml:"sum1:CuotaRecargoEquivalencia,omitempty"`
}

// newInvoiceRegistration creates a new VeriFactu registration for an invoice.
func newInvoiceRegistration(inv *bill.Invoice, ts time.Time, r IssuerRole, s *Software) (*InvoiceRegistration, error) {
	tf, err := getTaxExtKey(inv, verifactu.ExtKeyDocType)
	if err != nil {
		return nil, err
	}

	desc := newDescription(inv)
	dg, err := newDesglose(inv)
	if err != nil {
		return nil, err
	}

	reg := &InvoiceRegistration{
		NS:        SUM1, // to remove during sending
		IDVersion: CurrentVersion,
		IDFactura: &IDFactura{
			IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
			NumSerieFactura:        invoiceNumber(inv.Series, inv.Code),
			FechaExpedicionFactura: inv.IssueDate.Time().Format("02-01-2006"),
		},
		NombreRazonEmisor:        inv.Supplier.Name,
		TipoFactura:              tf,
		DescripcionOperacion:     desc,
		Desglose:                 dg,
		CuotaTotal:               newTotalTaxes(inv).String(),
		ImporteTotal:             newImporteTotal(inv).String(),
		SistemaInformatico:       s,
		FechaHoraHusoGenRegistro: formatDateTimeZone(ts),
		TipoHuella:               TipoHuella,
	}

	// Prepare the customer, but only if there are enough details, otherwise
	// we consider this to be a simplified or B2C invoice.
	if p := newParty(inv.Customer); p != nil {
		reg.Destinatarios = []*Destinatario{
			{
				IDDestinatario: p,
			},
		}
	} else {
		reg.FacturaSinIdentifDestinatarioArt61d = "S"
	}

	if inv.Tax.Ext[verifactu.ExtKeyDocType].In(correctiveCodes...) {
		k, err := getTaxExtKey(inv, verifactu.ExtKeyCorrectionType)
		if err != nil {
			return nil, err
		}
		reg.TipoRectificativa = k

		list := make([]*FacturaRectificada, len(inv.Preceding))
		taxes := new(tax.Total)
		for i, ref := range inv.Preceding {
			if ref.Tax != nil {
				taxes = taxes.Merge(ref.Tax)
			}
			list[i] = &FacturaRectificada{
				IDFactura: IDFactura{
					IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
					NumSerieFactura:        invoiceNumber(ref.Series, ref.Code),
					FechaExpedicionFactura: ref.IssueDate.Time().Format("02-01-2006"),
				},
			}
		}
		reg.FacturasRectificadas = list
		if k == "S" {
			// only include in substituted documents
			reg.ImporteRectificacion = newImporteRectificacion(taxes)
		}
	}

	// F3 covers the special use-case of full invoices that replace a
	// previous simplified document. This is the only time the "FacturaSustituida"
	// field is used.
	if reg.TipoFactura == "F3" {
		if inv.Preceding != nil {
			subs := make([]*FacturaSustituida, 0, len(inv.Preceding))
			for _, ref := range inv.Preceding {
				subs = append(subs, &FacturaSustituida{
					IDFactura: IDFactura{
						IDEmisorFactura:        inv.Supplier.TaxID.Code.String(),
						NumSerieFactura:        invoiceNumber(ref.Series, ref.Code),
						FechaExpedicionFactura: ref.IssueDate.Time().Format("02-01-2006"),
					},
				})
			}
			reg.FacturasSustituidas = subs
		}
	}

	if r == IssuerRoleThirdParty {
		reg.EmitidaPorTerceroODestinatario = "T"
		reg.Tercero = newParty(inv.Supplier)
	}

	// Flag for operations with totals over 100,000,000€. Added with optimism.
	if inv.Totals.TotalWithTax.Compare(num.MakeAmount(100000000, 0)) == 1 {
		reg.Macrodato = "S"
	}

	return reg, nil
}

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return code.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

func newDescription(inv *bill.Invoice) string {
	for _, note := range inv.Notes {
		if note.Key == org.NoteKeyGeneral {
			return note.Text
		}
	}

	var desc string
	// Iterate over invoice lines to build a description
	for i, line := range inv.Lines {
		// Only add an item name if it exists
		if line != nil && line.Item != nil && line.Item.Name != "" {
			// If the description is too long, we need to stop the loop
			if len(desc)+len(line.Item.Name)+3 > 500 {
				// If the description is not empty, add an ellipsis
				// This could happen if the item name length > 488
				if desc != "" {
					desc = desc[:len(desc)-2] + "..."
				}
				break
			}
			// Add the name and a comma if not the last line
			desc += line.Item.Name
			if i < len(inv.Lines)-1 {
				desc += ", "
			} else {
				desc += "."
			}
		}
	}

	if desc == "" {
		desc += "Sin descripción"
	}

	return desc
}

func newImporteTotal(inv *bill.Invoice) num.Amount {
	if inv.Totals.Taxes == nil {
		// This is likely to be wrong as all Spanish invoices need to account
		// for tax, even if exempt.
		return inv.Totals.Total
	}
	t := num.MakeAmount(0, 2)
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			// We need to recalculate the total based on the taxable bases
			for _, rate := range category.Rates {
				t = t.Add(rate.Base)
			}
			t = t.Add(category.Amount)
		}
	}
	return t
}

func newImporteRectificacion(taxes *tax.Total) *ImporteRectificacion {
	zero := currency.EUR.Def().Zero()
	ir := &ImporteRectificacion{
		BaseRectificada:         zero,
		CuotaRectificada:        zero,
		CuotaRecargoRectificado: zero,
	}
	for _, cat := range taxes.Categories {
		if cat.Code == tax.CategoryVAT {
			for _, rate := range cat.Rates {
				ir.BaseRectificada = ir.BaseRectificada.Add(rate.Base)
			}
			ir.CuotaRectificada = ir.CuotaRectificada.Add(cat.Amount)
			if cat.Surcharge != nil {
				ir.CuotaRecargoRectificado = ir.CuotaRecargoRectificado.Add(*cat.Surcharge)
			}
		}
	}
	return ir
}

func newTotalTaxes(inv *bill.Invoice) num.Amount {
	totalTaxes := num.MakeAmount(0, 2)
	if inv.Totals.Taxes == nil {
		return totalTaxes
	}
	for _, category := range inv.Totals.Taxes.Categories {
		if !category.Retained {
			totalTaxes = totalTaxes.Add(category.Amount)
		}
	}
	return totalTaxes
}

func getTaxExtKey(inv *bill.Invoice, k cbc.Key) (string, error) {
	if inv.Tax == nil || inv.Tax.Ext == nil || inv.Tax.Ext[k].String() == "" {
		return "", validation.Errors{
			"tax": validation.Errors{
				"ext": validation.Errors{
					k.String(): errors.New("required"),
				},
			},
		}
	}
	return inv.Tax.Ext[k].String(), nil
}

// fingerprint will add a fingerprint to the regisration line using the previous
// chain data entry details.
func (r *InvoiceRegistration) fingerprint(prev *ChainData) {
	h := ""
	if prev == nil {
		r.Encadenamiento = &Encadenamiento{
			PrimerRegistro: "S",
		}
	} else {
		r.Encadenamiento = &Encadenamiento{
			RegistroAnterior: &RegistroAnterior{
				IDEmisorFactura:        prev.IDIssuer,
				NumSerieFactura:        prev.NumSeries,
				FechaExpedicionFactura: prev.IssueDate,
				Huella:                 prev.Fingerprint,
			},
		}
		h = prev.Fingerprint
	}

	f := []string{
		formatChainField("IDEmisorFactura", r.IDFactura.IDEmisorFactura),
		formatChainField("NumSerieFactura", r.IDFactura.NumSerieFactura),
		formatChainField("FechaExpedicionFactura", r.IDFactura.FechaExpedicionFactura),
		formatChainField("TipoFactura", r.TipoFactura),
		formatChainField("CuotaTotal", r.CuotaTotal),
		formatChainField("ImporteTotal", r.ImporteTotal),
		formatChainField("Huella", h),
		formatChainField("FechaHoraHusoGenRegistro", r.FechaHoraHusoGenRegistro),
	}
	st := strings.Join(f, "&")
	hash := sha256.New()
	hash.Write([]byte(st))

	r.Huella = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
}

// ChainData provides the details for this registration entry.
func (r *InvoiceRegistration) ChainData() *ChainData {
	return &ChainData{
		IDIssuer:    r.IDFactura.IDEmisorFactura,
		NumSeries:   r.IDFactura.NumSerieFactura,
		IssueDate:   r.IDFactura.FechaExpedicionFactura,
		Fingerprint: r.Huella,
	}
}

// Bytes prepares an indendented XML document suitable for persistence.
func (r *InvoiceRegistration) Bytes() ([]byte, error) {
	return toBytesIndent(r)
}
