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
var keyWords []string

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
	if strings.Contains(RootClis.Keywords, " ,") {
		keyWords = strings.Split(RootClis.Keywords, " ,")
	} else {
		keyWords = strings.Split(RootClis.Keywords, ",")
	}
	cleanArgs := parseStdin(cmd)
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
	var sraFields []*extract.SraFields
	var pubMedFields []*extract.PubmedFields
	var lock sync.Mutex
	var pubmedJSON []parse.PubmedArticleJSON
	var sraJSON []parse.ExperimentPkgJSON
	if RootClis.Mode == "pubmed" && len(dat) > 0 {
		json.Unmarshal(dat, &pubmedJSON)
		for _, v := range pubmedJSON {
			lock.Lock()
			pubMedFields = append(pubMedFields, extract.GetSimplePubmedFields(&keyWords, &v, RootClis.CallCor))
			lock.Unlock()
		}
		dat, _ := json.MarshalIndent(pubMedFields, "", "    ")
		return &dat
	} else if RootClis.Mode == "sra" && len(dat) > 0 {
		json.Unmarshal(dat, &sraJSON)
		done := make(map[string]int)
		for _, v := range sraJSON {
			lock.Lock()
			sraFields = append(sraFields, extract.GetSimpleSraFields(&keyWords, &v, RootClis.CallCor, done))
			done[v.EXPERIMENT.TITLE+v.STUDY.DESCRIPTOR.STUDYTITLE] = 1
			lock.Unlock()
		}
		dat, _ := json.MarshalIndent(sraFields, "", "    ")
		return &dat
	} else if len(dat) > 0 {
		obj, _ := extract.GetPlainFields("", &dat, &keyWords, RootClis.CallCor)
		dat, _ := json.MarshalIndent(obj, "", "    ")
		return &dat
	}
	return nil
}

func init() {
	RootCmd.Flags().StringVarP(&RootClis.Keywords, "keywords", "w", "algorithm, tool, model, pipleline, method, database, workflow, dataset, bioinformatics, sequencing, http, github.com, gitlab.com, bitbucket.org", "Keywords to extracted from abstract.")
	RootCmd.Flags().BoolVarP(&RootClis.CallCor, "call-cor", "", false, "Wheather to calculate the corelated keywords, and return the sentence contains >=2 keywords.")
	RootCmd.Flags().StringVarP(&RootClis.Mode, "mode", "", "", "mode to extract information: plain,pubmed, or sra.")
}
