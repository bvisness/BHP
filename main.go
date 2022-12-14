package bhp

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/Masterminds/sprig"
)

func Run(dir string) {
	Must0(os.Chdir(dir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("\nERROR: %v\n\n", r)
				debug.PrintStack()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		filename := Must1(filepath.Rel("/", r.URL.Path))
		if _, err := os.Stat(filename); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				panic(err)
			}
		}

		t := template.New("root")
		t = t.Funcs(sprig.FuncMap())
		t = Must1(t.ParseGlob("*.html"))

		Must0(t.ExecuteTemplate(w, filename, nil))
	})

	log.Fatal(http.ListenAndServe(":8484", nil))
}

// Takes an (error) return and panics if there is an error.
// Helps avoid `if err != nil` in scripts. Use sparingly in real code.
func Must0(err error) {
	if err != nil {
		panic(err)
	}
}

// Takes a (something, error) return and panics if there is an error.
// Helps avoid `if err != nil` in scripts. Use sparingly in real code.
func Must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
