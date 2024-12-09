# GOBL to VeriFactu

Go library to convert [GOBL](https://github.com/invopop/gobl) invoices into VeriFactu declarations and send them to the AEAT (Agencia Estatal de Administración Tributaria) web services. This library assumes that clients will handle a local database of previous invoices in order to comply with the local requirements of chaining all invoices together.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [GNU Affero General Public License v3.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). For contributions to this library to be accepted, we will require you to accept transferring your copyright to Invopop Ltd.

## Source

The main resources used in this module include:

- [VeriFactu documentation](https://www.agenciatributaria.es/AEAT.desarrolladores/Desarrolladores/_menu_/Documentacion/Sistemas_Informaticos_de_Facturacion_y_Sistemas_VERI_FACTU/Sistemas_Informaticos_de_Facturacion_y_Sistemas_VERI_FACTU.html)
- [VeriFactu Ministerial Order](https://www.boe.es/diario_boe/txt.php?id=BOE-A-2024-22138)
- [Spanish VAT Law](https://www.boe.es/buscar/act.php?id=BOE-A-1992-28740).

## Usage

### Go Package

You must have first created a GOBL Envelope containing an Invoice that you'd like to send to the AEAT. For the document to be converted, the supplier contained in the invoice should have a "Tax ID" with the country set to `ES`.

The following is an example of how the GOBL VeriFactu package could be used:

```go
package main

import (
	"context"
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

	// VeriFactu requires a software definition to be provided. This is an example
	software := &verifactu.Software{
		NombreRazon:              "Company LTD",    // Company Name
		NIF:                      "B123456789",     // Software company's tax code
		NombreSistemaInformatico: "Software Name",  // Name of application
		IdSistemaInformatico:     "A1",             // Software ID
		Version:                  "1.0",            // Software version
		NumeroInstalacion:        "00001",          // Software installation number
	}

	// Load the certificate
	cert, err := xmldsig.LoadCertificate(
		"./path/to/certificate.p12",
		"password",
	)
	if err != nil {
		panic(err)
	}

	// Create the client with the software and certificate
	opts := []verifactu.Option{
		verifactu.WithCertificate(cert),
		verifactu.WithSupplierIssuer(), // The issuer can be either the supplier, the customer or a third party
		verifactu.InTesting(),          // Use the testing environment, as the production endpoint is not yet published 
	}

	tc, err := verifactu.New(software, opts...)
	if err != nil {
		panic(err)
	}

	// Convert the GOBL envelope to a VeriFactu document
	td, err := tc.Convert(env)
	if err != nil {
		panic(err)
	}

	// Prepare the previous document chain data
	previous, err := os.ReadFile("./path/to/previous_invoice.json")
	if err != nil {
		panic(err)
	}

	prev := new(doc.ChainData)
	if err := json.Unmarshal([]byte(previous), prev); err != nil {
		panic(err)
	}

	// Create the document fingerprint based on the previous document chain
	err = tc.Fingerprint(td, prev)
	if err != nil {
		panic(err)
	}

	// Add the QR code to the document
	if err := tc.AddQR(td, env); err != nil {
		panic(err)
	}

	// Send the document to the tax agency
	err = tc.Post(ctx, td)
	if err != nil {
		panic(err)
	}

	// Print the data to be used as previous document chain for the next invoice
	// Persist the data somewhere to be used by the next invoice
	cd, err := json.Marshal(td.ChainData())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated document with fingerprint: \n%s\n", string(cd))

}
```

### Command Line

The GOBL VeriFactu package tool also includes a command line helper. You can install manually in your Go environment with:

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

CERTIFICATE_PATH=xxxxxxxxx
CERTIFICATE_PASSWORD=xxxxxxxxx
```
To convert a document to XML, run:

```bash
gobl.verifactu convert ./test/data/sample-invoice.json
```
This function will output the XML to the terminal, or to a file if a second argument is provided. The output file will not include a fingerprint, and therefore will not be able to be submitted to the tax agency.

To submit to the tax agency testing environment:

```bash
gobl.verifactu send ./test/data/sample-invoice.json ./test/data/previous-invoice-info.json
```
Now, the output file will include a fingerprint, linked to the previous document, and will be submitted to the tax agency. An example for a previous file would look like this:

```json
{
	"emisor": "B12345678",
	"serie": "SAMPLE-001",
	"fecha": "2024-11-28",
	"huella": "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
}
```

## Tags and Extensions

In order to provide the supplier specific data required by VeriFactu, invoices need to include a bit of extra data. We've managed to simplify these into specific cases.

### Invoice Tags

Invoice tax tags can be added to invoice documents in order to reflect a special situation. The following schemes are supported:

- `simplified` - a retailer operating under a simplified tax regime (regimen simplificado) that must indicate that all of their sales are under this scheme. This implies that all operations in the invoice will have the `FacturaSinIdentifDestinatarioArt61d` tag set to `S` and the `TipoFactura` field set to `F2` in case of a regular invoice and `R5` in case of a corrective invoice.
- `substitution` - A simplified invoice that is being replaced by a standard invoice. Called a `Factura en Sustitución de Facturas Simplificadas` in VeriFactu. The `TipoFactura` field will be set to `F3`.

### Tax Extensions

The following extensions must be added to the document:

- `es-verifactu-doc-type` – defines the type of invoice being sent. In most cases this will be set automatically by GOBL, but it must be present. These are the valid values:
  - `F1` - Standard invoice.
  - `F2` - Simplified invoice.
  - `F3` - Invoice in substitution of simplified invoices.
  - `R1` - Rectified invoice based on law and Article 80.1, 80.2 and 80.6 in the Spanish VAT Law.
  - `R2` - Rectified invoice based on law and Article 80.3.
  - `R3` - Rectified invoice based on law and Article 80.4.
  - `R4` - Rectified invoice based on law and other reasons.
  - `R5` - Rectified invoice based on simplified invoices.


- `es-verifactu-tax-classification` - combines the tax classification and exemption codes used in VeriFactu. Must be included in each line item, or an error will be raised. These are the valid values:
  - `S1` - Subject and not exempt - Without reverse charge
  - `S2` - Subject and not exempt - With reverse charge. Known as `Inversión del Sujeto Pasivo` in Spanish VAT Law
  - `N1` - Not subject - Articles 7, 14, others
  - `N2` - Not subject - Due to location rules
  - `E1` - Exempt pursuant to Article 20 of the VAT Law
  - `E2` - Exempt pursuant to Article 21 of the VAT Law
  - `E3` - Exempt pursuant to Article 22 of the VAT Law
  - `E4` - Exempt pursuant to Articles 23 and 24 of the VAT Law
  - `E5` - Exempt pursuant to Article 25 of the VAT Law
  - `E6` - Exempt for other reasons

As a small consideration GOBL's tax internal tax framework differentiates between `exempt` and `zero-rated` taxes. In VeriFactu, GOBL `zero-rated` taxes refer to `Exenciones` (values `E1` to `E6` in the list above) and `exempt` taxes refer to `No Sujeto` (values `N1` and `N2` in the list above).

## Limitations

- VeriFactu allows more than one customer per invoice, but GOBL only has one possible customer.
- Invoices must have a note of type general that will be used as a general description of the invoice. If an invoice is missing this info, it will be rejected with an error.
- VeriFactu supports sending more than one invoice at a time (up to 1000). However, this module only currently supports 1 invoice at a time.
- VeriFactu requires a valid certificate to be provided, even when using the testing environment. It is the same certificate needed to access the AEAT's portal.
- When cancelling invoices, this module assumes the party issuing the cancellation is the same as the party that issued the original invoice. In the context of the app this would always be true, but VeriFactu does allow for a different issuer.

## Testing

This library includes a set of tests that can be used to validate the conversion and submission process. To run the tests, use the following command:

```bash
go test
```

Some sample test data is available in the `./test` directory. To update the JSON documents and regenerate the XML files for testing, use the following command:

```bash
go test --update
```