package bef93

type direction uint8

const (
	dirRight direction = iota
	dirDown
	dirLeft
	dirUp
	dirEND
)

type opcode byte

// https://github.com/catseye/Befunge-93/blob/master/doc/Befunge-93.markdown#appendix-a-command-summary
const (
	opAdd        = '+'  // Addition: Pop a and b, then push a+b
	opSub        = '-'  // Subtraction: Pop a and b, then push b-a
	opMul        = '*'  // Multiplication: Pop a and b, then push a*b
	opDiv        = '/'  // Integer division: Pop a and b, then push b/a, rounded towards 0.
	opMod        = '%'  // Modulo: Pop a and b, then push the remainder of the integer division of b/a.
	opNot        = '!'  // Logical NOT: Pop a value. If the value is zero, push 1; otherwise, push zero.
	opGt         = '`'  // Greater than: Pop a and b, then push 1 if b>a, otherwise zero.
	opRight      = '>'  // Start moving right
	opLeft       = '<'  // Start moving left
	opUp         = '^'  // Start moving up
	opDown       = 'v'  // Start moving down
	opRand       = '?'  // Start moving in a random cardinal direction
	opRif        = '_'  // Pop a value; move right if value=0, left otherwise
	opDif        = '|'  // Pop a value; move down if value=0, up otherwise
	opStr        = '"'  // Start string mode: push each character's ASCII value all the way up to the next "
	opDup        = ':'  // Duplicate value on top of the stack
	opSwp        = '\\' // Swap two values on top of the stack
	opPop        = '$'  // Pop value from the stack and discard it
	opPopWrtInt  = '.'  // Pop value and output as an integer followed by a space
	opPopWrtChr  = ','  // Pop value and output as ASCII character
	opSkip       = '#'  // Bridge: Skip next cell
	opPut        = 'p'  // A "put" call (a way to store a value for later use). Pop y, x, and v, then change the character at (x,y) in the program to the character with ASCII value v
	opGet        = 'g'  // A "get" call (a way to retrieve data in storage). Pop y and x, then push ASCII value of the character at that position in the program
	opReadNr     = '&'  // Ask user for a number and push it
	opReadChr    = '~'  // Ask user for a character and push its ASCII value
	opEnd        = '@'  // End program
	opWhitespace = ' '
	opEND
)
