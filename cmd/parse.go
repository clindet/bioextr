package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/openbiox/ligo/extract"
	"github.com/openbiox/ligo/flag"
	cio "github.com/openbiox/ligo/io"
	"github.com/openbiox/ligo/stringo"
	"github.com/spf13/cobra"
)

var stdin []byte
var keyWords []string
var cleanArgs []string
var keyWordsPat string

func parseStdin(cmd *cobra.Command) {
	var err error
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		stdin, err = ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func simpleExtr(cmd *cobra.Command, args []string) {
	if strings.Contains(RootClis.Keywords, " ,") {
		keyWords = strings.Split(RootClis.Keywords, " ,")
	} else {
		keyWords = strings.Split(RootClis.Keywords, ",")
	}
	if RootClis.KeywordsFile != "" {
		of, _ := os.Open(RootClis.KeywordsFile)
		keyWordsArr, _ := ioutil.ReadAll(of)
		keyWords = stringo.StrSplit(string(keyWordsArr), "\r\n|\n|\r|\t", 10000000)
	}
	keyWords = removeDuplicatesAndEmpty(keyWords)
	keyWordsPat = strings.Join(keyWords, "|")
	parseStdin(cmd)
	var wg sync.WaitGroup
	sem := make(chan struct{}, RootClis.Thread)
	if len(stdin) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			defer fmt.Println(string(*parseJSON(stdin, "")))
		}()
		RootClis.HelpFlags = false
	}
	if RootClis.ListFile != "" {
		cleanArgs = append(cleanArgs, cio.ReadLines(RootClis.ListFile)...)
	}
	if len(cleanArgs) > 0 {
		for _, v := range cleanArgs {
			wg.Add(1)
			go func(v string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				defer fmt.Println(string(*parseJSON(nil, v)))
			}(v)
		}
		RootClis.HelpFlags = false
	}
	wg.Wait()
}

func parseJSON(dat []byte, infile string) *[]byte {
	var sraFields []extract.SraFields
	var pubMedFields []extract.PubmedFields
	if len(dat) == 0 && infile == "" {
		return nil
	}
	if RootClis.Mode == "pubmed" {
		pubMedFields, _ = extract.GetSimplePubmedFields(infile, &dat, &keyWordsPat, RootClis.CallCor, RootClis.CallURLs, RootClis.KeepAbs, RootClis.Thread)
		dat2, _ := json.MarshalIndent(pubMedFields, "", "    ")
		return &dat2
	} else if RootClis.Mode == "sra" {
		sraFields, _ = extract.GetSimpleSraFields(infile, &dat, &keyWordsPat, RootClis.CallCor, RootClis.CallURLs, RootClis.KeepAbs, RootClis.Thread)
		dat2, _ := json.MarshalIndent(sraFields, "", "    ")
		return &dat2
	} else if RootClis.Mode == "bigd" {
		articleFields, _ := extract.GetBigdFields(infile, &dat, &keyWordsPat, RootClis.CallCor, RootClis.CallURLs, RootClis.KeepAbs, RootClis.Thread)
		dat2, _ := json.MarshalIndent(articleFields, "", "    ")
		return &dat2
	} else {
		obj, _ := extract.GetPlainFields(infile, &dat, &keyWordsPat, RootClis.CallCor, RootClis.CallURLs, RootClis.Thread)
		dat2, _ := json.MarshalIndent(obj, "", "    ")
		return &dat2
	}
}

func init() {
	RootCmd.Flags().StringVarP(&RootClis.Keywords, "keywords", "w", "algorithm, tool, model, pipleline, method, database, workflow, dataset, bioinformatics, sequencing, http, github.com, gitlab.com, bitbucket.org", "Keywords to extracted from abstract.")
	RootCmd.Flags().StringVarP(&RootClis.KeywordsFile, "keywords-file", "", "", "Keywors in file, one colum in a file")

	RootCmd.Flags().BoolVarP(&RootClis.CallURLs, "call-urls", "", false, "Wheather to extract all URLs")
	RootCmd.Flags().BoolVarP(&RootClis.KeepAbs, "keep-abs", "", false, "Wheather to keep abstract field")
	RootCmd.Flags().BoolVarP(&RootClis.CallCor, "call-cor", "", false, "Wheather to calculate the corelated keywords, and return the sentence contains >=2 keywords.")
	RootCmd.Flags().StringVarP(&RootClis.Mode, "mode", "", "", "mode to extract information: plain,pubmed, or sra.")
}
