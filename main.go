package main

import (
	"fmt"
	"os"
	"net/http"
	"flag"
	"log"
	"strconv"
	"path/filepath"
	"html/template"
	"sort"
)

func filesize(bytes int64) string {
	switch {
	case bytes < 1000:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1000:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1000:
		return fmt.Sprintf("%.2fM", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.2fG", float64(bytes)/(1024*1024*1024))
}

var template_funcs = template.FuncMap{
	"filesize": filesize,
}

func get_files_and_dirs(path string) DirectoryListing {
	var files, dirs FileInfoSlice
	d, _ := os.Open(path)
	defer d.Close()
	fi, _ := d.Readdir(-1)
	for _, fi := range fi {
		if fi.Mode()&os.ModeType == 0 { // is regular file
			files = append(files, fi)
		} else if fi.IsDir() {
			dirs = append(dirs, fi)
		}
	}
	files.Sort()
	dirs.Sort()
	return DirectoryListing{Dirs: dirs, Files: files}
}

func serve_index_html_if_exists(w http.ResponseWriter, r *http.Request, path string) bool {
	if f, err := os.Open(filepath.Join(path, "index.html")); err == nil {
		defer f.Close()
		finfo, _ := f.Stat()
		if finfo.Mode()&os.ModeType == 0 { // is regular file
			http.ServeFile(w, r, path)
			return true
		}
	}
	return false
}

var working_dir string
var tpl *template.Template

type FileInfoSlice []os.FileInfo
func (p FileInfoSlice) Len() int           { return len(p) }
func (p FileInfoSlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p FileInfoSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p FileInfoSlice) Sort()              { sort.Sort(p) }

type DirectoryListing struct {
	Dirs FileInfoSlice
	Files FileInfoSlice
}

func handle_req(w http.ResponseWriter, req *http.Request) {
	log.Printf("- %s - \"%s %s\" - %s", req.RemoteAddr, req.Method, req.URL, req.UserAgent())
	path := filepath.Join(working_dir, filepath.Clean(req.URL.Path))

	if f, err := os.Open(path); err == nil {
		defer f.Close()
		finfo, _ := f.Stat()

		switch {
		case finfo.IsDir():
			if serve_index_html_if_exists(w, req, path) {
				return
			}
			w.Header().Add("Content-Type", "text/html; charset=utf-8")

			err := tpl.Execute(w, get_files_and_dirs(path))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case finfo.Mode()&os.ModeType == 0: // is regular file
			http.ServeFile(w, req, path)
		}
	} else {
		http.NotFound(w, req)
	}
}

func main() {
	cwd, _ := os.Getwd()

	default_port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		default_port = 80
	}

	port := flag.Int("p", default_port, "port, default - environment variable PORT or 80, if not set")
	ip := flag.String("i", "0.0.0.0", "interface address, default - 0.0.0.0")
	flag.StringVar(&working_dir, "d", cwd, "if not specified serves from current working dir")
	flag.Parse()

	if working_dir != cwd {
		working_dir = filepath.Join(cwd, working_dir)
	}

	addr := fmt.Sprintf("%s:%d", *ip, *port)

	fmt.Printf("Serving \"%s\" directory\n", working_dir)
	fmt.Printf("Listening on %s ...\n", addr)

	tpl = template.Must(template.New("tpl").Funcs(template_funcs).Parse(tpl_src))

	http.HandleFunc("/", handle_req)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

const tpl_src = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="content-type" content="text/html; charset=UTF-8" />
	<title>directory listing</title>
	<style>
html{font-family:Verdana,sans-serif;font-size:15px;color:#444;cursor:default;-webkit-touch-callout:none;-webkit-user-select:none;-khtml-user-select:none;-moz-user-select:none;-ms-user-select:none;-o-user-select:none;user-select:none}
table{border-collapse:collapse;width:100%;margin:0 auto}
a{color:inherit;text-decoration:none;display:block;padding-left:30px}
td{border-bottom:#ccc 1px solid;}
td.size{text-align:right;padding:0 10px}
td.dir,td.file,td.parent{height:24px;background-repeat:no-repeat;background-position:3px}
td.dir{background-image:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAABC0lEQVR42mNkAILZebqrgVQIAyrYxi/0LzCs4eovBjyAEWrAf11NOQYmJjCX4f9/BoaHT14zfPz07StQxR8MXf8Z/jMy/F+WMulKNuOJmb42l68+OMzLx4OpDmQSFgAS//b12z8dLQV7lt+///a6x1YzcPAIMZACfnx5x/Tw5KJelr9//pr9+/2J4dOLZyQZwMLOwwDSy/L7z1+GX1/fM5AKfv15zwDSy/L771+gSb8ZyAEgvSx//vxj+PfnF1kGgPSCvfD3L5ku+AN3wW8KXfCHYhdQEgaUxgLIGf8o8sLfv6fev39nxsfLS5LmT58/M4D0svz6/b/4xq27vcBcZ0aSCf8ZTv1nZCwGAHEKqc3QXHeWAAAAAElFTkSuQmCC")}
td.file{background-image:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAArUlEQVR42u2TIQ7EIBBFp2QSBAKBqMWsZG+xN9reoqeqrWNRYLC9QltgmSbrS1O5P5mQEN7n/xC6eZ5f+76POecnNIgxZhHx3U3TZLXWpu972jwF18tgWRaIMX5wXVejlIK6tgQAYrz3BlNKQNOqH4e1P5RS4IqIxW3bLhsQe1S4anBPhb/BTQbHeyI2wzSsgt5ae/wFSnJm6CwxxKIQYgghjM65R0sCznmQUg5fD6e7S/20U8EAAAAASUVORK5CYII=")}
td.parent{background-image:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAA/klEQVR42mNgoBX4//8/IxDLALE0iE2OAUIqKioFIAxik6qZdf369QaMjIznQBjEBomRYoAMNzf3HCDzNAiD2CAxYjXzent7BwGZZ4B4PxSfAYmB5AhpZvr48aMKMzPzViB3BxBvh+IdIDGQHEgNPgPEpaWla6Ca16PhHSA5kBpcmjkbGhpsgIEG0rwMGwbJgdSA1GIzQJaLi2sWkLkQiOfhwAtBakBqsTr/8+fPGV+/fk3/8OFDOlBoKjIGiYMwSA0+bzCB8MOHD7WB3F5k/PPnT02YPMGofPDggSaQakXGnz59Uic6IT179kwDSNUgY5IMAAUSEKehYVma5FoAuw/AHFxK738AAAAASUVORK5CYII=")}
tr{vertical-align:center;}
tr:hover{background-color:#f9f9f9}
	</style>
</head>
<body>
	<table>
		<tr><td class="parent"><a href="javascript:history.back()">..</a></td><td class="size"></td></tr>
		{{range .Dirs}}<tr><td class="dir" colspan="2"><a href="{{.Name}}/">{{.Name}}</a></td></tr>{{end}}
		{{range .Files}}<tr><td class="file"><a href="{{.Name}}">{{.Name}}</a></td><td class="size"><a href="{{.Name}}/">{{filesize .Size}}</a></td></tr>{{end}}
	</table>
</body>
</html>`
