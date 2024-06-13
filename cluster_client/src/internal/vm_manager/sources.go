package vmmanager

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
  cont, err := os.ReadFile("src/internal/cluster_manager/sources/pool.xml")
  return string(cont), err
}

func (this *sources) getVolXml() (string, error) {
  cont, err := os.ReadFile("src/internal/cluster_manager/sources/vol.xml")
  var content string = strings.ReplaceAll(string(cont), "NUMERO", strconv.Itoa(this.volNum))
  this.volNum = this.volNum + 1
  return content, err
}

func (this *sources) getDomXml() (string, error) {
  cont, err := os.ReadFile("src/internal/cluster_manager/sources/dom.xml")
  var content string = strings.ReplaceAll(string(cont), "NUMERO", strconv.Itoa(this.domNum))
  this.domNum = this.domNum + 1
  return content, err
}
