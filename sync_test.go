package zzcache

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func BenchmarkCombinationParallel(b *testing.B) {
	NewGroup("fucker", 91929000)
	var cmpCnt int64 = 0

	// 并行测试cache是否线程安全
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache := GetGroup("fucker")
			for i:=0;i<10;i++ {
				atomic.AddInt64(&cmpCnt, 1) //原子加1操作
				key := fmt.Sprintf("%d%d",cmpCnt, i)
				cache.Set(key, i)
			}
		}
	})

	b.Logf("compare count:[%d]", cmpCnt)
	b.Logf("cache count:[%d]", GetGroup("fucker").Len())
}