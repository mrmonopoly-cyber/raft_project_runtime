package localfs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
)


type fs struct{
    lock sync.RWMutex
    rootDir string
    files map[string]string
}

func (this *fs) Create(file string) error{
    (*this).lock.Lock()
    defer (*this).lock.Unlock()

    var filePath string = this.getFilePath(file)
    var fd *os.File
    var err error

    fd,err = os.Create(filePath)
    if err == nil {
        fd.Close()
        (*this).files[file] = filePath
    }

    return err 
}
func (this *fs) Read(file string) ([]byte,error){
    (*this).lock.RLock()
    defer (*this).lock.Unlock()

    var resultBuffer bytes.Buffer
    var content []byte = make([]byte, 1024)
    var err error
    var fd *os.File 

    fd,err= this.searchFile(file)
    if err != nil {
        return nil,err
    }

    for err != io.EOF{
        _,err = fd.Read(content)
        if err != nil && err != io.EOF {
            fd.Close()
            return nil, errors.New("error reading file: " + file)
        }
        resultBuffer.Write(content)
    }
    fd.Close()

    return resultBuffer.Bytes(),nil
}

func (this *fs) Update(file string, data []byte) error{
    (*this).lock.Lock()
    defer (*this).lock.Unlock()

    var fd *os.File
    var err error

    fd,err = this.searchFile(file)
    if err != nil{
        return err
    }
    
    _,err = fd.WriteAt(data,0)
    fd.Close()
    return err 
}
func (this *fs) Delete(file string) error{
    (*this).lock.Lock()
    defer (*this).lock.Unlock()

    var err = os.Remove(this.getFilePath(file))
    if err == nil {
        delete(this.files,file)
    }

    return err
}

func (this *fs) Rename(file string, newName string) error{
    (*this).lock.Lock()
    defer (*this).lock.Unlock()

    var err = os.Rename(file,newName)
    if err == nil{
        delete(this.files,file)
        this.files[newName] = this.getFilePath(file)
    }

    return err
}

//utility

func (this *fs) getFilePath(file string) string{
    return this.rootDir + file
}

func (this *fs) searchFile(file string) (*os.File,error){
    var found bool
    var filePath string
    var err error
    var fd *os.File

    filePath,found = (*this).files[file]
    if !found{
        return nil,errors.New("file not found")
    }

    fd, err = os.Open(filePath)
    if err != nil{
        return nil,err
    }

    return fd,err
}
