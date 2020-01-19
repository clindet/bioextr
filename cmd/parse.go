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
	"github.com/openbiox/ligo/parse"
	"github.com/spf13/cobra"
)

var stdin []byte
var sraFields []*extract.SraFields
var pubMedFields []*extract.PubmedFields

func parseStdin(cmd *cobra.Command) []string {
	cleanArgs := []string{}
	var err error
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		stdin, err = ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
	}
	return cleanArgs
}

func simpleExtr(cmd *cobra.Command, args []string) {
	cleanArgs := parseStdin(cmd)
	var wg sync.WaitGroup
	sem := make(chan struct{}, RootClis.Thread)

	if len(stdin) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号
			defer func() { <-sem }() // 释放信号
			parseJSON(stdin)
		}()
		RootClis.HelpFlags = false
	}
	if len(cleanArgs) > 0 {
		for _, v := range cleanArgs {
			wg.Add(1)
			go func(v string) {
				defer wg.Done()
				sem <- struct{}{} // 获取信号
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
				parseJSON(input)
				defer func() { <-sem }() // 释放信号
			}(v)
		}
		RootClis.HelpFlags = false
	}
	wg.Wait()
}

func parseJSON(dat []byte) {
	var pubmedJson []parse.PubmedArticleJSON
	var sraJson []parse.ExperimentPkgJSON
	keyWords := []string{}
	if strings.Contains(RootClis.Keywords, " ,") {
		keyWords = strings.Split(RootClis.Keywords, " ,")
	} else {
		keyWords = strings.Split(RootClis.Keywords, ",")
	}
	if RootClis.Mode == "pubmed" && len(dat) > 0 {
		json.Unmarshal(stdin, &pubmedJson)
		for _, v := range pubmedJson {
			pubMedFields = append(pubMedFields, extract.GetSimplePubmedFields(&keyWords, &v, RootClis.CallCor))
		}
		dat, _ := json.MarshalIndent(pubMedFields, "", "    ")
		fmt.Println(string(dat))
	} else if RootClis.Mode == "sra" && len(dat) > 0 {
		json.Unmarshal(stdin, &sraJson)
		done := make(map[string]int)
		for _, v := range sraJson {
			sraFields = append(sraFields, extract.GetSimpleSraFields(&keyWords, &v, RootClis.CallCor, done))
			done[v.EXPERIMENT.TITLE+v.STUDY.DESCRIPTOR.STUDYTITLE] = 1
		}
		dat, _ := json.MarshalIndent(sraFields, "", "    ")
		fmt.Println(string(dat))
	}
}

func init() {
	RootCmd.Flags().StringVarP(&RootClis.Keywords, "keywords", "w", "algorithm, tool, model, pipleline, method, database, workflow, dataset, bioinformatics, sequencing, http, github.com, gitlab.com, bitbucket.org", "Keywords to extracted from abstract.")
	RootCmd.Flags().BoolVarP(&RootClis.CallCor, "call-cor", "", false, "Wheather to calculate the corelated keywords, and return the sentence contains >=2 keywords.")
	RootCmd.Flags().StringVarP(&RootClis.Mode, "mode", "", "", "mode to extract information (pubmed, sra).")
}
