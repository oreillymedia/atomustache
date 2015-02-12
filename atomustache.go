package atomustache

import(
  "io/ioutil"
  "strings"
  "log"
  "os"
  "fmt"
  "path/filepath"
)

// Struct used for rendering of mustache templates.
type Atomustache struct {
  root string
  ext string
  views map[string]*Template
  atomic map[string]*Template
  layouts map[string]*Template
}

// Create a new struct.
// root - location of templates directory
func New(root string) *Atomustache {
  r := Atomustache{
    root: root,
    ext: ".html",
    views: make(map[string]*Template),
    atomic: make(map[string]*Template),
    layouts: make(map[string]*Template),
  }
  r.loadLayouts()
  r.loadAtomic()
  r.loadViews()
  return &r
}

// Rendering
// ----------------------------------------------------

func (r *Atomustache) RenderView(view string, data ...interface{}) string {
  out := r.views[view].Render(data...)
  return out
}

func (r *Atomustache) RenderViewInLayout(view string, layout string, data ...interface{}) string {
  out := r.views[view].RenderInLayout(r.layouts[layout], data...)
  return out
}

// Load files into maps
// ----------------------------------------------------

func (r *Atomustache) readRelDir(path string) []os.FileInfo {
  items, err := ioutil.ReadDir(r.root + "/" + path)
  checkErr(err)
  return items
}

func (r *Atomustache) readRelFile(path string) string {
  buf, err := ioutil.ReadFile(r.root + "/" + path)
  checkErr(err)
  return string(buf)
}

func (r *Atomustache) loadLayouts() {
  files := r.readRelDir("layouts")
  for _,file := range files {
    if strings.HasSuffix(file.Name(), r.ext) {
      k := noExt(file.Name())
      v := r.readRelFile("layouts/" + file.Name())
      t, mErr := ParseString(v, nil)
      checkErr(mErr)
      r.layouts[k] = t
    }
  }
}

func (r *Atomustache) loadAtomic() {
  folders := r.readRelDir("atomic")
  for _,folder := range folders {
    r.folderToAtomic("atomic/" + folder.Name(), folder.Name())
  }
}

func (r *Atomustache) folderToAtomic(folder string, atomicType string) {
  items := r.readRelDir(folder)
  for _,item := range items {
    if item.IsDir() {
      r.folderToAtomic(folder + "/" + item.Name(), atomicType)
    } else if strings.HasSuffix(item.Name(), r.ext) {
      k := atomicType + "-" + noExt(item.Name())
      v := r.readRelFile(folder + "/" + item.Name())
      t, mErr := ParseString(string(v), r.atomic)
      checkErr(mErr)
      r.atomic[k] = t
    } 
  }
}

func (r *Atomustache) loadViews() {
  folders := r.readRelDir("views")
  for _,folder := range folders {
    if folder.IsDir() {
      files := r.readRelDir("views/" + folder.Name())
      for _,file := range files {
        if strings.HasSuffix(file.Name(), r.ext) {
          k := folder.Name() + "/" + noExt(file.Name())
          v := r.readRelFile("views/" + folder.Name() + "/" + file.Name())
          t, mErr := ParseString(string(v), r.atomic)
          checkErr(mErr)
          r.views[k] = t
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
  return filename[0:len(filename)-len(filepath.Ext(filename))]
}