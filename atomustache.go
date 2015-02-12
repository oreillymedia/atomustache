package atomustache

import(
  "io/ioutil"
  "strings"
  "log"
  "fmt"
  "path/filepath"
)

type Atomustache struct {
  root string
  ext string
  views map[string]*Template
  atomic map[string]*Template
  layouts map[string]*Template
}

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

func checkErr(err error) {
  if err != nil {
    fmt.Println("error")
    log.Fatal(err)
  }
}

func noExt(filename string) string {
  return filename[0:len(filename)-len(filepath.Ext(filename))]
}

func (r *Atomustache) loadLayouts() {
  layouts_root := r.root + "/layouts"
  files, _ := ioutil.ReadDir(layouts_root)
  for _,file := range files {
    if strings.HasSuffix(file.Name(), r.ext) {
      k := noExt(file.Name())
      v, _ := ioutil.ReadFile(layouts_root + "/" + file.Name())
      t, mErr := ParseString(string(v), nil)
      checkErr(mErr)
      r.layouts[k] = t
    }
  }
}

func (r *Atomustache) loadAtomic() {
  atomic_root := r.root + "/atomic"
  folders, _ := ioutil.ReadDir(atomic_root)
  for _,folder := range folders {
    r.folderToAtomic(atomic_root + "/" + folder.Name(), folder.Name())
  }
}

func (r *Atomustache) folderToAtomic(folder string, atomicType string) {
  items, _ := ioutil.ReadDir(folder)
  for _,item := range items {
    if item.IsDir() {
      r.folderToAtomic(folder + "/" + item.Name(), atomicType)
    } else if strings.HasSuffix(item.Name(), r.ext) {
      k := atomicType + "-" + noExt(item.Name())
      v, _ := ioutil.ReadFile(folder + "/" + item.Name())
      t, mErr := ParseString(string(v), r.atomic)
      checkErr(mErr)
      r.atomic[k] = t
    } 
  }
}

func (r *Atomustache) loadViews() {

  views_root := r.root + "/views"

  folders, _ := ioutil.ReadDir(views_root)
  for _,folder := range folders {
    if folder.IsDir() {
      files, _ := ioutil.ReadDir(views_root + "/" + folder.Name())
      for _,file := range files {
        if strings.HasSuffix(file.Name(), r.ext) {
          k := folder.Name() + "/" + noExt(file.Name())
          v, _ := ioutil.ReadFile(views_root + "/" + folder.Name() + "/" + file.Name())
          t, mErr := ParseString(string(v), r.atomic)
          checkErr(mErr)
          r.views[k] = t
        }
      }
    }
  }
  
}

func (r *Atomustache) RenderView(view string, data ...interface{}) string {
  out := r.views[view].Render(data...)
  return out
}

func (r *Atomustache) RenderViewInLayout(view string, layout string, data ...interface{}) string {
  out := r.views[view].RenderInLayout(r.layouts[layout], data...)
  return out
}