package OpCode

// OPCODE 操作指令类型
type OPCODE byte

// OpCode 脚本指令代码
const (
	PUSH0       OPCODE = 0x00 // An empty array of bytes is pushed onto the stack.
	PUSHF       OPCODE = PUSH0
	PUSHBYTES1  OPCODE = 0x01 // 0x01-0x4B The next opcode bytes is data to be pushed onto the stack
	PUSHBYTES75 OPCODE = 0x4B
	PUSHDATA1   OPCODE = 0x4C // The next OPCODE contains the number of bytes to be pushed onto the stack.
	PUSHDATA2   OPCODE = 0x4D // The next two bytes contain the number of bytes to be pushed onto the stack.
	PUSHDATA4   OPCODE = 0x4E // The next four bytes contain the number of bytes to be pushed onto the stack.
	PUSHM1      OPCODE = 0x4F // The number -1 is pushed onto the stack.
	PUSH1       OPCODE = 0x51 // The number 1 is pushed onto the stack.
	PUSHT       OPCODE = PUSH1
	PUSH2       OPCODE = 0x52 // The number 2 is pushed onto the stack.
	PUSH3       OPCODE = 0x53 // The number 3 is pushed onto the stack.
	PUSH4       OPCODE = 0x54 // The number 4 is pushed onto the stack.
	PUSH5       OPCODE = 0x55 // The number 5 is pushed onto the stack.
	PUSH6       OPCODE = 0x56 // The number 6 is pushed onto the stack.
	PUSH7       OPCODE = 0x57 // The number 7 is pushed onto the stack.
	PUSH8       OPCODE = 0x58 // The number 8 is pushed onto the stack.
	PUSH9       OPCODE = 0x59 // The number 9 is pushed onto the stack.
	PUSH10      OPCODE = 0x5A // The number 10 is pushed onto the stack.
	PUSH11      OPCODE = 0x5B // The number 11 is pushed onto the stack.
	PUSH12      OPCODE = 0x5C // The number 12 is pushed onto the stack.
	PUSH13      OPCODE = 0x5D // The number 13 is pushed onto the stack.
	PUSH14      OPCODE = 0x5E // The number 14 is pushed onto the stack.
	PUSH15      OPCODE = 0x5F // The number 15 is pushed onto the stack.
	PUSH16      OPCODE = 0x60 // The number 16 is pushed onto the stack.

	// Flow control
	NOP      OPCODE = 0x61 // Does nothing.
	JMP      OPCODE = 0x62
	JMPIF    OPCODE = 0x63
	JMPIFNOT OPCODE = 0x64
	CALL     OPCODE = 0x65
	RET      OPCODE = 0x66
	APPCALL  OPCODE = 0x67
	SYSCALL  OPCODE = 0x68
	TAILCALL OPCODE = 0x69

	// Stack
	DUPFROMALTSTACK OPCODE = 0x6A
	TOALTSTACK      OPCODE = 0x6B // Puts the input onto the top of the alt stack. Removes it from the main stack.
	FROMALTSTACK    OPCODE = 0x6C // Puts the input onto the top of the main stack. Removes it from the alt stack.
	XDROP           OPCODE = 0x6D
	XSWAP           OPCODE = 0x72
	XTUCK           OPCODE = 0x73
	DEPTH           OPCODE = 0x74 // Puts the number of stack items onto the stack.
	DROP            OPCODE = 0x75 // Removes the top stack item.
	DUP             OPCODE = 0x76 // Duplicates the top stack item.
	NIP             OPCODE = 0x77 // Removes the second-to-top stack item.
	OVER            OPCODE = 0x78 // Copies the second-to-top stack item to the top.
	PICK            OPCODE = 0x79 // The item n back in the stack is copied to the top.
	ROLL            OPCODE = 0x7A // The item n back in the stack is moved to the top.
	ROT             OPCODE = 0x7B // The top three items on the stack are rotated to the left.
	SWAP            OPCODE = 0x7C // The top two items on the stack are swapped.
	TUCK            OPCODE = 0x7D // The item at the top of the stack is copied and inserted before the second-to-top item.

	// Splice
	CAT    OPCODE = 0x7E // Concatenates two strings.
	SUBSTR OPCODE = 0x7F // Returns a section of a string.
	LEFT   OPCODE = 0x80 // Keeps only characters left of the specified point in a string.
	RIGHT  OPCODE = 0x81 // Keeps only characters right of the specified point in a string.
	SIZE   OPCODE = 0x82 // Returns the length of the input string.

	// Bitwise logic
	INVERT OPCODE = 0x83 // Flips all of the bits in the input.
	AND    OPCODE = 0x84 // Boolean and between each bit in the inputs.
	OR     OPCODE = 0x85 // Boolean or between each bit in the inputs.
	XOR    OPCODE = 0x86 // Boolean exclusive or between each bit in the inputs.
	EQUAL  OPCODE = 0x87 // Returns 1 if the inputs are exactly equal 0 otherwise.
	//OP_EQUALVERIFY OPCODE = 0x88 // Same as OP_EQUAL but runs OP_VERIFY afterward.
	//OP_RESERVED1 OPCODE = 0x89 // Transaction is invalid unless occuring in an unexecuted OP_IF branch
	//OP_RESERVED2 OPCODE = 0x8A // Transaction is invalid unless occuring in an unexecuted OP_IF branch

	// Arithmetic
	// Note: Arithmetic inputs are limited to signed 32-bit integers but may overflow their output.
	INC         OPCODE = 0x8B // 1 is added to the input.
	DEC         OPCODE = 0x8C // 1 is subtracted from the input.
	SIGN        OPCODE = 0x8D
	NEGATE      OPCODE = 0x8F // The sign of the input is flipped.
	ABS         OPCODE = 0x90 // The input is made positive.
	NOT         OPCODE = 0x91 // If the input is 0 or 1 it is flipped. Otherwise the output will be 0.
	NZ          OPCODE = 0x92 // Returns 0 if the input is 0. 1 otherwise.
	ADD         OPCODE = 0x93 // a is added to b.
	SUB         OPCODE = 0x94 // b is subtracted from a.
	MUL         OPCODE = 0x95 // a is multiplied by b.
	DIV         OPCODE = 0x96 // a is divided by b.
	MOD         OPCODE = 0x97 // Returns the remainder after dividing a by b.
	SHL         OPCODE = 0x98 // Shifts a left b bits preserving sign.
	SHR         OPCODE = 0x99 // Shifts a right b bits preserving sign.
	BOOLAND     OPCODE = 0x9A // If both a and b are not 0 the output is 1. Otherwise 0.
	BOOLOR      OPCODE = 0x9B // If a or b is not 0 the output is 1. Otherwise 0.
	NUMEQUAL    OPCODE = 0x9C // Returns 1 if the numbers are equal 0 otherwise.
	NUMNOTEQUAL OPCODE = 0x9E // Returns 1 if the numbers are not equal 0 otherwise.
	LT          OPCODE = 0x9F // Returns 1 if a is less than b 0 otherwise.
	GT          OPCODE = 0xA0 // Returns 1 if a is greater than b 0 otherwise.
	LTE         OPCODE = 0xA1 // Returns 1 if a is less than or equal to b 0 otherwise.
	GTE         OPCODE = 0xA2 // Returns 1 if a is greater than or equal to b 0 otherwise.
	MIN         OPCODE = 0xA3 // Returns the smaller of a and b.
	MAX         OPCODE = 0xA4 // Returns the larger of a and b.
	WITHIN      OPCODE = 0xA5 // Returns 1 if x is within the specified range (left-inclusive) 0 otherwise.

	// Crypto
	//RIPEMD160 OPCODE = 0xA6 // The input is hashed using RIPEMD-160.
	SHA1    OPCODE = 0xA7 // The input is hashed using SHA-1.
	SHA256  OPCODE = 0xA8 // The input is hashed using SHA-256.
	HASH160 OPCODE = 0xA9
	HASH256 OPCODE = 0xAA
	//因为这个hash函数可能仅仅是csharp 编译时专用的
	CSHARPSTRHASH32 OPCODE = 0xAB
	//这个是JAVA专用的
	JAVAHASH32 OPCODE = 0xAD

	CHECKSIG      OPCODE = 0xAC
	CHECKMULTISIG OPCODE = 0xAE

	// Array
	ARRAYSIZE OPCODE = 0xC0
	PACK      OPCODE = 0xC1
	UNPACK    OPCODE = 0xC2
	PICKITEM  OPCODE = 0xC3
	SETITEM   OPCODE = 0xC4
	NEWARRAY  OPCODE = 0xC5 //用作引用類型
	NEWSTRUCT OPCODE = 0xC6 //用作值類型

	SWITCH OPCODE = 0xD0

	// Exceptions
	THROW      OPCODE = 0xF0
	THROWIFNOT OPCODE = 0xF1
)
