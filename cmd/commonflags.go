package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	debugGoi          bool
	goc               bool
	spexMocker        bool
	httpMocker        bool
	cacheMocker       bool
	customInject      string
	debugInCISyncFile string
	buildFlags        string
	buildOutput       string
)

func addCommonFlags(cmdset *pflag.FlagSet) {
	cmdset.StringVar(&buildFlags, "buildflags", "", "specify the build flags")
	// bind to viper
	_ = viper.BindPFlags(cmdset)
}

func addBuildFlags(cmdset *pflag.FlagSet) {
	addCommonFlags(cmdset)
	// bind to viper
	_ = viper.BindPFlags(cmdset)
}
