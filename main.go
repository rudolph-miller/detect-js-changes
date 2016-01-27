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
	"sync"
)

type Config struct {
	Urls           []string
	TmpDir         string   `yaml:"tmp_dir"`
	IgnoreKeywords []string `yaml:"ignore_keywords"`
}

func setDefaultValue(config *Config) {
	if config.TmpDir == "" {
		config.TmpDir = "/tmp"
	}
}

func mergeDefaultConfig(config *Config, defaultConfig *Config) {
	if len(config.Urls) == 0 {
		config.Urls = defaultConfig.Urls
	}

	if config.TmpDir == "" {
		config.TmpDir = defaultConfig.TmpDir
	}

	if len(config.IgnoreKeywords) == 0 {
		config.IgnoreKeywords = defaultConfig.IgnoreKeywords
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
	if env != "default" {
		defaultConfig := parsed["default"]
		mergeDefaultConfig(&config, &defaultConfig)
	}
	setDefaultValue(&config)
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

func getFileName(index int) string {
	return "file_" + strconv.Itoa(index)
}

func formatResult(url string, result detect_js_changes.Result) string {
	switch result {
	case detect_js_changes.HasSomeChanges:
		return url + " has some changes"
	case detect_js_changes.HasNoChanges:
		return url + " has no changes"
	case detect_js_changes.HasIgnoredChanges:
		return url + " has ignored changes"
	}
	return ""
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
				hasSomeChange := false
				config := getConfig(configFile, env)
				urls := config.Urls
				dirs := getDownloadDirs(config)
				ignoreKeywords := config.IgnoreKeywords
				for _, dir := range dirs {
					files, err := ioutil.ReadDir(dir)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					if len(files) == 0 {
						fmt.Println("Please execute download twice")
						os.Exit(1)
					}
				}

				for index, url := range urls {
					fmt.Println("Detecting: " + url)
					filename := getFileName(index)
					file1 := path.Join(dirs[0], filename)
					file2 := path.Join(dirs[1], filename)
					result := detect_js_changes.Detect(file1, file2, ignoreKeywords)
					if result == detect_js_changes.HasSomeChanges {
						hasSomeChange = true
					}
					fmt.Println("Result: " + formatResult(url, result))
				}

				if hasSomeChange {
					os.Exit(1)
				}
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

				var wg sync.WaitGroup

				for index, url := range urls {
					file := getFileName(index)
					destination := path.Join(dir, file)
					wg.Add(1)

					go func(url string) {
						err := detect_js_changes.Download(url, destination)
						if err != nil {
							fmt.Println("Download error")
							fmt.Println(err)
							os.Exit(1)
						}
						msg := "Download: " + url + " as " + file
						fmt.Println(msg)
						wg.Done()
					}(url)
				}
				wg.Wait()
			},
		},
		{
			Name:  "reset",
			Usage: "resets downloaded JS files",
			Action: func(c *cli.Context) {
				config := getConfig(configFile, env)
				dirs := getDownloadDirs(config)
				args := c.Args()
				if len(args) > 0 {
					num, err := strconv.Atoi(args[0])
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					if num < 1 || num > 2 {
						fmt.Println("Please specify 1 or 2")
						os.Exit(1)
					}

					dir := dirs[num-1]
					detect_js_changes.Reset(dir)
					fmt.Println("Reset: " + dir)
				} else {
					for _, dir := range dirs {
						detect_js_changes.Reset(dir)
						fmt.Println("Reset: " + dir)
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
