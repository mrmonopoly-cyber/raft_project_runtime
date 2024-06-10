package cluster_manager

import (
	"os"
	"strconv"
	"strings"
)

type sources struct {
  volNum int
  domNum int
}



func newSources() *sources {
  return &sources{volNum: 2, domNum: 2}
}

func (this *sources) getPoolXml() (string, error) {
  cont, err := os.ReadFile("sources/pool.xml")
  return string(cont), err
}

func (this *sources) getVolXml() (string, error) {
  cont, err := os.ReadFile("sources/vol.xml")
  var content string = strings.ReplaceAll(string(cont), "NUMERO", strconv.Itoa(this.volNum))
  this.volNum++
  return content, err
}

func (this *sources) getDomXml() (string, error) {
  cont, err := os.ReadFile("sources/dom.xml")
  var content string = strings.ReplaceAll(string(cont), "NUMERO", strconv.Itoa(this.domNum))
  this.domNum++
  return content, err
}
