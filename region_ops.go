package goja

// binaryExpr 返回按运算与操作数种类特化的整数二元运算表达式
func binaryExpr(op intOp, x, y rexpr) rexpr {
	ar, ac, af := x.reg, x.c, x.f
	br, bc, bf := y.reg, y.c, y.f
	switch op {
	case intAdd:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a + b
				return v, v >= -maxInt && v <= maxInt
			}}
		}
	case intSub:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				v := a - b
				return v, v >= -maxInt && v <= maxInt
			}}
		}
	case intAnd:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(int32(a) & int32(b)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) & int32(b)), true
			}}
		}
	case intOr:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(int32(a) | int32(b)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) | int32(b)), true
			}}
		}
	case intXor:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(int32(a) ^ int32(b)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) ^ int32(b)), true
			}}
		}
	case intSal:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) << (uint32(b) & 0x1F)), true
			}}
		}
	case intSar:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(int32(a) >> (uint32(b) & 0x1F)), true
			}}
		}
	case intShr:
		switch {
		case x.kind == rReg && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := regs[br]
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b := bc
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rReg && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := regs[ar]
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := regs[br]
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b := bc
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rConst && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a := ac
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rReg:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := regs[br]
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rConst:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b := bc
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		case x.kind == rNode && y.kind == rNode:
			return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
				a, ok := af(regs)
				if !ok {
					return 0, false
				}
				b, ok := bf(regs)
				if !ok {
					return 0, false
				}
				return int64(uint32(a) >> (uint32(b) & 0x1F)), true
			}}
		}
	}
	return rexpr{kind: rNode, f: func(regs []int64) (int64, bool) {
		left, ok := x.eval(regs)
		if !ok {
			return 0, false
		}
		right, ok := y.eval(regs)
		if !ok {
			return 0, false
		}
		return regionApply(op, left, right)
	}}
}
