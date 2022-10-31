package inject

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/model"
	"git.garena.com/shopee/seller-server/seller-marketing/marketing-common/goi/pkg/util"

	log "github.com/sirupsen/logrus"
)

func injectImport(importConfig *model.ImportConfig, targetProjectDir string) error {
	targetFile := path.Join(targetProjectDir, importConfig.FileName)
	fileBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("read file %s: %s", targetFile, err)
	}

	fileContent := string(fileBytes)
	if strings.Contains(fileContent, importConfig.InjectCode) {
		return nil
	}

	fSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fSet, targetFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	log.Infof("targetFile: %v", targetFile)
	log.Infof("parsedFile.Name: %v", parsedFile.Name.Name)
	edit := util.NewBuffer(fileBytes)
	edit.Insert(int(parsedFile.Name.End())-1, fmt.Sprintf("; %s", importConfig.InjectCode))
	fd, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	_, err = fd.Write(edit.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func injectFunction(funcConfig *model.FunctionConfig, targetProjectDir string) error {
	targetFile := path.Join(targetProjectDir, funcConfig.FileName)
	fSet := token.NewFileSet()
	fileBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Read file %s: %s", targetFile, err)
	}
	fileContent := string(fileBytes)
	if !strings.Contains(fileContent, funcConfig.FuncName) {
		log.Fatalf("Func not find %s", funcConfig.FuncName)
		return nil
	}

	parsedFile, err := parser.ParseFile(fSet, targetFile, fileBytes, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse file: %s: %s", targetFile, err)
	}
	edit := util.NewBuffer(fileBytes)
	ast.Inspect(parsedFile, func(n ast.Node) bool {
		// Find Return Statements
		if ret, ok := n.(*ast.FuncDecl); ok {
			if ret.Name.Name == funcConfig.FuncName {
				//pos 和 end 都要减 1
				// 将代码注入到相对行位置，
				pos := endOfLine(fileBytes, int(ret.Pos())-1, int(ret.End())-1, funcConfig.RelativeLine)
				if pos != 0 {
					if fileBytes[pos-1] == '\n' {
						// 行首不加分号
						edit.Insert(pos, funcConfig.InjectCode)
					} else {
						// 非行首加分号
						edit.Insert(pos, fmt.Sprintf(";%s", funcConfig.InjectCode))
					}
				}
			}
			return false
		}
		return true
	})

	fd, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	if _, err = fd.Write(edit.Bytes()); err != nil {
		return err
	}
	return nil
}

// 找到相对行,找到行尾坐标
func endOfLine(fileBytes []byte, start, end, relativeLine int) int {
	if relativeLine >= 0 {
		// 向前找
		for i := start; i < end; i++ {
			if fileBytes[i] == '\n' {
				if relativeLine <= 0 {
					return i
				}
				relativeLine = relativeLine - 1
			}
		}
	} else {
		// 向后找
		for i := start; i >= 0; i-- {
			if fileBytes[i] == '\n' {
				relativeLine = relativeLine + 1
				if relativeLine >= 0 {
					return i
				}
			}
		}
	}
	return 0
}
