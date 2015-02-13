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
  Root string
  Ext string
  Views map[string]*Template
  Atomic map[string]*Template
  Layouts map[string]*Template
}

// Create a new struct.
// root - location of templates directory
func New(root string) *Atomustache {
  r := Atomustache{
    Root: root,
    Ext: ".html",
    Views: make(map[string]*Template),
    Atomic: make(map[string]*Template),
    Layouts: make(map[string]*Template),
  }
  r.loadLayouts()
  r.loadAtomic()
  r.loadViews()
  return &r
}

// Rendering
// ----------------------------------------------------

func (r *Atomustache) RenderView(view string, data ...interface{}) string {
  out := r.Views[view].Render(data...)
  return out
}

func (r *Atomustache) RenderViewInLayout(view string, layout string, data ...interface{}) string {
  out := r.Views[view].RenderInLayout(r.Layouts[layout], data...)
  return out
}

// Load files into maps
// ----------------------------------------------------

func (r *Atomustache) readRelDir(path string) []os.FileInfo {
  items, err := ioutil.ReadDir(r.Root + "/" + path)
  checkErr(err)
  return items
}

func (r *Atomustache) readRelFile(path string) string {
  buf, err := ioutil.ReadFile(r.Root + "/" + path)
  checkErr(err)
  return string(buf)
}

func (r *Atomustache) loadLayouts() {
  files := r.readRelDir("layouts")
  for _,file := range files {
    if strings.HasSuffix(file.Name(), r.Ext) {
      k := noExt(file.Name())
      v := r.readRelFile("layouts/" + file.Name())
      t, mErr := ParseString(v, nil)
      checkErr(mErr)
      r.Layouts[k] = t
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
    } else if strings.HasSuffix(item.Name(), r.Ext) {
      k := atomicType + "-" + noExt(item.Name())
      v := r.readRelFile(folder + "/" + item.Name())
      t, mErr := ParseString(string(v), r.Atomic)
      checkErr(mErr)
      r.Atomic[k] = t
    } 
  }
}

func (r *Atomustache) loadViews() {
  folders := r.readRelDir("views")
  for _,folder := range folders {
    if folder.IsDir() {
      files := r.readRelDir("views/" + folder.Name())
      for _,file := range files {
        if strings.HasSuffix(file.Name(), r.Ext) {
          k := folder.Name() + "/" + noExt(file.Name())
          v := r.readRelFile("views/" + folder.Name() + "/" + file.Name())
          t, mErr := ParseString(string(v), r.Atomic)
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
  return filename[0:len(filename)-len(filepath.Ext(filename))]
}