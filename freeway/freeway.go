package freeway

import (
  "io"
  "io/ioutil"
  "os"
  "errors"
  "bytes"
  "encoding/json"
)

func (f *Freeway) HasError() bool {
  return f.error != nil
}

func (f *Freeway) Fetch() error {
  if f == nil {
    return errors.New("Freeway struct should be allocated")
  }
  f.error = nil
  
  f.fetchSection()
  f.fetchInterchange()
  f.fetchLocation()
  f.fetchGeolocation()
  return f.error
}

func (f *Freeway) Save(path string, lint bool) error {
  if f == nil {
    return errors.New("Freeway struct should be allocated")
  }
  
  var file *os.File
  var bs []byte
  for action := 0; f.error == nil; action++ {
    switch action {
    case 0:
      file, f.error = os.OpenFile(path, os.O_WRONLY | os.O_CREATE, 0644)
    case 1:
      defer file.Close()
    case 2:
      if lint {
        bs, f.error = json.Marshal(f)
      } else {
        bs, f.error = json.MarshalIndent(f, "", "  ")
      }
    case 3:
      io.Copy(file, bytes.NewBuffer(bs))
    default:
      return f.error
    }
  }
  return f.error
}

func Load(path string) (f *Freeway, err error) {
  var file *os.File
  var bs []byte
  for action := 0; err == nil; action++ {
    switch action {
    case 0:
      file, err = os.OpenFile(path, os.O_RDONLY, 0644)
    case 1:
      defer file.Close()
    case 2:
      bs, err = ioutil.ReadAll(file)
    case 3:
      err = json.Unmarshal(bs, &f)
    default:
      return
    }
  }
  return
}
