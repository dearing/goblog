package main

import (
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
)

// Push a files contents up to the redis db after processing as markdown
func push(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	test := client.Set(filename, string(blackfriday.MarkdownCommon(body)))

	if test.Err() != nil {
		return test.Err()
	}

	if *verbose {
		log.Printf("pushed %s", filename)
	}

	return nil
}

// Reads a folder for files that are not folders itself (one level only) and pushes to the redis server as the folder + filename.
// TODO: needs a better naming scheme that will unfold when I get around to organizing data on the db
func pushall(folder string) error {

	files, _ := ioutil.ReadDir("articles")

	// for each file in the folder that, isn't a folder itself, push the parsed contents up
	for _, file := range files {
		if !file.IsDir() {
			err := push("articles/" + file.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

// Load an article to the db by reading it then parsing it as markdown and then pushing it up to the redis db by name
func pull(title string) (*Article, error) {

	key := fmt.Sprintf(title)

	if !client.Exists(key).Val() {
		if *verbose {
			log.Printf("not found: %v", key)
			return nil, errors.New(fmt.Sprintf("db does not contain key: %v", key))
		}
	}

	output := client.Get(key).Val()
	return &Article{Title: title, Body: template.HTML(output)}, nil
}
