/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

// Add k8s style header for api-reference html docs; And fix the alignment.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var (
	templateHTML = flag.String("templateHTML", "http://kubernetes.io/v1.0/index.html", "URL of the html file that has the k8s header")
)

//strip the <!DOCTYPE html> and <html> tag
func stripHTML(filename string) error {
	targetHTML, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer targetHTML.Close()
	target, err := html.Parse(targetHTML)
	if err != nil {
		return err
	}
	//node for <head>
	headNode := target.FirstChild.NextSibling.FirstChild

	err = targetHTML.Truncate(0)
	if err != nil {
		return err
	}
	_, err = targetHTML.Seek(0, 0)
	if err != nil {
		return err
	}
	for node := headNode; node != nil; node = node.NextSibling {
		err = html.Render(targetHTML, node)
		if err != nil {
			return err
		}
	}
	return nil
}

// Replace http://releases.k8s.io/HEAD/xxx.md
// or http://releases.k8s.io/v1.1.0/xxx.md
// http://kubernetes.io/v1.1/xxx.html
func fixLinks(fileName, version string) error {
	mdRE := regexp.MustCompile(`http://releases.k8s.io/[^/]+(.*?\.)md`)
	repl := []byte("http://kubernetes.io/" + version + "${1}" + "html")
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	file = mdRE.ReplaceAll(file, repl)

	if err = ioutil.WriteFile(fileName, file, 0666); err != nil {
		return err
	}
	return nil
}

func nitFix(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	//TODO: This really should be fixed when generating the original html files
	newContents := strings.Replace(string(file), "<h2 id=\"_paths\">Paths</h2>", "<h2 id=\"_paths\">Operations</h2>", 1)
	//Safari cannot parse *zoom and reloads the page forever
	newContents = strings.Replace(newContents, "*zoom: 1", "zoom: 1", -1)

	if err = ioutil.WriteFile(filename, []byte(newContents), 0666); err != nil {
		return err
	}
	return nil
}

// move the definitions.html and operations.html to _includes, with name like
// extensions-v1beta1-definitions.html.
func moveHTML(filePath string, version string, outputDir string) error {
	dirs := strings.Split(filePath, "/")
	var i int
	for i = 0; i < len(dirs); i++ {
		if dirs[i] == "api-reference" {
			break
		}
	}
	if i == len(dirs) {
		return fmt.Errorf("failed to find fold \"api-reference\"\n")
	}
	fileName := strings.Join(dirs[i+1:], "-")
	// version is needed to eliminate conflicts among releases.
	if err := os.MkdirAll(outputDir+"/../_includes/"+version, 0766); err != nil {
		return err
	}
	return os.Rename(filePath, outputDir+"/../_includes/"+version+"/"+fileName)
}

func processHTML(filePath string, version string, outputDir string) error {
	if err := stripHTML(filePath); err != nil {
		return err
	}
	if err := nitFix(filePath); err != nil {
		return err
	}
	if err := fixLinks(filePath, version); err != nil {
		return err
	}
	if err := moveHTML(filePath, version, outputDir); err != nil {
		return err
	}
	return nil
}
