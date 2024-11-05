package doc

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
)

// ValidationError is a simple wrapper around validation errors
type ValidationError struct {
	text string
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	return e.text
}

func validationErr(text string, args ...any) error {
	return &ValidationError{
		text: fmt.Sprintf(text, args...),
	}
}

func validate(inv *bill.Invoice, zone l10n.Code) error {
	if inv.Type == bill.InvoiceTypeCorrective {
		return validationErr("corrective invoices not supported, use credit or debit notes")
	}

	if inv.Supplier == nil || inv.Supplier.TaxID == nil {
		return nil // ignore
	}

	for _, l := range inv.Lines {
		if len(l.Charges) > 0 {
			return validationErr("charges are not supported")
		}
	}

	if len(inv.Charges) > 0 {
		return validationErr("charges are not supported")
	}

	return nil
}
