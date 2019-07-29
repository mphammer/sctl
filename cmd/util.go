/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
	//git "gopkg.in/src-d/go-git.v4"
)

// ReadFileData returns the information within a file as a byte slice
func ReadFileData(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read from file %s: %s", filepath, err)
	}
	return data, nil
}

// AddFileToGit ...
func AddFileToGit(data []byte, destination string) error {
	// GITFILENAME ...
	GITFILENAME := "personalized-on-prem-alert-final.yaml"
	// Verify correct location
	if !Exists(GITFILENAME) {
		return fmt.Errorf("Bad location - could not find %s", GITFILENAME)
	}
	// Verify you have the latest git
	if err := ExecCmd(exec.Command("git", "pull")); err != nil {
		return err
	}
	// Write to the file
	err := ioutil.WriteFile(GITFILENAME, data, 0666)
	if err != nil {
		return err
	}
	// Push the file's changes to git
	if err := ExecCmd(exec.Command("git", "add", fmt.Sprintf("%s", GITFILENAME))); err != nil {
		return err
	}
	if err := ExecCmd(exec.Command("git", "commit", "-m", fmt.Sprintf("updated at %s", time.Now().String()))); err != nil {
		return err
	}
	if err := ExecCmd(exec.Command("git", "push")); err != nil {
		return err
	}
	return nil
}

// ExecCmd ...
func ExecCmd(osCmd *exec.Cmd) error {
	osCmd.Stdin = os.Stdin
	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr
	err := osCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
