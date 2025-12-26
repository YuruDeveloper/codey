package styles

type Symbols struct {
	Pointer  string
	Cursor   string
	Selected string

	Success string
	Error   string
	Warning string
	Pending string
	Bullet  string
	Dot     string

	CheckOn  string
	CheckOff string

	ArrowRight string
	ArrowLeft  string
	ArrowUp    string
	ArrowDown  string

	Ellipsis string
}

var DefaultSymbols Symbols = Symbols{
	Pointer: ">",
	Cursor: "|",
	Selected: ">",

	Success: "v",
	Error: "x",
	Warning: "!",
	Pending: "o",
	Bullet: "*",
	Dot: "*",

	CheckOn: "[v]",
	CheckOff: "[ ]",

	ArrowRight: "->",
	ArrowLeft: "<-",
	ArrowUp: "^",
	ArrowDown: "v",
	
	Ellipsis: "...",
}