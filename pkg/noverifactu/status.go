package noverifactu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
)

// Status line keys for each NO VERI*FACTU event type. These map to the TipoEvento codes
// defined in the L2E list of the VeriFactu specification.
const (
	KeySystemStartup        cbc.Key = "system-startup"         // 01
	KeySystemShutdown       cbc.Key = "system-shutdown"        // 02
	KeyInvoiceAnomalyLaunch cbc.Key = "invoice-anomaly-launch" // 03
	KeyInvoiceAnomaly       cbc.Key = "invoice-anomaly"        // 04
	KeyEventAnomalyLaunch   cbc.Key = "event-anomaly-launch"   // 05
	KeyEventAnomaly         cbc.Key = "event-anomaly"          // 06
	KeyBackupRestoration    cbc.Key = "backup-restoration"     // 07
	KeyInvoiceExport        cbc.Key = "invoice-export"         // 08
	KeyEventExport          cbc.Key = "event-export"           // 09
	KeyEventSummary         cbc.Key = "event-summary"          // 10
	KeyOther                cbc.Key = "other"                  // 90
)

var statusRules = []*rules.Set{
	statusRuleSet(),
}

var validEventKeys = func() []any {
	keys := make([]any, 0, len(keyToComplement))
	for k := range keyToComplement {
		keys = append(keys, k)
	}
	return keys
}()

var keyToComplement = map[cbc.Key]any{
	KeySystemStartup:        nil,
	KeySystemShutdown:       nil,
	KeyInvoiceAnomalyLaunch: InvoiceAnomalyLaunch{},
	KeyInvoiceAnomaly:       InvoiceAnomaly{},
	KeyEventAnomalyLaunch:   EventAnomalyLaunch{},
	KeyEventAnomaly:         EventAnomaly{},
	KeyBackupRestoration:    nil,
	KeyInvoiceExport:        InvoiceExport{},
	KeyEventExport:          EventExport{},
	KeyEventSummary:         EventSummary{},
	KeyOther:                nil,
}

var isAnomalyDetection = is.Func("anomaly detection event", func(v any) bool {
	line, _ := v.(*bill.StatusLine)
	if line == nil {
		return false
	}
	return line.Key == KeyInvoiceAnomaly || line.Key == KeyEventAnomaly
})

var isOtherNote = is.Func("other note", func(v any) bool {
	note, _ := v.(*org.Note)
	return note != nil && note.Key == org.NoteKeyOther
})

var requiresComplement = is.Func("requires complement", func(v any) bool {
	line, _ := v.(*bill.StatusLine)
	if line == nil {
		return false
	}
	return keyToComplement[line.Key] != nil
})

var hasCorrectComplement = is.Func("correct complement schema", func(v any) bool {
	line, _ := v.(*bill.StatusLine)
	if line == nil || len(line.Complements) == 0 {
		return true // other rules handle missing info
	}

	s := schema.Lookup(keyToComplement[line.Key])
	if s == schema.UnknownID {
		return true // other rules handle unsupported keys
	}

	return line.Complements[0].Schema == s
})

func statusRuleSet() *rules.Set {
	return rules.For(new(bill.Status),
		rules.Field("lines",
			rules.Assert("01", "status must have exactly one line", is.Present, is.Length(1, 1)),
			rules.Each(
				rules.Field("key",
					rules.Assert("02", "line key is required", is.Present),
					rules.Assert("03", "line key must be a supported event type", is.In(validEventKeys...)),
				),
				rules.When(requiresComplement,
					rules.Field("complements",
						rules.Assert("04", "complement is required for this event type", is.Present, is.Length(1, 1)),
					),
					rules.Assert("05", "complement must correspond to the event type", hasCorrectComplement),
				),
				rules.When(isAnomalyDetection,
					rules.Field("description",
						rules.AssertIfPresent("06", "description must be 100 characters or less", is.Length(0, 100)),
					),
				),
			),
		),
		rules.Field("notes",
			rules.Each(
				rules.When(isOtherNote,
					rules.Field("text",
						rules.AssertIfPresent("07", "other note text must be 100 characters or less", is.Length(0, 100)),
					),
				),
			),
		),
	)
}
