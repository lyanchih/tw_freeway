package freeway

import (
  "log"
  "fmt"
  "sync"
  "errors"
  "io/ioutil"
  "net/http"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/html"
)

const (
  maxGokogiriGocurrency = 16
)

type htmlParseFunc func(string, *http.Response) (interface{}, error)

func get2xxResponse(url string) (*http.Response, error) {
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  
  if resp.StatusCode > 299 || resp.StatusCode < 200 {
    return nil, errors.New(fmt.Sprintf("Url %s get bad status code %d.", err.Error(), resp.StatusCode))
  }
  
  return resp, nil
}

func parseHtml(resp *http.Response) (*html.HtmlDocument, error) {
  defer resp.Body.Close()
  
  bs, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  
  return gokogiri.ParseHtml(bs)
}

func htmlParseChannel(urls []string, f htmlParseFunc) (<-chan interface{}) {
  wg := &sync.WaitGroup{}
  urlCh := make(chan string, len(urls))
  ch := make(chan interface{}, maxGokogiriGocurrency)
  
  for i := 0; i < cap(urls); i++ {
    wg.Add(1)
    
    go func() {
      defer wg.Done()
      for {
        select {
        case url, ok := <-urlCh:
          if !ok {
            return
          }
          
          res, err := get2xxResponse(url)
          if err != nil {
            log.Println(url, err)
            continue
          }
          
          obj, err := f(url, res)
          if err != nil {
            log.Println(url, err)
          } else {
            ch <- obj
          }
        }
      }
    }()
  }
  
  go func() {
    for _, url := range urls {
      urlCh <- url
    }
    close(urlCh)
    
    wg.Wait()
    close(ch)
  }()
  
  return ch
}
