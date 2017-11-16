// Copyright 2015 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var port = flag.Int("port", 8000, "port number")

func init() {
	flag.Parse()
}

var rootPath = ""

func init() {
	_, path, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(path), "..")
}

var jsDir = ""

func init() {
	var err error
	jsDir, err = ioutil.TempDir("", "ebiten")
	if err != nil {
		panic(err)
	}
}

func createJSIfNeeded(name string) (string, error) {
	out := filepath.Join(jsDir, name, "main.js")
	stat, err := os.Stat(out)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if (err != nil && os.IsNotExist(err)) || time.Now().Sub(stat.ModTime()) > 5*time.Second {
		target := "github.com/hajimehoshi/ebiten/examples/" + name
		out, err := exec.Command("gopherjs", "build", "--tags", "example", "-o", out, target).CombinedOutput()
		if err != nil {
			log.Print(string(out))
			return "", errors.New(string(out))
		}
	}
	return out, nil
}

func serveFile(w http.ResponseWriter, path, mime string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Type", mime)
	if _, err := io.Copy(w, f); err != nil {
		return err
	}
	return nil
}

func serveFileHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, filepath.Join(rootPath, r.URL.Path[1:]))
		return
	}
	if r.URL.RawQuery != "" {
		if err := serveFile(w, filepath.Join(rootPath, "index.html"), "text/html"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	apps := []string{}
	fs, err := ioutil.ReadDir(rootPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, f := range fs {
		if !f.IsDir() {
			continue
		}
		n := f.Name()
		if n[0] == '_' {
			continue
		}
		if n == "common" {
			continue
		}
		apps = append(apps, n)
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<ul>")
	for _, n := range apps {
		fmt.Fprintf(w, `<li><a href="/?%[1]s">%[1]s</a></li>`, template.HTMLEscapeString(n))
	}
	fmt.Fprintf(w, "</ul>")
}

func serveMainJS(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.Referer())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := u.RawQuery
	if name == "" {
		http.NotFound(w, r)
		return
	}
	out, err := createJSIfNeeded(name)
	if err != nil {
		t := template.JSEscapeString(template.HTMLEscapeString(err.Error()))
		js := `
window.onload = function() {
  document.body.innerHTML="<pre style='white-space: pre-wrap;'><code>` + t + `</code></pre>";
}`
		w.Header().Set("Content-Type", "text/javascript")
		fmt.Fprintf(w, js)
		return
	}
	if err := serveFile(w, out, "text/javascript"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveMainJSMap(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/main.js", serveMainJS)
	http.HandleFunc("/main.js.map", serveMainJSMap)
	http.HandleFunc("/", serveFileHandle)
	fmt.Printf("http://localhost:%d/\n", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
