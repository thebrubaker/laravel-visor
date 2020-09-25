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
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			log.Fatal(err)
		}
		// Get Ports
		databasePortOpen, databasePort := getAvailablePort(3306, 6306)
		appPortOpen, appPort := getAvailablePort(80, 8080)

		if !appPortOpen || !databasePortOpen {
			log.Fatal("💥 Could not get access to any port after trying several options")
		}

		// Docker-Compose File
		if fileNotExists("docker-compose.yaml") && fileNotExists("docker-compose.yml") {
			contents := []byte(getDockerCompose(appPort, databasePort))
			if err := replaceFile(".visor/docker-compose.yaml", contents); err != nil {
				log.Println(err)
				log.Fatal("💥 Could not write docker-compose.yaml file")
			}
		}

		// Docker-Compose Services Running
		// output, err := exec.Command("docker-compose", "--file", ".visor/docker-compose.yaml", "ps", "--filter", "status=running", "--services").Output()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// services := strings.Fields(string(output))
		// if len(services) == 0 {
		// 	logs = append(logs, "👌 Visor services not already running")
		// } else {
		// 	logs = append(logs, "👌 Visor services already running")
		// }

		fmt.Println("")
		fmt.Println("👉 running composer install...")
		command := exec.Command("docker-compose", "--project-directory", ".", "--file", ".visor/docker-compose.yaml", "run", "--rm", "php", "composer", "install")

		if verbose {
			fmt.Println("")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
		}

		command.Run()

		fmt.Println("👉 spinning up services...")
		command = exec.Command("docker-compose", "--project-directory", ".", "--file", ".visor/docker-compose.yaml", "up", "-d")

		if verbose {
			fmt.Println("")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
		}

		command.Run()

		// Print Tasks
		fmt.Println("")
		fmt.Println("💪 run `visor migrate` to migrate your database")
		fmt.Println("💪 run `visor down` to spin down your application and services")
		fmt.Println("💪 run `visor tinker` to jump into the php container")
		fmt.Println("")
		fmt.Printf("👌 Applicaton available at http://localhost:%s\n", strconv.Itoa(appPort))
		fmt.Printf("👌 Database available at mysql://root@127.0.0.1:%s/laravel_visor\n", strconv.Itoa(databasePort))
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")
}

func getAvailablePort(desiredPort int, backupPort int) (portIsOpen bool, openPort int) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(desiredPort)), time.Second)
	if err == nil {
		conn.Close()
		return true, desiredPort
	}
	maxPort := backupPort + 10
	for backupPort <= maxPort {
		conn, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(backupPort)))
		if err == nil {
			conn.Close()
			return true, backupPort
		}
		backupPort++
	}
	return false, backupPort
}

func getDockerCompose(apptPort int, databasePort int) string {
	return fmt.Sprintf(`version: "3"
services:
  php:
    image: lorisleiva/laravel-docker
    volumes:
      - ./:/var/www/
    environment:
      - "DB_HOST=mysql"
      - "DB_DATABASE=${DB_DATABASE}"
      - "DB_USERNAME=${DB_USERNAME}"
      - "DB_PASSWORD=${DB_PASSWORD}"
      - "REDIS_HOST=redis"
      - "REDIS_PORT=${REDIS_PORT}"
    command: php artisan serve --host=0.0.0.0 --port=80
    ports:
      - %s:80
  mysql:
    image: mysql:5.7
    volumes:
      - mysqldata:/var/lib/mysql
    environment:
      - "MYSQL_DATABASE=${DB_DATABASE}"
      - "MYSQL_USER=${DB_USERNAME}"
      - "MYSQL_PASSWORD=${DB_PASSWORD}"
    ports:
      - "%s:3306"
  redis:
    image: redis:4.0-alpine
    command: redis-server --appendonly yes  --port ${REDIS_PORT}
volumes:
  mysqldata:`, strconv.Itoa(apptPort), strconv.Itoa(databasePort))
}
