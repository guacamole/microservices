package files

import (
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"os"
	"path/filepath"
)

/*
Local is the implementation for storage interface which works with local disk in
our current machine
This can totally be replaced with a database or cloud service in futuer
 */
type Local struct {
	maxFileSize int //maximum number of bytes for teh file ,5MB in our case
	basePath string
}

// Storage defines the behavior for file operations
// Implementations may be of the time local disk, or cloud storage, etc
type Storage interface {
	Save(path string, file io.Reader) error
}

/*
NewLocal creates a new Local file system with given basePath
basePath is the base directory to save the files to
maxSize represents the maximum number of bytes that a file can be
returns Storage (local Storage)
 */
func NewLocal(basePath string,maxSize int) (Storage, error) {

	path,err := filepath.Abs(basePath)
	if err!= nil{
		return nil, err
	}

	return &Local{basePath: path},nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	//get the full path
	fp := l.fullPath(path)

	//make sure the directory at the mentioned path exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d,os.ModePerm)
	if err != nil{

		return xerrors.Errorf("Unable to create directory",err)
	}

	//if file exists, delete it

	_,err = os.Stat(fp)

	if err == nil{
		fmt.Println("came in Save inside fjnv func")
		err := os.Remove(fp)
		if err != nil{

			return xerrors.Errorf("Unable to remove teh file",err)

		}


	} else if !os.IsNotExist(err){

		//if anything other than a not exist error
		return xerrors.Errorf("Unable to get file info: %s\n",err)

	}

	//create a new file at the path

	f,err := os.Create(fp)
	if err != nil{
		xerrors.Errorf("Unable to create file",err)
	}
	defer f.Close()

	//save the content to teh file
	//make sure of not writing more than maxSize
	fmt.Println("starting to copy")
	_,err = io.Copy(f,contents)
	if err != nil {
		return xerrors.Errorf("unable to write to file",err)
	}
	fmt.Println("finished copying ")
	fmt.Println("came in Save func")

	return nil
}


//get the file at given path and return a reader
//
func (l *Local) Get(path string) (*os.File,error){
	//get full path
	fp:= l.fullPath(path)

	f,err := os.Open(fp)
	if err != nil{
		return nil,xerrors.Errorf("unable to open file",err)
	}

	return f,nil
}

//returns the absolute path
func (l *Local) fullPath(path string) string{
	//append the given path to base path
	return filepath.Join(l.basePath,path)
}