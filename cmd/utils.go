package cmd

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	cvrt "github.com/openbiox/ligo/convert"
	cio "github.com/openbiox/ligo/io"
	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/stringo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = clog.Logger
var logBash = clog.LoggerBash
var logEnv = log.WithFields(logrus.Fields{
	"prefix": "Env"})
var logPrefix string
var wd string

func setGlobalFlag(cmd *cobra.Command) {
	wd, _ = os.Getwd()
	cmd.PersistentFlags().IntVarP(&(RootClis.Verbose), "verbose", "", 1, "verbose level(0:no output, 1: basic level, 2: with env info")
	cmd.PersistentFlags().StringVarP(&(RootClis.TaskID), "task-id", "k", stringo.RandString(15), "task ID (default is random).")
	cmd.PersistentFlags().StringVarP(&(RootClis.LogDir), "log-dir", "", path.Join(wd, "_log"), "log dir.")
	cmd.PersistentFlags().BoolVarP(&(RootClis.SaveLog), "save-log", "s", false, "save log to file.")
	cmd.PersistentFlags().BoolVarP(&(RootClis.Clean), "clean", "", false, "remove log dir.")
	cmd.PersistentFlags().StringVarP(&RootClis.Out, "out", "o", "", "out specifies destination of the returned data (default to stdout or current woring directory).")
	cmd.PersistentFlags().IntVarP(&(RootClis.Thread), "thread", "t", 1, "thread to process.")

}
func initCmd(cmd *cobra.Command, args []string) {
	setLog()
	if RootClis.Verbose == 2 {
		logEnv.Infof("Prog: %s", cmd.CommandPath())
		logEnv.Infof("TaskID: %s", RootClis.TaskID)
		if RootClis.SaveLog && logPrefix != "" {
			logEnv.Infof("Log: %s.log", logPrefix)
		}
		if len(args) > 0 {
			logEnv.Infof("Args: %s", strings.Join(args, " "))
		}
		logEnv.Infof("Global: %v", cvrt.Struct2Map(RootClis))
	}
	if RootClis.Clean {
		cleanLog()
	}
}

func setLog() {
	var logCon io.Writer
	var logDir = RootClis.LogDir

	if RootClis.SaveLog {
		if logDir == "" {
			logDir = filepath.Join(os.TempDir(), "_log")
		}
		logPrefix = fmt.Sprintf("%s/%s", logDir, RootClis.TaskID)
		cio.CreateDir(logDir)
		logCon, _ = cio.Open(logPrefix + ".log")
	}
	clog.SetLogStream(log, RootClis.Verbose == 0, RootClis.SaveLog, &logCon)
}

func cleanLog() {
	RootClis.HelpFlags = false
	if err := os.RemoveAll(RootClis.LogDir); err != nil {
		log.Warn(err)
	}
}

func removeDuplicatesAndEmpty(a []string) (ret []string) {
	sort.Sort(sort.StringSlice(a))
	alen := len(a)
	for i := 0; i < alen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}
