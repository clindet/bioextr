package cmd

import (
	"fmt"
	"os"

	"code.sajari.com/docconv"
	"github.com/openbiox/ligo/io"
	"github.com/spf13/cobra"
)

// RootClisT is the bioctl global flags
type RootClisT struct {
	// version of bioctl
	Version string
	Verbose int
	SaveLog bool
	TaskID  string
	LogDir  string
	Clean   bool
	Out     string

	HelpFlags bool
}

var RootClis = RootClisT{
	Version:   version,
	Verbose:   1,
	HelpFlags: true,
}

var RootCmd = &cobra.Command{
	Use:   "pdf2plain [input.pdf]",
	Short: "A wrapper command line tool to convert pdf files to plain text.",
	Long:  `A wrapper command line tool to convert pdf files to plain text. More see here https://github.com/openbiox/bioextr.`,
	Run: func(cmd *cobra.Command, args []string) {
		if RootClis.Clean {
			RootClis.HelpFlags = false
		}
		if len(args) > 0 {
			initCmd(cmd, args)
			convertor(cmd, args)
			RootClis.HelpFlags = false
		}
		if RootClis.HelpFlags {
			cmd.Help()
		}
	},
}

func convertor(cmd *cobra.Command, args []string) {
	res, err := docconv.ConvertPath(args[0])
	if err != nil {
		log.Fatal(err)
	}
	if RootClis.Out != "" {
		if err := io.CreateFileParDir(RootClis.Out); err != nil {
			log.Warnln(err)
			return
		}
		con, _ := io.Open(RootClis.Out)
		fmt.Fprintf(con, res.Body)
		return
	}
	fmt.Println(res.Body)
}

// Execute main interface of bget
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if !RootCmd.HasFlags() && !RootCmd.HasSubCommands() {
			RootCmd.Help()
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	wd, _ = os.Getwd()
	RootCmd.Version = version
	setGlobalFlag(RootCmd)
	RootCmd.Example = `  pdf2plain _examples/Multi-omic_approaches_to_improve_outcome_for_T-cel.pdf`
}
