package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "goi",
	Short: "goi is a code injection tool for go language",
	Long: `	goi is a code injection tool for go language.
			Find more information at: 
			git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetReportCaller(true)
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				dirname, filename := filepath.Split(f.File)
				lastelem := filepath.Base(dirname)
				filename = filepath.Join(lastelem, filename)
				line := strconv.Itoa(f.Line)
				return "", "[" + filename + ":" + line + "]"
			},
		})
		if !debugGoi {
			// we only need log in debug mode
			log.SetLevel(log.FatalLevel)
			log.SetFormatter(&log.TextFormatter{
				DisableTimestamp: true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return "", ""
				},
			})
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if debugInCISyncFile != "" {
			f, err := os.Create(debugInCISyncFile)
			if err != nil {
				log.Fatalln(err)
			}
			defer func() {
				_ = f.Close()
			}()

			time.Sleep(5 * time.Second)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugGoi, "debug", false, "run goi in debug mode")
	rootCmd.PersistentFlags().BoolVar(&goc, "goc", false, "whether this project use goc to calculate code coverage")
	rootCmd.PersistentFlags().StringVar(&debugInCISyncFile, "debugcisyncfile", "", "internal use only, no explain")
	_ = rootCmd.PersistentFlags().MarkHidden("debugcisyncfile")
	_ = viper.BindPFlags(rootCmd.PersistentFlags())
}

// Execute the goi tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
