package wwwproxy

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"github.com/gorilla/mux"
)

type WWWProxy struct {
	redirects map[string]string
	router *mux.Router
}

func (p *WWWProxy) AddRedirect(node, target string) {
	p.redirects[node] = target
}

func (p *WWWProxy) RemoveRedirect(node string) {
	if _, ok := p.redirects[node]; ok {
		delete(p.redirects, node)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func (p *WWWProxy) node_redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Sprintf("%s/%s",vars["group"],vars["node"])
}

func (p *WWWProxy) show_all(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Sprintf("%s/%s",vars["group"],vars["node"])
}

func MakeWWWProxy() *WWWProxy {
	group_tmpl_str := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Group}}</title>
	</head>
	<body>
                <ul>
		{{range .Nodes}}<li><a href="{{.Group}}/{{ . }}">{{ . }}</a></li>{{end}}
                </ul>
	</body>
</html>`
	all_tmpl_str := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>All nodes</title>
	</head>
	<body>
                <ul>
		{{range .Groups}}<li>{{ . }}: <ul>{{range .Nodes}}<a href="{{.Group}}/{{ . }}">{{.Group}}/{{ . }}</a></li>{{end}}</ul></li>{{end}}
                </ul>
	</body>
</html>`
	p := &WWWProxy{
		redirects:make(map[string]string),
		router:mux.NewRouter(),
		group_tmpl:template.New("group").Parse(group_tmpl_str),
		all_tmpl:template.New("group").Parse(all_tmpl_str)}
	
	p.router.HandleFunc("/{group}/{name}", p.node_redirect).Methods("GET")
	p.router.HandleFunc("/{group}", p.show_group).Methods("GET")
	p.router.HandleFunc("/", p.show_all).Methods("GET")
	
}


func (p *WWWProxy) Run() {
	http.ListenAndServe(":8080", p.router)
}
