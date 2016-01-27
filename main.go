package main

import (
	"fmt"
	"github.com/Rudolph-Miller/detect-js-changes/detect_js_changes"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Urls   []string
	TmpDir string `yaml:"tmp_dir"`
}

func setDefaultConfig(config *Config) {
	if config.TmpDir == "" {
		config.TmpDir = "/tmp"
	}
}

func getConfig(file string, env string) *Config {
	data, _ := ioutil.ReadFile(file)
	parsed := make(map[string]Config)
	err := yaml.Unmarshal(data, &parsed)
	if err != nil {
		fmt.Println("Parse error")
		os.Exit(1)
	}
	config := parsed[env]
	setDefaultConfig(&config)
	return &config
}

func getDownloadDirs(config *Config) [2]string {
	tmpDir := config.TmpDir
	var result [2]string
	suffixes := [2]string{"1", "2"}
	for index, suffix := range suffixes {
		result[index] = tmpDir + "/detect_js_changes_download_" + suffix
	}
	return result
}

func main() {
	app := cli.NewApp()
	app.Name = "detect-js-changes"
	app.Usage = "detects JS changes"
	app.Version = "0.0.1"

	var env string
	var configFile string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "env, e",
			Usage:       "env for cofig file",
			EnvVar:      "ENV",
			Value:       "default",
			Destination: &env,
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "config file",
			Value:       "detect_config.yml",
			Destination: &configFile,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "detect",
			Usage: "detects changes",
			Action: func(c *cli.Context) {
				println("detect changes")
			},
		},
		{
			Name:  "download",
			Usage: "downloads JS files",
			Action: func(c *cli.Context) {
				config := getConfig(configFile, env)
				urls := config.Urls
				println(urls)
				println("download JS files")
			},
		},
		{
			Name:  "reset",
			Usage: "resets downloaded JS files",
			Action: func(c *cli.Context) {
				config := getConfig(configFile, env)
				dirs := getDownloadDirs(config)
				for _, dir := range dirs {
					detect_js_changes.Reset(dir)
				}
			},
		},
	}

	app.Run(os.Args)
}
