package freeway

import (
  "fmt"
  "regexp"
  "errors"
  "net/http"
)

const (
  locationUrl = "http://1968.freeway.gov.tw/traffic/index/fid/%s"
)

var (
  locLinkRe = regexp.MustCompile(`fid/(\d+)`)
  locNameRe = regexp.MustCompile(`^(.+)\([0-9.]+\) - (.+)\([0-9.]+\)$`)
  locIdRe = regexp.MustCompile(`loc/(\d+)$`)
)

func (f *Freeway) fetchLocation() {
  if f.error != nil {
    return
  }
  
  urls := make([]string, 0, len(f.Secs))
  for _, sec := range f.Secs {
    urls = append(urls, fmt.Sprintf(locationUrl, sec.Id))
  }
  
  ch := htmlParseChannel(urls, parseLocation)
  
  for {
    select {
    case obj, ok := <-ch:
      if !ok {
        return
      }
      
      locs, ok := obj.([]*Location)
      if !ok {
        continue
      }
      
      for _, loc := range locs {
        for _, name := range loc.insName {
          for _, in := range f.Ins {
            if in.Name != name || in.SecId != loc.SecId {
              continue
            }
            
            in.Locs = append(in.Locs, loc.Id)
            loc.Ins = append(loc.Ins, in.Id)
            break
          }
        }
        f.Locs = append(f.Locs, loc)
      }
    }
  }
}

func parseLocation(url string, resp *http.Response) (interface{}, error) {
  doc, err := parseHtml(resp)
  if err != nil {
    return nil, err
  }
  defer doc.Free()
  
  if !locLinkRe.MatchString(url) {
    return nil, errors.New("Link don't contain section id")
  }
  secId := locLinkRe.FindStringSubmatch(url)[1]
  
  rows, err := doc.Search("//tbody[@id='secs_body']/tr")
  if err != nil {
    return nil, errors.New("Html don't have any tr in tbody")
  }
  
  locs := make([]*Location, 0, len(rows))
  for i, row := range rows {
    tds, err := row.Search("td[@class='sec_name']")
    if err != nil || len(tds) == 0 {
      return nil, errors.New(fmt.Sprintf("%d'th td don't have sec_name class", i))
    }
    
    insName, err := parseInterchangeName(tds[0].Content())
    if err != nil {
      return nil, err
    }
    
    tds, err = row.Search("td[@class='sec_detail']/a/@href")
    if err != nil || len(tds) == 0 {
      return nil, errors.New(fmt.Sprintf("%d'th td don't have sec_detail href", i))
    }
    
    id, err := parseLocationId(tds[0].Content())
    if err != nil {
      return nil, err
    }
    
    locs = append(locs, &Location {
      id,
      secId,
      nil,
      nil,
      insName,
    })
  }
  
  return locs, nil
}

func parseInterchangeName(content string) ([]string, error) {
  if !locNameRe.MatchString(content) {
    return nil, errors.New("Location content " + content + " don't match name regexp format")
  }
  
  return locNameRe.FindStringSubmatch(content)[1:], nil
}

func parseLocationId(link string) (string, error) {
  if !locIdRe.MatchString(link) {
    return "", errors.New("Location link " + link + " don't match link regexp format")
  }
  
  return locIdRe.FindStringSubmatch(link)[1], nil
}
