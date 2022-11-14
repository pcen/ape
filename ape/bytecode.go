package ape

type Opcode byte

const (
	OpConstanti64 = iota + 1
	OpConstanti32
	OpConstanti16
	OpConstanti8

	OpNil
	OpTrue
	OpFalse

	OpPop

	OpGetLocal
	OpSetLocal

	OpGetGlobal
	OpDefineGlobal
	OpSetGlobal

	OpGetValue
	OpSetValue

	OpGetProperty
	OpSetProperty

	OpGetSuper

	OpEqual
	OpGreater
	OpLess

	OpAdd
	OpSubtract
	OpMultiply
	OpDivide

	OpNot

	OpNegate

	OpPrint

	OpJump
	OpJumpIfFalse
	OpLoop
	OpCall
	OpInvoke
	OpSuperInvoke
	OpClosure
	OpClosureUpValue
	OpReturn
	OpClass
	OpInherit
	OpMethod
)
