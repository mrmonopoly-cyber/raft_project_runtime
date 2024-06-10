package cluster_manager

import (
	"log"
	"libvirt.org/go/libvirt"
)

const (
  red string = "\x1b[31m"
  white string = "\x1b[0m"
)

type clusterManagerImpl struct {
  nodes []string
  sources sources
  libvirtMetaData
}

type libvirtMetaData struct {
  conn *libvirt.Connect
  pool *libvirt.StoragePool
  vols []*libvirt.StorageVol
  doms []*libvirt.Domain
}

/* 
 * 1. Create a new cluster with 5 nodes
*/
func (this *clusterManagerImpl) Init() {
  this.sources = *newSources()
  this.libvirtMetaData.vols = make([]*libvirt.StorageVol, 10)
  this.libvirtMetaData.doms = make([]*libvirt.Domain, 10)

  var errConn error

  this.libvirtMetaData.conn, errConn = libvirt.NewConnect("qemu:///system")
  if errConn != nil {
    log.Panicf("%sError connecting to the hypervisor: %s \n%s", red, errConn, white)
  }

  this.createPool(this.conn)
  
  this.createVol()

  this.createDom(this.conn)

  for _, d := range this.doms {
    d.Create()
  }
}

/*
 * Delete the entire cluster
*/
func (this *clusterManagerImpl) Terminate() bool {
  for _, el := range this.doms {
    el.Free()
  }

  for _, el := range this.vols {
    el.Free()
  }

  this.pool.Free()
  this.conn.Close()

  return true
}

/*
* Get the cluster configurarion: 
*  ask the leader for the configurarion
*/
func (this *clusterManagerImpl) GetConfig() {

}

/*
* Change the current configuration
*/
func (this *clusterManagerImpl) SetConfig() {
  
}

func (this *clusterManagerImpl) createPool(conn *libvirt.Connect) {
  var poolXml string
  var errPoolXml, errPool error

  poolXml, errPoolXml = this.sources.getPoolXml()
  if errPoolXml != nil {
    log.Println(errPoolXml)
  }

  this.pool, errPool = conn.StoragePoolDefineXML(poolXml, 0) 
  if errPool != nil {
    log.Panicf("%sError creating a storagePool: %s \n%s", red, errPool, white)
  }
}

func (this *clusterManagerImpl) createVol() {
  var volXml string
  var vol *libvirt.StorageVol
  var errVolXml, errVol error

  volXml, errVolXml = this.sources.getVolXml()
  if errVolXml != nil {
    log.Panicf("%s", errVolXml)
  }

  vol, errVol = this.pool.StorageVolCreateXML(volXml, 0)
  if errVol != nil {
    log.Panicf("%sError creating a volume: %s \n%s", red, errVol, white)
  }

  this.vols = append(this.vols, vol)
}

func (this *clusterManagerImpl) createDom(conn *libvirt.Connect) {
  var domXml string
  var errDomXml, errDom error
  var dom *libvirt.Domain

  domXml, errDomXml = this.sources.getDomXml()
  if errDomXml != nil {
    log.Panicf("%s", errDomXml)
  }

  dom, errDom = conn.DomainCreateXML(domXml, 0)
  if errDom != nil {
    log.Panicf("%sError creating a domain: %s \n%s", red, errDom, white)
  }

  this.libvirtMetaData.doms = append(this.libvirtMetaData.doms, dom)
}

