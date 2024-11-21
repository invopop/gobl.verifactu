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

type convertOpts struct {
	*rootOpts
}

func convert(o *rootOpts) *convertOpts {
	return &convertOpts{rootOpts: o}
}

func (c *convertOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [infile] [outfile]",
		Short: "Convert a GOBL JSON into a VeriFactu XML",
		RunE:  c.runE,
	}

	return cmd
}

func (c *convertOpts) runE(cmd *cobra.Command, args []string) error {
	input, err := openInput(cmd, args)
	if err != nil {
		return err
	}
	defer input.Close() // nolint:errcheck

	out, err := c.openOutput(cmd, args)
	if err != nil {
		return err
	}
	defer out.Close() // nolint:errcheck

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	env := new(gobl.Envelope)
	if err := json.Unmarshal(buf.Bytes(), env); err != nil {
		return fmt.Errorf("unmarshaling gobl envelope: %w", err)
	}

	vf, err := verifactu.New(c.software())
	if err != nil {
		return fmt.Errorf("creating verifactu client: %w", err)
	}

	doc, err := vf.Convert(env)
	if err != nil {
		return fmt.Errorf("converting to verifactu xml: %w", err)
	}

	data, err := doc.Envelop()
	if err != nil {
		return fmt.Errorf("generating verifactu xml: %w", err)
	}

	if _, err = out.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("writing verifactu xml: %w", err)
	}

	return nil
}
