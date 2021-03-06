/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize visor by creating a .visor directory and gitignore it",
	Run: func(cmd *cobra.Command, args []string) {
		initVisor()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initVisor() {
	if err := appendToFileIfMissing(".gitignore", ".visor"); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir(".visor", 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("")
	fmt.Println("👌 Docker is installed")
	fmt.Println("👌 Created .visor directory and added to .gitignore")
	fmt.Println("")

	fmt.Println("👉 downloading containers for php 7.4, redis 4.0 and mysql 5.7...")
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond) // Build our new spinner
	s.Start()

	exec.Command("docker", "pull", "lorisleiva/laravel-docker").Run()
	exec.Command("docker", "pull", "mysql:5.7").Run()
	exec.Command("docker", "pull", "redis:4.0-alpine").Run()

	s.Stop()

	fmt.Println("")
	fmt.Println("👌 Visor init success!")
	fmt.Println("")
}
