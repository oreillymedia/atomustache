package atomustache

import (
	"errors"
	"io/ioutil"
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
func New(styleguideFolder string, layoutsFolder string, viewsFolder string, ext string) (*Atomustache, error) {
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

	err1 := r.loadLayouts()
	if err1 != nil {
		return nil, err1
	}

	err2 := r.loadStyleguide()
	if err2 != nil {
		return nil, err2
	}

	err3 := r.loadViews()
	if err3 != nil {
		return nil, err3
	}

	return &r, nil
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

func (r *Atomustache) loadLayouts() error {

	files, err := ioutil.ReadDir(r.LayoutsFolder)
	if err != nil {
		return errors.New("Error reading layouts folder: " + err.Error())
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), r.Ext) {

			k := noExt(file.Name())
			path := r.LayoutsFolder + "/" + file.Name()

			v, fErr := ioutil.ReadFile(path)
			if fErr != nil {
				return errors.New("Error reading file (" + path + "): " + fErr.Error())
			}

			t, mErr := ParseString(string(v), nil)
			if mErr != nil {
				return errors.New("Error parsing string for file (" + path + "): " + mErr.Error())
			}

			r.Layouts[k] = t
		}
	}

	return nil
}

func (r *Atomustache) loadStyleguide() error {

	folders, err := ioutil.ReadDir(r.StyleguideFolder)
	if err != nil {
		return errors.New("Error reading styleguide folder: " + err.Error())
	}

	for _, folder := range folders {
		r.folderToAtomic(r.StyleguideFolder+"/"+folder.Name(), folder.Name())
	}

	return nil
}

func (r *Atomustache) loadViews() error {

	folders, err := ioutil.ReadDir(r.ViewsFolder)
	if err != nil {
		return errors.New("Error reading views folder: " + err.Error())
	}

	for _, folder := range folders {

		if folder.IsDir() {

			path := r.ViewsFolder + "/" + folder.Name()
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return errors.New("Error reading views folder (" + path + "): " + err.Error())
			}

			for _, file := range files {
				if strings.HasSuffix(file.Name(), r.Ext) {

					k := folder.Name() + "/" + noExt(file.Name())
					path := r.ViewsFolder + "/" + folder.Name() + "/" + file.Name()

					v, fErr := ioutil.ReadFile(path)
					if fErr != nil {
						return errors.New("Error reading view file (" + path + "): " + fErr.Error())
					}

					t, mErr := ParseString(string(v), r.Styleguide)
					if mErr != nil {
						return errors.New("Error parsing string from file (" + path + "): " + mErr.Error())
					}

					r.Views[k] = t
				}
			}
		}
	}

	return nil
}

// Helpers
// ----------------------------------------------------

func (r *Atomustache) folderToAtomic(folder string, atomicType string) error {

	items, err := ioutil.ReadDir(folder)
	if err != nil {
		return errors.New("Error reading atomic folder (" + folder + "): " + err.Error())
	}

	for _, item := range items {
		if item.IsDir() {
			r.folderToAtomic(folder+"/"+item.Name(), atomicType)
		} else if strings.HasSuffix(item.Name(), r.Ext) {
			k := atomicType + "-" + noExt(item.Name())
			path := folder + "/" + item.Name()
			v, fErr := ioutil.ReadFile(path)
			if fErr != nil {
				errors.New("Error reading atomic file (" + path + "): " + fErr.Error())
			}
			r.Styleguide[k] = string(v)
		}
	}

	return nil
}

func noExt(filename string) string {
	return filename[0 : len(filename)-len(filepath.Ext(filename))]
}
