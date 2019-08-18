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
	"os"
	"time"

	"github.com/joostvdg/jpb-mon-go/pkg/prometheus"
	"github.com/spf13/cobra"

	"github.com/joostvdg/jpb-mon-go/pkg/pipelinerun"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var host string
var job string
var runId int
var username string
var password string
var promEndpoint string
var push bool
var sleepTimeInMinutes time.Duration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jpb-mon-go",
	Short: "Jenkins Pipeline Binary - Monitoring",
	Long:  `This is a commandline tool for interacting with Prometheus and Jenkins Pipeline Runs`,
}

var getPipelineRunCommand = &cobra.Command{
	Use:   "get-run",
	Short: "Get a Jenkins Pipeline Run",
	Long:  `Retrieves a Pipeline Run from Jenkins' BlueOcean API'`,
	Run: func(cmd *cobra.Command, args []string) {
		pipelineRun := pipelinerun.GetPipelineRun(host, job, runId, username, password)
		if push {
			fmt.Println("> Push to Prometheus Push Gateway enabled, attempt to push metrics")
			pipelineRunMetadata := pipelinerun.PipelineRunMetadata{
				RunId:    runId,
				Instance: host,
				Job:      job,
			}
			prometheus.PushPipelineRunToGateway(promEndpoint, pipelineRun, pipelineRunMetadata)
		} else {
			fmt.Println("> Push to Prometheus disable, not doing anything with he pipeline runs")
		}
	},
}

var sleepCommand = &cobra.Command{
	Use:   "sleep",
	Short: "Sleeps given number of minutes",
	Long:  `Will sleep for the duration of given number in minutes`,
	Run: func(cmd *cobra.Command, args []string) {
		sleepTime := sleepTimeInMinutes * time.Minute
		fmt.Printf("Sleeping for %v\n", sleepTime)
		time.Sleep(sleepTime)
		fmt.Println("Returning")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jpb-mon-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	getPipelineRunCommand.Flags().StringVar(&host, "host", "", "Host to reach Jenkins at, http or https")
	getPipelineRunCommand.Flags().StringVar(&job, "job", "", "Job name to retrieve Pipeline Run from")
	getPipelineRunCommand.Flags().IntVar(&runId, "run", 0, "Job Run ID to retrieve Pipeline Run from")
	getPipelineRunCommand.Flags().StringVar(&username, "username", "", "Username for authentication with Jenkins")
	getPipelineRunCommand.Flags().StringVar(&password, "password", "", "Password for authentication with Jenkins")
	getPipelineRunCommand.Flags().StringVar(&promEndpoint, "prom", "http://prometheus-pushgateway.obs:9091", "Endpoint for the Prometheus Push Gateway")
	getPipelineRunCommand.Flags().BoolVar(&push, "push", false, "Push metrics to Prometheus Push Gateway")

	sleepCommand.Flags().DurationVar(&sleepTimeInMinutes, "sleep", 1, "Sleep time in minutes")

	rootCmd.AddCommand(getPipelineRunCommand)
	rootCmd.AddCommand(sleepCommand)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jpb-mon-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jpb-mon-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
