package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	for _, filename := range os.Args[1:] {
		snippets := extractSnippets("#_BEGIN_WRAP_TAG_", "#_END_WRAP_TAG_", filename)
		for _, snippet := range snippets {
			ioutil.WriteFile(
				filepath.Join("website", "pages", "partials", filepath.Base(filename), snippet.Identifier+filepath.Ext(filename)),
				snippet.Text,
				os.ModePerm,
			)
		}
	}
}

type snippet struct {
	Identifier string
	Text       []byte
	Closed     bool
}

func extractSnippets(beginPattern, endPattern, filename string) []snippet {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	snippets := []snippet{}
	for scanner.Scan() {
		line := scanner.Text()
		if identifier := matches(line, beginPattern); identifier != "" {
			snippets = append(snippets, snippet{
				Identifier: identifier,
			})
			continue
		}
		if identifier := matches(line, endPattern); identifier != "" {
			for _, snippet := range snippets {
				if snippet.Identifier == identifier {
					snippet.Closed = true
				}
			}
			continue
		}
		for _, snippet := range snippets {
			if snippet.Closed {
				continue
			}
			snippet.Text = append(snippet.Text, scanner.Bytes()...)
		}
	}
	return snippets
}

func matches(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return ""
	}
	return s[len(prefix)+1:]
}
