package freeway

import (
  "strconv"
)

type b2f struct {
  error
}

func (b b2f) convert(bs []byte) (f float64) {
  if b.error != nil {
    return
  }
  
  f, b.error = strconv.ParseFloat(string(bs), 32)
  return
}
