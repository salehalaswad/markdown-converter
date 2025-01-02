package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func RemoveMark(str string, mark string) (string, bool) {

	s := strings.Index(str, mark)

	if s == -1 {
		return str, false
	}
	switch mark {
	case "#":
		return str[2:], true
	case "##":
		return str[3:], true
	case "###":
		return str[4:], true
	case "####":
		return str[5:], true
	}
	new := str[s+len(mark):]

	e := strings.Index(new, mark)

	if e == -1 {
		return str, false
	}
	result := new[:e]

	return result, true
}
func renderHTML(w http.ResponseWriter, htmlStrings []string) {
	tpl, _ := template.ParseFiles("index.html")

	var htmlValues []template.HTML
	for _, n := range htmlStrings {
		htmlEncapsulate := template.HTML(n)
		htmlValues = append(htmlValues, htmlEncapsulate)
	}

	tpl.Execute(w, htmlValues)

}
func handler(w http.ResponseWriter, r *http.Request) {
	md, _ := os.ReadFile("index.md")
	scn := bufio.NewScanner(strings.NewReader(string(md)))

	var htmlValues []string
	for scn.Scan() {
		currentLine := scn.Text()

		bold, isBold := RemoveMark(currentLine, "**")
		if isBold {
			bold = "<strong>" + bold + "</strong>"
		}
		italic, isItalic := RemoveMark(bold, "*")
		if isItalic {
			italic = "<em>" + italic + "</em>"
		}

		headingL4, isHeadingL4 := RemoveMark(italic, "####")
		if isHeadingL4 {
			headingL4 = "<h4>" + headingL4 + "</h4>"
		}
		headingL3, isHeadingL3 := RemoveMark(headingL4, "###")
		if isHeadingL3 {
			headingL3 = "<h3>" + headingL3 + "</h3>"
		}

		headingL2, isHeadingL2 := RemoveMark(headingL3, "##")
		if isHeadingL2 {
			headingL2 = "<h2>" + headingL2 + "</h2>"
		}
		headingL1, isHeadingL1 := RemoveMark(headingL2, "#")
		if isHeadingL1 {
			headingL1 = "<h1>" + headingL1 + "</h1>"
		}

		htmlValues = append(htmlValues, headingL1)
	}
	renderHTML(w, htmlValues)

}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}
