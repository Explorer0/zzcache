package zzcache

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const dbName = "zzcache"
const paramLayerCount = 3
const defaultReplicaCnt = 5

type Server struct {
	address string
	dbName  string
	sync.RWMutex
	nodeHashMap *DistributeMap
	peerMap     map[string]Peer
}

func NewServer(addr string) *Server {
	return &Server{
		address:     addr,
		dbName:      dbName,
		nodeHashMap: NewDistributeMap(defaultReplicaCnt, nil),
		peerMap:     make(map[string]Peer),
	}
}

func (s *Server) Log(content string) {
	log.Printf("Server:[%s]:%s", s.address, content)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Log(fmt.Sprintf("method:[%s]", r.Method))
	s.Log(fmt.Sprintf("url:[%s]", r.URL.Path))

	if r.Method == http.MethodGet {
		s.get(w, r)
	}

	if r.Method == http.MethodPost {
		s.set(w, r)
	}

	return
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	keyList := strings.Split(r.URL.Path[1:], "/")

	// 参数校验
	if len(keyList) != paramLayerCount {
		http.Error(w, "wrong action in db", http.StatusNotFound)
		return
	}

	val, getErr := s.getAux(keyList[0], keyList[1], keyList[2])
	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(val)

	return
}

func (s *Server) getAux(db string, group string, key string) ([]byte, error) {
	if s.dbName != db{
		return nil, errors.New(fmt.Sprintf("base db exist in server. db:[%s]", db))
	}

	// 从key推导出节点地址
	nodeAddress := s.nodeHashMap.GetNode(key)
	// 如果key不存在当前节点，则从其他节点获取
	if nodeAddress != s.address {
		if peer, ok := s.peerMap[nodeAddress]; ok {
			s.Log(fmt.Sprintf("redirect to [%s]", peer.GetName()))
			return peer.Get(group, key)
		}
	}

	// 从当前节点获取
	dbMap := GetGroup(group)
	if dbMap == nil {
		return nil, errors.New(fmt.Sprintf("no such db. db:[%s]", group))
	}

	if val, ok := dbMap.Get(key); ok {
		return val.([]byte), nil
	} else {
		return nil, errors.New(fmt.Sprintf("not found this key:[%s]", key))
	}
}

func (s *Server) set(w http.ResponseWriter, r *http.Request) {
	keyList := strings.Split(r.URL.Path[1:], "/")

	// 参数校验
	if len(keyList) != paramLayerCount {
		http.Error(w, "wrong action in db", http.StatusNotFound)
		return
	}

	valBuf, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, fmt.Sprintf("read error. error:[%s]", readErr), http.StatusNotFound)
		return
	}

	setErr := s.setAux(keyList[1], keyList[2], valBuf)
	if setErr != nil {
		http.Error(w, fmt.Sprintf("read error. error:[%s]", readErr), http.StatusBadRequest)
		return
	}

	return
}

func (s *Server) setAux(group string, key string, value []byte) error {
	// 从key推导出节点地址
	nodeAddress := s.nodeHashMap.GetNode(key)
	// 如果key不存在当前节点，则写入刀其他节点
	if nodeAddress != s.address {
		if peer, ok := s.peerMap[nodeAddress]; ok {
			s.Log(fmt.Sprintf("redirect to [%s]", peer.GetName()))
			return peer.Set(group, key, value)
		}
	}

	// 写入当前节点
	dbMap := GetGroup(group)
	if dbMap == nil {
		return errors.New(fmt.Sprintf("no such db. db:[%s]", group))
	}

	return dbMap.Set(key, value)
}

// 添加服务节点
func (s *Server) AddServerNodes(peers ...string) {
	s.Lock()
	defer s.Unlock()

	s.nodeHashMap.AddNode(peers...)
	for _, peer := range peers {
		s.peerMap[peer] = &httpPeer{peer }
	}

	return
}

// 获取其他节点的Getter接口
func (s *Server) FetchPeer(node string) (Peer, bool) {
	if peer, ok := s.peerMap[node]; ok && s.address != node {
		return peer, true
	}

	return nil, false
}

type httpPeer struct {
	baseURL string
}

func (h *httpPeer) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v/%s/%v/%v",
		h.baseURL,
		dbName,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: [%v]", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: [%v]", err)
	}

	return bytes, nil
}

func (h *httpPeer) Set(group string, key string, value []byte) error {
	u := fmt.Sprintf(
		"%v/%s/%v/%v",
		h.baseURL,
		dbName,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	body := bytes.NewReader(value)
	res, postErr := http.Post(u, "application/x-www-form-urlencoded; charset=UTF-8", body)
	if postErr != nil {
		return postErr
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: [%v]", res.Status)
	}

	return nil
}

func (h *httpPeer) GetName() string {
	return h.baseURL
}
