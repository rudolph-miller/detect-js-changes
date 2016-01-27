package main

import (
	"errors"
	"fmt"
	"github.com/Rudolph-Miller/detect-js-changes/detect_js_changes"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
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
		dir := path.Join(tmpDir, "detect_js_changes_download_"+suffix)
		result[index] = dir
		os.MkdirAll(dir, 0777)
	}
	return result
}

func getAvailableDir(dirs [2]string) (string, error) {
	var result string
	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return result, err
		}
		if len(files) == 0 {
			result = dir
			break
		}
	}
	if result != "" {
		return result, nil
	} else {
		msg := "No available directory\nPleaze reset"
		return result, errors.New(msg)
	}
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
				dirs := getDownloadDirs(config)
				dir, err := getAvailableDir(dirs)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("Directory: " + dir)
				for index, url := range urls {
					file := "file_" + strconv.Itoa(index)
					destination := path.Join(dir, file)
					err := detect_js_changes.Download(url, destination)
					if err != nil {
						fmt.Println("Download error")
						fmt.Println(err)
						os.Exit(1)
					}
					msg := "Download: " + url + " as " + file
					fmt.Println(msg)
				}
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
					fmt.Println("Reset: " + dir)
				}
			},
		},
	}

	app.Run(os.Args)
}
