package goja

import (
	"sync/atomic"

	"github.com/Mag1cFall/goja/unistring"
)

// propCacheSize 是一个属性访问位置最多缓存的接收者形状数
const propCacheSize = 4

// propCacheAttaches 是一个属性访问位置最多加入缓存条目的次数，之后只使用已有条目
const propCacheAttaches = 64

// propCacheDepth 是读取缓存沿原型链查找持有者的最大深度
const propCacheDepth = 4

// propCache 是属性读写指令的内联缓存，条目不可变并以原子指针发布
type propCache struct {
	entries  atomic.Pointer[propEntries]
	attaches atomic.Uint32
}

// sharedPropSites 是进程内宿主与共享程序已分配的属性访问位置数
var sharedPropSites atomic.Int32

// sharedPropCache 返回宿主或共享程序第 site 个属性访问位置在本 agent 中的缓存，缓存随 agent 回收
func (a *Agent) sharedPropCache(site int32) *propCache {
	if int(site) >= len(a.propCaches) {
		a.propCaches = append(a.propCaches, make([]*propCache, int(sharedPropSites.Load())-len(a.propCaches))...)
	}
	c := a.propCaches[site]
	if c == nil {
		c = new(propCache)
		a.propCaches[site] = c
	}
	return c
}

// propEntries 是一组按接收者形状区分的缓存条目
type propEntries struct {
	n int
	e [propCacheSize]propEntry
}

// propEntry 是一种接收者形状上属性所在的槽位：holder 为 nil 时在接收者自身，否则在依次经过 protos 后的原型 holder 上；absent 表示原型链 protos 上都没有该属性；resolve 表示读取时按 String 对象通知解析 String
type propEntry struct {
	shape       *shape
	kind        objectKind
	resolve     bool
	absent      bool
	slot        int
	holder      *Object
	holderShape *shape
	protos      []protoLink
}

// protoLink 是原型链上位于持有者之前的一个原型及其形状
type protoLink struct {
	obj   *Object
	shape *shape
}

// lookup 按缓存返回 obj 的属性保存的值（数据值或 *valueProperty），属性不存在时返回 undefined，未命中时返回 nil
func (c *propCache) lookup(obj *Object) Value {
	es := c.entries.Load()
	if es == nil {
		return nil
	}
	var proto *Object
	protoKnown := false
next:
	for i := 0; i < es.n; i++ {
		e := &es.e[i]
		if obj.shape != e.shape || obj.kind != e.kind {
			continue
		}
		holder := obj
		if e.holder != nil || e.absent {
			if !protoKnown {
				proto, protoKnown = obj.self.proto(), true
			}
			p := proto
			for j, link := range e.protos {
				if j > 0 {
					p = p.self.proto()
				}
				if p != link.obj || p.shape != link.shape {
					continue next
				}
			}
			if len(e.protos) > 0 {
				p = p.self.proto()
			}
			if e.absent {
				if p != nil {
					continue
				}
				return _undefined
			}
			if p != e.holder || p.shape != e.holderShape {
				continue
			}
			holder = p
		}
		v := holder.slots[e.slot]
		if v != nil && e.resolve {
			obj.runtime.resolveBuiltin("String")
		}
		return v
	}
	return nil
}

// first 返回条目要求的接收者原型，自有属性的条目返回 nil
func (e *propEntry) first() *Object {
	if len(e.protos) > 0 {
		return e.protos[0].obj
	}
	return e.holder
}

// attachGet 为已按慢路径读取过的 obj 的 name 属性加入缓存条目，absent 为真时也缓存原型链上都没有该属性的结果
func (c *propCache) attachGet(obj *Object, name unistring.String, absent bool) {
	kind := obj.cacheKind()
	if c.attaches.Load() >= propCacheAttaches || !kind.cacheable(name) {
		return
	}
	s := obj.cacheShape()
	if s == nil {
		return
	}
	e := propEntry{shape: s, kind: kind, resolve: kind == kindString && name != "length"}
	if e.slot = s.lookup(name); e.slot < 0 {
		p := obj
		for depth := 0; ; depth++ {
			if p = p.self.proto(); p == nil && absent {
				e.absent = true
				break
			}
			if p == nil || depth == propCacheDepth {
				return
			}
			pk := p.cacheKind()
			if pk == kindString || !pk.cacheable(name) {
				return
			}
			ps := p.cacheShape()
			if ps == nil {
				return
			}
			if e.slot = ps.lookup(name); e.slot >= 0 {
				e.holder, e.holderShape = p, ps
				break
			}
			e.protos = append(e.protos, protoLink{obj: p, shape: ps})
		}
	}
	c.attach(e)
}

// set 按缓存把 val 写入 obj 已有的可写数据属性，未命中时返回 false
func (c *propCache) set(obj *Object, val Value) bool {
	es := c.entries.Load()
	if es == nil {
		return false
	}
	for i := 0; i < es.n; i++ {
		e := &es.e[i]
		if obj.shape == e.shape && obj.kind == e.kind {
			slot := &obj.slots[e.slot]
			switch (*slot).(type) {
			case nil, *valueProperty:
				return false
			}
			*slot = val
			return true
		}
	}
	return false
}

// attachSet 为赋值后 obj 的 name 自有数据属性加入缓存条目
func (c *propCache) attachSet(obj *Object, name unistring.String) {
	kind := obj.cacheKind()
	s := obj.shape
	if c.attaches.Load() >= propCacheAttaches || s == nil || s.dict || !kind.cacheableWrite(name) {
		return
	}
	if slot := s.lookup(name); slot >= 0 {
		switch obj.slots[slot].(type) {
		case nil, *valueProperty:
			return
		}
		c.attach(propEntry{shape: s, kind: kind, slot: slot})
	}
}

// attach 发布加入 e 后的条目组，形状、种类与接收者原型相同的旧条目被替换，条目数满时丢弃最早加入的条目
func (c *propCache) attach(e propEntry) {
	c.attaches.Add(1)
	next := &propEntries{n: 1}
	next.e[0] = e
	if old := c.entries.Load(); old != nil {
		for i := 0; i < old.n && next.n < propCacheSize; i++ {
			if old.e[i].shape != e.shape || old.e[i].kind != e.kind || old.e[i].first() != e.first() {
				next.e[next.n] = old.e[i]
				next.n++
			}
		}
	}
	c.entries.Store(next)
}

// cacheShape 返回用于属性缓存的形状，没有自有字符串键的对象取空形状，字典形状返回 nil
func (o *Object) cacheShape() *shape {
	if o.shape == nil {
		o.shape = o.runtime.rootShape()
	}
	if o.shape.dict {
		return nil
	}
	return o.shape
}
