// Package noverifactu provides GOBL types and validation for the NO VERI*FACTU modality.
package noverifactu

import "github.com/invopop/gobl/rules"

// Rules is the main rules set for validating NO VERI*FACTU documents
var Rules = rules.NewSet("NOVERIFACTU", append(statusRules, complementRules...)...)

// Validate checks the provided object against the NO VERI*FACTU rules and returns any
// faults found.
func Validate(obj any) rules.Faults {
	return Rules.Validate(obj)
}
