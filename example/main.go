package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"text/template"
)

var indexTextTemplate = template.Must(template.New("Index").Parse(`# Properties

Runtime: {{.Runtime}}
TARGETPLATFORM: {{.TARGETPLATFORM}}
GOOS: {{.GOOS}}
GOARCH: {{.GOARCH}}
URL: {{.URL}}
User: {{.User}}
UserEmail: {{.UserEmail}}
{{- range $kv := .Properties}}
{{- range .Values }}
{{$kv.Name}}: {{.}}
{{- end}}
{{- end}}

# Request Headers
{{ range $kv := .RequestHeaders}}
{{- range .Values }}
{{$kv.Name}}: {{.}}
{{- end}}
{{- end}}
`))

var indexTemplate = template.Must(template.New("Index").Parse(`<!DOCTYPE html>
<html>
<head>
<title>oauth2-proxy service example</title>
<style>
body {
	font-family: monospace;
	color: #555;
	background: #e6edf4;
	padding: 1.25rem;
	margin: 0;
}
table {
	background: #fff;
	border: .0625rem solid #c4cdda;
	border-radius: 0 0 .25rem .25rem;
	border-spacing: 0;
	margin-bottom: 1.25rem;
	padding: .75rem 1.25rem;
	text-align: left;
	white-space: pre;
}
table > caption {
	background: #f1f6fb;
	text-align: left;
	font-weight: bold;
	padding: .75rem 1.25rem;
	border: .0625rem solid #c4cdda;
	border-radius: .25rem .25rem 0 0;
	border-bottom: 0;
}
table td, table th {
	padding: .25rem;
}
table > tbody > tr:hover {
	background: #f1f6fb;
}
</style>
</head>
<body>
	<table>
		<caption>Properties</caption>
		<tbody>
			<tr>
				<th>Runtime</th>
				<td>{{.Runtime}}</td>
			</tr>
			<tr>
				<th>GOOS</th>
				<td>{{.GOOS}}</td>
			</tr>
			<tr>
				<th>GOARCH</th>
				<td>{{.GOARCH}}</td>
			</tr>
			<tr>
				<th>URL</th>
				<td>{{.URL}}</td>
			</tr>
			<tr>
				<th>User</th>
				<td>{{.User}}</td>
			</tr>
			<tr>
				<th>User Email</th>
				<td>{{.UserEmail}}</td>
			</tr>
			{{- range $kv := .Properties }}
			{{- range .Values }}
			<tr>
				<th>{{$kv.Name}}</th>
				<td>{{.}}</td>
			</tr>
			{{- end}}
			{{- end}}
		</tbody>
	</table>
	<table>
		<caption>Request Headers</caption>
		<tbody>
			{{- range $kv := .RequestHeaders }}
			{{- range .Values }}
			<tr>
				<th>{{$kv.Name}}</th>
				<td>{{.}}</td>
			</tr>
			{{- end}}
			{{- end}}
		</tbody>
	</table>
</body>
</html>
`))

type keyValue struct {
	Name   string
	Values []string
}

type keyValues []keyValue

func (a keyValues) Len() int      { return len(a) }
func (a keyValues) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a keyValues) Less(i, j int) bool {
	return strings.ToLower(a[i].Name) < strings.ToLower(a[j].Name)
}

func getProperties() keyValues {
	result := make(keyValues, 0)
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		name := parts[0]
		if !strings.HasPrefix(name, "EXAMPLE_") {
			continue
		}
		value := parts[1]
		result = append(result, keyValue{
			Name:   name,
			Values: []string{value},
		})
	}
	sort.Sort(result)
	return result
}

func getRequestHeaders(httpHeaders http.Header) keyValues {
	result := make(keyValues, 0, len(httpHeaders))
	for k := range httpHeaders {
		result = append(result, keyValue{
			Name:   k,
			Values: httpHeaders[k],
		})
	}
	sort.Sort(result)
	return result
}

type indexData struct {
	Runtime        string
	GOOS           string
	GOARCH         string
	URL            string
	User           string
	UserEmail      string
	Properties     keyValues
	RequestHeaders keyValues
}

func main() {
	log.SetFlags(0)

	log.Printf("%s", runtime.Version())
	log.Printf("GOOS=%s", runtime.GOOS)
	log.Printf("GOARCH=%s", runtime.GOARCH)
	//log.Printf("GOARM=%s", runtime.GOARM) // NB there is no GOARM.

	var listenAddress = flag.String("listen", "", "Listen address")

	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		log.Fatalf("\nERROR You MUST NOT pass any positional arguments")
	}

	if *listenAddress == "" || *listenAddress == "no" {
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		properties := getProperties()
		requestHeaders := getRequestHeaders(r.Header)
		user := r.Header.Get("X-Forwarded-User")
		userEmail := r.Header.Get("X-Forwarded-Email")

		var scheme string

		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}

		url := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.String())

		var t *template.Template
		var contentType string

		switch r.URL.Query().Get("format") {
		case "text":
			t = indexTextTemplate
			contentType = "text/plain"
		default:
			t = indexTemplate
			contentType = "text/html"
		}

		w.Header().Set("Content-Type", contentType)

		err := t.ExecuteTemplate(w, "Index", indexData{
			Runtime: runtime.Version(),
			GOOS:    runtime.GOOS,
			GOARCH:  runtime.GOARCH,
			//GOARM:          runtime.GOARM, // NB there is no GOARM.
			URL:            url,
			User:           user,
			UserEmail:      userEmail,
			Properties:     properties,
			RequestHeaders: requestHeaders,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Printf("Listening at http://%s\n", *listenAddress)

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatalf("Failed to ListenAndServe: %v", err)
	}
}
