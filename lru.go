package zzcache

import (
	"errors"
	"fmt"
)

type Value interface {
}

type ValueNode struct {
	Key string // 节点保存key以便于需要反查
	Value
	PrevNode *ValueNode
	NextNode *ValueNode
}

type lruCache struct {
	valueList *ValueNode
	valueMap  map[string]*ValueNode
	len       uint64
	capacity  uint64
}

func (c *lruCache) Get(key string) (interface{}, bool) {
	if item, ok := c.valueMap[key]; ok {
		// 从链表中删除当前value节点
		prevValueNode := item.PrevNode
		nextValueNode := item.NextNode
		headNode := c.valueList

		prevValueNode.NextNode = nextValueNode
		nextValueNode.PrevNode = prevValueNode

		// 添加当前value节点到链表头部
		item.NextNode = headNode.NextNode
		item.PrevNode = headNode
		headNode.NextNode.PrevNode = item
		headNode.NextNode = item

		return item.Value, true
	}

	return nil, false
}

func (c *lruCache) Delete(key string) (err error) {
	defer func() {
		fatalErr := recover()
		if fatalErr != nil {
			err = errors.New(fmt.Sprintf("delete error:[%s]", fatalErr))
			return
		}
	}()

	if item, ok := c.valueMap[key]; ok {
		// 从链表中删除当前value节点
		prevValueNode := item.PrevNode
		nextValueNode := item.NextNode
		prevValueNode.NextNode = nextValueNode
		nextValueNode.PrevNode = prevValueNode

		delete(c.valueMap, key)
		c.len -= 1
	}

	return
}

func (c *lruCache) Set(key string, value Value) (err error) {
	headNode := c.valueList
	var item *ValueNode

	defer func() {
		fatalErr := recover()
		if fatalErr != nil {
			err = errors.New(fmt.Sprintf("set error:[%s]", fatalErr))
			return
		}
	}()

	if _, ok := c.valueMap[key]; !ok { // 增加kv
		if c.len >= c.capacity { // 直接复用末尾value节点，减少内存分配
			item = headNode.PrevNode
			delete(c.valueMap, item.Key)
		} else {
			item = new(ValueNode)
			c.len += 1
		}

		c.valueMap[key] = item
		item.Key = key
		item.Value = value
	} else { // 更新kv
		item = c.valueMap[key]

		// 从链表中删除当前value节点
		prevValueNode := item.PrevNode
		nextValueNode := item.NextNode
		prevValueNode.NextNode = nextValueNode
		nextValueNode.PrevNode = prevValueNode

		item.Value = value
	}

	// 添加当前value节点到链表头部
	item.NextNode = headNode.NextNode
	item.PrevNode = headNode
	headNode.NextNode.PrevNode = item
	headNode.NextNode = item

	return
}

func (c *lruCache) Len() uint64 {
	return c.len
}

func NewLRU(capacity uint64) *lruCache {
	list := new(ValueNode)
	list.NextNode = list
	list.PrevNode = list

	return &lruCache{
		valueList: list,
		valueMap:  make(map[string]*ValueNode),
		len:       0,
		capacity:  capacity,
	}
}
