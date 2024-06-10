package localfs

import "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

type LocalFs interface{
    ApplyLogEntry(log *protobuf.LogEntry) ([]byte,error)
    GetRootDir() string
}

func NewFs(rootDir string) LocalFs{
    var res = &fs{
        rootDir: rootDir + "/",
    }

    return res
}
