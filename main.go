package main

import (
	"bufio"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	RawText string
	HTML    []template.HTML
}

const (
	H1 = iota
	H2
	H3
	H4
	H5
	H6
	Bold
	Italic
	Span
)

func converToHTML(before, target, after, color string, tag int) (result string) {

	switch tag {
	case Span:
		result = "<span style=\"color:" + color + " \" >" + target + "</span>"
	case Bold:
		result = "<strong style=\"color:" + color + " \" >" + target + "</strong>"
	case Italic:
		result = "<em style=\"color:" + color + " \" >" + target + "</em>"
	case H1:
		result = "<h1 style=\"color:" + color + " \" >" + target + "</h1>"
	case H2:
		result = "<h2 style=\"color:" + color + " \" >" + target + "</h2>"
	case H3:
		result = "<h3 style=\"color:" + color + " \" >" + target + "</h3>"
	case H4:
		result = "<h4 style=\"color:" + color + " \" >" + target + "</h4>"
	case H5:
		result = "<h5 style=\"color:" + color + " \" >" + target + "</h5>"
	case H6:
		result = "<h6 style=\"color:" + color + " \" >" + target + "</h6>"

	}
	return before + result + after
}
func span(rawText string) string {
	indexOfAttr := strings.Index(rawText, "[clr=")
	textAfterAttr := rawText[indexOfAttr:]
	textBeforeAttr := rawText[:indexOfAttr]
	indexOfPipe := strings.Index(textAfterAttr, "|")
	indexOfBracket := strings.Index(textAfterAttr, "]")
	clr := textAfterAttr[len("[clr="):indexOfPipe]
	targetText := textAfterAttr[indexOfPipe+1 : indexOfBracket]
	textAfterTarget := textAfterAttr[indexOfBracket+1:]
	rawText = converToHTML(textBeforeAttr, targetText, textAfterTarget, clr, Span)
	return rawText
}
func bold(rawText string) string {
	indexOfFirstDinkus := strings.Index(rawText, "**")
	textBeforeDinkus := rawText[:indexOfFirstDinkus]
	textAfterDinkus := rawText[indexOfFirstDinkus+2:]
	indexOfSecondDinkus := strings.Index(textAfterDinkus, "**")
	target := textAfterDinkus[:indexOfSecondDinkus]
	textAfterDinkus = textAfterDinkus[indexOfSecondDinkus+2:]
	clr := "black"
	if strings.HasPrefix(target, "clr=") {

		indexOfPipe := strings.Index(target, "|")
		clr = target[len("clr="):indexOfPipe]

		target = target[indexOfPipe+1:]
	}
	rawText = converToHTML(textBeforeDinkus, target, textAfterDinkus, clr, Bold)
	return rawText
}
func italic(rawText string) string {
	indexOfFirstDinkus := strings.Index(rawText, "*")
	textBeforeDinkus := rawText[:indexOfFirstDinkus]
	textAfterDinkus := rawText[indexOfFirstDinkus+1:]
	indexOfSecondDinkus := strings.Index(textAfterDinkus, "*")
	target := textAfterDinkus[:indexOfSecondDinkus]
	textAfterDinkus = textAfterDinkus[indexOfSecondDinkus+1:]
	clr := "black"
	if strings.HasPrefix(target, "clr=") {

		indexOfPipe := strings.Index(target, "|")
		clr = target[len("clr="):indexOfPipe]

		target = target[indexOfPipe+1:]
	}
	rawText = converToHTML(textBeforeDinkus, target, textAfterDinkus, clr, Italic)
	return rawText
}
func cutText(rawText string) string {

	for strings.Contains(rawText, "[clr=") {
		rawText = span(rawText)
	}
	if strings.HasPrefix(rawText, "#") {
		headingLevel := strings.Count(rawText, "#")
		clr := "black"
		rawText = rawText[headingLevel+1:]
		if strings.HasPrefix(rawText, "clr=") {

			indexOfPipe := strings.Index(rawText, "|")
			clr = rawText[len("clr="):indexOfPipe]

			rawText = rawText[indexOfPipe+1:]
		}
		switch headingLevel {
		case 1:
			rawText = converToHTML("", rawText, "", clr, H1)
		case 2:
			rawText = converToHTML("", rawText, "", clr, H2)
		case 3:
			rawText = converToHTML("", rawText, "", clr, H3)
		case 4:
			rawText = converToHTML("", rawText, "", clr, H4)
		case 5:
			rawText = converToHTML("", rawText, "", clr, H5)
		case 6:
			rawText = converToHTML("", rawText, "", clr, H6)
		}
	}
	for strings.Contains(rawText, "**") {
		rawText = bold(rawText)
	}
	for strings.Contains(rawText, "*") {
		rawText = italic(rawText)
	}

	return rawText
}

func renderHTML(w http.ResponseWriter, rawText string, htmlStrings []string) {
	tpl, _ := template.ParseFiles("edit.html")

	var htmlValues []template.HTML
	for _, n := range htmlStrings {
		htmlEncapsulate := template.HTML(n)
		htmlValues = append(htmlValues, htmlEncapsulate)
	}
	data := &Data{RawText: rawText, HTML: htmlValues}

	tpl.Execute(w, data)

}

func resultHandler(w http.ResponseWriter, r *http.Request) {

	md, _ := os.ReadFile("index.md")
	scn := bufio.NewScanner(strings.NewReader(string(md)))

	var htmlValues []string

	for scn.Scan() {

		currentLine := strings.ReplaceAll(scn.Text(), "\\n", "<br/>")

		cutted := cutText(currentLine)
		cutted += "<br/>"
		htmlValues = append(htmlValues, cutted)
	}

	renderHTML(w, string(md), htmlValues)

}

type Input struct {
	Data string
}

func writeHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var input Input
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	os.WriteFile("index.md", []byte(input.Data), 0600)

	w.WriteHeader(http.StatusOK)

}

func main() {
	http.HandleFunc("/", resultHandler)

	http.HandleFunc("/create", writeHandler)

	http.ListenAndServe(":8080", nil)
}
