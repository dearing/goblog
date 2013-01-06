package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ContentFolder  string
	TemplateFolder string
	Suffix         string
	WWWHost        string
	WWWRoot        string
	RedisHost      string
	RedisPass      string
	RedisDB        int64
	Verbose        bool
	EnableWWW      bool
}

// Load up a JSON config file.
func (c *Config) LoadConfig(path string) {

	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}

	if err := json.Unmarshal(f, &c); err != nil {
		log.Panicln(err)
	}

}

// Generate a default config in the current directory for the user to manipulate.
func (c *Config) GenerateConfig(path string) {

	c = &Config{
		ContentFolder:  "content",
		TemplateFolder: "templates",
		Suffix:         ".md",
		WWWHost:        ":9002",
		WWWRoot:        "example",
		RedisHost:      "localhost:6379",
		RedisPass:      "",
		RedisDB:        -1,
		Verbose:        true,
		EnableWWW:      false,
	}

	b, _ := json.MarshalIndent(c, "", "\t")
	ioutil.WriteFile(*conf, b, 0644)

}
