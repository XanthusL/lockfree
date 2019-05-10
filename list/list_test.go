package list

import (
	"sync"
	"testing"
	"time"
)

func TestPushBack(t *testing.T) {
	l := NewList()
	now := time.Now()
	var wg sync.WaitGroup
	total := 1000000
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(m int) {
			for k := m; k < total; k += 10 {
				l.PushBack(k)
			}
			wg.Done()
			t.Log(m, "done")
		}(i)
	}
	t.Log(time.Now().Sub(now).String())
	wg.Wait()
	all := make(map[int]struct{})
	l.Walk(func(v interface{}) {
		value := v.(int)
		if _, ok := all[value]; ok {
			t.Fatal(value, "repeated")
		} else {
			all[value] = struct{}{}
		}
	})
	for i := 0; i < total; i++ {
		_, ok := all[i]
		if !ok {
			t.Fatal(i, "missing")
		}
	}

}
