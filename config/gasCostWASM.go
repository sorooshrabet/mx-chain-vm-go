// Code generated by elrondapi generator. DO NOT EDIT.

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

package config

type WASMOpcodeCost struct {
	AtomicFence                   uint32
	AtomicNotify                  uint32
	Block                         uint32
	Br                            uint32
	BrIf                          uint32
	BrTable                       uint32
	Call                          uint32
	CallIndirect                  uint32
	Catch                         uint32
	CatchAll                      uint32
	DataDrop                      uint32
	Delegate                      uint32
	Drop                          uint32
	ElemDrop                      uint32
	Else                          uint32
	End                           uint32
	F32Abs                        uint32
	F32Add                        uint32
	F32Ceil                       uint32
	F32Const                      uint32
	F32ConvertI32S                uint32
	F32ConvertI32U                uint32
	F32ConvertI64S                uint32
	F32ConvertI64U                uint32
	F32Copysign                   uint32
	F32DemoteF64                  uint32
	F32Div                        uint32
	F32Eq                         uint32
	F32Floor                      uint32
	F32Ge                         uint32
	F32Gt                         uint32
	F32Le                         uint32
	F32Load                       uint32
	F32Lt                         uint32
	F32Max                        uint32
	F32Min                        uint32
	F32Mul                        uint32
	F32Ne                         uint32
	F32Nearest                    uint32
	F32Neg                        uint32
	F32ReinterpretI32             uint32
	F32Sqrt                       uint32
	F32Store                      uint32
	F32Sub                        uint32
	F32Trunc                      uint32
	F32x4Abs                      uint32
	F32x4Add                      uint32
	F32x4Ceil                     uint32
	F32x4ConvertI32x4S            uint32
	F32x4ConvertI32x4U            uint32
	F32x4DemoteF64x2Zero          uint32
	F32x4Div                      uint32
	F32x4Eq                       uint32
	F32x4ExtractLane              uint32
	F32x4Floor                    uint32
	F32x4Ge                       uint32
	F32x4Gt                       uint32
	F32x4Le                       uint32
	F32x4Lt                       uint32
	F32x4Max                      uint32
	F32x4Min                      uint32
	F32x4Mul                      uint32
	F32x4Ne                       uint32
	F32x4Nearest                  uint32
	F32x4Neg                      uint32
	F32x4PMax                     uint32
	F32x4PMin                     uint32
	F32x4ReplaceLane              uint32
	F32x4Splat                    uint32
	F32x4Sqrt                     uint32
	F32x4Sub                      uint32
	F32x4Trunc                    uint32
	F64Abs                        uint32
	F64Add                        uint32
	F64Ceil                       uint32
	F64Const                      uint32
	F64ConvertI32S                uint32
	F64ConvertI32U                uint32
	F64ConvertI64S                uint32
	F64ConvertI64U                uint32
	F64Copysign                   uint32
	F64Div                        uint32
	F64Eq                         uint32
	F64Floor                      uint32
	F64Ge                         uint32
	F64Gt                         uint32
	F64Le                         uint32
	F64Load                       uint32
	F64Lt                         uint32
	F64Max                        uint32
	F64Min                        uint32
	F64Mul                        uint32
	F64Ne                         uint32
	F64Nearest                    uint32
	F64Neg                        uint32
	F64PromoteF32                 uint32
	F64ReinterpretI64             uint32
	F64Sqrt                       uint32
	F64Store                      uint32
	F64Sub                        uint32
	F64Trunc                      uint32
	F64x2Abs                      uint32
	F64x2Add                      uint32
	F64x2Ceil                     uint32
	F64x2ConvertI64x2S            uint32
	F64x2ConvertI64x2U            uint32
	F64x2ConvertLowI32x4S         uint32
	F64x2ConvertLowI32x4U         uint32
	F64x2Div                      uint32
	F64x2Eq                       uint32
	F64x2ExtractLane              uint32
	F64x2Floor                    uint32
	F64x2Ge                       uint32
	F64x2Gt                       uint32
	F64x2Le                       uint32
	F64x2Lt                       uint32
	F64x2Max                      uint32
	F64x2Min                      uint32
	F64x2Mul                      uint32
	F64x2Ne                       uint32
	F64x2Nearest                  uint32
	F64x2Neg                      uint32
	F64x2PMax                     uint32
	F64x2PMin                     uint32
	F64x2PromoteLowF32x4          uint32
	F64x2ReplaceLane              uint32
	F64x2Splat                    uint32
	F64x2Sqrt                     uint32
	F64x2Sub                      uint32
	F64x2Trunc                    uint32
	GlobalGet                     uint32
	GlobalSet                     uint32
	I16x8Abs                      uint32
	I16x8Add                      uint32
	I16x8AddSatS                  uint32
	I16x8AddSatU                  uint32
	I16x8AddSaturateS             uint32
	I16x8AddSaturateU             uint32
	I16x8AllTrue                  uint32
	I16x8AnyTrue                  uint32
	I16x8Bitmask                  uint32
	I16x8Eq                       uint32
	I16x8ExtAddPairwiseI8x16S     uint32
	I16x8ExtAddPairwiseI8x16U     uint32
	I16x8ExtMulHighI8x16S         uint32
	I16x8ExtMulHighI8x16U         uint32
	I16x8ExtMulLowI8x16S          uint32
	I16x8ExtMulLowI8x16U          uint32
	I16x8ExtendHighI8x16S         uint32
	I16x8ExtendHighI8x16U         uint32
	I16x8ExtendLowI8x16S          uint32
	I16x8ExtendLowI8x16U          uint32
	I16x8ExtractLaneS             uint32
	I16x8ExtractLaneU             uint32
	I16x8GeS                      uint32
	I16x8GeU                      uint32
	I16x8GtS                      uint32
	I16x8GtU                      uint32
	I16x8LeS                      uint32
	I16x8LeU                      uint32
	I16x8Load8x8S                 uint32
	I16x8Load8x8U                 uint32
	I16x8LtS                      uint32
	I16x8LtU                      uint32
	I16x8MaxS                     uint32
	I16x8MaxU                     uint32
	I16x8MinS                     uint32
	I16x8MinU                     uint32
	I16x8Mul                      uint32
	I16x8NarrowI32x4S             uint32
	I16x8NarrowI32x4U             uint32
	I16x8Ne                       uint32
	I16x8Neg                      uint32
	I16x8Q15MulrSatS              uint32
	I16x8ReplaceLane              uint32
	I16x8RoundingAverageU         uint32
	I16x8Shl                      uint32
	I16x8ShrS                     uint32
	I16x8ShrU                     uint32
	I16x8Splat                    uint32
	I16x8Sub                      uint32
	I16x8SubSatS                  uint32
	I16x8SubSatU                  uint32
	I16x8SubSaturateS             uint32
	I16x8SubSaturateU             uint32
	I16x8WidenHighI8x16S          uint32
	I16x8WidenHighI8x16U          uint32
	I16x8WidenLowI8x16S           uint32
	I16x8WidenLowI8x16U           uint32
	I32Add                        uint32
	I32And                        uint32
	I32AtomicLoad                 uint32
	I32AtomicLoad16U              uint32
	I32AtomicLoad8U               uint32
	I32AtomicRmw16AddU            uint32
	I32AtomicRmw16AndU            uint32
	I32AtomicRmw16CmpxchgU        uint32
	I32AtomicRmw16OrU             uint32
	I32AtomicRmw16SubU            uint32
	I32AtomicRmw16XchgU           uint32
	I32AtomicRmw16XorU            uint32
	I32AtomicRmw8AddU             uint32
	I32AtomicRmw8AndU             uint32
	I32AtomicRmw8CmpxchgU         uint32
	I32AtomicRmw8OrU              uint32
	I32AtomicRmw8SubU             uint32
	I32AtomicRmw8XchgU            uint32
	I32AtomicRmw8XorU             uint32
	I32AtomicRmwAdd               uint32
	I32AtomicRmwAnd               uint32
	I32AtomicRmwCmpxchg           uint32
	I32AtomicRmwOr                uint32
	I32AtomicRmwSub               uint32
	I32AtomicRmwXchg              uint32
	I32AtomicRmwXor               uint32
	I32AtomicStore                uint32
	I32AtomicStore16              uint32
	I32AtomicStore8               uint32
	I32AtomicWait                 uint32
	I32Clz                        uint32
	I32Const                      uint32
	I32Ctz                        uint32
	I32DivS                       uint32
	I32DivU                       uint32
	I32Eq                         uint32
	I32Eqz                        uint32
	I32Extend16S                  uint32
	I32Extend8S                   uint32
	I32GeS                        uint32
	I32GeU                        uint32
	I32GtS                        uint32
	I32GtU                        uint32
	I32LeS                        uint32
	I32LeU                        uint32
	I32Load                       uint32
	I32Load16S                    uint32
	I32Load16U                    uint32
	I32Load8S                     uint32
	I32Load8U                     uint32
	I32LtS                        uint32
	I32LtU                        uint32
	I32Mul                        uint32
	I32Ne                         uint32
	I32Or                         uint32
	I32Popcnt                     uint32
	I32ReinterpretF32             uint32
	I32RemS                       uint32
	I32RemU                       uint32
	I32Rotl                       uint32
	I32Rotr                       uint32
	I32Shl                        uint32
	I32ShrS                       uint32
	I32ShrU                       uint32
	I32Store                      uint32
	I32Store16                    uint32
	I32Store8                     uint32
	I32Sub                        uint32
	I32TruncF32S                  uint32
	I32TruncF32U                  uint32
	I32TruncF64S                  uint32
	I32TruncF64U                  uint32
	I32TruncSatF32S               uint32
	I32TruncSatF32U               uint32
	I32TruncSatF64S               uint32
	I32TruncSatF64U               uint32
	I32WrapI64                    uint32
	I32Xor                        uint32
	I32x4Abs                      uint32
	I32x4Add                      uint32
	I32x4AllTrue                  uint32
	I32x4AnyTrue                  uint32
	I32x4Bitmask                  uint32
	I32x4DotI16x8S                uint32
	I32x4Eq                       uint32
	I32x4ExtAddPairwiseI16x8S     uint32
	I32x4ExtAddPairwiseI16x8U     uint32
	I32x4ExtMulHighI16x8S         uint32
	I32x4ExtMulHighI16x8U         uint32
	I32x4ExtMulLowI16x8S          uint32
	I32x4ExtMulLowI16x8U          uint32
	I32x4ExtendHighI16x8S         uint32
	I32x4ExtendHighI16x8U         uint32
	I32x4ExtendLowI16x8S          uint32
	I32x4ExtendLowI16x8U          uint32
	I32x4ExtractLane              uint32
	I32x4GeS                      uint32
	I32x4GeU                      uint32
	I32x4GtS                      uint32
	I32x4GtU                      uint32
	I32x4LeS                      uint32
	I32x4LeU                      uint32
	I32x4Load16x4S                uint32
	I32x4Load16x4U                uint32
	I32x4LtS                      uint32
	I32x4LtU                      uint32
	I32x4MaxS                     uint32
	I32x4MaxU                     uint32
	I32x4MinS                     uint32
	I32x4MinU                     uint32
	I32x4Mul                      uint32
	I32x4Ne                       uint32
	I32x4Neg                      uint32
	I32x4ReplaceLane              uint32
	I32x4Shl                      uint32
	I32x4ShrS                     uint32
	I32x4ShrU                     uint32
	I32x4Splat                    uint32
	I32x4Sub                      uint32
	I32x4TruncSatF32x4S           uint32
	I32x4TruncSatF32x4U           uint32
	I32x4TruncSatF64x2SZero       uint32
	I32x4TruncSatF64x2UZero       uint32
	I32x4WidenHighI16x8S          uint32
	I32x4WidenHighI16x8U          uint32
	I32x4WidenLowI16x8S           uint32
	I32x4WidenLowI16x8U           uint32
	I64Add                        uint32
	I64And                        uint32
	I64AtomicLoad                 uint32
	I64AtomicLoad16U              uint32
	I64AtomicLoad32U              uint32
	I64AtomicLoad8U               uint32
	I64AtomicRmw16AddU            uint32
	I64AtomicRmw16AndU            uint32
	I64AtomicRmw16CmpxchgU        uint32
	I64AtomicRmw16OrU             uint32
	I64AtomicRmw16SubU            uint32
	I64AtomicRmw16XchgU           uint32
	I64AtomicRmw16XorU            uint32
	I64AtomicRmw32AddU            uint32
	I64AtomicRmw32AndU            uint32
	I64AtomicRmw32CmpxchgU        uint32
	I64AtomicRmw32OrU             uint32
	I64AtomicRmw32SubU            uint32
	I64AtomicRmw32XchgU           uint32
	I64AtomicRmw32XorU            uint32
	I64AtomicRmw8AddU             uint32
	I64AtomicRmw8AndU             uint32
	I64AtomicRmw8CmpxchgU         uint32
	I64AtomicRmw8OrU              uint32
	I64AtomicRmw8SubU             uint32
	I64AtomicRmw8XchgU            uint32
	I64AtomicRmw8XorU             uint32
	I64AtomicRmwAdd               uint32
	I64AtomicRmwAnd               uint32
	I64AtomicRmwCmpxchg           uint32
	I64AtomicRmwOr                uint32
	I64AtomicRmwSub               uint32
	I64AtomicRmwXchg              uint32
	I64AtomicRmwXor               uint32
	I64AtomicStore                uint32
	I64AtomicStore16              uint32
	I64AtomicStore32              uint32
	I64AtomicStore8               uint32
	I64AtomicWait                 uint32
	I64Clz                        uint32
	I64Const                      uint32
	I64Ctz                        uint32
	I64DivS                       uint32
	I64DivU                       uint32
	I64Eq                         uint32
	I64Eqz                        uint32
	I64Extend16S                  uint32
	I64Extend32S                  uint32
	I64Extend8S                   uint32
	I64ExtendI32S                 uint32
	I64ExtendI32U                 uint32
	I64GeS                        uint32
	I64GeU                        uint32
	I64GtS                        uint32
	I64GtU                        uint32
	I64LeS                        uint32
	I64LeU                        uint32
	I64Load                       uint32
	I64Load16S                    uint32
	I64Load16U                    uint32
	I64Load32S                    uint32
	I64Load32U                    uint32
	I64Load8S                     uint32
	I64Load8U                     uint32
	I64LtS                        uint32
	I64LtU                        uint32
	I64Mul                        uint32
	I64Ne                         uint32
	I64Or                         uint32
	I64Popcnt                     uint32
	I64ReinterpretF64             uint32
	I64RemS                       uint32
	I64RemU                       uint32
	I64Rotl                       uint32
	I64Rotr                       uint32
	I64Shl                        uint32
	I64ShrS                       uint32
	I64ShrU                       uint32
	I64Store                      uint32
	I64Store16                    uint32
	I64Store32                    uint32
	I64Store8                     uint32
	I64Sub                        uint32
	I64TruncF32S                  uint32
	I64TruncF32U                  uint32
	I64TruncF64S                  uint32
	I64TruncF64U                  uint32
	I64TruncSatF32S               uint32
	I64TruncSatF32U               uint32
	I64TruncSatF64S               uint32
	I64TruncSatF64U               uint32
	I64Xor                        uint32
	I64x2Abs                      uint32
	I64x2Add                      uint32
	I64x2AllTrue                  uint32
	I64x2AnyTrue                  uint32
	I64x2Bitmask                  uint32
	I64x2Eq                       uint32
	I64x2ExtMulHighI32x4S         uint32
	I64x2ExtMulHighI32x4U         uint32
	I64x2ExtMulLowI32x4S          uint32
	I64x2ExtMulLowI32x4U          uint32
	I64x2ExtendHighI32x4S         uint32
	I64x2ExtendHighI32x4U         uint32
	I64x2ExtendLowI32x4S          uint32
	I64x2ExtendLowI32x4U          uint32
	I64x2ExtractLane              uint32
	I64x2GeS                      uint32
	I64x2GtS                      uint32
	I64x2LeS                      uint32
	I64x2Load32x2S                uint32
	I64x2Load32x2U                uint32
	I64x2LtS                      uint32
	I64x2Mul                      uint32
	I64x2Ne                       uint32
	I64x2Neg                      uint32
	I64x2ReplaceLane              uint32
	I64x2Shl                      uint32
	I64x2ShrS                     uint32
	I64x2ShrU                     uint32
	I64x2Splat                    uint32
	I64x2Sub                      uint32
	I64x2TruncSatF64x2S           uint32
	I64x2TruncSatF64x2U           uint32
	I8x16Abs                      uint32
	I8x16Add                      uint32
	I8x16AddSatS                  uint32
	I8x16AddSatU                  uint32
	I8x16AddSaturateS             uint32
	I8x16AddSaturateU             uint32
	I8x16AllTrue                  uint32
	I8x16AnyTrue                  uint32
	I8x16Bitmask                  uint32
	I8x16Eq                       uint32
	I8x16ExtractLaneS             uint32
	I8x16ExtractLaneU             uint32
	I8x16GeS                      uint32
	I8x16GeU                      uint32
	I8x16GtS                      uint32
	I8x16GtU                      uint32
	I8x16LeS                      uint32
	I8x16LeU                      uint32
	I8x16LtS                      uint32
	I8x16LtU                      uint32
	I8x16MaxS                     uint32
	I8x16MaxU                     uint32
	I8x16MinS                     uint32
	I8x16MinU                     uint32
	I8x16Mul                      uint32
	I8x16NarrowI16x8S             uint32
	I8x16NarrowI16x8U             uint32
	I8x16Ne                       uint32
	I8x16Neg                      uint32
	I8x16Popcnt                   uint32
	I8x16ReplaceLane              uint32
	I8x16RoundingAverageU         uint32
	I8x16Shl                      uint32
	I8x16ShrS                     uint32
	I8x16ShrU                     uint32
	I8x16Shuffle                  uint32
	I8x16Splat                    uint32
	I8x16Sub                      uint32
	I8x16SubSatS                  uint32
	I8x16SubSatU                  uint32
	I8x16SubSaturateS             uint32
	I8x16SubSaturateU             uint32
	I8x16Swizzle                  uint32
	If                            uint32
	LocalGet                      uint32
	LocalSet                      uint32
	LocalTee                      uint32
	LocalAllocate                 uint32
	LocalsUnmetered               uint32
	Loop                          uint32
	MaxMemoryGrow                 uint32
	MaxMemoryGrowDelta            uint32
	MemoryAtomicNotify            uint32
	MemoryAtomicWait32            uint32
	MemoryAtomicWait64            uint32
	MemoryCopy                    uint32
	MemoryFill                    uint32
	MemoryGrow                    uint32
	MemoryInit                    uint32
	MemorySize                    uint32
	Nop                           uint32
	RefFunc                       uint32
	RefIsNull                     uint32
	RefNull                       uint32
	Rethrow                       uint32
	Return                        uint32
	ReturnCall                    uint32
	ReturnCallIndirect            uint32
	Select                        uint32
	TableCopy                     uint32
	TableFill                     uint32
	TableGet                      uint32
	TableGrow                     uint32
	TableInit                     uint32
	TableSet                      uint32
	TableSize                     uint32
	Throw                         uint32
	Try                           uint32
	TypedSelect                   uint32
	Unreachable                   uint32
	Unwind                        uint32
	V128And                       uint32
	V128AndNot                    uint32
	V128AnyTrue                   uint32
	V128Bitselect                 uint32
	V128Const                     uint32
	V128Load                      uint32
	V128Load16Lane                uint32
	V128Load16Splat               uint32
	V128Load16x4S                 uint32
	V128Load16x4U                 uint32
	V128Load32Lane                uint32
	V128Load32Splat               uint32
	V128Load32Zero                uint32
	V128Load32x2S                 uint32
	V128Load32x2U                 uint32
	V128Load64Lane                uint32
	V128Load64Splat               uint32
	V128Load64Zero                uint32
	V128Load8Lane                 uint32
	V128Load8Splat                uint32
	V128Load8x8S                  uint32
	V128Load8x8U                  uint32
	V128Not                       uint32
	V128Or                        uint32
	V128Store                     uint32
	V128Store16Lane               uint32
	V128Store32Lane               uint32
	V128Store64Lane               uint32
	V128Store8Lane                uint32
	V128Xor                       uint32
	V16x8LoadSplat                uint32
	V32x4LoadSplat                uint32
	V64x2LoadSplat                uint32
	V8x16LoadSplat                uint32
	V8x16Shuffle                  uint32
	V8x16Swizzle                  uint32
}
