package zzcache

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func([]byte) uint32

type DistributeMap struct {
	hash         HashFunc
	nodeHashList []int          // 存储节点名称哈希环
	nodeMap      map[int]string // 存储哈希值->节点名称
	replicaCnt   int
}

func (m *DistributeMap) AddNode(nodeNames ...string) {
	for _, nodeName := range nodeNames {
		for i := 0; i < m.replicaCnt; i++ {
			nodeHash := int(m.hash([]byte(strconv.Itoa(i) + nodeName)))
			m.nodeMap[nodeHash] = nodeName
			m.nodeHashList = append(m.nodeHashList, nodeHash)
		}
	}

	sort.Ints(m.nodeHashList)
}

func (m *DistributeMap) GetNode(key string) string {
	if key == "" {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.nodeHashList), func(i int) bool {
		return m.nodeHashList[i] >= hash
	})

	return m.nodeMap[m.nodeHashList[idx%len(m.nodeHashList)]]
}

func NewDistributeMap(replicaCnt int, fn HashFunc) *DistributeMap {
	m := &DistributeMap{
		hash:       fn,
		nodeMap:    make(map[int]string),
		replicaCnt: replicaCnt,
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (m *DistributeMap) format() string {
	return fmt.Sprintf("%v\n", m.nodeMap)
}
