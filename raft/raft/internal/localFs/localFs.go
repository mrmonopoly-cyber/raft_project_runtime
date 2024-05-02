package localfs

import "raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

type LocalFs interface{
    ApplyLogEntry(log *protobuf.LogEntry) error
}

func NewFs(rootDir string) LocalFs{
    return &fs{
        rootDir: rootDir,
    }
}
