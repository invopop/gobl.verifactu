// Package main provides the command line interface to the VeriFactu package.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
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
		Short: "Convert a GOBL JSON into a TickeBAI XML",
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

	// tbai, err := ticketbai.New(&ticketbai.Software{})
	// if err != nil {
	// 	return fmt.Errorf("creating ticketbai client: %w", err)
	// }

	// doc, err := tbai.NewDocument(env)
	// if err != nil {
	// 	panic(err)
	// }

	// data, err := doc.BytesIndent()
	// if err != nil {
	// 	return fmt.Errorf("generating ticketbai xml: %w", err)
	// }

	// if _, err = out.Write(append(data, '\n')); err != nil {
	// 	return fmt.Errorf("writing ticketbai xml: %w", err)
	// }

	return nil
}
