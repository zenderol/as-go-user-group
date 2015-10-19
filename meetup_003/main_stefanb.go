package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	shortenedUrls = map[string]string{}
	letterRunes   = []rune("abcdefghijklmnopqrstuvwxyz")
	hostPart      = "http://localhost:8080/short/"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getNextUniqueKey() string {
	n := 5
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.New("index").Parse(`
		<html>
			<body>
				<form action="submit">
					<label for="url">Url:</label>
					<input type="text" name="url">
					<button type="submit">Send</button>
				</form>
			</body>
		</html>
	`)

	tpl.Execute(w, nil)
}

func handleShorteningRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formData := r.Form
	key := getNextUniqueKey()

	shortenedUrls[key] = formData.Get("url")

	http.Redirect(w, r, "/urls", 302)
}

func listUrls(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.New("urls").Parse(`
		<html>
			<body>
				<ul>
					{{ range $key, $value := . }}
					   <li><a href="/short/{{ $key }}">click me to get to {{ $value }}</a><br/>
					       Or copy'n paste this link: http://localhost:8080/short/{{ $key }}
					   </li>
					{{ end }}
				</ul>
				<br/><br/>
				<div>
					<a href="/">back to home</a>
				</div>
			</body>
		</html>
	`)
	tpl.Execute(w, shortenedUrls)
}

func handleShortened(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["url"]
	targetUrl := shortenedUrls[key]

	http.Redirect(w, r, targetUrl, 302)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/submit", handleShorteningRequest)
	r.HandleFunc("/urls", listUrls)
	r.HandleFunc("/short/{url}", handleShortened)
	log.Fatal(http.ListenAndServe(":8080", r))
}
