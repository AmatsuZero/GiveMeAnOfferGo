package Core

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// weakRef -- See https://research.swtch.com/interfaces
type weakRef struct {
	t uintptr // interface type
	d uintptr // interface data
}

// newWeakRef -- create a new weakRef from target object/interface{}
func newWeakRef(v interface{}) *weakRef {
	i := (*[2]uintptr)(unsafe.Pointer(&v))
	w := &weakRef{^i[0], ^i[1]}
	runtime.SetFinalizer(&i[1], func(_ *uintptr) {
		atomic.StoreUintptr(&w.d, uintptr(0))
		atomic.StoreUintptr(&w.t, uintptr(0))
	})
	return w
}

// IsAlive -- check if the object referenced by the weakRef has already been GC'd
func (w *weakRef) IsAlive() bool {
	return atomic.LoadUintptr(&w.d) != 0
}

// GetTarget -- return a target object/interface{} from weakRef
func (w *weakRef) GetTarget() (v interface{}) {
	t := atomic.LoadUintptr(&w.t)
	d := atomic.LoadUintptr(&w.d)
	if d != 0 {
		i := (*[2]uintptr)(unsafe.Pointer(&v))
		i[0] = ^t
		i[1] = ^d
	}
	return
}

func (w *weakRef) Delete() {
	runtime.SetFinalizer(&w.d, nil)
}

type WeakMapKVRelation int

const (
	WeakMapStrongToStrong WeakMapKVRelation = iota
	WeakMapWeakToStrong
	WeakMapWeakToWeak
	WeakMapStrongToWeak
)

type WeakMap struct {
	weakToStrongStore   map[*weakRef]interface{}
	weakToWeakStore     map[*weakRef]*weakRef
	strongToStrongStore map[interface{}]interface{}
	strongToWeakStore   map[interface{}]*weakRef
	strategy            WeakMapKVRelation
}

func NewWeakMap(strategy WeakMapKVRelation) *WeakMap {
	wm := WeakMap{strategy: strategy}
	switch strategy {
	case WeakMapStrongToStrong:
		wm.strongToStrongStore = map[interface{}]interface{}{}
	case WeakMapWeakToWeak:
		wm.weakToWeakStore = map[*weakRef]*weakRef{}
	case WeakMapStrongToWeak:
		wm.strongToWeakStore = map[interface{}]*weakRef{}
	case WeakMapWeakToStrong:
		wm.weakToStrongStore = map[*weakRef]interface{}{}
	}
	return &wm
}

func (wm *WeakMap) Get(key interface{}) (interface{}, bool) {
	if key == nil {
		return nil, false
	}
	switch wm.strategy {
	case WeakMapStrongToStrong:
		val, ok := wm.strongToStrongStore[key]
		return val, ok
	case WeakMapStrongToWeak:
		val, ok := wm.strongToWeakStore[key]
		if !ok {
			return nil, ok
		}
		if !val.IsAlive() {
			val.Delete()
			delete(wm.strongToWeakStore, key)
			return nil, false
		}
		ret := val.GetTarget()
		return ret, ret != nil
	case WeakMapWeakToStrong:
		wKey := wm.checkWeakKey(key)
		if !wKey.IsAlive() {
			delete(wm.weakToStrongStore, wKey)
			wKey.Delete()
			return nil, false
		}
		val, ok := wm.weakToStrongStore[wKey]
		if !ok { // 由于不存在值，去除临时创建的变量
			wKey.Delete()
		}
		return val, ok
	case WeakMapWeakToWeak:
		wKey := wm.checkWeakKey(key)
		if !wKey.IsAlive() {
			delete(wm.weakToStrongStore, wKey)
			wKey.Delete()
			return nil, false
		}
		val, ok := wm.weakToWeakStore[wKey]
		if !ok {
			wKey.Delete()
			return nil, ok
		}
		if !val.IsAlive() {
			delete(wm.weakToStrongStore, wKey)
			val.Delete()
			wKey.Delete()
			return nil, false
		}
		ret := val.GetTarget()
		return ret, ret != nil
	}
	return nil, false
}

func (wm *WeakMap) Set(k, v interface{}) {
	if k == nil || v == nil {
		return
	}
	switch wm.strategy {
	case WeakMapStrongToStrong:
		wm.strongToStrongStore[k] = v
	case WeakMapWeakToWeak:
		key := wm.checkWeakKey(k)
		pre, ok := wm.weakToWeakStore[key]
		if ok {
			if pre.GetTarget() == v {
				return
			}
			pre.Delete()
		}
		wm.weakToWeakStore[key] = newWeakRef(v)
	case WeakMapStrongToWeak:
		// 在赋值之前，先检查一下是否已经有值，如果已有，需要把之前的删除
		val, ok := wm.Get(k)
		if ok {
			if val.(*weakRef).GetTarget() == v {
				return // 值没有变化，直接返回
			}
			val.(*weakRef).Delete() // 释放掉之前的值
		}
		wm.strongToWeakStore[k] = newWeakRef(v)
	case WeakMapWeakToStrong:
		key := wm.checkWeakKey(k)
		wm.weakToStrongStore[key] = v
	}
}

func (wm *WeakMap) getRefValue(key interface{}) (*weakRef, bool) {
	wKey := wm.checkWeakKey(key)
	if !wKey.IsAlive() {
		delete(wm.weakToStrongStore, wKey)
		return nil, false
	}
	return nil, false
}

func (wm *WeakMap) RemoveAll() {
	switch wm.strategy { // Golang 本身没有自带清除方法，完全依靠GC ……
	case WeakMapStrongToStrong:
		wm.strongToStrongStore = map[interface{}]interface{}{}
	case WeakMapStrongToWeak:
		for _, v := range wm.strongToWeakStore {
			v.Delete()
		}
		wm.strongToWeakStore = map[interface{}]*weakRef{}
	case WeakMapWeakToStrong:
		for k := range wm.weakToStrongStore {
			k.Delete()
		}
		wm.weakToStrongStore = map[*weakRef]interface{}{}
	case WeakMapWeakToWeak:
		for k, v := range wm.weakToWeakStore {
			k.Delete()
			v.Delete()
		}
		wm.weakToWeakStore = map[*weakRef]*weakRef{}
	}
}

func (wm *WeakMap) Remove(key interface{}) {
	switch wm.strategy {
	case WeakMapStrongToStrong:
		delete(wm.strongToStrongStore, key)
	case WeakMapWeakToWeak:
		k := wm.checkWeakKey(key)
		v, ok := wm.weakToWeakStore[k]
		if ok {
			delete(wm.weakToWeakStore, k)
			k.Delete()
			v.Delete()
		} else {
			k.Delete()
		}
	case WeakMapWeakToStrong:
		k := wm.checkWeakKey(key)
		delete(wm.weakToStrongStore, k)
		k.Delete()
	case WeakMapStrongToWeak:
		v, ok := wm.strongToWeakStore[key]
		if ok {
			delete(wm.strongToWeakStore, key)
			v.Delete()
		}
	}
}

/// 检查是否已经有过类似的同样的 Key, 避免重新创建弱引用Key
func (wm *WeakMap) checkWeakKey(k interface{}) *weakRef {
	if wm.strategy == WeakMapWeakToWeak {
		for k := range wm.weakToWeakStore {
			if k.GetTarget() == k {
				return k
			}
		}
	} else if wm.strategy == WeakMapWeakToStrong {
		for k := range wm.weakToStrongStore {
			if k.GetTarget() == k {
				return k
			}
		}
	}
	return newWeakRef(k)
}
