// Package main provides the command line interface to the VeriFactu package.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/invopop/gobl.verifactu/doc"
	"github.com/spf13/cobra"
)

type sendTestOpts struct {
	*rootOpts
	previous string
}

func sendTest(o *rootOpts) *sendTestOpts {
	return &sendTestOpts{rootOpts: o}
}

func (c *sendTestOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendTest [infile]",
		Short: "Sends the GOBL invoiceFactu service",
		RunE:  c.runE,
	}

	f := cmd.Flags()
	c.prepareFlags(f)

	f.StringVar(&c.previous, "prev", "", "Previous document fingerprint to chain with")

	return cmd
}

func (c *sendTestOpts) runE(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return fmt.Errorf("unmarshaling gobl envelope: %w", err)
	}

	opts := []verifactu.Option{
		verifactu.WithThirdPartyIssuer(),
		verifactu.InTesting(),
	}

	tc, err := verifactu.New(c.software(), opts...)
	if err != nil {
		return err
	}

	td, err := tc.Convert(env)
	if err != nil {
		return err
	}

	c.previous = `{
		"emisor": "B123456789",
		"serie": "FACT-001", 
		"fecha": "2024-11-11",
		"huella": "abc123def456"
	}`
	prev := new(doc.ChainData)
	if err := json.Unmarshal([]byte(c.previous), prev); err != nil {
		return err
	}

	err = tc.Fingerprint(td, prev)
	if err != nil {
		return err
	}

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

	if _, err = out.Write(append(convOut, '\n')); err != nil {
		return fmt.Errorf("writing verifactu xml: %w", err)
	}

	err = tc.Post(cmd.Context(), td)
	if err != nil {
		return err
	}
	fmt.Println("made it!")

	data, err := json.Marshal(td.ChainData())
	if err != nil {
		return err
	}
	fmt.Printf("Generated document with fingerprint: \n%s\n", string(data))

	// TEMP

	return nil
}
