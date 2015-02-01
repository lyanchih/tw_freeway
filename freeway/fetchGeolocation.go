package freeway

import (
  "io"
  "fmt"
  "bufio"
  "errors"
  "regexp"
  "net/http"
)

const (
  geolocationUrl = "http://1968.freeway.gov.tw/section/getlocationinfo/?loc=%s"
)

var (
  geoLinkRe = regexp.MustCompile(`loc=(\d+)`)
  geoRe = regexp.MustCompile(`fromX=(\d+\.\d+);\s*fromY=(\d+\.\d+);\s*toX=(\d+\.\d+);\s*toY=(\d+\.\d+)`)
)
func (f *Freeway) fetchGeolocation() {
  if f.error != nil {
    return
  }
  
  urls := make([]string, 0, len(f.Locs))
  for _, loc := range f.Locs {
    urls = append(urls, fmt.Sprintf(geolocationUrl, loc.Id))
  }
  
  ch := htmlParseChannel(urls, parseGeolocation)
  
  for {
    select {
    case obj, ok := <- ch:
      if !ok {
        return
      }
      
      geo, ok := obj.(*Geolocation)
      if !ok {
        continue
      }
      
      for _, loc := range f.Locs {
        if loc.Id == geo.id {
          loc.Geolocation = geo
          break
        }
      }
    }
  }
}

func parseGeolocation(url string, resp *http.Response) (geo interface{}, err error) {
  defer resp.Body.Close()
  
  if !geoLinkRe.MatchString(url) {
    return nil, errors.New("Link can't contain location id")
  }
  locId := geoLinkRe.FindStringSubmatch(url)[1]
  
  b := bufio.NewReader(resp.Body)
  var line []byte
  for {
    line, _, err = b.ReadLine()
    if err != nil {
      if err == io.EOF {
        err = errors.New("Html can't find location's latitide and longitude")
      }
      break
    }
    
    if !geoRe.Match(line) {
      continue
    }
    
    res := geoRe.FindSubmatch(line)
    
    var b b2f
    fromX := b.convert(res[1])
    fromY := b.convert(res[2])
    toX := b.convert(res[3])
    toY := b.convert(res[4])
    if b.error != nil {
      return nil, errors.New(fmt.Sprintf("Geolocation %v parse float fail", res[1:]))
    }
    
    return &Geolocation{
      locId,
      fromX, fromY,
      toX, toY,
    }, nil
  }
  
  return nil, err
}
