package confmetadata

type ConfMetadata interface{
    GetConfig() map[string]string
    GetNumNodesInConf() uint
}
