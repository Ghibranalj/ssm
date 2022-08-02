package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Service struct {
	Endpoint     string   `yaml:"endpoint"`
	Email        string   `yaml:"email"`
	Env          string   `yaml:"env"`
	Destinations []string `yaml:"destinations"`
	Bc           []string `yaml:"cc"`
	Bcc          []string `yaml:"bcc"`
	BodyFormat   string   `yaml:"bodyFormat"`
	Cors         string   `yaml:"cors"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Port     string    `yaml:"port"`
}

type cliOpts struct {
	configFile string
	envFile    string
}

var opts cliOpts

func main() {
	file, err := ioutil.ReadFile(opts.configFile)
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(os.Stderr, "Error loading yaml file: %v\n", err)
	}
	config := Config{}
	err = yaml.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s file: %v\n", opts.configFile, err)
		os.Exit(1)
	}
	err = godotenv.Load(opts.envFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}
}
func init() {
	confFile := flag.String("c", "config.yml", "config file")
	envFile := flag.String("e", ".env", "env file")

	flag.Parse()
	opts.configFile = *confFile
	opts.envFile = *envFile
}
