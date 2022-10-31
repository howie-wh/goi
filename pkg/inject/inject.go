package inject

import (
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/model"
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/util"
)

type Injection struct {
	SourceProjectDir string // 被注入工程目录
	TargetProjectDir string // 注入目标工程目录
	GoPath           string
	BuildArgs        string
	ResourceDir      string // 存放下载资源的临时目录
	Config           *model.InjectConfig
}

func (i *Injection) Inject() error {
	if err := i.injectImports(); err != nil {
		return err
	}
	if err := i.injectFunctions(); err != nil {
		return err
	}
	if err := i.tidyDep(); err != nil {
		return err
	}
	return nil
}

func (i *Injection) injectImports() error {
	for _, importConfig := range i.Config.ImportConfigs {
		if err := injectImport(importConfig, i.TargetProjectDir); err != nil {
			return err
		}
	}
	return nil
}

func (i *Injection) injectFunctions() error {
	for _, funcConfig := range i.Config.FunctionConfigs {
		if err := injectFunction(funcConfig, i.TargetProjectDir); err != nil {
			return err
		}
	}
	return nil
}

func (i *Injection) tidyDep() error {
	return util.ExecuteCmd("go mod tidy && go mod vendor", i.TargetProjectDir)
}
