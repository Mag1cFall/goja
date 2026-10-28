package goja

import (
	"sync/atomic"
)

// region 是编译为 Go 闭包的整数循环：循环头到回跳之间的指令只读写栈上局部变量与整数常量，循环内的值以 int64 保存在寄存器中，离开循环或遇到无法按整数表示的结果时写回栈并交还解释器
type region struct {
	entry, end int
	fallback   instruction
	// blocks 是以基本块起点相对 entry 的位置为下标的语句序列，非起点为 nil
	blocks [][]rstmt
	// slots 是区域读写的局部变量编号，寄存器 i 对应 slots[i]，其后依次是栈上临时值的寄存器与各局部变量本次执行中是否写过的标记
	slots []int
	// read 标记寄存器在区域内被读取
	read  []bool
	temps int
}

const (
	// rNext 表示语句执行完毕，继续执行块内下一条语句
	rNext = -1
	// rDeopt 表示语句的结果无法按整数表示，从语句起点交还解释器
	rDeopt = -2
)

// rstmt 是区域内的一条语句：从 restart 处开始、到写入局部变量、跳转或刷新临时值为止的一段原指令
type rstmt struct {
	run func(vm *vm, regs []int64) int
	// restart 与 stack 是语句开始时的指令位置与栈上的值，交还解释器时据此恢复
	restart int
	stack   []rexpr
	// exit 是跳出区域时栈上临时值的个数
	exit int
}

// rexpr 是区域内的整数表达式：寄存器、常量或由子表达式组成的闭包
type rexpr struct {
	kind uint8
	reg  int
	c    int64
	f    func(regs []int64) (int64, bool)
}

// rReg、rConst 与 rNode 是表达式的种类：寄存器、常量与闭包
const (
	rReg uint8 = iota
	rConst
	rNode
)

// eval 返回表达式的值，结果无法按整数表示时返回 false
func (e *rexpr) eval(regs []int64) (int64, bool) {
	switch e.kind {
	case rReg:
		return regs[e.reg], true
	case rConst:
		return e.c, true
	}
	return e.f(regs)
}

// exec 在局部变量都是整数时执行区域，否则执行循环头处原来的指令
func (r *region) exec(vm *vm) {
	if regionYield(vm) || !r.enter(vm) {
		r.fallback.exec(vm)
		return
	}
	regs := vm.regionRegs
	pc := r.entry
	for {
		stmts := r.blocks[pc-r.entry]
		for i := range stmts {
			s := &stmts[i]
			res := s.run(vm, regs)
			if res == rNext {
				continue
			}
			if res == rDeopt {
				r.leave(vm, regs, s.restart, 0)
				for k := range s.stack {
					v, _ := s.stack[k].eval(regs)
					vm.push(vm.intValue(v))
				}
				return
			}
			pc = res
			if pc < r.entry || pc > r.end {
				r.leave(vm, regs, pc, s.exit)
				return
			}
			break
		}
	}
}

// regionYield 判断解释循环需要取回控制：有中断请求或 JS 性能分析正在采样
func regionYield(vm *vm) bool {
	return atomic.LoadUint32(&vm.interrupted) != 0 || vm.profTracker != nil || atomic.LoadInt32(&globalProfiler.enabled) == 1
}

// enter 检查区域读写的局部变量当前都是整数并装入寄存器，不满足时返回 false
func (r *region) enter(vm *vm) bool {
	need := 2*len(r.slots) + r.temps
	if cap(vm.regionRegs) < need {
		vm.regionRegs = make([]int64, need)
	}
	vm.regionRegs = vm.regionRegs[:need]
	regs := vm.regionRegs
	clear(regs[len(r.slots)+r.temps:])
	for i, l := range r.slots {
		index, ok := regionSlot(vm, l)
		if !ok {
			return false
		}
		v := vm.stack[index]
		if v == nil {
			return false
		}
		if r.read[i] {
			n, ok := v.(valueInt)
			if !ok {
				return false
			}
			regs[i] = int64(n)
		}
	}
	return true
}

// regionSlot 返回局部变量编号 l 在栈上的位置，参数不存在时返回 false
func regionSlot(vm *vm, l int) (int, bool) {
	if l <= 0 {
		arg := -l
		if arg > vm.args {
			return 0, false
		}
		return vm.sb + arg, true
	}
	return vm.sb + vm.args + l, true
}

// leave 把本次执行中写过的局部变量寄存器写回栈，压入 depth 个临时值后从 pc 继续解释执行
func (r *region) leave(vm *vm, regs []int64, pc, depth int) {
	dirty := len(r.slots) + r.temps
	for i, l := range r.slots {
		if regs[dirty+i] != 0 {
			index, _ := regionSlot(vm, l)
			vm.stack[index] = vm.intValue(regs[i])
		}
	}
	base := len(r.slots)
	for k := 0; k < depth; k++ {
		vm.push(vm.intValue(regs[base+k]))
	}
	vm.pc = pc
}

// regionBuilder 在编译期把一段原指令翻译为区域闭包
type regionBuilder struct {
	code      []instruction
	entry     int
	end       int
	slotIndex map[int]int
	r         *region
	depth     []int
	leader    []bool
}

// findRegions 返回 code 中可编译为整数区域的循环，键为循环头位置
func findRegions(code []instruction) map[int]*region {
	var regions map[int]*region
	ends := map[int]int{}
	for e, ins := range code {
		if j, ok := ins.(jump); ok && j < 0 {
			if _, ok := code[e+int(j)].(loopHead); ok {
				ends[e+int(j)] = e
			}
		}
	}
	for h, end := range ends {
		b := &regionBuilder{code: code, entry: h, end: end, slotIndex: map[int]int{}, r: &region{entry: h, end: end}}
		if !b.build() {
			continue
		}
		if regions == nil {
			regions = map[int]*region{}
		}
		regions[h] = b.r
	}
	return regions
}

// slot 登记局部变量编号 l 的寄存器，read 时标记区域会读取它
func (b *regionBuilder) slot(l int, read bool) {
	i, ok := b.slotIndex[l]
	if !ok {
		i = len(b.r.slots)
		b.slotIndex[l] = i
		b.r.slots = append(b.r.slots, l)
		b.r.read = append(b.r.read, false)
	}
	if read {
		b.r.read[i] = true
	}
}

// isCompare 判断指令是否为比较运算
func isCompare(ins instruction) bool {
	switch ins.(type) {
	case _op_lt, _op_lte, _op_gt, _op_gte, _op_eq, _op_neq, _op_strict_eq, _op_strict_neq:
		return true
	}
	return false
}

// stackEffect 返回区域支持的指令执行后的栈深度变化，不支持时返回 false
func stackEffect(code []instruction, pc int) (int, bool) {
	switch ins := code[pc].(type) {
	case loadStackLex:
		return 1, true
	case loadVal:
		_, ok := ins.v.(valueInt)
		return 1, ok
	case storeStackLex, loopHead, jump, _inc, _dec, _bnot:
		return 0, true
	case storeStackLexP, _pop, jneP, jeqP:
		return -1, true
	case _add, _sub, _mul, _mod, _and, _or, _xor, _sal, _sar, _shr:
		return -1, true
	}
	if isCompare(code[pc]) && pc+1 < len(code) {
		switch code[pc+1].(type) {
		case jneP, jeqP:
			return -1, true
		}
	}
	return 0, false
}

// build 计算区域内每条指令前的栈深度、划分基本块并生成语句闭包，有不支持的指令或深度不一致时返回 false
func (b *regionBuilder) build() bool {
	n := b.end - b.entry + 1
	b.depth = make([]int, n)
	b.leader = make([]bool, n)
	for i := range b.depth {
		b.depth[i] = -1
	}
	b.depth[0] = 0
	b.leader[0] = true
	work := []int{b.entry}
	maxDepth := 0
	set := func(pc, d int, branch bool) bool {
		if pc < b.entry || pc > b.end {
			return true
		}
		if branch {
			b.leader[pc-b.entry] = true
		}
		cur := b.depth[pc-b.entry]
		if cur < 0 {
			b.depth[pc-b.entry] = d
			work = append(work, pc)
			return true
		}
		return cur == d
	}
	for len(work) > 0 {
		pc := work[len(work)-1]
		work = work[:len(work)-1]
		d := b.depth[pc-b.entry]
		effect, ok := stackEffect(b.code, pc)
		if !ok || d+effect < 0 {
			return false
		}
		if d+1 > maxDepth {
			maxDepth = d + 1
		}
		switch ins := b.code[pc].(type) {
		case jump:
			if !set(pc+int(ins), d, true) {
				return false
			}
			continue
		case jneP:
			if !set(pc+int(ins), d-1, true) || !set(pc+1, d-1, true) {
				return false
			}
			continue
		case jeqP:
			if !set(pc+int(ins), d-1, true) || !set(pc+1, d-1, true) {
				return false
			}
			continue
		}
		if isCompare(b.code[pc]) {
			if d < 2 || pc+1 > b.end {
				return false
			}
			var off int
			switch j := b.code[pc+1].(type) {
			case jneP:
				off = int(j)
			case jeqP:
				off = int(j)
			}
			if b.depth[pc+1-b.entry] != -1 {
				return false
			}
			b.depth[pc+1-b.entry] = -2
			if !set(pc+1+off, d-2, true) || !set(pc+2, d-2, true) {
				return false
			}
			continue
		}
		if !set(pc+1, d+effect, false) {
			return false
		}
	}
	for pc := b.entry; pc <= b.end; pc++ {
		if b.depth[pc-b.entry] < 0 {
			continue
		}
		switch ins := b.code[pc].(type) {
		case loadStackLex:
			b.slot(int(ins), true)
		case storeStackLex:
			b.slot(int(ins), false)
		case storeStackLexP:
			b.slot(int(ins), false)
		}
	}
	b.r.temps = maxDepth + 1
	b.r.blocks = make([][]rstmt, n)
	for pc := b.entry; pc <= b.end; pc++ {
		if !b.leader[pc-b.entry] || b.depth[pc-b.entry] < 0 {
			continue
		}
		stmts, ok := b.compileBlock(pc)
		if !ok {
			return false
		}
		b.r.blocks[pc-b.entry] = stmts
	}
	return true
}

// temp 返回栈深度 k 处临时值的寄存器
func (b *regionBuilder) temp(k int) int {
	return len(b.r.slots) + k
}

// compileBlock 把从 start 开始的基本块翻译为语句序列
func (b *regionBuilder) compileBlock(start int) ([]rstmt, bool) {
	d := b.depth[start-b.entry]
	st := make([]rexpr, d)
	for k := range st {
		st[k] = rexpr{kind: rReg, reg: b.temp(k)}
	}
	var stmts []rstmt
	restart := start
	snapshot := func() []rexpr { return append([]rexpr(nil), st...) }
	restartStack := snapshot()
	emit := func(run func(vm *vm, regs []int64) int, exit int) {
		stmts = append(stmts, rstmt{run: run, restart: restart, stack: restartStack, exit: exit})
	}
	// flush 返回把栈上尚未存入对应临时寄存器的值写入该寄存器的目标与表达式，keepConst 时常量保留在栈上
	flush := func(keepConst bool) ([]int, []rexpr) {
		var targets []int
		var values []rexpr
		for k := range st {
			e := st[k]
			if e.kind == rReg && e.reg == b.temp(k) || keepConst && e.kind == rConst {
				continue
			}
			targets = append(targets, b.temp(k))
			values = append(values, e)
			st[k] = rexpr{kind: rReg, reg: b.temp(k)}
		}
		return targets, values
	}
	for pc := start; ; pc++ {
		if pc > b.end {
			return nil, false
		}
		if pc != start && b.leader[pc-b.entry] {
			targets, values := flush(false)
			marks := b.marks(targets)
			next := pc
			emit(func(vm *vm, regs []int64) int {
				if !assign(regs, targets, values, marks) {
					return rDeopt
				}
				return next
			}, len(st))
			return stmts, true
		}
		ins := b.code[pc]
		switch ins := ins.(type) {
		case loadStackLex:
			st = append(st, rexpr{kind: rReg, reg: b.slotIndex[int(ins)]})
			continue
		case loadVal:
			st = append(st, rexpr{kind: rConst, c: int64(ins.v.(valueInt))})
			continue
		case _add, _sub, _mul, _mod, _and, _or, _xor, _sal, _sar, _shr:
			op, _ := intOpOf(ins)
			x, y := st[len(st)-2], st[len(st)-1]
			st = append(st[:len(st)-2], binaryExpr(op, x, y))
			continue
		case _inc, _dec, _bnot:
			st[len(st)-1] = unaryExpr(ins, st[len(st)-1])
			continue
		case _pop:
			st = st[:len(st)-1]
			continue
		case storeStackLexP, storeStackLex:
			var l int
			keep := false
			switch s := ins.(type) {
			case storeStackLexP:
				l = int(s)
			case storeStackLex:
				l, keep = int(s), true
			}
			v := st[len(st)-1]
			st = st[:len(st)-1]
			targets, values := flush(true)
			targets = append(targets, b.slotIndex[l])
			values = append(values, v)
			if keep {
				targets = append(targets, b.temp(len(st)))
				values = append(values, v)
				st = append(st, rexpr{kind: rReg, reg: b.temp(len(st))})
			}
			marks := b.marks(targets)
			emit(func(vm *vm, regs []int64) int {
				if !assign(regs, targets, values, marks) {
					return rDeopt
				}
				return rNext
			}, 0)
			restart, restartStack = pc+1, snapshot()
			continue
		case loopHead:
			if targets, values := flush(true); len(targets) > 0 {
				marks := b.marks(targets)
				emit(func(vm *vm, regs []int64) int {
					if !assign(regs, targets, values, marks) {
						return rDeopt
					}
					return rNext
				}, 0)
			}
			restart, restartStack = pc, snapshot()
			emit(func(vm *vm, regs []int64) int {
				if regionYield(vm) {
					return rDeopt
				}
				p := vm.prg
				p.countWarmUp()
				if !p.host && p.warmUp >= baselineInterpreterWarmUpThreshold {
					p.baseline = true
				}
				return rNext
			}, 0)
			restart, restartStack = pc+1, snapshot()
			continue
		case jump:
			targets, values := flush(false)
			marks := b.marks(targets)
			target := pc + int(ins)
			emit(func(vm *vm, regs []int64) int {
				if !assign(regs, targets, values, marks) {
					return rDeopt
				}
				return target
			}, len(st))
			return stmts, true
		case jneP, jeqP:
			cond := st[len(st)-1]
			st = st[:len(st)-1]
			targets, values := flush(false)
			marks := b.marks(targets)
			var target int
			onTrue := false
			switch j := ins.(type) {
			case jneP:
				target = pc + int(j)
			case jeqP:
				target, onTrue = pc+int(j), true
			}
			next := pc + 1
			emit(func(vm *vm, regs []int64) int {
				c, ok := cond.eval(regs)
				if !ok || !assign(regs, targets, values, marks) {
					return rDeopt
				}
				if (c != 0) == onTrue {
					return target
				}
				return next
			}, len(st))
			return stmts, true
		}
		if isCompare(ins) {
			op, _ := intOpOf(ins)
			x, y := st[len(st)-2], st[len(st)-1]
			st = st[:len(st)-2]
			targets, values := flush(false)
			marks := b.marks(targets)
			var target int
			onTrue := false
			switch j := b.code[pc+1].(type) {
			case jneP:
				target = pc + 1 + int(j)
			case jeqP:
				target, onTrue = pc+1+int(j), true
			}
			next := pc + 2
			emit(func(vm *vm, regs []int64) int {
				left, ok := x.eval(regs)
				if !ok {
					return rDeopt
				}
				right, ok := y.eval(regs)
				if !ok || !assign(regs, targets, values, marks) {
					return rDeopt
				}
				if op.compare(left, right) == onTrue {
					return target
				}
				return next
			}, len(st))
			return stmts, true
		}
		return nil, false
	}
}

// marks 返回 targets 中局部变量寄存器对应的写入标记寄存器
func (b *regionBuilder) marks(targets []int) []int {
	var marks []int
	for _, t := range targets {
		if t < len(b.r.slots) {
			marks = append(marks, len(b.r.slots)+b.r.temps+t)
		}
	}
	return marks
}

// assign 先求出全部表达式的值再写入目标寄存器并标记写过的局部变量，有表达式无法按整数表示时不写入并返回 false
func assign(regs []int64, targets []int, values []rexpr, marks []int) bool {
	switch len(targets) {
	case 0:
		return true
	case 1:
		v, ok := values[0].eval(regs)
		if !ok {
			return false
		}
		regs[targets[0]] = v
		for _, m := range marks {
			regs[m] = 1
		}
		return true
	}
	var buf [8]int64
	results := buf[:0]
	for i := range values {
		v, ok := values[i].eval(regs)
		if !ok {
			return false
		}
		results = append(results, v)
	}
	for i, t := range targets {
		regs[t] = results[i]
	}
	for _, m := range marks {
		regs[m] = 1
	}
	return true
}

// unaryExpr 返回自增、自减或按位取反的表达式
func unaryExpr(ins instruction, x rexpr) rexpr {
	switch ins.(type) {
	case _bnot:
		return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
			v, ok := x.eval(regs)
			return int64(^int32(v)), ok
		}}
	}
	delta := int64(1)
	if _, ok := ins.(_dec); ok {
		delta = -1
	}
	return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
		v, ok := x.eval(regs)
		v += delta
		return v, ok && v >= -maxInt && v <= maxInt
	}}
}

// regionApply 返回整数运算的结果，结果不能按 valueInt 表示时返回 false
func regionApply(op intOp, left, right int64) (int64, bool) {
	var v int64
	switch op {
	case intAdd:
		v = left + right
	case intSub:
		v = left - right
	case intMul:
		if left == 0 || right == 0 {
			if left < 0 || right < 0 {
				return 0, false
			}
			return 0, true
		}
		v = left * right
		if v/left != right {
			return 0, false
		}
	case intMod:
		if right == 0 {
			return 0, false
		}
		v = left % right
		if v == 0 && left < 0 {
			return 0, false
		}
	case intAnd:
		return int64(int32(left) & int32(right)), true
	case intOr:
		return int64(int32(left) | int32(right)), true
	case intXor:
		return int64(int32(left) ^ int32(right)), true
	case intSal:
		return int64(int32(left) << (uint32(right) & 0x1F)), true
	case intSar:
		return int64(int32(left) >> (uint32(right) & 0x1F)), true
	case intShr:
		return int64(uint32(left) >> (uint32(right) & 0x1F)), true
	}
	if v < -maxInt || v > maxInt {
		return 0, false
	}
	return v, true
}
