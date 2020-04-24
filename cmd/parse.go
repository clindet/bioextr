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
	"github.com/openbiox/ligo/stringo"
	"github.com/spf13/cobra"
)

var stdin []byte
var keyWords []string
var cleanArgs []string

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
	parseStdin(cmd)
	var wg sync.WaitGroup
	sem := make(chan struct{}, RootClis.Thread)

	if len(stdin) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			defer fmt.Println(string(*parseJSON(stdin)))
		}()
		RootClis.HelpFlags = false
	}
	if len(cleanArgs) > 0 {
		for _, v := range cleanArgs {
			wg.Add(1)
			go func(v string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				var input []byte
				var con *os.File
				var err error
				if con, err = os.Open(v); err != nil {
					log.Warnln(err)
					return
				}
				if input, err = ioutil.ReadAll(con); err != nil {
					log.Warnln(err)
					return
				}
				defer fmt.Println(string(*parseJSON(input)))
			}(v)
		}
		RootClis.HelpFlags = false
	}
	wg.Wait()
}

func parseJSON(dat []byte) *[]byte {
	var sraFields []extract.SraFields
	var pubMedFields []extract.PubmedFields
	var keyWordsPat string
	if RootClis.Mode == "pubmed" && len(dat) > 0 {
		keyWordsPat = strings.Join(keyWords, "|")
		pubMedFields, _ = extract.GetSimplePubmedFields("", &dat, &keyWordsPat, RootClis.CallCor, RootClis.Thread)
		dat2, _ := json.MarshalIndent(pubMedFields, "", "    ")
		return &dat2
	} else if RootClis.Mode == "sra" && len(dat) > 0 {
		keyWordsPat = strings.Join(keyWords, "|")
		sraFields, _ = extract.GetSimpleSraFields("", &dat, &keyWordsPat, RootClis.CallCor, RootClis.Thread)
		dat2, _ := json.MarshalIndent(sraFields, "", "    ")
		return &dat2
	} else if len(dat) > 0 {
		keyWordsPat = strings.Join(keyWords, "|")
		obj, _ := extract.GetPlainFields("", &dat, &keyWordsPat, RootClis.CallCor, RootClis.Thread)
		dat2, _ := json.MarshalIndent(obj, "", "    ")
		return &dat2
	}
	return nil
}

func init() {
	RootCmd.Flags().StringVarP(&RootClis.Keywords, "keywords", "w", "algorithm, tool, model, pipleline, method, database, workflow, dataset, bioinformatics, sequencing, http, github.com, gitlab.com, bitbucket.org", "Keywords to extracted from abstract.")
	RootCmd.Flags().StringVarP(&RootClis.KeywordsFile, "keywords-file", "", "", "Keywors in file, one colum in a file")

	RootCmd.Flags().BoolVarP(&RootClis.CallCor, "call-cor", "", false, "Wheather to calculate the corelated keywords, and return the sentence contains >=2 keywords.")
	RootCmd.Flags().StringVarP(&RootClis.Mode, "mode", "", "", "mode to extract information: plain,pubmed, or sra.")
}
