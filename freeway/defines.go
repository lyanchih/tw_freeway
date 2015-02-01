package freeway

type Freeway struct {
  Secs []*Section `json:"sections"`
  Ins []*Interchange `json:"interchanges"`
  Locs []*Location `json:"locations"`
  error
}

type Section struct {
  Id string `json:"id"`
  Name string `json:"name"`
}

type Interchange struct {
  Id string `json:"id"`
  Name string `json:"name"`
  SecId string `json:"section_id"`
  Locs []string `json:"locations"`
}

type Location struct {
  Id string `json:"id"`
  SecId string `json:"section_id"`
  Ins []string `json:"interchanges"`
  *Geolocation
  insName []string
}

type Geolocation struct {
  id string
  FromX float64 `json:"from_x"`
  FromY float64 `json:"from_y"`
  ToX float64 `json:"to_x"`
  ToY float64 `json:"to_y"`
}
