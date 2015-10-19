package templating

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var templates *template.Template

func InitTemplates() {
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"showDate":     func(date time.Time) string { return date.Format("Jan 2, 2006") },
		"showDateTime": func(date time.Time) string { return date.Format(time.RFC850) },
		"showISODate":  func(date time.Time) string { return date.Format("2006-01-02") },
		"minus":        func(a, b int) int { return a - b },
		"add":          func(a, b int) int { return a + b },
		"xpToLevel":    func(xp int) int { return xp/100 + 1 },
		"fixNewLines": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			s = regexp.MustCompile("\r?\n").ReplaceAllString(s, "<br>")
			return template.HTML(s)
		},
		"boldItalics": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			imageTags := regexp.MustCompile(`&lt;img\s+src=&#34;(.*?)&#34;&gt;`)
			s = imageTags.ReplaceAllString(s, `<img src="$1" style="max-width:570px;">`)
			unescapeTags := regexp.MustCompile("&lt;(/?(b|i|pre|u|sub|sup|strike|marquee))&gt;")
			s = unescapeTags.ReplaceAllString(s, "<$1>")
			s = regexp.MustCompile("\r?\n").ReplaceAllString(s, "<br>")
			return template.HTML(s)
		},
	}).ParseGlob("./templates/*.tmpl.html"))
}

func RenderPage(w http.ResponseWriter, template string, data interface{}) {
	err := templates.ExecuteTemplate(w, template+".tmpl.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrorPage(w http.ResponseWriter, statusCode int) {
	errorMessage := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	http.Error(w, errorMessage, statusCode)
}
