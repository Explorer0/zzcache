package zzcache

import "sync"

type call struct {
	val []byte
	err error
}

type RequestCache struct {
	sync.Mutex
	reqMap map[string]*call
}

func (rc *RequestCache) Do(key string, fn func() ([]byte, error)) ([]byte, error) {
	rc.Lock()
	if rc.reqMap == nil {
		rc.reqMap = make(map[string]*call)
	}

	if c, ok := rc.reqMap[key]; ok {
		rc.Unlock()
		return c.val, c.err
	}

	c := new(call)
	c.val, c.err = fn()
	rc.reqMap[key] = c
	rc.Unlock()

	rc.Lock()
	delete(rc.reqMap, key)
	rc.Unlock()

	return c.val, c.err
}

func NewRequestCache() *RequestCache {
	return &RequestCache{
		reqMap: make(map[string]*call),
	}
}
