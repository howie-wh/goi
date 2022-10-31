package cmd

import (
	"os"

	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/build"
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/config"
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/inject"
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/model"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "inject mock for your project and execute go build command",
	Long: `	Build command will copy the project code and its necessary dependencies to a temporary directory,
			then do mock for the target, binaries will be generated to their original place.`,
	Example: ``,

	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Fail to build: %v", err)
		}
		runBuild(args, wd)
	},
}

func init() {
	addBuildFlags(buildCmd.Flags())
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "Specifies the output execution file ： -o xxx")
	buildCmd.Flags().StringVarP(&customInject, "custom-inject", "c", "", "Specifying a configuration file： -c xxx.json")
	buildCmd.PersistentFlags().BoolVar(&spexMocker, "spex-mocker", false, "inject default spex-mocker code: --spex-mocker")
	buildCmd.PersistentFlags().BoolVar(&httpMocker, "http-mocker", false, "inject default http-mocker: --http-mocker")
	buildCmd.PersistentFlags().BoolVar(&cacheMocker, "cache-mocker", false, "inject default cache-mocker: --cache-mocker")
	rootCmd.AddCommand(buildCmd)
}

func runBuild(args []string, wd string) {
	goiBuild, err := build.NewBuild(buildFlags, args, wd, buildOutput)
	//goiBuild.InjectedProjectDir, _ = filepath.Abs(injectedProject)
	if err != nil {
		log.Fatalf("Fail to build: %v", err)
	}
	// remove GocBuild directory if needed
	defer func() {
		if err = goiBuild.Clean(); err != nil {
			log.Fatalf("Fail to goi build clean: %v", err)
		}
	}()

	injectConfig := getInjectConfig()
	if injectConfig == nil {
		log.Fatalf("Fail to get inject config")
	}
	injection := &inject.Injection{
		BuildArgs:        buildFlags,
		GoPath:           goiBuild.GocBuild.NewGOPATH,
		TargetProjectDir: goiBuild.GocBuild.TmpDir,
		Config:           injectConfig,
		ResourceDir:      goiBuild.TmpDir,
	}
	err = injection.Inject()
	if err != nil {
		log.Fatalf("Fail to Inject Code: %v", err)
	}

	// do install in the temporary directory
	if goc {
		err = goiBuild.Build(debugGoi)
	} else {
		err = goiBuild.GocBuild.Build()
	}
	if err != nil {
		log.Fatalf("Fail to build: %v", err)
	}
}

func getInjectConfig() *model.InjectConfig {
	if customInject != "" {
		return config.GetCustomInjectConfig(customInject)
	}

	injectConfig := &model.InjectConfig{}
	if spexMocker {
		if spexInjectConfig := config.GetDefaultSpexInjectConfig(); spexInjectConfig != nil {
			config.MergeInjectConfig(injectConfig, spexInjectConfig)
		}
	}
	if httpMocker {
		if httpInjectConfig := config.GetDefaultHTTPInjectConfig(); httpInjectConfig != nil {
			config.MergeInjectConfig(injectConfig, httpInjectConfig)
		}
	}
	if cacheMocker {
		if cacheInjectConfig := config.GetDefaultCacheInjectConfig(); cacheInjectConfig != nil {
			config.MergeInjectConfig(injectConfig, cacheInjectConfig)
		}
	}

	return injectConfig
}
