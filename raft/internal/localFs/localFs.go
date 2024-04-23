package localfs

import "raft/pkg/rpcEncoding/out/protobuf"

type LocalFs interface{
    ApplyLogEntry(log *protobuf.LogEntry) error
}

func NewFs(rootDir string) LocalFs{
    return &fs{
        rootDir: rootDir,
    }
}
