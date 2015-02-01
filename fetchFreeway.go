package main

import (
  "log"
  "./freeway"
)

func main() {
  f := new(freeway.Freeway)
  f.Fetch()
  f.Save(freewayJsonPath, freewayJsonLint)
  if f.HasError() {
    log.Fatal(f.Error())
  }
}
