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
		License: "XYZ",        // provided by tax agency
		NIF:     "B123456789", // Software company's tax code
		Name:    "Invopop",    // Name of application
		Version: "v0.1.0",     // Software version
	}


	// Instantiate the TicketBAI client with sofrward config
	// and specific zone.
	c, err := verifactu.New(soft,
		verifactu.WithSupplierIssuer(),  // The issuer is the invoice's supplier
		verifactu.InTesting(),           // Use the tax agency testing environment
	)
	if err != nil {
		panic(err)
	}

	// Create a new Veri*Factu document:
	doc, err := c.Convert(env)
	if err != nil {
		panic(err)
	}

	// Create the document fingerprint
	// Assume here that we don't have a previous chain data object.
	if err = c.Fingerprint(doc, nil); err != nil {
		panic(err)
	}

	// Sign the document:
	if err := c.AddQR(doc, env); err != nil {
		panic(err)
	}

	// Create the XML output
	bytes, err := doc.BytesIndent()
	if err != nil {
		panic(err)
	}

	// Do something with the output, you probably want to store
	// it somewhere.
	fmt.Println("Document created:\n", string(bytes))

	// Grab and persist the Chain Data somewhere so you can use this
	// for the next call to the Fingerprint method.
	cd := doc.ChainData()

	// Send to Veri*Factu, if rejected, you'll want to fix any
	// issues and send in a new XML document. The original
	// version should not be modified.
	if err := c.Post(ctx, doc); err != nil {
		panic(err)
	}

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
SOFTWARE_COMPANY_NAME="Invopop S.L."
SOFTWARE_NAME="Invopop"
SOFTWARE_ID_SISTEMA_INFORMATICO="IP"
SOFTWARE_NUMERO_INSTALACION="12345678"
SOFTWARE_VERSION="1.0"
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

<!-- ### Tax Tags

Invoice tax tags can be added to invoice documents in order to reflect a special situation. The following schemes are supported:

- `simplified-scheme` - a retailer operating under a simplified tax regime (regimen simplificado) that must indicate that all of their sales are under this scheme. This implies that all operations in the invoice will have the `OperacionEnRecargoDeEquivalenciaORegimenSimplificado` tag set to `S`.
- `reverse-charge` - B2B services or goods sold to a tax registered EU member who will pay VAT on the suppliers behalf. Implies that all items will be classified under the `TipoNoExenta` value of `S2`.
- `customer-rates` - B2C services, specifically for the EU digital goods act (2015) which imply local taxes will be applied. All items will specify the `DetalleNoSujeta` cause of `RL`.

## Tax Extensions

The following extension can be applied to each line tax:

- `es-tbai-product` – allows to correctly group the invoice's lines taxes in the TicketBAI breakdowns (a.k.a. desgloses). These are the valid values:

  - `services` - indicates that the product being sold is a service (as opposed to a physical good). Services are accounted in the `DesgloseTipoOperacion > PrestacionServicios` breakdown of invoices to foreign customers. By default, all items are considered services.
  - `goods` - indicates that the product being sold is a physical good. Products are accounted in the `DesgloseTipoOperacion > Entrega` breakdown of invoices to foreign customers.
  - `resale` - indicates that a line item is sold without modification from a provider under the Equalisation Charge scheme. (This implies that the `OperacionEnRecargoDeEquivalenciaORegimenSimplificado` tag will be set to `S`).

- `es-tbai-exemption` - identifies the specific TicketBAI reason code as to why taxes should not be applied to the line according to the whole set of exemptions or not-subject scenarios defined in the law. It has to be set along with the tax rate value of `exempt`. These are the valid values:
  - `E1` – Exenta por el artículo 20 de la Norma Foral del IVA
  - `E2` – Exenta por el artículo 21 de la Norma Foral del IVA
  - `E3` – Exenta por el artículo 22 de la Norma Foral del IVA
  - `E4` – Exenta por el artículo 23 y 24 de la Norma Foral del IVA
  - `E5` – Exenta por el artículo 25 de la Norma Foral del IVA
  - `E6` – Exenta por otra causa
  - `OT` – No sujeto por el artículo 7 de la Norma Foral de IVA / Otros supuestos
  - `RL` – No sujeto por reglas de localización (\*)

_(\*) As noted elsewhere, `RL` will be set automatically set in invoices using the `customer-rates` tax tag. It can also be set explicitly using the `es-tbai-exemption` extension in invoices not using that tag._

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

All generate XML documents will be validated against the TicketBAI XSD documents. -->
