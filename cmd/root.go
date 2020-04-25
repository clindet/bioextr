package cmd

import (
	"fmt"
	"os"

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
	Thread  int

	// call corelation
	Keywords     string
	KeywordsFile string
	CallCor      bool

	// type
	Mode      string
	ListFile  string
	HelpFlags bool
}

var RootClis = RootClisT{
	Version:   version,
	Verbose:   1,
	HelpFlags: true,
}

var RootCmd = &cobra.Command{
	Use:   "bioextr [filename]",
	Short: "A simple command line tool to extract information from text and json files.",
	Long:  `A simple command line tool to extract information from text and json files. More see here https://github.com/openanno/bioextr.`,
	Run: func(cmd *cobra.Command, args []string) {
		if RootClis.Clean {
			initCmd(cmd, args)
			RootClis.HelpFlags = false
		}
		simpleExtr(cmd, args)
		if RootClis.HelpFlags {
			cmd.Help()
		}
	},
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
	RootCmd.Example = `  # extract from pubmed abstract
  bget api ncbi -q "Galectins control MTOR and AMPK in response to lysosomal damage to induce autophagy OR MTOR-independent autophagy induced by interrupted endoplasmic reticulum-mitochondrial Ca2+ communication: a dead end in cancer cells. OR The PARK10 gene USP24 is a negative regulator of autophagy and ULK1 protein stability OR Coordinate regulation of autophagy and the ubiquitin proteasome system by MTOR." | bioctl cvrt --xml2json pubmed - | bioextr --mode pubmed -w 'MTOR,AMPK,autophagy' --call-cor -
	
  # extract from sra json
  bget api ncbi -d 'sra' -q PRJNA527714 | bioctl cvrt --xml2json sra - | bioextr --mode sra --call-cor -w "Chromatin,mouse" -`
}
