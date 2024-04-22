package model

type Model interface {
  Show() (map[string]string, error)
}
