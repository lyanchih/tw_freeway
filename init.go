package main

var (
  freewayJsonDir = getenv("FreewayJsonDir", "./")
  freewayJsonName = getenv("FreewayJsonName", "freeway.json")
  freewayJsonPath = freewayJsonDir + freewayJsonName
  freewayJsonLint = getenvBool("FreewayJsonLint", true)
)
