package vmmanager

import (
	"log"
	c "raft/client/src/internal/cluster"
	"raft/client/src/internal/utility"
	"slices"
	"strings"
	"time"

	"libvirt.org/go/libvirt"
)

const (
  red string = "\x1b[31m"
  white string = "\x1b[0m"
)

type vMManagerImpl struct {
  sources sources
  libvirtMetaData
}

type libvirtMetaData struct {
  conn *libvirt.Connect
  pool *libvirt.StoragePool
  vols []*libvirt.StorageVol
  doms []utility.Pair[*libvirt.Domain, []libvirt.DomainInterface]
}

/* 
 * 1. Create a new cluster with 5 nodes
*/
func (this *vMManagerImpl) Init() c.Cluster {
  var errConn error
  var IPs []utility.Pair[string,string]

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

  IPs = this.getIPs()

  return c.NewCluster(IPs)
}

/*
 * Delete the entire cluster
*/
func (this *vMManagerImpl) Terminate() bool {

  for _, d := range this.doms {
    d.Fst.Destroy()
    d.Fst.Undefine()
  }
  this.doms = make([]utility.Pair[*libvirt.Domain, []libvirt.DomainInterface], 0)

  for _, sv := range this.vols {
    name, _ := sv.GetName()
    if (!strings.Contains(name, "iso")) {
       sv.Delete(0)
    }
  }
  this.vols = make([]*libvirt.StorageVol, 0)

  this.pool.Destroy()
  this.pool.Undefine()
  this.pool = nil

  this.conn.Close()
  this.conn = nil

  return true
}

/*
* Change the current configuration, adding a new node
*/
func (this *vMManagerImpl) AddNode() c.Cluster {
  var IPs []utility.Pair[string,string]

  this.createVol()
  this.createDom()
  time.Sleep(50*time.Second)
  
  IPs = this.getIPs()

  return c.NewCluster(IPs)
}

/* 
 * Change the current configuration, removing a new node
*/
func (this *vMManagerImpl) RemoveNode(IP string) {
  var domInterfaces [][]libvirt.DomainInterface = getDomInterfaces(this.doms)

  for _, v := range domInterfaces {
    
  }

}

/*
 * Perform a clean up of storages, pools and domains
*/
func (this *vMManagerImpl) cleanUp() {
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
func (this *vMManagerImpl) createPool() {
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
func (this *vMManagerImpl) createVol() {
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
func (this *vMManagerImpl) createDom() {
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
func (this *vMManagerImpl) getIPs() []utility.Pair[string,string] {
  var domInterfaces [][]libvirt.DomainInterface = getDomInterfaces(this.doms)
  var IPs []utility.Pair[string,string]

  for _, v := range domInterfaces {
    IPs = append(IPs, utility.Pair[string, string]{Fst: v[0].Addrs[0].Addr, Snd: v[1].Addrs[0].Addr})
  }

  return IPs
}

/*
 * Collect domain interfaces of all domain
*/
func getDomInterfaces(doms []*libvirt.Domain) [][]libvirt.DomainInterface {
  var domInterfaces [][]libvirt.DomainInterface

  for _, d := range doms  {
    var err error
    var domInt []libvirt.DomainInterface
    domInt, err = d.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
    if err != nil {
      log.Printf("Error retreiving IPs: %s \n", err)
    }
    domInterfaces = append(domInterfaces, domInt)
  }

  return domInterfaces
}
