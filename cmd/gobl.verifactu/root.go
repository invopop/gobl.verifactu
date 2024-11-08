package main

import (
	"io"
	"os"

	verifactu "github.com/invopop/gobl.verifactu"
	"github.com/spf13/cobra"
)

type rootOpts struct {
	swNIF         string
	swCompanyName string
	swVersion     string
	swLicense     string
	production    bool
}

func root() *rootOpts {
	return &rootOpts{}
}

func (o *rootOpts) cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           name,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(versionCmd())
	cmd.AddCommand(convert(o).cmd())

	return cmd
}

func (o *rootOpts) outputFilename(args []string) string {
	if len(args) >= 2 && args[1] != "-" {
		return args[1]
	}
	return ""
}

func openInput(cmd *cobra.Command, args []string) (io.ReadCloser, error) {
	if inFile := inputFilename(args); inFile != "" {
		return os.Open(inFile)
	}
	return io.NopCloser(cmd.InOrStdin()), nil
}

func (o *rootOpts) openOutput(cmd *cobra.Command, args []string) (io.WriteCloser, error) {
	if outFile := o.outputFilename(args); outFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		return os.OpenFile(outFile, flags, os.ModePerm)
	}
	return writeCloser{cmd.OutOrStdout()}, nil
}

func (o *rootOpts) software() *verifactu.Software {
	return &verifactu.Software{
		NombreRazon: o.swCompanyName,
		NIF:         o.swNIF,
		Version:     o.swVersion,
	}
}

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }
