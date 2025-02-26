/*
Copyright © 2025 EspressoCake
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Retrieve ports/connection string for running Bloodhound Community Edition containers",
	Long:  `Retrieve ports/connection string for running containers via SSH or assumed to be local`,
	Run:   retrievePorts,
}

var (
	uniqueName      string
	localConnection bool
)

func init() {
	portsCmd.Flags().StringVarP(&uniqueName, "prefix", "p", "", "Prefix for project containers to query")
	portsCmd.MarkFlagRequired("prefix")

	portsCmd.Flags().BoolVarP(&localConnection, "local", "l", true, "Boolean if the connection is meant to be local or remote")

	rootCmd.AddCommand(portsCmd)
}

func retrievePorts(cmd *cobra.Command, args []string) {
	dockerNeedle, err := cmd.Flags().GetString("prefix")
	if err != nil || len(dockerNeedle) == 0 {
		log.Fatal(err)
	}

	localInstance, err := cmd.Flags().GetBool("local")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string][]string)

	for _, container := range containers {
		name := container.Names
		ports := container.Ports

		currentRegex := regexp.MustCompile(`neo4j\-inst\-` + dockerNeedle)

		if currentRegex.Match([]byte(name[0])) {
			for _, port := range ports {
				if port.IP == "0.0.0.0" && port.PublicPort != 0 {
					_, ok := results[currentRegex.FindString(name[0])]
					if !ok {
						results[currentRegex.FindString(name[0])] = []string{}
					}

					if localInstance {
						switch port.PrivatePort {
						case 7687:
							results[currentRegex.FindString(name[0])] = append(results[currentRegex.FindString(name[0])], fmt.Sprintf("%-25s %d", "NEO4J_Database_Port:", port.PublicPort))
						case 7474:
							results[currentRegex.FindString(name[0])] = append(results[currentRegex.FindString(name[0])], fmt.Sprintf("%-25s %d", "NEO4J_Web_Port:", port.PublicPort))
						case 8080:
							results[currentRegex.FindString(name[0])] = append(results[currentRegex.FindString(name[0])], fmt.Sprintf("%-25s %d", "Bloodhound_Web_Port:", port.PublicPort))
						}
					} else {
						results[currentRegex.FindString(name[0])] = append(results[currentRegex.FindString(name[0])], fmt.Sprintf("-L %d:localhost:%d", port.PrivatePort, port.PublicPort))
					}
				}
			}
		}
	}

	if len(results) == 0 {
		fmt.Println("Didn't return any results, please modify your query.")

		return
	}

	for item := range results {
		fmt.Println(item)
		fmt.Println(strings.Repeat("=", len(item)))

		switch localInstance {
		case false:
			fmt.Printf("ssh %s username@SERVER_IP\n\n", strings.Join(results[item], " "))
		case true:
			fmt.Printf("%s\n", strings.Join(results[item], "\n"))
		}
	}

}
