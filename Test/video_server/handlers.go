package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func CreateUser(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	io.WriteString(w,"create user handle")
}


func Login(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	uname:=p.ByName("user_name")
	io.WriteString(w,uname)
}