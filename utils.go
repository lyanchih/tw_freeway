package main

import (
  "os"
  "strconv" 
)

func getenv(key, defaultValue string) (ret string) {
  ret = os.Getenv(key)
  if len(ret) == 0 {
    ret = defaultValue
  }
  return
}

func getenvBool(key string, defaultValue bool) (ret bool) {
  s := os.Getenv(key)
  if len(s) == 0 {
    return defaultValue
  }
  
  ret, err := strconv.ParseBool(s)
  if err != nil {
    ret = defaultValue
  }
  return ret
}
