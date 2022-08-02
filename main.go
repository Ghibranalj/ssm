package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
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
	Cors         []string `yaml:"cors"`
}

type Config struct {
	Services []Service `yaml:"services"`
	Port     string    `yaml:"port"`
}

type cliOpts struct {
	configFile string
	envFile    string
}

var (
	opts   cliOpts
	router *mux.Router
)

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

	router = mux.NewRouter()
	for _, service := range config.Services {
		GenerateEndpoint(service)
	}
	http.Handle("/", router)

	port := config.Port
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
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

func GenerateEndpoint(service Service) {

	endpoint := service.Endpoint
	corsString := strings.Join(service.Cors, ",")

	// TODO get the body format from the config & get the field names

	bodyTemplate := template.New("body")
	bodyTemplate, err := bodyTemplate.Parse(service.BodyFormat)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing body format in service %s: %v\n", service.Endpoint, err)
		os.Exit(1)
	}

	router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", corsString)

		body := json.NewDecoder(r.Body)

		var data map[string]interface{}
		err := body.Decode(&data)
		if err != nil {
			log.Printf("Error parsing body: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing body: %s\n", err.Error())
			return
		}
		log.Printf("Received request: %v\n", data)

		sb := &strings.Builder{}
		err = bodyTemplate.Execute(sb, data)
		if err != nil {
			log.Printf("Error executing template: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error executing template: %s\n", err.Error())
			return
		}
		log.Printf("Sending email: %s\n", sb.String())
	}).Methods("POST")

	log.Printf("Endpoint %s created succesfully\n", endpoint)
}
