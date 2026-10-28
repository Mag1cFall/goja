package goja

// fuse 系列指令把同一程序中相邻两条指令合为一次分派；第一条指令顺序执行到下一条时才执行第二条，第二条指令保留在原位置供跳转与异常处理

type fuseLoadStackLexLoadVal struct {
	first  loadStackLex
	second loadVal
}

func (f *fuseLoadStackLexLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexLoadStackLex struct {
	first  loadStackLex
	second loadStackLex
}

func (f *fuseLoadStackLexLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseStoreStackLexPLoadStackLex struct {
	first  storeStackLexP
	second loadStackLex
}

func (f *fuseStoreStackLexPLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseJnePLoadStackLex struct {
	first  jneP
	second loadStackLex
}

func (f *fuseJnePLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexGetPropSource struct {
	first  loadStackLex
	second getPropSource
}

func (f *fuseLoadStackLexGetPropSource) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValOpEq struct {
	first  loadVal
	second _op_eq
}

func (f *fuseLoadValOpEq) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValOr struct {
	first  loadVal
	second _or
}

func (f *fuseLoadValOr) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOpEqJneP struct {
	first  _op_eq
	second jneP
}

func (f *fuseOpEqJneP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseGetPropSourceLoadStackLex struct {
	first  getPropSource
	second loadStackLex
}

func (f *fuseGetPropSourceLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValLoadStackLex struct {
	first  loadVal
	second loadStackLex
}

func (f *fuseLoadValLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoopHeadLoadStackLex struct {
	first  loopHead
	second loadStackLex
}

func (f *fuseLoopHeadLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseCheckBaseGetPropCallee struct {
	first  checkBase
	second getPropCallee
}

func (f *fuseCheckBaseGetPropCallee) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOrStoreStackLexP struct {
	first  _or
	second storeStackLexP
}

func (f *fuseOrStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexCall struct {
	first  loadStackLex
	second call
}

func (f *fuseLoadStackLexCall) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStashLexLoadVal struct {
	first  loadStashLex
	second loadVal
}

func (f *fuseLoadStashLexLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexGetElem struct {
	first  loadStackLex
	second _getElem
}

func (f *fuseLoadStackLexGetElem) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseStoreStackLexPJump struct {
	first  storeStackLexP
	second jump
}

func (f *fuseStoreStackLexPJump) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOrSub struct {
	first  _or
	second _sub
}

func (f *fuseOrSub) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStashLexLoadStashLex struct {
	first  loadStashLex
	second loadStashLex
}

func (f *fuseLoadStashLexLoadStashLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValShr struct {
	first  loadVal
	second _shr
}

func (f *fuseLoadValShr) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOpLtJneP struct {
	first  _op_lt
	second jneP
}

func (f *fuseOpLtJneP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseXorStoreStackLexP struct {
	first  _xor
	second storeStackLexP
}

func (f *fuseXorStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseAddStoreStackLexP struct {
	first  _add
	second storeStackLexP
}

func (f *fuseAddStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseGetPropCalleeLoadStackLex struct {
	first  getPropCallee
	second loadStackLex
}

func (f *fuseGetPropCalleeLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValAnd struct {
	first  loadVal
	second _and
}

func (f *fuseLoadValAnd) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStashLexMul struct {
	first  loadStashLex
	second _mul
}

func (f *fuseLoadStashLexMul) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexStoreStackLexP struct {
	first  loadStackLex
	second storeStackLexP
}

func (f *fuseLoadStackLexStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseAddLoadVal struct {
	first  _add
	second loadVal
}

func (f *fuseAddLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fusePopLoadStackLex struct {
	first  _pop
	second loadStackLex
}

func (f *fusePopLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseMulAdd struct {
	first  _mul
	second _add
}

func (f *fuseMulAdd) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOrAdd struct {
	first  _or
	second _add
}

func (f *fuseOrAdd) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexBnot struct {
	first  loadStackLex
	second _bnot
}

func (f *fuseLoadStackLexBnot) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseJneLoadStackLex struct {
	first  jne
	second loadStackLex
}

func (f *fuseJneLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexCheckBase struct {
	first  loadStackLex
	second checkBase
}

func (f *fuseLoadStackLexCheckBase) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseJnePLoadVal struct {
	first  jneP
	second loadVal
}

func (f *fuseJnePLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadUndefLoadStashLex struct {
	first  _loadUndef
	second loadStashLex
}

func (f *fuseLoadUndefLoadStashLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOrLoadStackLex struct {
	first  _or
	second loadStackLex
}

func (f *fuseOrLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseSalOr struct {
	first  _sal
	second _or
}

func (f *fuseSalOr) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexSal struct {
	first  loadStackLex
	second _sal
}

func (f *fuseLoadStackLexSal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValSal struct {
	first  loadVal
	second _sal
}

func (f *fuseLoadValSal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexAdd struct {
	first  loadStackLex
	second _add
}

func (f *fuseLoadStackLexAdd) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseOpStrictEqJneP struct {
	first  _op_strict_eq
	second jneP
}

func (f *fuseOpStrictEqJneP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseGetElemStoreStackLexP struct {
	first  _getElem
	second storeStackLexP
}

func (f *fuseGetElemStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseBnotLoadVal struct {
	first  _bnot
	second loadVal
}

func (f *fuseBnotLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStashLexLoadStackLex struct {
	first  loadStashLex
	second loadStackLex
}

func (f *fuseLoadStashLexLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadUndefRet struct {
	first  _loadUndef
	second _ret
}

func (f *fuseLoadUndefRet) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValMod struct {
	first  loadVal
	second _mod
}

func (f *fuseLoadValMod) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValOpStrictEq struct {
	first  loadVal
	second _op_strict_eq
}

func (f *fuseLoadValOpStrictEq) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexRet struct {
	first  loadStackLex
	second _ret
}

func (f *fuseLoadStackLexRet) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseSubLoadStackLex struct {
	first  _sub
	second loadStackLex
}

func (f *fuseSubLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexInc struct {
	first  loadStackLex
	second _inc
}

func (f *fuseLoadStackLexInc) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseIncStoreStackLexP struct {
	first  _inc
	second storeStackLexP
}

func (f *fuseIncStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseStoreStackLexPLoadVal struct {
	first  storeStackLexP
	second loadVal
}

func (f *fuseStoreStackLexPLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseInitStackPLoadStackLex struct {
	first  initStackP
	second loadStackLex
}

func (f *fuseInitStackPLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseSubLoadVal struct {
	first  _sub
	second loadVal
}

func (f *fuseSubLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValSar struct {
	first  loadVal
	second _sar
}

func (f *fuseLoadValSar) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexXor struct {
	first  loadStackLex
	second _xor
}

func (f *fuseLoadStackLexXor) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValStoreStackLexP struct {
	first  loadVal
	second storeStackLexP
}

func (f *fuseLoadValStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseJnePop struct {
	first  jne
	second _pop
}

func (f *fuseJnePop) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValLoadStashLex struct {
	first  loadVal
	second loadStashLex
}

func (f *fuseLoadValLoadStashLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValXor struct {
	first  loadVal
	second _xor
}

func (f *fuseLoadValXor) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseSalLoadStackLex struct {
	first  _sal
	second loadStackLex
}

func (f *fuseSalLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValOpLt struct {
	first  loadVal
	second _op_lt
}

func (f *fuseLoadValOpLt) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseAddXor struct {
	first  _add
	second _xor
}

func (f *fuseAddXor) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseShrLoadStackLex struct {
	first  _shr
	second loadStackLex
}

func (f *fuseShrLoadStackLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseShrOr struct {
	first  _shr
	second _or
}

func (f *fuseShrOr) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValAdd struct {
	first  loadVal
	second _add
}

func (f *fuseLoadValAdd) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseJnePLoadStashLex struct {
	first  jneP
	second loadStashLex
}

func (f *fuseJnePLoadStashLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexJne struct {
	first  loadStackLex
	second jne
}

func (f *fuseLoadStackLexJne) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseStoreStackLexPLoopHead struct {
	first  storeStackLexP
	second loopHead
}

func (f *fuseStoreStackLexPLoopHead) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseGetPropSourceJneP struct {
	first  getPropSource
	second jneP
}

func (f *fuseGetPropSourceJneP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseDupGetPropSource struct {
	first  _dup
	second getPropSource
}

func (f *fuseDupGetPropSource) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fusePopLoadUndef struct {
	first  _pop
	second _loadUndef
}

func (f *fusePopLoadUndef) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexDup struct {
	first  loadStackLex
	second _dup
}

func (f *fuseLoadStackLexDup) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseMulLoadStashLex struct {
	first  _mul
	second loadStashLex
}

func (f *fuseMulLoadStashLex) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseSubStoreStackLexP struct {
	first  _sub
	second storeStackLexP
}

func (f *fuseSubStoreStackLexP) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseMulSub struct {
	first  _mul
	second _sub
}

func (f *fuseMulSub) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseGetElemLoadVal struct {
	first  _getElem
	second loadVal
}

func (f *fuseGetElemLoadVal) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValOpGt struct {
	first  loadVal
	second _op_gt
}

func (f *fuseLoadValOpGt) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadValSub struct {
	first  loadVal
	second _sub
}

func (f *fuseLoadValSub) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

type fuseLoadStackLexMul struct {
	first  loadStackLex
	second _mul
}

func (f *fuseLoadStackLexMul) exec(vm *vm) {
	pc, prg := vm.pc, vm.prg
	f.first.exec(vm)
	if vm.pc == pc+1 && vm.prg == prg {
		f.second.exec(vm)
	}
}

// fusePair 返回 a、b 的合并指令，没有对应的合并指令时返回 nil
func fusePair(a, b instruction) instruction {
	switch a := a.(type) {
	case loadStackLex:
		switch b := b.(type) {
		case loadVal:
			return &fuseLoadStackLexLoadVal{a, b}
		case loadStackLex:
			return &fuseLoadStackLexLoadStackLex{a, b}
		case getPropSource:
			return &fuseLoadStackLexGetPropSource{a, b}
		case call:
			return &fuseLoadStackLexCall{a, b}
		case _getElem:
			return &fuseLoadStackLexGetElem{a, b}
		case storeStackLexP:
			return &fuseLoadStackLexStoreStackLexP{a, b}
		case _bnot:
			return &fuseLoadStackLexBnot{a, b}
		case checkBase:
			return &fuseLoadStackLexCheckBase{a, b}
		case _sal:
			return &fuseLoadStackLexSal{a, b}
		case _add:
			return &fuseLoadStackLexAdd{a, b}
		case _ret:
			return &fuseLoadStackLexRet{a, b}
		case _inc:
			return &fuseLoadStackLexInc{a, b}
		case _xor:
			return &fuseLoadStackLexXor{a, b}
		case jne:
			return &fuseLoadStackLexJne{a, b}
		case _dup:
			return &fuseLoadStackLexDup{a, b}
		case _mul:
			return &fuseLoadStackLexMul{a, b}
		}
	case storeStackLexP:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseStoreStackLexPLoadStackLex{a, b}
		case jump:
			return &fuseStoreStackLexPJump{a, b}
		case loadVal:
			return &fuseStoreStackLexPLoadVal{a, b}
		case loopHead:
			return &fuseStoreStackLexPLoopHead{a, b}
		}
	case jneP:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseJnePLoadStackLex{a, b}
		case loadVal:
			return &fuseJnePLoadVal{a, b}
		case loadStashLex:
			return &fuseJnePLoadStashLex{a, b}
		}
	case loadVal:
		switch b := b.(type) {
		case _op_eq:
			return &fuseLoadValOpEq{a, b}
		case _or:
			return &fuseLoadValOr{a, b}
		case loadStackLex:
			return &fuseLoadValLoadStackLex{a, b}
		case _shr:
			return &fuseLoadValShr{a, b}
		case _and:
			return &fuseLoadValAnd{a, b}
		case _sal:
			return &fuseLoadValSal{a, b}
		case _mod:
			return &fuseLoadValMod{a, b}
		case _op_strict_eq:
			return &fuseLoadValOpStrictEq{a, b}
		case _sar:
			return &fuseLoadValSar{a, b}
		case storeStackLexP:
			return &fuseLoadValStoreStackLexP{a, b}
		case loadStashLex:
			return &fuseLoadValLoadStashLex{a, b}
		case _xor:
			return &fuseLoadValXor{a, b}
		case _op_lt:
			return &fuseLoadValOpLt{a, b}
		case _add:
			return &fuseLoadValAdd{a, b}
		case _op_gt:
			return &fuseLoadValOpGt{a, b}
		case _sub:
			return &fuseLoadValSub{a, b}
		}
	case _op_eq:
		switch b := b.(type) {
		case jneP:
			return &fuseOpEqJneP{a, b}
		}
	case getPropSource:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseGetPropSourceLoadStackLex{a, b}
		case jneP:
			return &fuseGetPropSourceJneP{a, b}
		}
	case loopHead:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseLoopHeadLoadStackLex{a, b}
		}
	case checkBase:
		switch b := b.(type) {
		case getPropCallee:
			return &fuseCheckBaseGetPropCallee{a, b}
		}
	case _or:
		switch b := b.(type) {
		case storeStackLexP:
			return &fuseOrStoreStackLexP{a, b}
		case _sub:
			return &fuseOrSub{a, b}
		case _add:
			return &fuseOrAdd{a, b}
		case loadStackLex:
			return &fuseOrLoadStackLex{a, b}
		}
	case loadStashLex:
		switch b := b.(type) {
		case loadVal:
			return &fuseLoadStashLexLoadVal{a, b}
		case loadStashLex:
			return &fuseLoadStashLexLoadStashLex{a, b}
		case _mul:
			return &fuseLoadStashLexMul{a, b}
		case loadStackLex:
			return &fuseLoadStashLexLoadStackLex{a, b}
		}
	case _op_lt:
		switch b := b.(type) {
		case jneP:
			return &fuseOpLtJneP{a, b}
		}
	case _xor:
		switch b := b.(type) {
		case storeStackLexP:
			return &fuseXorStoreStackLexP{a, b}
		}
	case _add:
		switch b := b.(type) {
		case storeStackLexP:
			return &fuseAddStoreStackLexP{a, b}
		case loadVal:
			return &fuseAddLoadVal{a, b}
		case _xor:
			return &fuseAddXor{a, b}
		}
	case getPropCallee:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseGetPropCalleeLoadStackLex{a, b}
		}
	case _pop:
		switch b := b.(type) {
		case loadStackLex:
			return &fusePopLoadStackLex{a, b}
		case _loadUndef:
			return &fusePopLoadUndef{a, b}
		}
	case _mul:
		switch b := b.(type) {
		case _add:
			return &fuseMulAdd{a, b}
		case loadStashLex:
			return &fuseMulLoadStashLex{a, b}
		case _sub:
			return &fuseMulSub{a, b}
		}
	case jne:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseJneLoadStackLex{a, b}
		case _pop:
			return &fuseJnePop{a, b}
		}
	case _loadUndef:
		switch b := b.(type) {
		case loadStashLex:
			return &fuseLoadUndefLoadStashLex{a, b}
		case _ret:
			return &fuseLoadUndefRet{a, b}
		}
	case _sal:
		switch b := b.(type) {
		case _or:
			return &fuseSalOr{a, b}
		case loadStackLex:
			return &fuseSalLoadStackLex{a, b}
		}
	case _op_strict_eq:
		switch b := b.(type) {
		case jneP:
			return &fuseOpStrictEqJneP{a, b}
		}
	case _getElem:
		switch b := b.(type) {
		case storeStackLexP:
			return &fuseGetElemStoreStackLexP{a, b}
		case loadVal:
			return &fuseGetElemLoadVal{a, b}
		}
	case _bnot:
		switch b := b.(type) {
		case loadVal:
			return &fuseBnotLoadVal{a, b}
		}
	case _sub:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseSubLoadStackLex{a, b}
		case loadVal:
			return &fuseSubLoadVal{a, b}
		case storeStackLexP:
			return &fuseSubStoreStackLexP{a, b}
		}
	case _inc:
		switch b := b.(type) {
		case storeStackLexP:
			return &fuseIncStoreStackLexP{a, b}
		}
	case initStackP:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseInitStackPLoadStackLex{a, b}
		}
	case _shr:
		switch b := b.(type) {
		case loadStackLex:
			return &fuseShrLoadStackLex{a, b}
		case _or:
			return &fuseShrOr{a, b}
		}
	case _dup:
		switch b := b.(type) {
		case getPropSource:
			return &fuseDupGetPropSource{a, b}
		}
	}
	return nil
}

// fuse 把可编译的整数循环换为区域入口，再自左向右把可合并的指令序列替换为超级指令或相邻指令对的合并指令，每条指令最多参与一次合并，合并不跨过区域入口
func (p *Program) fuse() {
	regions := findRegions(p.code)
	covers := func(i, n int) bool {
		for k := i + 1; k < i+n; k++ {
			if regions[k] != nil {
				return true
			}
		}
		return false
	}
	for i := 0; i+1 < len(p.code); i++ {
		if super, n := fuseSuper(p.code, i); super != nil && !covers(i, n) {
			p.code[i] = super
			i += n - 1
			continue
		}
		if fused := fusePair(p.code[i], p.code[i+1]); fused != nil && !covers(i, 2) {
			p.code[i] = fused
			i++
		}
	}
	for h, r := range regions {
		r.fallback = p.code[h]
		p.code[h] = r
	}
}
