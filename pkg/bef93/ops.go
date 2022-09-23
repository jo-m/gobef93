package bef93

type direction uint8

const (
	dirRight direction = iota
	dirDown
	dirLeft
	dirUp
	dirEND
)

type opcode rune

// https://github.com/catseye/Befunge-93/blob/master/doc/Befunge-93.markdown#appendix-a-command-summary

const (
	opAdd        opcode = '+'  // Addition: Pop a and b, then push a+b
	opSub        opcode = '-'  // Subtraction: Pop a and b, then push b-a
	opMul        opcode = '*'  // Multiplication: Pop a and b, then push a*b
	opDiv        opcode = '/'  // Integer division: Pop a and b, then push b/a, rounded towards 0.
	opMod        opcode = '%'  // Modulo: Pop a and b, then push the remainder of the integer division of b/a.
	opNot        opcode = '!'  // Logical NOT: Pop a value. If the value is zero, push 1; otherwise, push zero.
	opGt         opcode = '`'  // Greater than: Pop a and b, then push 1 if b>a, otherwise zero.
	opRight      opcode = '>'  // Start moving right
	opLeft       opcode = '<'  // Start moving left
	opUp         opcode = '^'  // Start moving up
	opDown       opcode = 'v'  // Start moving down
	opRand       opcode = '?'  // Start moving in a random cardinal direction
	opRif        opcode = '_'  // Pop a value; move right if value=0, left otherwise
	opDif        opcode = '|'  // Pop a value; move down if value=0, up otherwise
	opStr        opcode = '"'  // Start string mode: push each character's ASCII value all the way up to the next "
	opDup        opcode = ':'  // Duplicate value on top of the stack
	opSwp        opcode = '\\' // Swap two values on top of the stack
	opPop        opcode = '$'  // Pop value from the stack and discard it
	opPopWrtInt  opcode = '.'  // Pop value and output as an integer followed by a space
	opPopWrtChr  opcode = ','  // Pop value and output as ASCII character
	opSkip       opcode = '#'  // Bridge: Skip next cell
	opPut        opcode = 'p'  // A "put" call (a way to store a value for later use). Pop y, x, and v, then change the character at (x,y) in the program to the character with ASCII value v
	opGet        opcode = 'g'  // A "get" call (a way to retrieve data in storage). Pop y and x, then push ASCII value of the character at that position in the program
	opReadNr     opcode = '&'  // Ask user for a number and push it
	opReadChr    opcode = '~'  // Ask user for a character and push its ASCII value
	opEnd        opcode = '@'  // End program
	opWhitespace opcode = ' '
)
