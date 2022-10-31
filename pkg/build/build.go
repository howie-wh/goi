/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package build

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	gocBuild "github.com/qiniu/goc/pkg/build"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GoiBuild struct {
	GocBuild *gocBuild.Build
	TmpDir   string
}

func NewBuild(buildFlags string, args []string, workingDir string, outputDir string) (*GoiBuild, error) {
	childBuild, err := gocBuild.NewBuild(buildFlags, args, workingDir, outputDir)
	if err != nil {
		return nil, err
	}
	//b.TmpDir
	goiBuild := &GoiBuild{
		GocBuild: childBuild,
	}

	goiBuild.TmpDir = filepath.Join(os.TempDir(), tmpFolderName(childBuild.WorkingDir))
	if err = os.RemoveAll(goiBuild.TmpDir); err != nil {
		return nil, err
	}
	if err = os.MkdirAll(goiBuild.TmpDir, os.ModePerm); err != nil {
		return nil, err
	}
	return goiBuild, nil
}

// Clean clears up the temporary workspace
func (b *GoiBuild) Clean() error {
	if err := b.GocBuild.Clean(); err != nil {
		return err
	}
	if !viper.GetBool("debug") {
		return os.RemoveAll(b.TmpDir)
	}
	return nil
}

// Build calls 'go build' tool to do building
func (b *GoiBuild) Build(debugGoi bool) error {
	log.Infof("Goc building in temp...")
	debugFlag := ""
	if debugGoi {
		debugFlag = "--debug"
	}
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("goc build %v --buildflags=\"%v\" -o %v %v", debugFlag, b.GocBuild.BuildFlags, b.GocBuild.Target, b.GocBuild.Packages))
	cmd.Dir = b.GocBuild.TmpWorkingDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if b.GocBuild.NewGOPATH != "" {
		// Change to temp GOPATH for go install command
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%v", b.GocBuild.NewGOPATH))
	}

	log.Infof("goc build cmd is: %v", cmd.Args)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	log.Infof("Goc build exit successful.")
	return nil
}

// tmpFolderName uses the first six characters of the input path's SHA256 checksum
// as the suffix.
func tmpFolderName(path string) string {
	sum := sha256.Sum256([]byte(path))
	h := fmt.Sprintf("%x", sum[:6])

	return "goi-build-" + h
}
