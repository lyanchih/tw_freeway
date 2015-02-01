package freeway

import (
  "net/http"
  "net/textproto"
)

const (
  sectionUrl = "http://1968.freeway.gov.tw/common/getfrees?id=sec_selt&df=1"
)

func (f *Freeway) fetchSection() {
  if f.error != nil {
    return
  }
  
  ch := htmlParseChannel([]string{sectionUrl}, parseSection)
  
  for {
    select {
    case obj, ok := <-ch:
      if !ok {
        return
      }
      f.Secs, ok = obj.([]*Section)
    }
  }
}

func parseSection(url string, resp *http.Response) (interface{}, error) {
  doc, err := parseHtml(resp)
  if err != nil {
    return nil, err
  }
  
  secs, err := doc.Search("//select[@id='sec_selt']/option")
  if err != nil {
    return nil, err
  }
  
  arr := make([]*Section, 0, len(secs))
  for _, sec := range secs {
    arr = append(arr, &Section{
      sec.Attr("value"),
      textproto.TrimString(sec.Content()),
    })
  }
  return arr, nil
}
