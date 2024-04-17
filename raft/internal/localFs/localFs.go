package localfs

type LocalFs interface{
    Create(file string) error
    Read(file string) ([]byte,error)
    Update(file string, data []byte) error
    Delete(file string) error
    Rename(file string, newName string) error
}

func NewFs(rootDir string) LocalFs{
    return &fs{
        rootDir: rootDir,
    }
}
