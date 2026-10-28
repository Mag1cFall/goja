package goja

import (
	"hash/maphash"
	"math"
	"sort"

	"github.com/Mag1cFall/goja/unistring"
)

// shape 是对象自有字符串键的布局：同一 Runtime 中按相同顺序加入相同键的对象共享同一形状，第 i 个键的值保存在 Object.slots[i]
type shape struct {
	names    []unistring.String
	table    shapeTable
	child    *shape
	children map[unistring.String]*shape
	searches uint8
	extended bool
	hasIdx   bool
	dict     bool
	sorted   int
}

// shapeTableMin 是形状按需建立键索引的最少键数
const shapeTableMin = 8

// shapeTable 是形状键名的开放寻址索引，按键名哈希定位，元素为键的下标加 1，0 为空位
type shapeTable []int32

// shapeSeed 是形状索引的键名哈希种子
var shapeSeed = maphash.MakeSeed()

// newShapeTable 为 names 建立装载率不超过一半的索引
func newShapeTable(names []unistring.String) shapeTable {
	size := 16
	for size < 2*len(names) {
		size <<= 1
	}
	t := make(shapeTable, size)
	for i, name := range names {
		t.insert(name, i)
	}
	return t
}

// insert 记录 name 位于下标 i
func (t shapeTable) insert(name unistring.String, i int) {
	mask := uint64(len(t) - 1)
	h := maphash.String(shapeSeed, string(name)) & mask
	for t[h] != 0 {
		h = (h + 1) & mask
	}
	t[h] = int32(i + 1)
}

// find 返回 name 在 names 中的下标，不存在时返回 -1
func (t shapeTable) find(names []unistring.String, name unistring.String) int {
	mask := uint64(len(t) - 1)
	for h := maphash.String(shapeSeed, string(name)) & mask; ; h = (h + 1) & mask {
		e := t[h]
		if e == 0 {
			return -1
		}
		if names[e-1] == name {
			return int(e - 1)
		}
	}
}

// shapeTreeMax 是共享形状的最多键数，再加键时对象转为字典形状
const shapeTreeMax = 128

// objectKind 区分具名属性语义不同的对象，决定属性缓存可以按形状处理哪些键
type objectKind uint8

const (
	kindUnknown objectKind = iota
	kindNone
	kindPlain
	kindArray
	kindString
	kindFunction
	kindArguments
	kindTemplated
	kindRegExp
)

// lookup 返回 name 的槽位，不存在时返回 -1
func (s *shape) lookup(name unistring.String) int {
	if s.table == nil && len(s.names) > shapeTableMin {
		if s.dict || s.searches >= 4 {
			s.buildTable()
		} else {
			s.searches++
		}
	}
	if s.table != nil {
		return s.table.find(s.names, name)
	}
	for i := len(s.names) - 1; i >= 0; i-- {
		if s.names[i] == name {
			return i
		}
	}
	return -1
}

// buildTable 建立键到槽位的索引
func (s *shape) buildTable() {
	s.table = newShapeTable(s.names)
}

// transition 返回加入 name 后的共享形状，尚未出现过时返回 nil
func (s *shape) transition(name unistring.String) *shape {
	if c := s.child; c != nil && c.names[len(c.names)-1] == name {
		return c
	}
	return s.children[name]
}

// addChild 创建加入 name 后的共享形状，子形状在父形状之后加入的键数组上追加，父形状的键索引移交给子形状
func (s *shape) addChild(name unistring.String) *shape {
	n := len(s.names)
	var names []unistring.String
	if !s.extended && n < cap(s.names) {
		names = s.names[:n+1]
		names[n] = name
	} else {
		names = append(s.names[:n:n], name)
	}
	s.extended = true
	c := &shape{names: names, hasIdx: s.hasIdx || strToArrayIdx(name) != math.MaxUint32}
	if s.table != nil {
		c.table, s.table, s.searches = s.table, nil, 0
		if 2*len(names) > len(c.table) {
			c.table = newShapeTable(names)
		} else {
			c.table.insert(name, n)
		}
	}
	if s.child == nil {
		s.child = c
	} else {
		if s.children == nil {
			s.children = make(map[unistring.String]*shape)
		}
		s.children[name] = c
	}
	return c
}

// toDict 返回键与顺序和 s 相同、只属于一个对象的字典形状
func (s *shape) toDict() *shape {
	return &shape{names: append([]unistring.String(nil), s.names...), hasIdx: s.hasIdx, dict: true}
}

// dictAdd 在字典形状末尾加入 name
func (s *shape) dictAdd(name unistring.String) {
	s.names = append(s.names, name)
	if s.table != nil {
		if 2*len(s.names) > len(s.table) {
			s.table = newShapeTable(s.names)
		} else {
			s.table.insert(name, len(s.names)-1)
		}
	}
	if strToArrayIdx(name) != math.MaxUint32 {
		s.hasIdx = true
	}
}

// dictRemove 从字典形状删除第 i 个键，键数组换为新数组以保持进行中的枚举不变
func (s *shape) dictRemove(i int) {
	names := make([]unistring.String, 0, len(s.names)-1)
	names = append(names, s.names[:i]...)
	s.names = append(names, s.names[i+1:]...)
	s.table = nil
	s.hasIdx = hasIdxName(s.names)
	if i < s.sorted {
		s.sorted--
	}
}

// rootShape 返回本 Runtime 中没有自有字符串键的形状
func (r *Runtime) rootShape() *shape {
	if r.emptyShape == nil {
		r.emptyShape = &shape{}
	}
	return r.emptyShape
}

// templateShape 返回本 Runtime 中按模板键顺序排列的形状
func (r *Runtime) templateShape(t *objectTemplate) *shape {
	if s := r.templateShapes[t]; s != nil {
		return s
	}
	n := len(t.propNames)
	s := &shape{names: t.propNames[:n:n], hasIdx: hasIdxName(t.propNames)}
	if r.templateShapes == nil {
		r.templateShapes = make(map[*objectTemplate]*shape)
	}
	r.templateShapes[t] = s
	return s
}

// ownIndex 返回自有字符串键 name 的槽位，不存在时返回 -1
func (o *baseObject) ownIndex(name unistring.String) int {
	if s := o.val.shape; s != nil {
		return s.lookup(name)
	}
	return -1
}

// ownValue 返回自有字符串键 name 保存的值，不存在或尚未生成时返回 nil
func (o *baseObject) ownValue(name unistring.String) Value {
	if i := o.ownIndex(name); i >= 0 {
		return o.val.slots[i]
	}
	return nil
}

// hasIdxKeys 判断对象是否有整数自有字符串键
func (o *baseObject) hasIdxKeys() bool {
	s := o.val.shape
	return s != nil && s.hasIdx
}

// toDict 把对象换为只属于它的字典形状并返回该形状
func (o *baseObject) toDict() *shape {
	obj := o.val
	if obj.shape == nil {
		obj.shape = &shape{dict: true}
	} else if !obj.shape.dict {
		obj.shape = obj.shape.toDict()
	}
	return obj.shape
}

// keyOrder 返回自有字符串键的枚举顺序：整数键按数值升序在前，其余键按加入顺序；有整数键的对象换为字典形状并按此顺序重排
func (o *baseObject) keyOrder() []unistring.String {
	s := o.val.shape
	if s == nil {
		return nil
	}
	if !s.hasIdx || s.dict && s.sorted == len(s.names) {
		return s.names
	}
	s = o.toDict()
	slots := o.val.slots
	order := make([]int, len(s.names))
	for i := range order {
		order[i] = i
	}
	sort.SliceStable(order, func(a, b int) bool {
		x, y := strToArrayIdx(s.names[order[a]]), strToArrayIdx(s.names[order[b]])
		return x != math.MaxUint32 && (y == math.MaxUint32 || x < y)
	})
	names := make([]unistring.String, len(order))
	values := make([]Value, len(order), cap(slots))
	for i, j := range order {
		names[i], values[i] = s.names[j], slots[j]
	}
	s.names, s.table, s.sorted = names, nil, len(names)
	o.val.slots = values
	return names
}

// reorder 按 names 给出的顺序重排自有字符串键与槽位，共享形状的对象换为按该顺序加键得到的共享形状
func (o *baseObject) reorder(names []unistring.String) {
	obj := o.val
	s := obj.shape
	values := make([]Value, len(names), cap(obj.slots))
	for i, name := range names {
		values[i] = obj.slots[s.lookup(name)]
	}
	obj.slots = values
	if s.dict {
		s.names, s.table, s.sorted = names, nil, 0
		return
	}
	if len(names) > shapeTreeMax {
		obj.shape = &shape{names: names, hasIdx: hasIdxName(names), dict: true}
		return
	}
	s = obj.runtime.rootShape()
	for _, name := range names {
		if c := s.transition(name); c != nil {
			s = c
		} else {
			s = s.addChild(name)
		}
	}
	obj.shape = s
}

// hasIdxName 判断 names 中是否有整数键
func hasIdxName(names []unistring.String) bool {
	for _, name := range names {
		if strToArrayIdx(name) != math.MaxUint32 {
			return true
		}
	}
	return false
}

// appendable 判断在数组末尾追加 n 个元素是否只需追加到元素表：数组稠密、可扩展、长度可写且原型链上没有整数键
func (a *arrayObject) appendable(n int) bool {
	if !a.extensible || !a.lengthProp.writable || a.length != uint32(len(a.values)) || uint64(a.length)+uint64(n) > math.MaxUint32 {
		return false
	}
	for p := a.prototype; p != nil; p = p.self.proto() {
		switch p.self.(type) {
		case *baseObject, *templatedObject, *templatedArrayObject:
			if p.shape != nil && p.shape.hasIdx {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// kindOf 返回 self 的属性缓存种类，具名属性不全部按形状读写的对象为 kindNone
func kindOf(self objectImpl) objectKind {
	switch self.(type) {
	case *baseObject, *errorObject, *primitiveValueObject, *mapObject, *setObject, *weakMapObject, *weakSetObject,
		*Promise, *dateObject, *arrayIterObject, *mapIterObject, *setIterObject, *generatorObject,
		*asyncGeneratorObject:
		return kindPlain
	case *arrayObject, *sparseArrayObject, *templatedArrayObject:
		return kindArray
	case *stringObject:
		return kindString
	case *funcObject, *arrowFuncObject, *methodFuncObject, *classFuncObject, *nativeFuncObject, *boundFuncObject,
		*generatorFuncObject, *generatorMethodFuncObject, *asyncFuncObject, *asyncArrowFuncObject, *asyncMethodFuncObject,
		*asyncGeneratorFuncObject, *asyncGeneratorMethodFuncObject:
		return kindFunction
	case *argumentsObject:
		return kindArguments
	case *regexpObject:
		return kindRegExp
	case *templatedObject, *templatedFuncObject:
		return kindTemplated
	}
	return kindNone
}

// cacheKind 返回对象的属性缓存种类，首次调用时按对象实现确定
func (o *Object) cacheKind() objectKind {
	if o.kind == kindUnknown {
		o.kind = kindOf(o.self)
	}
	return o.kind
}

// cacheable 判断 kind 种类对象的 name 键是否完全由形状决定
func (kind objectKind) cacheable(name unistring.String) bool {
	if kind <= kindNone || strToArrayIdx(name) != math.MaxUint32 {
		return false
	}
	switch kind {
	case kindArray:
		return name != "length"
	case kindFunction:
		return name != "length" && name != "name" && name != "prototype"
	}
	return true
}

// cacheableWrite 判断 kind 种类对象对 name 键的赋值是否只改写槽位
func (kind objectKind) cacheableWrite(name unistring.String) bool {
	return kind.cacheable(name) && (kind != kindRegExp || name != "exec")
}
