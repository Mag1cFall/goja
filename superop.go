package goja

// super 系列指令把取两个操作数的指令、整数二元运算及其后的条件跳转或局部变量写入合为一次分派；两个操作数都是整数时直接算出结果，否则顺序执行被合并的原指令，原指令保留在原位置供跳转与异常处理

// intOp 是超级指令的整数二元运算
type intOp uint8

const (
	intAdd intOp = iota
	intSub
	intMul
	intMod
	intAnd
	intOr
	intXor
	intSal
	intSar
	intShr
	intLt
	intLte
	intGt
	intGte
	intEq
	intNeq
)

// intOpOf 返回二元运算指令对应的整数运算
func intOpOf(op instruction) (intOp, bool) {
	switch op.(type) {
	case _add:
		return intAdd, true
	case _sub:
		return intSub, true
	case _mul:
		return intMul, true
	case _mod:
		return intMod, true
	case _and:
		return intAnd, true
	case _or:
		return intOr, true
	case _xor:
		return intXor, true
	case _sal:
		return intSal, true
	case _sar:
		return intSar, true
	case _shr:
		return intShr, true
	case _op_lt:
		return intLt, true
	case _op_lte:
		return intLte, true
	case _op_gt:
		return intGt, true
	case _op_gte:
		return intGte, true
	case _op_eq, _op_strict_eq:
		return intEq, true
	case _op_neq, _op_strict_neq:
		return intNeq, true
	}
	return 0, false
}

// compare 返回整数比较的结果
func (op intOp) compare(left, right int64) bool {
	switch op {
	case intLt:
		return left < right
	case intLte:
		return left <= right
	case intGt:
		return left > right
	case intGte:
		return left >= right
	case intEq:
		return left == right
	default:
		return left != right
	}
}

// apply 返回整数运算的结果值
func (op intOp) apply(vm *vm, left, right int64) Value {
	switch op {
	case intAdd:
		return vm.intValue(left + right)
	case intSub:
		return vm.intValue(left - right)
	case intMul:
		return vm.mulInt(left, right)
	case intMod:
		return vm.modInt(left, right)
	case intAnd:
		return vm.intValue(int64(int32(left) & int32(right)))
	case intOr:
		return vm.intValue(int64(int32(left) | int32(right)))
	case intXor:
		return vm.intValue(int64(int32(left) ^ int32(right)))
	case intSal:
		return vm.intValue(int64(int32(left) << (uint32(right) & 0x1F)))
	case intSar:
		return vm.intValue(int64(int32(left) >> (uint32(right) & 0x1F)))
	case intShr:
		return vm.intValue(int64(uint32(left) >> (uint32(right) & 0x1F)))
	}
	if op.compare(left, right) {
		return valueTrue
	}
	return valueFalse
}

// operand 是超级指令的操作数：常量或栈上局部变量
type operand struct {
	value Value
	lex   loadStackLex
	local bool
}

// operandOf 返回取值指令对应的操作数
func operandOf(ins instruction) (operand, bool) {
	switch ins := ins.(type) {
	case loadStackLex:
		return operand{lex: ins, local: true}, true
	case loadVal:
		return operand{value: ins.v}, true
	}
	return operand{}, false
}

// intValue 返回操作数当前的整数值
func (o *operand) intValue(vm *vm) (int64, bool) {
	v := o.value
	if o.local {
		if o.lex <= 0 {
			arg := int(-o.lex)
			if arg > vm.args {
				return 0, false
			}
			v = vm.stack[vm.sb+arg]
		} else {
			v = vm.stack[vm.sb+vm.args+int(o.lex)]
		}
	}
	i, ok := v.(valueInt)
	return int64(i), ok
}

// superBinary 合并取两个操作数的指令与一条二元运算
type superBinary struct {
	left, right operand
	op          intOp
	seq         []instruction
}

func (s *superBinary) exec(vm *vm) {
	if left, ok := s.left.intValue(vm); ok {
		if right, ok := s.right.intValue(vm); ok {
			vm.push(s.op.apply(vm, left, right))
			vm.pc += len(s.seq)
			return
		}
	}
	runSequence(vm, s.seq)
}

// superCompareJump 合并取两个操作数的指令、一条比较与随后弹出条件的跳转
type superCompareJump struct {
	left, right operand
	op          intOp
	offset      int
	onTrue      bool
	seq         []instruction
}

func (s *superCompareJump) exec(vm *vm) {
	if left, ok := s.left.intValue(vm); ok {
		if right, ok := s.right.intValue(vm); ok {
			if s.op.compare(left, right) == s.onTrue {
				vm.pc += len(s.seq) - 1 + s.offset
			} else {
				vm.pc += len(s.seq)
			}
			return
		}
	}
	runSequence(vm, s.seq)
}

// superBinaryStore 合并取两个操作数的指令、一条二元运算与随后弹出结果的局部变量写入
type superBinaryStore struct {
	left, right operand
	op          intOp
	target      int
	seq         []instruction
}

func (s *superBinaryStore) exec(vm *vm) {
	if left, ok := s.left.intValue(vm); ok {
		if right, ok := s.right.intValue(vm); ok {
			var p *Value
			if s.target < 0 {
				p = &vm.stack[vm.sb-s.target]
			} else {
				p = &vm.stack[vm.sb+vm.args+s.target]
			}
			if *p != nil {
				*p = s.op.apply(vm, left, right)
				vm.pc += len(s.seq)
				return
			}
		}
	}
	runSequence(vm, s.seq)
}

// runSequence 依次执行被合并的原指令，执行离开顺序位置时停止
func runSequence(vm *vm, seq []instruction) {
	pc, prg := vm.pc, vm.prg
	for k, ins := range seq {
		if vm.pc != pc+k || vm.prg != prg {
			return
		}
		ins.exec(vm)
	}
}

// fuseSuper 返回从 code[i] 起的超级指令与它覆盖的指令数，不匹配时返回 nil
func fuseSuper(code []instruction, i int) (instruction, int) {
	if i+2 >= len(code) {
		return nil, 0
	}
	left, ok := operandOf(code[i])
	if !ok {
		return nil, 0
	}
	right, ok := operandOf(code[i+1])
	if !ok {
		return nil, 0
	}
	op, ok := intOpOf(code[i+2])
	if !ok {
		return nil, 0
	}
	end := i + 3
	var tail instruction
	if end < len(code) {
		switch next := code[end].(type) {
		case jneP, jeqP:
			if op >= intLt {
				tail = next
			}
		case storeStackLexP:
			tail = next
		}
	}
	if tail != nil {
		end++
	}
	seq := append([]instruction(nil), code[i:end]...)
	switch next := tail.(type) {
	case jneP:
		return &superCompareJump{left: left, right: right, op: op, offset: int(next), seq: seq}, len(seq)
	case jeqP:
		return &superCompareJump{left: left, right: right, op: op, offset: int(next), onTrue: true, seq: seq}, len(seq)
	case storeStackLexP:
		return &superBinaryStore{left: left, right: right, op: op, target: int(next), seq: seq}, len(seq)
	}
	return &superBinary{left: left, right: right, op: op, seq: seq}, len(seq)
}
