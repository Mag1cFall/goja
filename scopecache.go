package goja

import (
	"sync/atomic"

	"github.com/Mag1cFall/goja/unistring"
)

// dynCache 是按名称查找变量的指令的作用域缓存，条目不可变并以原子指针发布
type dynCache struct {
	entry    atomic.Pointer[dynEntry]
	attaches atomic.Uint32
}

// dynEntry 是一次查找自当前作用域起依次经过的各层作用域的 scope 与目标作用域中变量的槽位（含类型位）；absent 表示 scopes 为整条作用域链且都没有该变量
type dynEntry struct {
	scopes []*map[unistring.String]uint32
	idx    uint32
	absent bool
}

// find 按缓存查找变量：hit 为真时 st 与 idx 为变量所在作用域与槽位，st 为 nil 表示变量不在作用域链上
func (c *dynCache) find(vm *vm) (st *stash, idx uint32, hit bool) {
	e := c.entry.Load()
	if e == nil {
		return nil, 0, false
	}
	st = vm.stash
	last := len(e.scopes) - 1
	for i, scope := range e.scopes {
		if st == nil || st.scope != scope || st.obj != nil {
			return nil, 0, false
		}
		if i < last || e.absent {
			st = st.outer
		}
	}
	if e.absent {
		return nil, 0, st == nil
	}
	return st, e.idx, true
}

// attach 在 name 于自当前作用域起第 depth 层的作用域 target 中找到后加入缓存条目，target 为 nil 时 depth 为作用域链的长度
func (c *dynCache) attach(vm *vm, target *stash, depth int, name unistring.String) {
	if c.attaches.Load() >= propCacheAttaches || depth == 0 && target == nil {
		return
	}
	e := &dynEntry{absent: target == nil}
	if target != nil {
		idx, ok := target.names[name]
		if !ok || target.scope == nil {
			return
		}
		e.idx = idx
		depth++
	}
	e.scopes = make([]*map[unistring.String]uint32, depth)
	st := vm.stash
	for i := range e.scopes {
		if st.obj != nil {
			return
		}
		e.scopes[i] = st.scope
		st = st.outer
	}
	c.attaches.Add(1)
	c.entry.Store(e)
}
