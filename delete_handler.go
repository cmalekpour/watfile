package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	/* Security checks */
	request_id_t := strings.TrimSpace(r.FormValue("id"))
	if len(request_id_t) == 0 {
		http.Redirect(w, r, CONF_DOMAIN, 303)
		return
	}

	request_commands := strings.Split(request_id_t, "/")
	request_id := strings.Split(request_commands[0], ".")[0]
	if len(request_id) == 0 {
		http.Redirect(w, r, CONF_DOMAIN, 303)
		return
	}

	exists, _ := Exists(UPLOAD_DIR + request_id + "/")
	if exists == false {
		http.Redirect(w, r, CONF_DOMAIN, 303)
		return
	}

	files_t, _ := ioutil.ReadDir(UPLOAD_DIR + request_id + "/")

	filename := ""
	for a := range files_t {
		if files_t[a].Name() != "." && files_t[a].Name() != ".." {
			filename = files_t[a].Name()
			break
		}
	}

	if len(filename) == 0 {
		http.Redirect(w, r, CONF_DOMAIN, 303)
		return
	}

	filedat_t, _ := ioutil.ReadFile(UPLOAD_DIR + request_id + "/" + filename)
	md5_t := md5.New()
	md5_t.Write(filedat_t)
	md5_s := fmt.Sprintf("%x", md5_t.Sum(nil))

	delete_id := request_commands[1]
	exists, _ = Exists(DELETE_DIR + delete_id + "/" + request_id)
	if exists {
		os.RemoveAll(DELETE_DIR + delete_id)
		os.RemoveAll(FORCEDL_DIR + request_id)
		os.RemoveAll(HASH_DIR + md5_s)
		os.RemoveAll(UPLOAD_DIR + request_id)
	}
	http.Redirect(w, r, CONF_DOMAIN, 303)
	log.Printf("[LOG] File %s deleted\n", request_id)
	return
}
