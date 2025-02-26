/*
Copyright © 2025 EspressoCake
*/
package cmd

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var (
	//go:embed templates/bloodhound.json.tmpl
	jsonConfiguration embed.FS

	//go:embed templates/config.tmpl
	configurationTemplate embed.FS
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Create configurations for projects",
	Long:  `Generates respective configruation(s) for Docker instances to deploy Bloodhound Community Edition via docker-compose.yml files`,
	Run:   execute,
}

var (
	path string
	name string
)

func init() {
	generateCmd.Flags().StringVarP(&path, "path", "p", "", "Filepath on system to desired root directory. Default will be the current working directory")
	generateCmd.Flags().StringVarP(&name, "name", "n", "", "Name for project, in lowercase")
	generateCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(generateCmd)
}

func execute(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("path")
	name, _ := cmd.Flags().GetString("name")

	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = currentPath
	}

	createPathAndGenerateConfigurations(filepath.Join(path, fmt.Sprintf("neo4j-inst-%s", strings.ToLower(name))), strings.ToLower(name))
}

func createPathAndGenerateConfigurations(path string, name string) {
	type (
		configuration struct {
			Codename   string
			Password   string
			CEPassword string
			OSPath     string
		}
	)

	currentConfiguration := configuration{
		Codename:   name,
		Password:   passwordGenerator(),
		CEPassword: passwordGenerator(),
		OSPath:     path,
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/bloodhound.json", currentConfiguration.OSPath)); os.IsNotExist(err) {
		currentTemplate := template.Must(template.ParseFS(jsonConfiguration, "templates/bloodhound.json.tmpl"))

		f, err := os.Create(fmt.Sprintf("%s/bloodhound.json", currentConfiguration.OSPath))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = currentTemplate.Execute(f, currentConfiguration)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(fmt.Sprintf("%s/docker-compose.yml", currentConfiguration.OSPath)); os.IsNotExist(err) {
		currentTemplate := template.Must(template.ParseFS(configurationTemplate, "templates/config.tmpl"))

		f, err := os.Create(fmt.Sprintf("%s/docker-compose.yml", currentConfiguration.OSPath))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = currentTemplate.Execute(f, currentConfiguration)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Password for Bloodhound CE Web Server:  %s\n", currentConfiguration.CEPassword)
		fmt.Printf("Current password for your Neo4j is:     %s\n", currentConfiguration.Password)
		fmt.Printf("Go to the following directory:          %s\n", currentConfiguration.OSPath)
		fmt.Printf("Run the following:                      %s\n", "docker compose up -d OR docker-compose up -d")
	}
}

func passwordGenerator() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	characterOptions := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, 12)
	for i := range b {
		b[i] = characterOptions[rand.Intn(len(characterOptions))]
	}
	return string(b)
}
