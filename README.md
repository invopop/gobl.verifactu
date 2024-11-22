# GOBL to Veri*Factu

Go library to convert [GOBL](https://github.com/invopop/gobl) invoices into Veri*Factu declarations and send them to the AEAT (Agencia Estatal de Administración Tributaria) web services. This library assumes that clients will handle a local database of previous invoices in order to comply with the local requirements of chaining all invoices together.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [GNU Affero General Public License v3.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). For contributions to this library to be accepted, we will require you to accept transferring your copyright to Invopop Ltd.

## Source

The main resources for Veri*Factu can be found in the AEAT website and include:

- [Veri*Factu documentation](https://www.agenciatributaria.es/AEAT.desarrolladores/Desarrolladores/_menu_/Documentacion/Sistemas_Informaticos_de_Facturacion_y_Sistemas_VERI_FACTU/Sistemas_Informaticos_de_Facturacion_y_Sistemas_VERI_FACTU.html)
- [Veri*Factu Ministerial Order](https://www.boe.es/diario_boe/txt.php?id=BOE-A-2024-22138)

## Usage

### Go Package

You must have first created a GOBL Envelope containing an Invoice that you'd like to send to the AEAT. For the document to be converted, the supplier contained in the invoice should have a "Tax ID" with the country set to `ES`.

The following is an example of how the GOBL Veri*Factu package could be used:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/gobl"
	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/invopop/xmldsig"
)

func main() {
	ctx := context.Background()

	// Load sample envelope:
	data, _ := os.ReadFile("./test/data/sample-invoice.json")

	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err != nil {
		panic(err)
	}

	// Prepare software configuration:
	soft := &verifactu.Software{
		NombreRazon:              "Company LTD",    // Company Name
		NIF:                      "B123456789",     // Software company's tax code
		NombreSistemaInformatico: "Software Name",  // Name of application
		IdSistemaInformatico:     "A1",             // Software ID
		Version:                  "1.0",            // Software version
		NumeroInstalacion:        "00001",          // Software installation number
	}

	// Load the certificate
	cert, err := xmldsig.LoadCertificate(c.cert, c.password)
	if err != nil {
		return err
	}

	// Create the client with the software and certificate
	opts := []verifactu.Option{
		verifactu.WithCertificate(cert),
		verifactu.WithSupplierIssuer(), // The issuer can be either the supplier, the
		verifactu.InTesting(),
	}


	tc, err := verifactu.New(c.software(), opts...)
	if err != nil {
		return err
	}

	// Convert the GOBL envelope to a Veri*Factu document
	td, err := tc.Convert(env)
	if err != nil {
		return err
	}

	// Prepare the previous document chain data
	c.previous = `{
		"emisor": "B85905495",
		"serie": "SAMPLE-001", 
		"fecha": "11-11-2024",
		"huella": "13EC0696104D1E529667184C6CDFC67D08036BCA4CD1B7887DE9C6F8F7EEC69C"
		}`
	prev := new(doc.ChainData)
	if err := json.Unmarshal([]byte(c.previous), prev); err != nil {
		return err
	}

	// Create the document fingerprint based on the previous document chain
	err = tc.Fingerprint(td, prev)
	if err != nil {
		return err
	}

	// Add the QR code to the document
	if err := tc.AddQR(td, env); err != nil {
		return err
	}

	out, err := c.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	convOut, err := td.BytesIndent()
	if err != nil {
		return fmt.Errorf("generating verifactu xml: %w", err)
	}


	err = tc.Post(cmd.Context(), td)
	if err != nil {
		return err
	}

	data, err := json.Marshal(td.ChainData())
	if err != nil {
		return err
	}
	fmt.Printf("Generated document with fingerprint: \n%s\n", string(data))

	return nil

}
```

## Command Line

The GOBL Veri*Factu package tool also includes a command line helper. You can install manually in your Go environment with:

```bash
go install github.com/invopop/gobl.verifactu
```

We recommend using a `.env` file to prepare configuration settings, although all parameters can be set using command line flags. Heres an example:

```
SOFTWARE_COMPANY_NIF=B85905495
SOFTWARE_COMPANY_NAME=Invopop S.L.
SOFTWARE_NAME=gobl.verifactu
SOFTWARE_VERSION=1.0
SOFTWARE_ID_SISTEMA_INFORMATICO=A1
SOFTWARE_NUMERO_INSTALACION=00001

CERTIFICATE_PATH=./xxxxxxxxx.p12
CERTIFICATE_PASSWORD=xxxxxxxx

```

To convert a document to XML, run:

```bash
gobl.verifactu convert ./test/data/sample-invoice.json
```

To submit to the tax agency testing environment:

```bash
gobl.verifactu send ./test/data/sample-invoice.json
```

## Limitations

- Veri*Factu allows more than one customer per invoice, but GOBL only has one possible customer.

- Invoices should have a note of type general that will be used as a general description of the invoice. If an invoice is missing this info, it will be rejected with an error.

- Currently VeriFactu supportts sending more than one invoice at a time (up to 1000). However, this module only currently supports 1 invoice at a time.

## Tags, Keys and Extensions

In order to provide the supplier specific data required by Veri*Factu, invoices need to include a bit of extra data. We've managed to simplify these into specific cases.

### Tax Tags

Invoice tax tags can be added to invoice documents in order to reflect a special situation. The following schemes are supported:

- `simplified-scheme` - a retailer operating under a simplified tax regime (regimen simplificado) that must indicate that all of their sales are under this scheme. This implies that all operations in the invoice will have the `FacturaSinIdentifDestinatarioArt61d` tag set to `S`.
- `reverse-charge` - B2B services or goods sold to a tax registered EU member who will pay VAT on the suppliers behalf. Implies that all items will be classified under the `TipoNoExenta` value of `S2`.

## Tax Extensions

The following extension can be applied to each line tax:

- `es-verifactu-doc-type` – defines the type of invoice being sent. In most cases this will be set automatically by the GOBL add-on. These are the valid values:

  - `F1` - Standard invoice.
  - `F2` - Simplified invoice.
  - `F3` - Invoice in substitution of simplified invoices.
  - `R1` - Rectified invoice based on law and Article 80.1, 80.2 and 80.6 in the Spanish VAT Law ([LIVA](https://www.boe.es/buscar/act.php?id=BOE-A-1992-28740)).
  - `R2` - Rectified invoice based on law and Article 80.3.
  - `R3` - Rectified invoice based on law and Article 80.4.
  - `R4` - Rectified invoice based on law and other reasons.
  - `R5` - Rectified invoice based on simplified invoices.

- `es-verifactu-tax-classification` - combines the tax classification and exemption codes used in Veri*Factu. These are the valid values:

  - `S1` - Subject and not exempt - Without reverse charge
  - `S2` - Subject and not exempt - With reverse charge
  - `N1` - Not subject - Articles 7, 14, others
  - `N2` - Not subject - Due to location rules
  - `E1` - Exempt pursuant to Article 20 of the VAT Law
  - `E2` - Exempt pursuant to Article 21 of the VAT Law
  - `E3` - Exempt pursuant to Article 22 of the VAT Law
  - `E4` - Exempt pursuant to Articles 23 and 24 of the VAT Law
  - `E5` - Exempt pursuant to Article 25 of the VAT Law
  - `E6` - Exempt for other reasons


### Use-Cases

Under what situations should the TicketBAI system be expected to function:

- B2B & B2C: regular national invoice with VAT. Operation with minimal data.
- B2B Provider to Retailer: Include equalisation surcharge VAT rates
- B2B Retailer: Same as regular invoice, except with invoice lines that include `ext[es-tbai-product] = resale` when the goods being provided are being sold without modification (recargo de equivalencia), very much related to the next point.
- B2B Retailer Simplified: Include the simplified scheme key. (This implies that the `OperacionEnRecargoDeEquivalenciaORegimenSimplificado` tag will be set to `S`).
- EU B2B: Reverse charge EU export, scheme: reverse-charge taxes calculated, but not applied to totals. By default all line items assumed to be services. Individual lines can use the `ext[es-tbai-product] = goods` value to identify when the line is a physical good. Operations like this are normally assigned the TipoNoExenta value of S2. If however the service or goods are exempt of tax, each line's tax `ext[exempt]` field can be used to identify a reason.
- EU B2C Digital Goods: use tax tag `customer-rates`, that applies VAT according to customer location. In TicketBAI, these cases are "not subject" to tax, and thus should have the cause RL (por reglas de localización).

## Test Data

Some sample test data is available in the `./test` directory. To update the JSON documents and regenerate the XML files for testing, use the following command:

```bash
go test ./examples_test.go --update
```

All generate XML documents will be validated against the Veri*Factu XSD documents.
