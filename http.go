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

	if !strings.HasPrefix(r.URL.Path, s.dbName) {
		panic("Server serving unexpected path: " + r.URL.Path)
	}

	keyList := strings.Split(r.URL.Path, "/")
	if len(keyList) != paramLayerCount {
		panic("wrong action in db. wrong path:["+r.URL.Path+"]")
	}

	dbMap := GetGroup(keyList[1])
	if dbMap == nil {
		http.Error(w, fmt.Sprintf("no such db:[%s]", keyList[1]), http.StatusNotFound)
		return
	}

	if val, ok := dbMap.Get(keyList[2]); ok {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte(val.(string)))
	} else {
		http.Error(w, fmt.Sprintf("not found this key:[%s]", keyList[2]), http.StatusNotFound)
	}

	return
}