package cmd

import (
	"fmt"
	"github.com/evandro-slv/arcentry-docker/api"
	"github.com/evandro-slv/arcentry-docker/stats"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/evandro-slv/arcentry-docker/parser"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arcentry-docker",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := os.Open(cfgFile)
		bytes, err := ioutil.ReadAll(file)

		if err != nil {
			log.Fatal(err)
		}

		config, err := parser.ReadConfig(bytes)

		if err != nil {
			log.Fatal(err)
		}

		interval := 5 * time.Second

		if config.Config.Watch.Interval != "" {
			interval, err = parser.ParseDuration(config.Config.Watch.Interval)

			if err != nil {
				fmt.Println(fmt.Sprintf("Error parsing duration from YAML file: %s, using default of '5s'", err))
			}
		}

		fmt.Println(fmt.Sprintf("watching for changes in %s", interval))

		dataArr := make(map[string][][]float64)
		pos := make(map[string]int)

		for true {
			st, err := stats.GetStats()

			if err != nil {
				log.Fatal(err)
			}

			req := parser.Request{
				Objects: make(map[string]interface{}),
			}

			for _, container := range config.Config.Containers {
				for k, v := range container {
					for _, stat := range st {
						if stat != nil && stat.Container == k {
							if v.Chart.Cpu != "" {
								key := v.Chart.Cpu
								cpu, err := strconv.ParseFloat(strings.Replace(strings.Replace(stat.Cpu, "%", "", -1), ",", ".", -1), 64)

								if err != nil {
									log.Fatal(err)
								}

								dataArr[key] = append(dataArr[key], []float64{float64(pos[key]), cpu})
								pos[key] += 1

								if len(dataArr[key]) > 5 {
									dataArr[key] = dataArr[key][1:] // Dequeue
								}

								req.Objects[key] = parser.ChartDocument{
									Data:       dataArr[key],
									YAxisSpace: "custom",
									MinY:       0,
									MaxY:       100,
									//Data: [][]float64 {{float64(1), float64(10)}, {float64(2), float64(20)}},
								}
							}

							if v.Text.Cpu != "" {
								key := v.Text.Cpu
								cpu := stat.Cpu

								req.Objects[key] = parser.TextDocument{
									Text: fmt.Sprint(cpu),
								}
							}

							if v.Chart.Memory != "" {
								key := v.Chart.Memory
								mem, err := strconv.ParseFloat(strings.Replace(strings.Replace(stat.Memory.Percent, "%", "", -1), ",", ".", -1), 64)

								if err != nil {
									log.Fatal(err)
								}

								dataArr[key] = append(dataArr[key], []float64{float64(pos[key]), mem})
								pos[key] += 1

								if len(dataArr[key]) > 5 {
									dataArr[key] = dataArr[key][1:] // Dequeue
								}

								req.Objects[key] = parser.ChartDocument{
									Data:       dataArr[key],
									YAxisSpace: "custom",
									MinY:       0,
									MaxY:       100,
									//Data: [][]float64 {{float64(1), float64(10)}, {float64(2), float64(20)}},
								}
							}

							if v.Text.Memory != "" {
								key := v.Text.Memory
								cpu := stat.Memory.Percent

								req.Objects[key] = parser.TextDocument{
									Text: fmt.Sprint(cpu),
								}
							}
						}
					}
				}
			}

			reqs, err := parser.ParseRequest(&req)

			if err != nil {
				log.Fatal(err)
			}

			err = api.UpdateDoc(config.Config.Arcentry, reqs)

			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(interval)
		}
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.arcentry-docker.yaml)")

	//err := rootCmd.MarkPersistentFlagRequired("config")

	//if err != nil {
	//	log.Fatal(err)
	//}
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

		// Search config in home directory with name ".arcentry-docker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".arcentry-docker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
