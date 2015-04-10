package atomustache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Struct used for rendering of mustache templates.
type Atomustache struct {
	StyleguideFolder string
	LayoutsFolder    string
	ViewsFolder      string
	Ext              string
	Views            map[string]*Template
	Styleguide       map[string]string
	Layouts          map[string]*Template
}

// Create a new struct.
// root - location of templates directory
func New(styleguideFolder string, layoutsFolder string, viewsFolder string, ext string) *Atomustache {
	r := Atomustache{
		StyleguideFolder: styleguideFolder,
		LayoutsFolder:    layoutsFolder,
		ViewsFolder:      viewsFolder,
		Ext:              ext,
		Views:            make(map[string]*Template),

		// the atomic templates are saved as strings
		// because we need all of the partials before
		// parsing into template to avoid "partial not found"
		Styleguide: make(map[string]string),
		Layouts:    make(map[string]*Template),
	}
	r.loadLayouts()
	r.loadStyleguide()
	r.loadViews()
	return &r
}

// Rendering
// ----------------------------------------------------

func (r *Atomustache) RenderView(view string, data ...interface{}) (string, error) {
	if r.Views[view] == nil {
		return "", errors.New("no such view: " + view)
	}
	out := r.Views[view].Render(data...)
	return out, nil
}

func (r *Atomustache) RenderViewInLayout(view string, layout string, data ...interface{}) (string, error) {
	if r.Views[view] == nil {
		return "", errors.New("no such view: " + view)
	}
	out := r.Views[view].RenderInLayout(r.Layouts[layout], data...)
	return out, nil
}

// Load files into maps
// ----------------------------------------------------

func (r *Atomustache) readDir(path string) []os.FileInfo {
	items, err := ioutil.ReadDir(path)
	checkErr(err)
	return items
}

func (r *Atomustache) readFile(path string) string {
	buf, err := ioutil.ReadFile(path)
	checkErr(err)
	return string(buf)
}

func (r *Atomustache) loadLayouts() {
	files := r.readDir(r.LayoutsFolder)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), r.Ext) {
			k := noExt(file.Name())
			v := r.readFile(r.LayoutsFolder + "/" + file.Name())
			t, mErr := ParseString(v, nil)
			checkErr(mErr)
			r.Layouts[k] = t
		}
	}
}

func (r *Atomustache) loadStyleguide() {
	folders := r.readDir(r.StyleguideFolder)
	for _, folder := range folders {
		r.folderToAtomic(r.StyleguideFolder+"/"+folder.Name(), folder.Name())
	}
}

func (r *Atomustache) folderToAtomic(folder string, atomicType string) {
	items := r.readDir(folder)
	for _, item := range items {
		if item.IsDir() {
			r.folderToAtomic(folder+"/"+item.Name(), atomicType)
		} else if strings.HasSuffix(item.Name(), r.Ext) {
			k := atomicType + "-" + noExt(item.Name())
			v := r.readFile(folder + "/" + item.Name())
			r.Styleguide[k] = v
		}
	}
}

func (r *Atomustache) loadViews() {
	folders := r.readDir(r.ViewsFolder)
	for _, folder := range folders {
		if folder.IsDir() {
			files := r.readDir(r.ViewsFolder + "/" + folder.Name())
			for _, file := range files {
				if strings.HasSuffix(file.Name(), r.Ext) {
					k := folder.Name() + "/" + noExt(file.Name())
					v := r.readFile(r.ViewsFolder + "/" + folder.Name() + "/" + file.Name())
					t, mErr := ParseString(string(v), r.Styleguide)
					checkErr(mErr)
					r.Views[k] = t
				}
			}
		}
	}

}

// Helpers
// ----------------------------------------------------

func checkErr(err error) {
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
}

func noExt(filename string) string {
	return filename[0 : len(filename)-len(filepath.Ext(filename))]
}
