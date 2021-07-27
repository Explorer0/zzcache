package zzcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const dbName = "zzcache"
const paramLayerCount = 3

type Server struct {
	address string
	dbName string
}

func NewServer(addr string) *Server {
	return &Server{
		address: addr,
		dbName: dbName,
	}
}

func (s *Server) Log(content string) {
	log.Printf("Server:[%s]:%s", s.address, content)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Log(fmt.Sprintf("url:[%s]", r.URL.Path))
	keyList := strings.Split(r.URL.Path[1:], "/")
	if len(keyList) != paramLayerCount {
		http.Error(w, fmt.Sprintf("wrong action in db. wrong path:[%s]", r.URL.Path), http.StatusBadRequest)
		return
	}

	if s.dbName != keyList[0] {
		http.Error(w, fmt.Sprintf("no such db exist in server. wrong path:[%s]", r.URL.Path), http.StatusBadRequest)
		return
	}

	dbMap := GetGroup(keyList[1])
	if dbMap == nil {
		http.Error(w, fmt.Sprintf("no such db:[%s]", keyList[1]), http.StatusNotFound)
		return
	}

	if val, ok := dbMap.Get(keyList[2]); ok {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(val.([]byte))
	} else {
		http.Error(w, fmt.Sprintf("not found this key:[%s]", keyList[2]), http.StatusNotFound)
	}

	return
}