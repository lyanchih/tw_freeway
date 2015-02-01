package freeway

import (
  "fmt"
  "errors"
  "regexp"
  "net/http"
  "net/textproto"
)

const (
  interchangeFromUrl = "http://1968.freeway.gov.tw/common/getnodsecs/fid/%s?id=from_selt"
  interchangeToUrl = "http://1968.freeway.gov.tw/common/getnodsecs/fid/%s?lc=1&id=end_selt"
)

var inLinkRe = regexp.MustCompile(`fid/(\d+)`)

func (f *Freeway) fetchInterchange() {
  if f.error != nil {
    return
  }
  
  urls := make([]string, 0, len(f.Secs)<<1)
  for _, sec := range f.Secs {
    urls = append(urls, fmt.Sprintf(interchangeFromUrl, sec.Id), fmt.Sprintf(interchangeToUrl, sec.Id))
  }
  ch := htmlParseChannel(urls, parseInterchange)
  
  for {
    select {
    case obj, ok := <-ch:
      if !ok {
        return
      }
      ins, ok := obj.([]*Interchange)
      if !ok {
        continue
      }
      
      for _, in := range ins {
        exist := false
        for i := 0; !exist && i < len(f.Ins); i++ {
          exist = f.Ins[i].Id == in.Id
        }
        
        if !exist {
          f.Ins = append(f.Ins, in)
        }
      }
    }
  }
}

func parseInterchange(url string, resp *http.Response) (interface{}, error){
  doc, err := parseHtml(resp)
  if err != nil {
    return nil, err
  }
  defer doc.Free()
  
  if !inLinkRe.MatchString(url) {
    return nil, errors.New("Link don't contain section id")
  }
  
  secId := inLinkRe.FindStringSubmatch(url)[1]
  ins, err := doc.Search("//select/option")
  if err != nil {
    return nil, err
  }
  
  arr := make([]*Interchange, 0, len(ins))
  for _, in := range ins {
    arr = append(arr, &Interchange{
      fmt.Sprintf("%s.%s", secId, in.Attr("value")),
      textproto.TrimString(in.Content()),
      secId,
      nil,
    })
  }
  return arr, nil
}
