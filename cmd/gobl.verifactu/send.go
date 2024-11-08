// Package main provides the command line interface to the VeriFactu package.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/spf13/cobra"
)

type sendOpts struct {
	*rootOpts
	previous string
}

func send(o *rootOpts) *sendOpts {
	return &sendOpts{rootOpts: o}
}

func (c *sendOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [infile]",
		Short: "Sends the GOBL invoice to the VeriFactu service",
		RunE:  c.runE,
	}

	f := cmd.Flags()

	f.StringVar(&c.previous, "prev", "", "Previous document fingerprint to chain with")
	f.BoolVarP(&c.production, "production", "p", false, "Production environment")

	return cmd
}

func (c *sendOpts) runE(cmd *cobra.Command, args []string) error {
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
	}

	if c.production {
		opts = append(opts, verifactu.InProduction())
	} else {
		opts = append(opts, verifactu.InTesting())
	}

	tc, err := verifactu.New(c.software())
	if err != nil {
		return err
	}

	td, err := tc.Convert(env)
	if err != nil {
		return err
	}

	err = tc.Fingerprint(td, c.previous)
	if err != nil {
		return err
	}

	// if err := tc.Sign(td, env); err != nil {
	// 	return err
	// }

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
