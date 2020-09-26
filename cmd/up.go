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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		appPortOpen, appPort := getAvailablePort(80, 8080)
		databasePortOpen, databasePort := getAvailablePort(viper.GetInt("database.port"), 6306)

		if !appPortOpen || !databasePortOpen {
			fmt.Println("💥 Could not get access to any port after trying several options")
			os.Exit(0)
		}

		// Docker-Compose File
		contents := []byte(getDockerCompose(appPort, databasePort))
		if err := replaceFile(".visor/docker-compose.yaml", contents); err != nil {
			log.Println(err)
			log.Fatal("💥 Could not write .visor/docker-compose.yaml")
		}

		contents = []byte(getServerConfig())
		if err := replaceFile(".visor/default.conf", contents); err != nil {
			log.Println(err)
			log.Fatal("💥 Could not write visor/default.conf")
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
		//
		if directoryNotExists("vendor") {
			fmt.Println("")
			fmt.Println("👉 running composer install...")
			install := exec.Command("docker-compose", "--project-directory", ".", "--file", ".visor/docker-compose.yaml", "run", "--rm", "php", "composer", "install")

			if verbose {
				fmt.Println("")
				install.Stdout = os.Stdout
				install.Stderr = os.Stderr
			}

			install.Run()
		}

		fmt.Println("👉 spinning up services...")
		up := exec.Command("docker-compose", "--project-directory", ".", "--file", ".visor/docker-compose.yaml", "up", "-d")

		if verbose {
			fmt.Println("")
			up.Stdout = os.Stdout
			up.Stderr = os.Stderr
		}

		up.Run()

		appConnectionString := getAppConnectionString(appPort)
		databaseConnectionString := getDatabaseConnectionString(databasePort)

		// Print Tasks
		fmt.Println("")
		fmt.Println("💪 run `visor migrate` to migrate your database")
		fmt.Println("💪 run `visor down` to spin down your application and services")
		fmt.Println("💪 run `visor tinker` to jump into the php container")
		fmt.Println("")
		fmt.Printf("👌 Applicaton available at %s\n", appConnectionString)
		fmt.Printf("👌 Database available at %s\n", databaseConnectionString)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")
}

func getAppConnectionString(port int) string {
	return fmt.Sprintf("http://localhost:%s", strconv.Itoa(port))
}

func getDatabaseConnectionString(port int) string {
	return fmt.Sprintf("mysql://root@127.0.0.1:%s/%s", strconv.Itoa(port), viper.GetString("database.database"))
}

func getAvailablePort(desiredPort int, backupPort int) (portIsOpen bool, openPort int) {
	conn, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(desiredPort)))
	if err == nil {
		conn.Close()
		return true, desiredPort
	}
	maxPort := backupPort + 10
	for backupPort <= maxPort {
		conn, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(backupPort)))
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
  nginx:
    image: nginx:alpine
    tty: true
    ports:
      - "%s:80"
    volumes:
      - ./public:/var/www/public
      - ./.visor/default.conf:/etc/nginx/conf.d/default.conf
  php:
    image: cyberduck/php-fpm-laravel:7.4
    volumes:
      - ./:/var/www/
    environment:
      - "DB_HOST=mysql"
      - "DB_DATABASE=${DB_DATABASE}"
      - "DB_PORT=3306"
      - "DB_USERNAME=root"
      - "DB_PASSWORD="
      - "REDIS_HOST=redis"
      - "REDIS_PORT=6379"
    tty: true
    expose:
      - 9000
    command: php-fpm
  mysql:
    image: mysql:5.7
    volumes:
      - mysqldata:/var/lib/mysql
    environment:
      - "MYSQL_ALLOW_EMPTY_PASSWORD=true"
      - "MYSQL_DATABASE=${DB_DATABASE}"
    ports:
      - "%s:3306"
  redis:
    image: redis:4.0-alpine
    command: redis-server --appendonly yes
volumes:
  mysqldata:
`, strconv.Itoa(apptPort), strconv.Itoa(databasePort))
}

func getServerConfig() string {
	return `server {
    listen 80;
    server_name example.com;
    root /var/www/public;

    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-XSS-Protection "1; mode=block";
    add_header X-Content-Type-Options "nosniff";

    index index.php;

    charset utf-8;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    error_page 404 /index.php;

    location ~ \.php$ {
        fastcgi_pass php:9000;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.(?!well-known).* {
        deny all;
    }
}`
}
