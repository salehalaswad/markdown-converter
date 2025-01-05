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

func RemoveMark(str string, mark string) (pre string, result string, suf string, found bool, color string) {

	color = "black"
	pre = ""
	suf = ""
	s := strings.Index(str, mark)
	colored := strings.Contains(str, "[clr=")
	if colored {

		begin := strings.Index(str, "[clr=")
		pre = str[:begin]
		end := strings.Index(str, "|")
		clr := str[begin+len("[clr=") : end]
		endWord := strings.Index(str, "]")
		suf = str[endWord+len("]"):]
		word := str[end+len("|") : endWord]
		return pre, word, suf, true, clr
	}
	if s == -1 {
		return pre, str, suf, false, color
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
		pre = str[:s]

		new := str[s+len(mark):]
		e := strings.Index(new, mark)

		if e == -1 {
			return pre, str, suf, false, color
		}
		result = new[:e]
		suf = str[strings.Index(str, result)+len(result)+len(mark):]

	}

	if len(result) > 4 {

		if result[0:4] == "clr=" {
			i := strings.Index(result, "|")
			color = result[4:i]
			result = result[i+1:]

		}
	}
	return pre, result, suf, true, color
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
		pre, paragraph, suf, isColored, color := RemoveMark(currentLine, "")
		if isColored {

			paragraph = pre + "<span style=\"color:" + color + "\">" + paragraph + "</span>" + suf
		}

		_, headingL4, _, isHeadingL4, color := RemoveMark(paragraph, "####")
		if isHeadingL4 {
			headingL4 = "<h4 style=\"color:" + color + "\">" + headingL4 + "</h4>"
		}

		pre, headingL3, suf, isHeadingL3, color := RemoveMark(headingL4, "###")
		if isHeadingL3 {
			headingL3 = pre + "<h3 style=\"color:" + color + "\">" + headingL3 + "</h3>" + suf
		}

		pre, headingL2, suf, isHeadingL2, color := RemoveMark(headingL3, "##")
		if isHeadingL2 {
			headingL2 = pre + "<h2 style=\"color:" + color + "\">" + headingL2 + "</h2>" + suf
		}
		pre, headingL1, suf, isHeadingL1, color := RemoveMark(headingL2, "#")
		if isHeadingL1 {
			headingL1 = pre + "<h1 style=\"color:" + color + "\">" + headingL1 + "</h1>" + suf
		}

		pre, bold, suf, isBold, color := RemoveMark(headingL1, "**")
		if isBold {
			bold = pre + "<strong style=\"color:" + color + "\">" + bold + "</strong>" + suf

		}
		pre, italic, suf, isItalic, color := RemoveMark(bold, "*")
		if isItalic {

			italic = pre + "<em style=\"color:" + color + "\">" + italic + "</em>" + suf
		}
		htmlValues = append(htmlValues, italic+"<br/>")

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
