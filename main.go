// Copyright 2016 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/yosssi/gohtml"
)

var (
	overwrite = flag.Bool("w", false, "overwrite when -l option is specified")
	output    = flag.String("o", "", "save file with this name")
)

type InputType int

const (
	FILE InputType = iota
	URL
)

func main() {
	flag.Parse()

	typ, err := inputType(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	var data string
	switch typ {
	case FILE:
		data, err = openLocal(flag.Arg(0))
	case URL:
		data, err = openURL(flag.Arg(0))
	}

	var file *os.File
	if err == nil {
		switch {
		case *overwrite && typ == FILE:
			file, err = os.Create(flag.Arg(0))
		case *overwrite && typ == URL:
			filename := path.Base(flag.Arg(0))
			if filename == "/" || filename == "." {
				filename = "result.html"
			}
			if path.Ext(filename) == "" {
				filename += ".html"
			}
			file, err = os.Create(filename)
		case *output != "":
			file, err = os.Create(*output)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	result := gohtml.Format(data)
	if file != nil {
		file.Write([]byte(result))
	} else {
		os.Stdout.Write([]byte(result))
	}
}

func inputType(input string) (InputType, error) {
	matched, err := regexp.MatchString(`https?://*`, input)
	if err != nil {
		return FILE, err
	}
	if matched {
		return URL, nil
	}
	return FILE, nil
}

func openLocal(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func openURL(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
