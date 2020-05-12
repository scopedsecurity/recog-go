package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	recog "github.com/hdm/recog-go"
)

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}

func fingerprint(fingerprints []recog.FingerprintDB, text string) {
	for _, term := range strings.Fields(text) {
		for _, fdb := range fingerprints {
			matches := fdb.MatchAll(term)
			for _, match := range matches {
				if len(match.Errors) > 0 {
					fmt.Printf("%q", match.Errors)
				}
				if match.Matched {
					j, _ := json.Marshal(match.Values)
					fmt.Printf("%s\n", j)
				}
			}
		}
	}
}

func main() {
	var files []string
	if len(os.Args) < 2 {
		log.Fatalf("missing: recog xml directory")
	}

	err := filepath.Walk(os.Args[1], visit(&files))
	if err != nil {
		log.Fatal(err)
	}

	var fingerprints []recog.FingerprintDB
	for _, file := range files {
		fdb, err := recog.LoadFingerprintDBFromFile(file)
		if err != nil {
			log.Fatalf("error loading fingerprints from %s: %s", file, err)
		}
		fingerprints = append(fingerprints, fdb)
	}

	var text string

	text = strings.Join(os.Args[2:], " ")
	if len(text) < 1 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text = scanner.Text()
			fingerprint(fingerprints, text)
		}
	} else {
		fingerprint(fingerprints, text)
	}
}
