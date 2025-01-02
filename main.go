package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func RemoveMark(str string, mark string) (result string, found bool, color string) {

	s := strings.Index(str, mark)
	color = "black"
	if s == -1 {
		return str, false, color
	}
	switch mark {
	case "#":
		result = str[2:]
		// return result, true
	case "##":
		result = str[3:]
		// return result, true
	case "###":
		result = str[4:]
		// return result, true
	case "####":
		result = str[5:]
		// return result, true
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

		bold, isBold, color := RemoveMark(currentLine, "**")
		if isBold {
			bold = "<strong style=\"color:" + color + "\">" + bold + "</strong>"
		}
		italic, isItalic, color := RemoveMark(bold, "*")
		if isItalic {
			italic = "<em style=\"color:" + color + "\">" + italic + "</em>"
		}

		headingL4, isHeadingL4, color := RemoveMark(italic, "####")
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

		htmlValues = append(htmlValues, headingL1)
	}
	renderHTML(w, htmlValues)

}

func main() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}
