package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", handleIndex)

	log.Println("Listening on 'http://127.0.0.1:8080'")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Failed to start server: ", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	entries, err := os.ReadDir("img")
	if err != nil {
		log.Fatalln("Could not read 'img' directory: ", err)
	}

	srcs := []string{}
	for _, entry := range entries {
		if !isImage(entry) {
			continue
		}

		multiplier := parseMultiplier(entry)

		for range multiplier {
			srcs = append(srcs, path.Join("img", entry.Name()))
		}
	}

	output := "<div class=\"container\">"
	for i, src := range srcs {
		if i > 1 && i%8 == 0 {
			output += "</div>"
			output += "<div class=\"container\">"
		}

		output += fmt.Sprintf("<img src=\"%s\" alt=\"\">", src)
	}
	output += "</div>"

	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatalln("Failed to parse 'index.html': ", err)
	}

	if err := t.Execute(w, template.HTML(output)); err != nil {
		log.Fatalln("Failed to write html: ", err)
	}
}

func isImage(entry os.DirEntry) bool {
	return !entry.IsDir() && slices.Contains([]string{".png", ".jpg", ".webp"}, path.Ext(entry.Name()))
}

func parseMultiplier(entry os.DirEntry) int {
	re := regexp.MustCompile(`X(\d+) - `)
	matches := re.FindStringSubmatch(entry.Name())
	if len(matches) <= 1 {
		return 1
	}

	multiplier, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		log.Fatalf("Unexpected symbol '%s' while parsing multiplier: %s\n", matches[1], err)
	}

	return int(multiplier)
}
