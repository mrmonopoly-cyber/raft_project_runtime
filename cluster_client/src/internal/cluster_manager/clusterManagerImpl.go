package cluster_manager

import (
	"log"
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strings"
	"time"

	"libvirt.org/go/libvirt"
)

const (
  red string = "\x1b[31m"
  white string = "\x1b[0m"
)

type clusterManagerImpl struct {
  config protobuf.ClusterConf
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
  var errConn error
  var privIPs, publIPs []string

  this.sources = *newSources()

  this.libvirtMetaData.conn, errConn = libvirt.NewConnect("qemu:///system")
  if errConn != nil {
    log.Panicf("%sError connecting to the hypervisor: %s \n%s", red, errConn, white)
  }

  this.cleanUp()
  
  this.createPool()

  for i := 0; i < 5; i++ {
    this.createVol()
    this.createDom()
  }

  time.Sleep(50 * time.Second) // IPs aren't immediately available since VMs take time to actually start

  publIPs, privIPs = this.getIPs()
  log.Println("public IPs", publIPs)
  log.Println("private IPs", privIPs)

}

/*
 * Delete the entire cluster
*/
func (this *clusterManagerImpl) Terminate() bool {

  for _, d := range this.doms {
    d.Destroy()
    d.Undefine()
  }

  for _, sv := range this.vols {
    name, _ := sv.GetName()
    if (!strings.Contains(name, "iso")) {
       sv.Delete(0)
    }
  }

  this.pool.Destroy()
  this.pool.Undefine()

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

/*
 * Perform a clean up of storages, pools and domains
*/
func (this *clusterManagerImpl) cleanUp() {
  var existingVols []libvirt.StorageVol
  var existingDom []libvirt.Domain
  var existingPools []libvirt.StoragePool
  var errVols, errDom, errPool error

  existingPools, errPool = this.conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_DIR)
  if errPool != nil {
    log.Println(errPool)
  }

  for _, p := range existingPools {

    existingVols, errVols = p.ListAllStorageVolumes(0)
    if errVols != nil {
      log.Println("err listing vols ", errVols)
    }

    if len(existingVols) > 0 {
      for _, sv := range existingVols {
        name, _ := sv.GetName()
        if (!strings.Contains(name, "iso")) {
          sv.Delete(0)
        }
      }
    }

    p.Destroy()
    p.Undefine()
  }

  existingDom, errDom = this.conn.ListAllDomains(0)
  if errDom != nil {
    log.Println(errDom)
  }

  if len(existingDom) > 0 {
    for _, d := range existingDom {
      d.Destroy()
      d.Undefine()
    }
  }

}

/*
 * Create a storage pool, necessary to create storages
*/
func (this *clusterManagerImpl) createPool() {
  var existingPools []libvirt.StoragePool
  var poolXml string
  var errPoolXml, errPool, err error

  existingPools, err = this.conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_DIR)
  if err != nil {
    log.Println(err)
  }

  if len(existingPools) < 1 {

      poolXml, errPoolXml = this.sources.getPoolXml()
      if errPoolXml != nil {
        log.Println(errPoolXml)
      }

      this.pool, errPool = this.conn.StoragePoolDefineXML(poolXml, 0) 
      if errPool != nil {
        log.Panicf("%sError creating a storagePool: %s \n%s", red, errPool, white)
      }
  } else {
      this.pool = &existingPools[0]
  }
  this.pool.SetAutostart(true)
  this.pool.Create(0)
}

/*
 * Create a storage volume 
*/
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

/* 
 * Create a domain (this is where a VM is actually created)
*/
func (this *clusterManagerImpl) createDom() {
  var domXml string
  var errDomXml, errDom error
  var dom *libvirt.Domain

  domXml, errDomXml = this.sources.getDomXml()
  if errDomXml != nil {
    log.Panicf("%s", errDomXml)
  }

  dom, errDom = this.conn.DomainCreateXML(domXml, 0)
  if errDom != nil {
    log.Panicf("%sError creating a domain: %s \n%s", red, errDom, white)
  }

  this.doms = append(this.doms, dom)

  dom.Create()

}

/*
 * Collect private and public ips from DomainInterfaces
*/
func (this *clusterManagerImpl) getIPs() ([]string, []string) {
  var domInterfaces [][]libvirt.DomainInterface
  var publIPs, privIPs []string

  for _, d := range this.doms  {
    var err error
    var domInt []libvirt.DomainInterface
    domInt, err = d.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
    if err != nil {
      log.Printf("Error retreiving IPs: %s \n", err)
    }
    domInterfaces = append(domInterfaces, domInt)
  }

  for _, v := range domInterfaces {
    publIPs = append(publIPs, v[0].Addrs[0].Addr)
    privIPs = append(privIPs, v[1].Addrs[0].Addr)
  }

  return publIPs, privIPs
}
