package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"io"
	"net/http"
	"path/filepath"
	"product_images/files"
	"strconv"
)

//Files is a handler for reading and writing files
type Files struct{
	log hclog.Logger
	store files.Storage
}

// NewFiles creates a new file handler
func NewFiles(l hclog.Logger,s files.Storage) *Files{
	return &Files{l,s}
}

//UploadREST implements the http handler interface
func (f *Files) UploadREST(rw http.ResponseWriter,r *http.Request){

	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("handle POST","ID",id,"Filename",fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters
	f.SaveFiles(id,fn,r.Body,rw)


}

//multipart
func (f *Files) UploadMultipart(rw http.ResponseWriter,r *http.Request){

	err := r.ParseMultipartForm(128 * 1024)
	if err!= nil{

		f.log.Error("Bad request error",err)
		http.Error(rw,"Expected multipart form data",http.StatusBadRequest)
	}

	id,iderr:= strconv.Atoi(r.FormValue("id"))
	if iderr != nil{
		f.log.Error("Bad Request error",err,"id",id)
		http.Error(rw,"invalid id value, should be int",http.StatusBadRequest)
		return
	}

	ff,mh,err := r.FormFile("file")
	if err != nil{
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}

f.SaveFiles(r.FormValue("id"),mh.Filename,ff,rw)

}

func (f *Files) SaveFiles(id , path string,r io.ReadCloser,rw http.ResponseWriter ){

	fp := filepath.Join(id,path)
	err := f.store.Save(fp,r)
	if err != nil {
		f.log.Error("Unable to save file", "error",err)
		http.Error(rw,"Unable to save file",http.StatusBadRequest)
	}
}
func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}