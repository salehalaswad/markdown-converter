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

func RemoveMark(str string, mark string) (result string, found bool, color string) {

	s := strings.Index(str, mark)
	color = "black"
	if s == -1 {
		return str, false, color
	}
	switch mark {
	case "#":
		result = str[2:]

	case "##":
		result = str[3:]

	case "###":
		result = str[4:]

	case "####":
		result = str[5:]

	default:
		new := str[s+len(mark):]

		e := strings.Index(new, mark)

		if e == -1 {
			return str, false, color
		}
		result = new[:e]
	}

	if result[0:4] == "clr=" {
		i := strings.Index(result, "|")
		color = result[4:i]
		result = result[i+1:]

	}

	return result, true, color
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
		currentLine := scn.Text()

		headingL4, isHeadingL4, color := RemoveMark(currentLine, "####")
		if isHeadingL4 {
			headingL4 = "<h4 style=\"color:" + color + "\">" + headingL4 + "</h4>"
		}
		headingL3, isHeadingL3, color := RemoveMark(headingL4, "###")
		if isHeadingL3 {
			headingL3 = "<h3 style=\"color:" + color + "\">" + headingL3 + "</h3>"
		}

		headingL2, isHeadingL2, color := RemoveMark(headingL3, "##")
		if isHeadingL2 {
			headingL2 = "<h2 style=\"color:" + color + "\">" + headingL2 + "</h2>"
		}
		headingL1, isHeadingL1, color := RemoveMark(headingL2, "#")
		if isHeadingL1 {
			headingL1 = "<h1 style=\"color:" + color + "\">" + headingL1 + "</h1>"
		}
		bold, isBold, color := RemoveMark(headingL1, "**")
		if isBold {
			bold = "<strong style=\"color:" + color + "\">" + bold + "</strong>"
		}
		italic, isItalic, color := RemoveMark(bold, "*")
		if isItalic {
			italic = "<em style=\"color:" + color + "\">" + italic + "</em>"
		}

		htmlValues = append(htmlValues, italic)

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
