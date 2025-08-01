package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"regexp"
	_ "embed"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var myPrompt = ":) "
//go:embed README.md
var readmeContent string

var mainWordRegex = regexp.MustCompile(`:\s+main\s+[^;]*;`) // Regex to detect ": main ... ;"

var ErrBreak = fmt.Errorf("break")
var ErrContinue = fmt.Errorf("continue")

// Error represents a custom error with a code and message.
type Error struct {
	Code    int
	Message string
}

// Predefined errors
var errors = []Error{
	{Code: 1, Message: "stack underflow"},
	{Code: 2, Message: "division by zero"},
	{Code: 3, Message: "type error: expected a number, got %T"},
	{Code: 4, Message: "type error: expected a boolean, got %T"},
	{Code: 5, Message: "type error: expected a string, got %T"},
	{Code: 6, Message: "type error: expected a code block, got %T"},
	{Code: 7, Message: "type error: '+' requires two numbers or two strings, got %T and %T"},
	{Code: 8, Message: "type error: '+' requires two numbers or two strings, got %T"},
	{Code: 9, Message: "invalid arguments for if"},
	{Code: 10, Message: "index: not inside a loop"},
	{Code: 11, Message: "variable name '%s' conflicts with an existing command"},
	{Code: 12, Message: "variable names starting with '_' are reserved for internal use: %s"},
	{Code: 13, Message: "internal variable %s can only be set to a boolean value"},
	{Code: 14, Message: "cannot modify internal variable: %s"},
	{Code: 15, Message: "undefined variable: %s"},
	{Code: 16, Message: "invalid function definition"},
	{Code: 17, Message: "word names starting with '_' are reserved for internal use: %s"},
	{Code: 18, Message: "word name '%s' conflicts with an existing command"},
	{Code: 19, Message: "delete: missing variable or word name"},
	{Code: 20, Message: "cannot delete internal variable or word: %s"},
	{Code: 21, Message: "undefined variable or word: %s"},
	{Code: 22, Message: "undefined %s: %s"},
	{Code: 23, Message: "see: missing variable or word name"},
	{Code: 24, Message: "edit: missing name"},
	{Code: 25, Message: "edit: undefined word or variable: %s"},
	{Code: 26, Message: "edit: word '%s' not found"},
	{Code: 27, Message: "edit: variable '%s' is not a code block"},
	{Code: 28, Message: "edit: variable '%s' not found"},
	{Code: 29, Message: "error executing '%s'\n%w"},
	{Code: 30, Message: "undefined word: %s"},
	{Code: 31, Message: "error executing word '%s'\n%w"},
	{Code: 32, Message: "variable '%s' contains non-string elements in block"},
	{Code: 33, Message: "error executing variable as word '%s'\n%w"},
	{Code: 34, Message: "unrecognized token: %s"},
	{Code: 35, Message: "unmatched ')'"},
	{Code: 36, Message: "unmatched quote"},
	{Code: 37, Message: "unmatched '('"},
	{Code: 38, Message: "expected '{' to start block"},
	{Code: 39, Message: "unmatched '{'"},
	{Code: 40, Message: "failed to get home directory\n%w"},
	{Code: 41, Message: "failed to create .rpn directory\n%w"},
	{Code: 42, Message: "failed to marshal interpreter state\n%w"},
	{Code: 43, Message: "failed to write state to file %s\n%w"},
	{Code: 44, Message: "failed to read state from file %s\n%w"},
	{Code: 45, Message: "failed to unmarshal interpreter state\n%w"},
	{Code: 46, Message: "failed to read RPN file %s\n%w"},
	{Code: 47, Message: "failed to open file %s for export\n%w"},
	{Code: 48, Message: "failed to read RPN directory %s\n%w"},
	{Code: 49, Message: "while: condition must evaluate to a boolean or number, got %T"},
	{Code: 50, Message: "'%s' is low level defined into the core"},
	{Code: 51, Message: "execution interrupted by user"},
	{Code: 52, Message: "factorial: expected a non-negative integer, got %v"},
	{Code: 53, Message: "editfile: missing filename"},
	{Code: 54, Message: "editfile: file not found: %s"},
	{Code: 55, Message: "variable names cannot contain spaces: '%s'"},
	{Code: 56, Message: "cannot define local variable '%s' in global scope"},
	{Code: 57, Message: "semicolon out of context"},
	{Code: 58, Message: "invalid character input: %s"},
	{Code: 59, Message: "string bounds out of range"},
}

// History variables
var history []string
var historyIndex int = -1
var historyFile = ".rpn_history"
var rpnDir = ".polish"

const majorVersion = "0"
const appName = "Polish"

var version string // This will be set by ldflags during build

var _version float64 // Internal read-only variable for version as float

func init() {
	// Parse the 'version' string which is expected to be in "MAJOR.COMMITS-HASH" format
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		majorStr := parts[0]
		commitsAndHash := parts[1]

		major, err := strconv.ParseFloat(majorStr, 64)
		if err != nil {
			_version = 0.0 // Default or error value
			return
		}

		commitParts := strings.Split(commitsAndHash, "-")
		var commitsStr string
		if len(commitParts) >= 1 {
			commitsStr = commitParts[0]
		} else {
			_version = major // If commitCount is not available, use major version only
			return
		}

		commits, err := strconv.ParseFloat(commitsStr, 64)
		if err != nil {
			_version = major // If parsing fails, use major version only
			return
		}

		// Combine major and commit count as 0.XX
		_version = major + (commits / 100.0)
	} else {
		// Fallback if version string is not in the expected format
		// Use the majorVersion constant as a base
		major, err := strconv.ParseFloat(majorVersion, 64)
		if err != nil {
			_version = 0.0
		} else {
			_version = major
		}
	}
}

// loadHistory loads command history from the history file.
func loadHistory() {
	home, err := os.UserHomeDir()
	if err != nil {
		return // Don't worry if we can't find home
	}
	rpnPath := filepath.Join(home, rpnDir)
	path := filepath.Join(rpnPath, historyFile)
	file, err := os.Open(path)
	if err != nil {
		return // It's okay if the file doesn't exist yet
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
	historyIndex = len(history)
}

// saveHistory saves command history to the history file.
func saveHistory() {
	if len(history) == 0 {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return // Can't save if we can't find home
	}
	rpnPath := filepath.Join(home, rpnDir)
	// Ensure the .rpn directory exists
	if err := os.MkdirAll(rpnPath, 0755); err != nil {
		return // Can't save if we can't create the directory
	}
	path := filepath.Join(rpnPath, historyFile)
	file, err := os.Create(path)
	if err != nil {
		return // Can't save if we can't create the file
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range history {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}

// SetFocusables sets the focusable elements for the interpreter.
func (i *Interpreter) SetFocusables(elements ...tview.Primitive) {
	i.focusables = elements
	i.currentFocus = 0 // Default to the first element
}

// CycleFocus moves the focus to the next element in the focusables slice.
func (i *Interpreter) CycleFocus() {
	if len(i.focusables) == 0 {
		return
	}

	// Helper function to set border color
	setBorderColor := func(p tview.Primitive, color tcell.Color) {
		switch v := p.(type) {
		case *tview.InputField:
			v.SetBorderColor(color)
		case *tview.TextView:
			v.SetBorderColor(color)
		case *tview.Table:
			v.SetBorderColor(color)
		case *tview.TextArea:
			v.SetBorderColor(color)
		}
	}

	// Reset border of the old focused item
	setBorderColor(i.focusables[i.currentFocus], tview.Styles.BorderColor)

	i.currentFocus = (i.currentFocus + 1) % len(i.focusables)
	i.app.SetFocus(i.focusables[i.currentFocus])

	// Set border of the new focused item
	setBorderColor(i.focusables[i.currentFocus], tcell.ColorYellow)
}

// Interpreter holds the state of our RPN calculator.
type Interpreter struct {
	stack       []interface{}
	opcodes     map[string]func(*Interpreter) error
	variables   map[string]interface{} // Global variables
	scopeStack  []map[string]interface{}
	words       map[string][]string
	interrupted chan struct{}

	outputView      io.Writer
	angleModeView   *tview.TextView   // New field for angle mode display
	clockView       *tview.TextView   // New field for clock display
	variablesTable  *tview.Table      // New field for variables display
	stackTable      *tview.Table      // New field for stack display
	wordsTable      *tview.Table      // New field for words display
	suggestions     []string          // New field for tab completion suggestions
	suggestionIndex int               // New field for current suggestion index
	inputField      *tview.InputField // New field for input field access
	loopIndex       float64           // New field to store current loop index
	clrEdit         bool              // Flag to clear or not the inputText field

	editingFile     bool   // Flag to indicate if currently in file editing mode
	currentEditFile string // Stores the path of the file being edited
	fileDirty       bool   // Flag to indicate if the file has been modified
	app             *tview.Application // Reference to the tview application
	appFlex         *tview.Flex        // Reference to the main application flex layout

	originalPrompt string      // Stores the original prompt of the input field
	inputChan      chan string // Channel to communicate input from prompt command
	promptActive   bool        // Flag to indicate if prompt command is active
	promptMutex    sync.Mutex  // Mutex to protect promptActive

	focusables      []tview.Primitive // Slice to hold focusable elements
	currentFocus    int               // Index of the currently focused element
}

// newError creates a new error with a code and formatted message.
func (i *Interpreter) newError(code int, args ...interface{}) error {
	i.variables["_error"] = true
	for _, e := range errors {
		if e.Code == code {
			i.variables["_last_error"] = float64(code)
			return fmt.Errorf(fmt.Sprintf("[red]error %d: %s[white]\n", e.Code, e.Message), args...)
		}
	}
	return fmt.Errorf("[red]unknown error code: %d[white]\n", code)
}

// InterpreterState represents the savable state of the interpreter.
type InterpreterState struct {
	Stack     []interface{}          `json:"stack"`
	Variables map[string]interface{} `json:"variables"`
	Words     map[string][]string    `json:"words"`
}

// saveState saves the current interpreter state to a file.
func (i *Interpreter) saveState(filename string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return i.newError(40, err)
	}
	rpnPath := filepath.Join(home, rpnDir)
	if err := os.MkdirAll(rpnPath, 0755); err != nil {
		return i.newError(41, err)
	}
	fullPath := filepath.Join(rpnPath, filename)

	state := InterpreterState{
		Stack:     i.stack,
		Variables: i.variables,
		Words:     i.words,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return i.newError(42, err)
	}

	err = ioutil.WriteFile(fullPath, data, 0644)
	if err != nil {
		return i.newError(43, fullPath, err)
	}
	return nil
}

func updateAngleAndEchoModeView(i *Interpreter) {
	if i.angleModeView == nil {
		return
	}
	i.angleModeView.Clear()
	mode := "RAD"
	if val, ok := i.variables["_degree_mode"].(bool); ok && val {
		mode = "DEG"
	}
	echoStatus := "OFF"
	if val, ok := i.variables["_echo_mode"].(bool); ok && val {
		echoStatus = "ON"
	}
	fmt.Fprintf(i.angleModeView, "%s | ECHO %s", mode, echoStatus)
}

// updateClockView updates the clock display with the current time.
func updateClockView(i *Interpreter) {
	if i.clockView == nil {
		return
	}
	currentTime := time.Now().Format("15:04:05") // HH:MM:SS
	i.clockView.SetText(currentTime)
}

// loadState loads the interpreter state from a file.
func (i *Interpreter) loadState(filename string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return i.newError(40, err)
	}
	rpnPath := filepath.Join(home, rpnDir)
	fullPath := filepath.Join(rpnPath, filename)

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return i.newError(44, fullPath, err)
	}

	var state InterpreterState
	err = json.Unmarshal(data, &state)
	if err != nil {
		return i.newError(45, err)
	}

	i.stack = state.Stack
	// Merge loaded variables into existing ones, preserving internal variables
	for k, v := range state.Variables {
		// Only overwrite if it's not an internal variable (starts with '_')
		if !strings.HasPrefix(k, "_") {
			i.variables[k] = v
		}
	}
	if state.Words == nil {
		i.words = make(map[string][]string)
	} else {
		i.words = state.Words
	}
	// Update the angle mode display after loading state
	updateAngleAndEchoModeView(i)
	// Update the variables view after loading state
	showVarsValue := true
	if val, ok := i.variables["_vars_value"].(bool); ok {
		showVarsValue = val
	}
	hideInternalVars := true
	if val, ok := i.variables["_hidden_vars"].(bool); ok {
		hideInternalVars = val
	}
	if i.variablesTable != nil {
		if outputTextView, ok := i.outputView.(*tview.TextView); ok {
			updateVariablesView(i.variablesTable, i.variables, showVarsValue, hideInternalVars, outputTextView)
		}
	}
	// Update the stack view after loading state
	showStackType := false
	if val, ok := i.variables["_stack_type"].(bool); ok {
		showStackType = val
	}
	if i.stackTable != nil {
		updateStackView(i.stackTable, i.stack, showStackType)
	}
	return nil
}

// NewInterpreter creates a new interpreter instance with all opcodes registered.
func NewInterpreter(app *tview.Application, appFlex *tview.Flex, outputView io.Writer, angleModeView *tview.TextView, clockView *tview.TextView, variablesTable *tview.Table, stackTable *tview.Table, inputField *tview.InputField) *Interpreter {
	interp := &Interpreter{
		stack:       make([]interface{}, 0),
		opcodes:     make(map[string]func(*Interpreter) error),
		variables:   make(map[string]interface{}),
		scopeStack:  make([]map[string]interface{}, 0),
		words:       make(map[string][]string),
		interrupted: make(chan struct{}, 1),

		outputView:      outputView,
		angleModeView:   angleModeView,
		clockView:       clockView,
		variablesTable:  variablesTable,
		stackTable:      stackTable,
		wordsTable:      tview.NewTable().SetBorders(false), // Initialize wordsTable
		suggestions:     []string{},                         // Initialize empty suggestions
		suggestionIndex: -1,                                 // No suggestion selected initially
		inputField:      inputField,
		app:             app,
		appFlex:         appFlex,
		originalPrompt:  inputField.GetLabel(), // Store the initial prompt
		inputChan:       make(chan string, 1), // Initialize the channel
		promptActive:    false,                // Not active initially
	}
	interp.scopeStack = append(interp.scopeStack, interp.variables) // Add global scope
	interp.clrEdit = true
	interp.variables["_echo_mode"] = true
	interp.variables["_degree_mode"] = false
	interp.variables["_vars_value"] = true
	interp.variables["_stack_type"] = false
	interp.variables["_hidden_vars"] = false
	interp.variables["_exit_save"] = false
	interp.variables["_last_error"] = float64(0)
	interp.variables["_error"] = false
	interp.variables["_last_x"] = nil // Initialize _last_x
	interp.loopIndex = -1 // Initialize loop index to -1 (no active loop)

	// Add _version to internal variables
	interp.variables["_version"] = _version

	interp.registerOpcodes()
	return interp
}

// push adds a value to the stack.
func (i *Interpreter) push(v interface{}) {
	i.stack = append(i.stack, v)
}

// pop removes and returns a value from the stack.
func (i *Interpreter) pop() (interface{}, error) {
	if len(i.stack) == 0 {
		return 0, i.newError(1)
	}
	val := i.stack[len(i.stack)-1]
	i.stack = i.stack[:len(i.stack)-1]
	i.variables["_last_x"] = val // Store the popped value in _last_x
	return val, nil
}

// popFloat pops a value and asserts it's a float64.
func (i *Interpreter) popFloat() (float64, error) {
	val, err := i.pop()
	if err != nil {
		return 0, err
	}
	f, ok := val.(float64)
	if !ok {
		// Allow bool to be converted to float64
		if b, ok := val.(bool); ok {
			if b {
				return 1.0, nil
			}
			return 0.0, nil
		}
		return 0, i.newError(3, val)
	}
	return f, nil
}

// popBool pops a value and asserts it's a bool.
func (i *Interpreter) popBool() (bool, error) {
	val, err := i.pop()
	if err != nil {
		return false, err
	}
	b, ok := val.(bool)
	if !ok {
		// Allow float64 to be converted to bool
		if f, ok := val.(float64); ok {
			return f != 0, nil
		}
		return false, i.newError(4, val)
	}
	return b, nil
}

// popString pops a value and asserts it's a string.
func (i *Interpreter) popString() (string, error) {
	val, err := i.pop()
	if err != nil {
		return "", err
	}
	s, ok := val.(string)
	if !ok {
		return "", i.newError(5, val)
	}
	return s, nil
}

// popBlock pops a value and asserts it's a code block ([]string).
func (i *Interpreter) popBlock() ([]string, error) {
	val, err := i.pop()
	if err != nil {
		return nil, err
	}
	block, ok := val.([]string)
	if !ok {
		return nil, i.newError(6, val)
	}
	return block, nil
}

// registerOpcodes maps string commands to their functions.
func (i *Interpreter) registerOpcodes() {
	// Null words treated directly into the parser
	i.opcodes[":"] = nil
	i.opcodes[";"] = nil
	i.opcodes["delete"] = nil
	i.opcodes["edit"] = nil
	i.opcodes["editfile"] = nil
	i.opcodes["see"] = nil
	i.opcodes["("] = nil
	i.opcodes[")"] = nil
	i.opcodes["{"] = nil
	i.opcodes["}"] = nil
	i.opcodes["\""] = nil

	// Arithmetic & String Concat
	i.opcodes["+"] = func(i *Interpreter) error {
		b, err := i.pop()
		if err != nil {
			return err
		}
		a, err := i.pop()
		if err != nil {
			return err
		}

		switch aVal := a.(type) {
		case float64:
			if bVal, ok := b.(float64); ok {
				i.push(aVal + bVal)
			} else {
				return i.newError(7, a, b)
			}
		case string:
			if bVal, ok := b.(string); ok {
				i.push(aVal + bVal)
			} else {
				return i.newError(7, a, b)
			}
		default:
			return i.newError(8, a)
		}
		return nil
	}

	i.opcodes["-"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a - b)
		return nil
	}
	i.opcodes["*"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a * b)
		return nil
	}
	i.opcodes["/"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		if b == 0 {
			return i.newError(2)
		}
		i.push(a / b)
		return nil
	}
			i.opcodes["mod"] = func(i *Interpreter) error {
			b, err := i.popFloat()
			if err != nil {
				return err
			}
			a, err := i.popFloat()
			if err != nil {
				return err
			}
			i.push(math.Mod(a, b))
			return nil
		}

		i.opcodes["%"] = func(i *Interpreter) error {
			total, err := i.popFloat()
			if err != nil {
				return err
			}
			value, err := i.popFloat()
			if err != nil {
				return err
			}
			if total == 0 {
				return i.newError(2) // Division by zero
			}
			i.push((value / total) * 100)
			return nil
		}

	// Math functions
	i.opcodes["sqrt"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Sqrt(a))
		return nil
	}
	i.opcodes["pow"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Pow(a, b))
		return nil
	}
	i.opcodes["nroot"] = func(i *Interpreter) error {
		n, err := i.popFloat()
		if err != nil {
			return err
		}
		x, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Pow(x, 1/n))
		return nil
	}
	i.opcodes["sq"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a * a)
		return nil
	}
	i.opcodes["sin"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to radians
			a = a * math.Pi / 180
		}
		i.push(math.Sin(a))
		return nil
	}
	i.opcodes["cos"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to radians
			a = a * math.Pi / 180
		}
		i.push(math.Cos(a))
		return nil
	}
	i.opcodes["tan"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to radians
			a = a * math.Pi / 180
		}
		i.push(math.Tan(a))
		return nil
	}
	i.opcodes["log"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Log10(a))
		return nil
	}
	i.opcodes["pow10"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Pow(10, a))
		return nil
	}
	i.opcodes["exp"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Exp(a))
		return nil
	}
	i.opcodes["ln"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Log(a))
		return nil
	}

	i.opcodes["factorial"] = func(i *Interpreter) error {
		val, err := i.popFloat()
		if err != nil {
			return err
		}
		if val < 0 || val != math.Trunc(val) {
			return i.newError(52, val)
		}
		result := float64(1)
		for k := float64(1); k <= val; k++ {
			result *= k
		}
		i.push(result)
		return nil
	}
	i.opcodes["int"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(math.Trunc(a))
		return nil
	}
	i.opcodes["frac"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a - math.Trunc(a))
		return nil
	}
	i.opcodes["asin"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		res := math.Asin(a)
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to degrees
			res = res * 180 / math.Pi
		}
		i.push(res)
		return nil
	}
	i.opcodes["acos"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		res := math.Acos(a)
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to degrees
			res = res * 180 / math.Pi
		}
		i.push(res)
		return nil
	}
	i.opcodes["atan"] = func(i *Interpreter) error {
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		res := math.Atan(a)
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to degrees
			res = res * 180 / math.Pi
		}
		i.push(res)
		return nil
	}
	i.opcodes["atan2"] = func(i *Interpreter) error {
		y, err := i.popFloat()
		if err != nil {
			return err
		}
		x, err := i.popFloat()
		if err != nil {
			return err
		}
		res := math.Atan2(y, x)
		if val, ok := i.variables["_degree_mode"].(bool); ok && val { // If in degrees mode, convert to degrees
			res = res * 180 / math.Pi
		}
		i.push(res)
		return nil
	}

	// Stack manipulation
	i.opcodes["dup"] = func(i *Interpreter) error {
		a, err := i.pop()
		if err != nil {
			return err
		}
		i.push(a)
		i.push(a)
		return nil
	}
	i.opcodes["drop"] = func(i *Interpreter) error {
		_, err := i.pop()
		return err
	}
	i.opcodes["swap"] = func(i *Interpreter) error {
		b, err := i.pop()
		if err != nil {
			return err
		}
		a, err := i.pop()
		if err != nil {
			return err
		}
		i.push(b)
		i.push(a)
		return nil
	}

	i.opcodes["rot"] = func(i *Interpreter) error {
		if len(i.stack) < 3 {
			return i.newError(1) // Stack underflow
		}
		c, err := i.pop()
		if err != nil {
			return err
		}
		b, err := i.pop()
		if err != nil {
			return err
		}
		a, err := i.pop()
		if err != nil {
			return err
		}
		i.push(b)
		i.push(c)
		i.push(a)
		return nil
	}

	// Comparison operators
	i.opcodes["=="] = func(i *Interpreter) error {
		b, err := i.pop()
		if err != nil {
			return err
		}
		a, err := i.pop()
		if err != nil {
			return err
		}
		equal := false
		switch aVal := a.(type) {
		case float64:
			if bVal, ok := b.(float64); ok {
				equal = aVal == bVal
			}
		case string:
			if bVal, ok := b.(string); ok {
				equal = aVal == bVal
			}
		case bool:
			if bVal, ok := b.(bool); ok {
				equal = aVal == bVal
			}
		}
		i.push(equal)
		return nil
	}
	i.opcodes["!="] = func(i *Interpreter) error {
		b, err := i.pop()
		if err != nil {
			return err
		}
		a, err := i.pop()
		if err != nil {
			return err
		}
		equal := false
		switch aVal := a.(type) {
		case float64:
			if bVal, ok := b.(float64); ok {
				equal = aVal == bVal
			}
		case string:
			if bVal, ok := b.(string); ok {
				equal = aVal == bVal
			}
		case bool:
			if bVal, ok := b.(bool); ok {
				equal = aVal == bVal
			}
		}
		i.push(!equal)
		return nil
	}
	i.opcodes[">"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a > b)
		return nil
	}
	i.opcodes["<"] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a < b)
		return nil
	}
	i.opcodes[">="] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a >= b)
		return nil
	}
	i.opcodes["<="] = func(i *Interpreter) error {
		b, err := i.popFloat()
		if err != nil {
			return err
		}
		a, err := i.popFloat()
		if err != nil {
			return err
		}
		i.push(a <= b)
		return nil
	}

	// Control flow
	i.opcodes["if"] = func(i *Interpreter) error {
		// Pop the 'then' or 'else' block first.
		block1, err := i.popBlock()
		if err != nil {
			return err
		}

		// Pop the next item. It could be the 'then' block or the condition.
		next, err := i.pop()
		if err != nil {
			return err
		}

		if block2, ok := next.([]string); ok {
			// This is an if-else. Stack: condition {then} {else} if
			// block1 is {else}, block2 is {then}
			condition, err := i.popBool()
			if err != nil {
				return err
			}

			if condition {
				err := i.execute(block2) // then block
				if err == ErrBreak || err == ErrContinue {
					return err // Re-propagate break/continue
				}
				return err
			} else {
				err := i.execute(block1) // else block
				if err == ErrBreak || err == ErrContinue {
					return err // Re-propagate break/continue
				}
				return err
			}
		} else if condition, ok := next.(bool); ok {
			// This is a simple if. Stack: condition {then} if
			// block1 is {then}, condition is the boolean
			if condition {
				err := i.execute(block1) // then block
				if err == ErrBreak || err == ErrContinue {
					return err // Re-propagate break/continue
				}
				return err
			}
			return nil // no else block, do nothing
		} else if condition, ok := next.(float64); ok {
			// This is a simple if. Stack: condition {then} if
			// block1 is {then}, condition is the number
			if condition != 0 {
				err := i.execute(block1) // then block
				if err == ErrBreak || err == ErrContinue {
					return err // Re-propagate break/continue
				}
				return err
			}
			return nil // no else block, do nothing
		} else {
			return i.newError(9)
		}
	}
	i.opcodes["loop"] = func(i *Interpreter) error {
		block, err := i.popBlock()
		if err != nil {
			return err
		}
		count, err := i.popFloat()
		if err != nil {
			return err
		}
		for j := 0; j < int(count); j++ {
			// Check for interruption inside the loop
			select {
			case <-i.interrupted:
				return i.newError(51)
			default:
				// Continue execution
			}
			i.loopIndex = float64(j) // Set current loop index
			if err := i.execute(block); err != nil {
				if err == ErrBreak {
					break
				} else if err == ErrContinue {
					continue
				} else {
					return err
				}
			}
			time.Sleep(time.Millisecond) // Yield to other goroutines
		}
		i.loopIndex = -1 // Reset loop index after loop completes
		return nil
	}

	i.opcodes["while"] = func(i *Interpreter) error {
		bodyBlock, err := i.popBlock()
		if err != nil {
			return err
		}
		conditionBlock, err := i.popBlock()
		if err != nil {
			return err
		}

		for {
			// Check for interruption inside the loop
			select {
			case <-i.interrupted:
				return i.newError(51)
			default:
				// Continue execution
			}
			// Execute the condition block
			if err := i.execute(conditionBlock); err != nil {
				return err
			}

			// Pop the result of the condition
			condResult, err := i.pop()
			if err != nil {
				return err
			}

			// Evaluate the condition
			condition := false
			switch v := condResult.(type) {
			case bool:
				condition = v
			case float64:
				condition = v != 0
			default:
				return i.newError(49, condResult)
			}

			if !condition {
				break // Exit loop if condition is false
			}

			// Execute the body block
			if err := i.execute(bodyBlock); err != nil {
				if err == ErrBreak {
					break
				} else if err == ErrContinue {
					continue
				} else {
					return err
				}
			}
			time.Sleep(time.Millisecond) // Yield to other goroutines
		}
		return nil
	}

	i.opcodes["break"] = func(i *Interpreter) error {
		return ErrBreak
	}

	i.opcodes["continue"] = func(i *Interpreter) error {
		return ErrContinue
	}

	// Loop index
	i.opcodes["index"] = func(i *Interpreter) error {
		if i.loopIndex == -1 {
			return i.newError(10)
		}
		i.push(i.loopIndex)
		return nil
	}

	// Storage
	i.opcodes["store"] = func(i *Interpreter) error {
		name, err := i.popString()
		if err != nil {
			return err
		}
		if strings.Contains(name, " ") {
			return i.newError(55, name)
		}
		val, err := i.pop()
		if err != nil {
			return err
		}

		if strings.HasPrefix(name, "$") {
			if len(i.scopeStack) > 1 {
				currentScope := i.scopeStack[len(i.scopeStack)-1]
				currentScope[name] = val
			} else {
				return i.newError(56, name)
			}
			return nil
		}

		// Prevent defining variables with the same name as an opcode
		if _, exists := i.opcodes[name]; exists {
			return i.newError(11, name)
		}

		// Handle internal variables (names starting with '_')
		if strings.HasPrefix(name, "_") {
			// Check if it's an attempt to create a *new* internal variable
			if _, exists := i.variables[name]; !exists {
				return i.newError(12, name)
			}

			// If it's an *existing* internal variable, apply type protection
			switch name {
			case "_echo_mode", "_degree_mode", "_vars_value", "_stack_type", "_hidden_vars", "_exit_save":
				if _, ok := val.(bool); !ok {
					return i.newError(13, name)
				}
			default:
				return i.newError(14, name)
			}
		}

		i.variables[name] = val
		return nil
	}
	i.opcodes["load"] = func(i *Interpreter) error {
		name, err := i.popString()
		if err != nil {
			return err
		}

		if strings.HasPrefix(name, "$") {
			for j := len(i.scopeStack) - 1; j >= 1; j-- {
				if val, ok := i.scopeStack[j][name]; ok {
					i.push(val)
					return nil
				}
			}
			return i.newError(15, name)
		}

		val, ok := i.variables[name]
		if !ok {
			return i.newError(15, name)
		}
		i.push(val)
		return nil
	}

	// String manipulation
	i.opcodes["len"] = func(i *Interpreter) error {
		s, err := i.popString()
		if err != nil {
			return err
		}
		i.push(float64(len(s)))
		return nil
	}
	i.opcodes["mid"] = func(i *Interpreter) error {
		length, err := i.popFloat()
		if err != nil {
			return err
		}
		start, err := i.popFloat()
		if err != nil {
			return err
		}
		s, err := i.popString()
		if err != nil {
			return err
		}
		if int(start) >= 0 && int(start) < len(s) && int(start+length) >=0 && int(start+length) <= len(s) {
			i.push(s[int(start):int(start+length)])
			return nil
		} else {
			return i.newError(59)
		}
	}

	// Output
	i.opcodes["."] = func(i *Interpreter) error {
		val, err := i.pop()
		if err != nil {
			return err
		}
		fmt.Fprint(i.outputView, val)
		return nil
	}
	i.opcodes["print"] = i.opcodes["."]
	i.opcodes["cr"] = func(i *Interpreter) error {
		fmt.Fprintln(i.outputView)
		return nil
	}
	i.opcodes["cls"] = func(i *Interpreter) error {
		if tv, ok := i.outputView.(*tview.TextView); ok {
			tv.Clear()
		}
		return nil
	}

	// State management
	i.opcodes["save"] = func(i *Interpreter) error {
		filename, err := i.popString()
		if err != nil {
			return err
		}
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}
		return i.saveState(filename)
	}
	i.opcodes["restore"] = func(i *Interpreter) error {
		filename, err := i.popString()
		if err != nil {
			return err
		}
		if !strings.HasSuffix(filename, ".json") {
			filename += ".json"
		}
		return i.loadState(filename)
	}

	i.opcodes["import"] = func(i *Interpreter) error {
		filename, err := i.popString()
		if err != nil {
			return err
		}
		if !strings.HasSuffix(filename, ".rpn") {
			filename += ".rpn"
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return i.newError(40, err)
		}
		rpnPath := filepath.Join(home, rpnDir)
		fullPath := filepath.Join(rpnPath, filename)

		content, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return i.newError(46, fullPath, err)
		}

		// Check if the imported content defines a "main" word
		mainDefinedInImport := mainWordRegex.MatchString(string(content))

		err = i.Eval(string(content)) // Execute the content of the imported file
		if err != nil {
			return err
		}

		// After importing, execute a "main" word if it was defined in the imported file
		if mainDefinedInImport {
			if newMainWord, mainExistsAfter := i.words["main"]; mainExistsAfter {
				err := i.execute(newMainWord)
				if err != nil {
					return i.newError(31, "main", err) // Error executing word 'main'
				}
			}
		}
		return nil
	}

	i.opcodes["export"] = func(i *Interpreter) error {
		filename, err := i.popString()
		if err != nil {
			return err
		}
		wordName, err := i.popString()
		if err != nil {
			return err
		}

		wordDef, ok := i.words[wordName]
		if !ok {
			return i.newError(30, wordName)
		}

		if strings.HasPrefix(wordName, "_") {
			return i.newError(17, wordName)
		}

		if !strings.HasSuffix(filename, ".rpn") {
			filename += ".rpn"
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return i.newError(40, err)
		}
		rpnPath := filepath.Join(home, rpnDir)
		if err := os.MkdirAll(rpnPath, 0755); err != nil {
			return i.newError(41, err)
		}
		fullPath := filepath.Join(rpnPath, filename)

		file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return i.newError(47, fullPath, err)
		}
		defer file.Close()

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		comment := fmt.Sprintf("( %s - %s )", wordName, timestamp)
		fmt.Fprintln(file, comment)

		fmt.Fprintln(file, formatWord(wordName, wordDef))

		fmt.Fprintln(file, "")

		return nil
	}

	i.opcodes["list"] = func(i *Interpreter) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return i.newError(40, err)
		}
		rpnPath := filepath.Join(home, rpnDir)

		files, err := ioutil.ReadDir(rpnPath)
		if err != nil {
			return i.newError(48, rpnPath, err)
		}

		fmt.Fprintln(i.outputView, "RPN files in ~/.polish:")
		found := false
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".rpn") {
				fmt.Fprintln(i.outputView, "  -", file.Name())
				found = true
			}
		}

		if !found {
			fmt.Fprintln(i.outputView, "  (No .rpn files found)")
		}

		return nil
	}

	// Prompt command
	i.opcodes["prompt"] = func(i *Interpreter) error {
		promptMsg, err := i.popString()
		if err != nil {
			return err
		}

		// Store original prompt and set new one
		i.originalPrompt = i.inputField.GetLabel()

		i.promptMutex.Lock()
		i.promptActive = true
		i.promptMutex.Unlock()

		// Update UI on the main goroutine
		i.app.QueueUpdateDraw(func() {
			i.inputField.SetLabel(promptMsg + " ")
			i.inputField.SetText("")
		})

		// Block until input is received from the UI thread.
		input := <-i.inputChan

		// Push the received input onto the stack.
		i.push(input)

		// The prompt is no longer active.
		i.promptMutex.Lock()
		i.promptActive = false
		i.promptMutex.Unlock()

		return nil
	}

	i.opcodes["words"] = func(i *Interpreter) error {
		var allWords []string

		// Get core commands
		for opcode := range i.opcodes {
			allWords = append(allWords, opcode)
		}

		// Get user-defined words
		for word := range i.words {
			allWords = append(allWords, word)
		}

		// Get variables
		for variable := range i.variables {
			allWords = append(allWords, variable)
		}

		sort.Strings(allWords)
		// Display in color
		for _, word := range allWords {
			if _, exists := i.opcodes[word]; exists {
                        	word = "[yellow]" + word + "[white]"
			}
			if _, exists := i.variables[word]; exists {
				if strings.HasPrefix(word, "_") {
					word = "[green]" + word + "[white]"
				} else {
					word = "[blue]" + word + "[white]"
				}
			}
			fmt.Fprint(i.outputView, word + " ")
                }
		fmt.Fprint(i.outputView, "\n")
		// fmt.Fprintln(i.outputView, strings.Join(allWords, " "))
		return nil
	}

	// Time and Date
	i.opcodes["time"] = func(i *Interpreter) error {
		i.push(time.Now().Format("15:04:05"))
		return nil
	}
	i.opcodes["date"] = func(i *Interpreter) error {
		i.push(time.Now().Format("2006-01-02"))
		return nil
	}
	i.opcodes["year"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Year()))
		return nil
	}
	i.opcodes["month"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Month()))
		return nil
	}
	i.opcodes["day"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Day()))
		return nil
	}
	i.opcodes["hour"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Hour()))
		return nil
	}
	i.opcodes["minute"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Minute()))
		return nil
	}
	i.opcodes["second"] = func(i *Interpreter) error {
		i.push(float64(time.Now().Second()))
		return nil
	}

	// Stack manipulation
	i.opcodes["clear"] = func(i *Interpreter) error {
		i.stack = make([]interface{}, 0)
		return nil
	}

	i.opcodes["free"] = func(i *Interpreter) error {
		newVariables := make(map[string]interface{})
		for name, value := range i.variables {
			if strings.HasPrefix(name, "_") {
				newVariables[name] = value
			}
		}
		i.variables = newVariables
		// Update the variables view after freeing user-defined variables
		showVarsValue := true
		if val, ok := i.variables["_vars_value"].(bool); ok {
			showVarsValue = val
		}
		hideInternalVars := true
		if val, ok := i.variables["_hidden_vars"].(bool); ok {
			hideInternalVars = val
		}
		if i.variablesTable != nil {
			if outputTextView, ok := i.outputView.(*tview.TextView); ok {
				updateVariablesView(i.variablesTable, i.variables, showVarsValue, hideInternalVars, outputTextView)
			}
		}
		return nil
	}

	i.opcodes["forget"] = func(i *Interpreter) error {
		i.words = make(map[string][]string, 0)
		if i.wordsTable != nil {
			updateWordsView(i.wordsTable, i.words)
		}
		return nil
	}

	// Constants
	i.opcodes["pi"] = func(i *Interpreter) error {
		i.push(math.Pi)
		return nil
	}
	i.opcodes["e"] = func(i *Interpreter) error {
		i.push(math.E)
		return nil
	}
	// Phi (Golden Ratio) - not directly in math package, calculate it
	i.opcodes["phi"] = func(i *Interpreter) error {
		i.push((1 + math.Sqrt(5)) / 2)
		return nil
	}

	// Random number generation
	i.opcodes["rand"] = func(i *Interpreter) error {
		i.push(rand.Float64())
		return nil
	}

	// Stack depth
	i.opcodes["depth"] = func(i *Interpreter) error {
		i.push(float64(len(i.stack)))
		return nil
	}

	// Boolean operations
	i.opcodes["true"] = func(i *Interpreter) error {
		i.push(true)
		return nil
	}
	i.opcodes["false"] = func(i *Interpreter) error {
		i.push(false)
		return nil
	}
	i.opcodes["and"] = func(i *Interpreter) error {
		b, err := i.popBool()
		if err != nil {
			return err
		}
		a, err := i.popBool()
		if err != nil {
			return err
		}
		i.push(a && b)
		return nil
	}
	i.opcodes["or"] = func(i *Interpreter) error {
		b, err := i.popBool()
		if err != nil {
			return err
		}
		a, err := i.popBool()
		if err != nil {
			return err
		}
		i.push(a || b)
		return nil
	}
	i.opcodes["not"] = func(i *Interpreter) error {
		a, err := i.popBool()
		if err != nil {
			return err
		}
		i.push(!a)
		return nil
	}
	i.opcodes["xor"] = func(i *Interpreter) error {
		b, err := i.popBool()
		if err != nil {
			return err
		}
		a, err := i.popBool()
		if err != nil {
			return err
		}
		i.push(a != b)
		return nil
	}

	i.opcodes["set"] = func(i *Interpreter) error {
		name, err := i.popString()
		if err != nil {
			return err
		}

		// Prevent defining variables with the same name as an opcode
		if _, exists := i.opcodes[name]; exists {
			return i.newError(11, name)
		}

		// Handle internal variables (names starting with '_')
		if strings.HasPrefix(name, "_") {
			// Check if it's an attempt to create a *new* internal variable
			if _, exists := i.variables[name]; !exists {
				return i.newError(12, name)
			}

			// If it's an *existing* internal variable, apply type protection
			switch name {
			case "_echo_mode", "_degree_mode", "_vars_value", "_stack_type", "_hidden_vars", "_exit_save":
				// These are boolean flags, so allow setting them to true
			default:
				return i.newError(14, name)
			}
		}

		i.variables[name] = true
		return nil
	}

	i.opcodes["unset"] = func(i *Interpreter) error {
		name, err := i.popString()
		if err != nil {
			return err
		}

		// Prevent defining variables with the same name as an opcode
		if _, exists := i.opcodes[name]; exists {
			return i.newError(11, name)
		}

		// Handle internal variables (names starting with '_')
		if strings.HasPrefix(name, "_") {
			// Check if it's an attempt to create a *new* internal variable
			if _, exists := i.variables[name]; !exists {
				return i.newError(12, name)
			}

			// If it's an *existing* internal variable, apply type protection
			switch name {
			case "_echo_mode", "_degree_mode", "_vars_value", "_stack_type", "_hidden_vars", "_exit_save":
				// These are boolean flags, so allow setting them to false
			default:
				return i.newError(14, name)
			}
		}

		i.variables[name] = false
		return nil
	}

	i.opcodes["toggle"] = func(i *Interpreter) error {
		name, err := i.popString()
		if err != nil {
			return err
		}

		// Prevent toggling variables with the same name as an opcode
		if _, exists := i.opcodes[name]; exists {
			return i.newError(11, name)
		}

		val, ok := i.variables[name]
		if !ok {
			return i.newError(15, name)
		}

		b, isBool := val.(bool)
		if !isBool {
			return i.newError(4, val)
		}

		// Handle internal variables (names starting with '_')
		if strings.HasPrefix(name, "_") {
			// Only allow toggling of specific internal boolean variables
			switch name {
			case "_echo_mode", "_degree_mode", "_vars_value", "_stack_type", "_hidden_vars", "_exit_save":
				// These are boolean flags, so allow toggling them
			default:
				return i.newError(14, name)
			}
		}

		i.variables[name] = !b
		return nil
	}

	// String case conversion
	i.opcodes["upper"] = func(i *Interpreter) error {
		s, err := i.popString()
		if err != nil {
			return err
		}
		i.push(strings.ToUpper(s))
		return nil
	}
	i.opcodes["lower"] = func(i *Interpreter) error {
		s, err := i.popString()
		if err != nil {
			return err
		}
		i.push(strings.ToLower(s))
		return nil
	}

	i.opcodes["val"] = func(i *Interpreter) error {
		s, err := i.popString()
		if err != nil {
			return err
		}
		// Try to parse as a boolean
		if s == "true" {
			i.push(true)
		} else if s == "false" {
			i.push(false)
		} else {
			// Try to parse as a float
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				i.push(f)
			} else {
				// If it's neither, push the original string back
				i.push(s)
			}
		}
		return nil
	}

	i.opcodes["str"] = func(i *Interpreter) error {
		v, err := i.pop()
		if err != nil {
			return err
		}
		i.push(fmt.Sprintf("%v", v))
		return nil
	}

	// UTF-8 commands
	i.opcodes["code"] = func(i *Interpreter) error {
		s, err := i.popString()
		if err != nil {
			return err
		}
		if len(s) == 0 {
			return i.newError(58, "empty string")
		}
		runeValue := []rune(s)[0]
		i.push(float64(runeValue))
		return nil
	}

	i.opcodes["char"] = func(i *Interpreter) error {
		code, err := i.popFloat()
		if err != nil {
			return err
		}
		runeValue := rune(code)
		i.push(string(runeValue))
		return nil
	}

	i.opcodes["emit"] = func(i *Interpreter) error {
		code, err := i.popFloat()
		if err != nil {
			return err
		}
		runeValue := rune(code)
		fmt.Fprint(i.outputView, string(runeValue))
		return nil
	}

	// Math functions
	i.opcodes["inv"] = func(i *Interpreter) error {
		x, err := i.popFloat()
		if err != nil {
			return err
		}
		if x == 0 {
			return i.newError(2) // Division by zero
		}
		i.push(1 / x)
		return nil
	}

	// Help
	i.opcodes["help"] = func(i *Interpreter) error {
		fmt.Fprintln(i.outputView, CleanMarkdown(readmeContent))
		return nil
	}
}

func (i *Interpreter) enterFileEditMode(filePath string) {
	i.editingFile = true
	i.currentEditFile = filePath
	i.fileDirty = false // Reset dirty flag

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// if file does not exist, create it
		if os.IsNotExist(err) {
			content = []byte{}
		} else {
			fmt.Fprintf(i.outputView, "[red]Error reading file %s: %v[white]\n", filePath, err)
			return
		}
	}

	textArea := tview.NewTextArea()
	textArea.SetText(string(content), false)
	textArea.SetBorder(true)

	updateTitle := func() {
		dirtyFlag := ""
		if i.fileDirty {
			dirtyFlag = "*"
		}
		fromRow, fromCol, _, _ := textArea.GetCursor()
		// The row and column are 0-indexed, so we add 1 for display.
		title := fmt.Sprintf("File %s%s | L: %d, C: %d | Ctrl-S to Save | Esc to Cancel", filepath.Base(filePath), dirtyFlag, fromRow+1, fromCol+1)
		textArea.SetTitle(title)
	}

	updateTitle() // Initial title update

	textArea.SetChangedFunc(func() {
		i.fileDirty = true
		updateTitle()
	})

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// We update the title on every key press.
		// This is not ideal, but it's the simplest way to keep the cursor position updated.
		defer updateTitle()

		switch event.Key() {
		case tcell.KeyCtrlS:
			i.exitFileEditMode(textArea, true)
			return nil
		case tcell.KeyEsc:
			i.exitFileEditMode(textArea, false)
			return nil
		}
		return event
	})

	// Replace input field with text area
	i.appFlex.RemoveItem(i.inputField)
	i.appFlex.AddItem(textArea, 0, 1, true)
	i.app.SetFocus(textArea)
	// Update focusables
	i.SetFocusables(textArea, i.outputView.(tview.Primitive), i.stackTable, i.variablesTable, i.wordsTable)
}

func (i *Interpreter) exitFileEditMode(textArea *tview.TextArea, save bool) {
	i.editingFile = false
	i.fileDirty = false

	// Get content before removing the text area
	content := textArea.GetText()

	// Restore input field
	i.appFlex.RemoveItem(textArea)
	i.appFlex.AddItem(i.inputField, 3, 0, true)
	i.app.SetFocus(i.inputField)

	// Restore focusables
	i.SetFocusables(i.inputField, i.outputView.(tview.Primitive), i.angleModeView, i.stackTable, i.variablesTable, i.wordsTable)

	if save {
		err := ioutil.WriteFile(i.currentEditFile, []byte(content), 0644)
		if err != nil {
			fmt.Fprintf(i.outputView, "[red]Error writing file %s: %v[white]\n", i.currentEditFile, err)
		} else {
			fmt.Fprintf(i.outputView, "[green]File '%s' saved.[white]\n", filepath.Base(i.currentEditFile))
			// Execute the content
			go i.app.QueueUpdateDraw(func() {
				if err := i.Eval(content); err != nil {
					fmt.Fprintln(i.outputView, err.Error())
				}
				// Update views
				showStackType := false
				if val, ok := i.variables["_stack_type"].(bool); ok {
					showStackType = val
				}
				updateStackView(i.stackTable, i.stack, showStackType)
			})
		}
	} else {
		fmt.Fprintf(i.outputView, "[yellow]Edit cancelled.[white]\n")
	}

	i.currentEditFile = ""
	i.inputField.SetText("") // Clear input field after editing
}

// compareStringSlices compares two string slices for equality.
func compareStringSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// execute runs a sequence of tokens through the interpreter.
func (i *Interpreter) execute(tokens []string) error {
	i.variables["_error"] = false
	i.clrEdit = true
	commentLevel := 0
	for j := 0; j < len(tokens); j++ {
		// Check for interruption
		select {
		case <-i.interrupted:
			return i.newError(51)
		default:
			// Continue execution
		}

		token := tokens[j]

		if commentLevel > 0 {
			if token == "(" {
				commentLevel++
			} else if token == ")" {
				commentLevel--
			}
			continue
		}

		if token == "(" {
			commentLevel++
			continue
		}

		// Prioritize quoted strings as literals
		if len(token) > 1 && token[0] == '"' && token[len(token)-1] == '"' {
			i.push(token[1 : len(token)-1]) // Push the unquoted string
			continue                        // Move to next token
		}

		// Handle function definition
		if token == ":" {
			if len(tokens) < j+3 {
				return i.newError(16)
			}
			wordName := tokens[j+1]

			// Prevent defining words starting with "_"
			if strings.HasPrefix(wordName, "_") {
				return i.newError(17, wordName)
			}

			// Prevent defining words with the same name as an opcode
			if _, exists := i.opcodes[wordName]; exists {
				return i.newError(18, wordName)
			}

			var wordDef []string
			j += 2
			for ; j < len(tokens); j++ {
				if tokens[j] == ";" {
					break
				}
				wordDef = append(wordDef, tokens[j])
			}
			i.words[wordName] = wordDef
			if i.wordsTable != nil {
				updateWordsView(i.wordsTable, i.words)
			}
			continue
		}

		// Handle end of function definition or standalone semicolon
		if token == ";" {
			return i.newError(57) // Return an error for out-of-context semicolons
		}

		// Handle code blocks for control flow
		if token == "{" {
			block, end, err := i.parseBlock(tokens, j)
			if err != nil {
				return err
			}
			i.push(block)
			j = end
			continue
		}

		// Handle delete command
		if token == "delete" {
			if len(tokens) < j+2 {
				return i.newError(19)
			}

			name := tokens[j+1]

			var targetType string
			var targetName string

			if strings.HasPrefix(name, "word:") {
				targetType = "word"
				targetName = strings.TrimPrefix(name, "word:")
			} else if strings.HasPrefix(name, "var:") {
				targetType = "variable"
				targetName = strings.TrimPrefix(name, "var:")
			} else {
				targetType = "any"
				targetName = name
			}

			if strings.HasPrefix(targetName, "_") {
				return i.newError(20, targetName)
			}

			// Check if it's an internal command
                        if _, exists := i.opcodes[targetName]; exists {
                                if targetType == "any" || targetType == "word" {
                                        return i.newError(50, targetName)
                                }
                        }

			deleted := false
			if targetType == "word" || targetType == "any" {
				if _, ok := i.words[targetName]; ok {
					delete(i.words, targetName)
					if i.wordsTable != nil {
						updateWordsView(i.wordsTable, i.words)
					}
					deleted = true
				}
			}

			if !deleted && (targetType == "variable" || targetType == "any") {
				if _, ok := i.variables[targetName]; ok {
					delete(i.variables, targetName)
					showVarsValue := true
					if val, ok := i.variables["_vars_value"].(bool); ok {
						showVarsValue = val
					}
					hideInternalVars := true
					if val, ok := i.variables["_hidden_vars"].(bool); ok {
						hideInternalVars = val
					}
					if i.variablesTable != nil {
						if outputTextView, ok := i.outputView.(*tview.TextView); ok {
							updateVariablesView(i.variablesTable, i.variables, showVarsValue, hideInternalVars, outputTextView)
						}
					}
					deleted = true
				}
			}

			if !deleted {
				if targetType == "any" {
					return i.newError(21, name)
				}
				return i.newError(22, targetType, targetName)
			}
			j++ // consume the name token
			continue
		}

		// Handle see command
		if token == "see" {
			if len(tokens) < j+2 {
				return i.newError(23)
			}
			name := tokens[j+1]

			var targetType string
			var targetName string

			if strings.HasPrefix(name, "word:") {
				targetType = "word"
				targetName = strings.TrimPrefix(name, "word:")
			} else if strings.HasPrefix(name, "var:") {
				targetType = "variable"
				targetName = strings.TrimPrefix(name, "var:")
			} else {
				targetType = "any"
				targetName = name
			}

			// Check if it's an internal command
			if _, exists := i.opcodes[targetName]; exists {
				if targetType == "any" || targetType == "word" {
					return i.newError(50, targetName)
				}
			}

			found := false
			if targetType == "word" || targetType == "any" {
				if wordDef, ok := i.words[targetName]; ok {
					fmt.Fprintln(i.outputView, formatWord(targetName, wordDef))
					found = true
				}
			}

			if !found && (targetType == "variable" || targetType == "any") {
				if varVal, ok := i.variables[targetName]; ok {
					var formattedValue string
					switch v := varVal.(type) {
					case []string:
						formattedValue = "{ " + strings.Join(v, " ") + " }"
					case []interface{}:
						strSlice := make([]string, len(v))
						for i, val := range v {
							strSlice[i] = fmt.Sprintf("%v", val)
						}
						formattedValue = "{ " + strings.Join(strSlice, " ") + " }"
					default:
						formattedValue = fmt.Sprintf("%v", v)
					}
					fmt.Fprintln(i.outputView, formattedValue)
					found = true
				}
			}

			if !found {
				if targetType == "any" {
					return i.newError(21, name)
				}
				return i.newError(22, targetType, targetName)
			}
			j++ // consume the name token
			continue
		}

		// Handle editfile command
		if token == "editfile" {
			if i.inputField == nil {
				// Silently ignore in non-interactive mode
				j++ // consume the name token
				continue
			}
			if len(tokens) < j+2 {
				return i.newError(53)
			}
			filename := tokens[j+1]
			// If the name is a quoted string, unquote it
			if len(filename) > 1 && filename[0] == '"' && filename[len(filename)-1] == '"' {
				filename = filename[1 : len(filename)-1]
			}
			// Add .rpn extension if not present
			if filepath.Ext(filename) == "" {
				filename += ".rpn"
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return i.newError(40, err)
			}
			rpnPath := filepath.Join(home, rpnDir)
			fullPath := filepath.Join(rpnPath, filename)

			// No need to check for existence, enterFileEditMode handles it
			i.enterFileEditMode(fullPath)
			i.clrEdit = false
			j++ // consume the name token
			return nil // Return to avoid further execution in the line
		}

		// Handle edit command
		if token == "edit" {
			if i.inputField == nil {
				// Silently ignore in non-interactive mode
				j++ // consume the name token
				continue
			}
			if len(tokens) < j+2 {
				return i.newError(24)
			}
			name := tokens[j+1]
			// If the name is a quoted string, unquote it
			if len(name) > 1 && name[0] == '"' && name[len(name)-1] == '"' {
				name = name[1 : len(name)-1]
			}
			var editString string
			var targetType string

			// Check if it's an internal command
                        if _, exists := i.opcodes[name]; exists {
                        	return i.newError(50, name)
                        }

			if strings.HasPrefix(name, "var:") {
				name = strings.TrimPrefix(name, "var:")
				targetType = "variable"
			} else if strings.HasPrefix(name, "word:") {
				name = strings.TrimPrefix(name, "word:")
				targetType = "word"
			} else {
				// Default behavior: check words first, then variables
				if _, ok := i.words[name]; ok {
					targetType = "word"
				} else if _, ok := i.variables[name]; ok {
					targetType = "variable"
				} else {
					return i.newError(25, name)
				}
			}

			switch targetType {
			case "word":
				if wordDef, ok := i.words[name]; ok {
					editString = ": " + name + " " + strings.Join(wordDef, " ") + " ;"
				} else {
					return i.newError(26, name)
				}
			case "variable":
				if varVal, ok := i.variables[name]; ok {
					// ... (existing logic for formatting variable for editing)
					if blockStr, isStringSlice := varVal.([]string); isStringSlice {
						editString = "{ " + strings.Join(blockStr, " ") + " } \"" + name + "\" store"
					} else if blockIface, isInterfaceSlice := varVal.([]interface{}); isInterfaceSlice {
						convertedBlock := make([]string, len(blockIface))
						for k, v := range blockIface {
							if s, isString := v.(string); isString {
								convertedBlock[k] = s
							} else {
								convertedBlock[k] = fmt.Sprintf("%v", v)
							}
						}
						editString = "{ " + strings.Join(convertedBlock, " ") + " } \"" + name + "\" store"
					} else {
						return i.newError(27, name)
					}
				} else {
					return i.newError(28, name)
				}
			}

			i.inputField.SetText(editString)
			i.clrEdit = false
			j += 1 // Consume the name token
			continue
		}

		// Main token processing logic
		if op, exists := i.opcodes[token]; exists { // Check for opcodes
			if err := op(i); err != nil {
				if err == ErrBreak || err == ErrContinue {
					return err
				}
				return i.newError(29, token, err)
			}
		} else if strings.HasPrefix(token, "word:") { // Explicitly execute a word
			wordName := strings.TrimPrefix(token, "word:")
			if wordDef, exists := i.words[wordName]; exists {
				if err := i.execute(wordDef); err != nil {
					return i.newError(31, wordName, err)
				}
			} else {
				return i.newError(30, wordName)
			}
		} else if strings.HasPrefix(token, "var:") { // Explicitly execute a variable (as a code block) or push its value
			varName := strings.TrimPrefix(token, "var:")
			if val, exists := i.variables[varName]; exists {
				if blockStr, ok := val.([]string); ok {
					if err := i.execute(blockStr); err != nil {
						return i.newError(33, varName, err)
					}
				} else if blockIface, ok := val.([]interface{}); ok {
					convertedBlock := make([]string, len(blockIface))
					for k, v := range blockIface {
						if s, isString := v.(string); isString {
							convertedBlock[k] = s
						} else {
							return i.newError(32, varName)
						}
					}
					if err := i.execute(convertedBlock); err != nil {
						return i.newError(33, varName, err)
					}
				} else {
					i.push(val) // Push the variable's value if not a code block
				}
			} else {
				return i.newError(15, varName)
			}
		} else if wordDef, exists := i.words[token]; exists { // Check for user-defined words (default)
			i.scopeStack = append(i.scopeStack, make(map[string]interface{}))
			err := i.execute(wordDef)
			i.scopeStack = i.scopeStack[:len(i.scopeStack)-1]
			if err != nil {
				return i.newError(31, token, err)
			}
		} else if val, exists := i.variables[token]; exists { // Check for variables
			// If the variable holds a block, execute it
			if blockStr, ok := val.([]string); ok {
				i.scopeStack = append(i.scopeStack, make(map[string]interface{}))
				err := i.execute(blockStr)
				i.scopeStack = i.scopeStack[:len(i.scopeStack)-1]
				if err != nil {
					return i.newError(33, token, err)
				}
			} else if blockIface, ok := val.([]interface{}); ok {
				// If it's []interface{}, try to convert it to []string
				convertedBlock := make([]string, len(blockIface))
				for k, v := range blockIface {
					if s, isString := v.(string); isString {
						convertedBlock[k] = s
					} else {
						return i.newError(32, token)
					}
				}
				i.scopeStack = append(i.scopeStack, make(map[string]interface{}))
				err := i.execute(convertedBlock)
				i.scopeStack = i.scopeStack[:len(i.scopeStack)-1]
				if err != nil {
					return i.newError(33, token, err)
				}
			} else if blockStr, ok := val.(string); ok && strings.HasPrefix(blockStr, "{") && strings.HasSuffix(blockStr, "}") {
				// If it's a string that looks like a block, tokenize and execute it
				innerContent := strings.TrimSpace(blockStr[1 : len(blockStr)-1])
				tempTokens, err := i.tokenize(innerContent)
				if err != nil {
					return i.newError(33, token, fmt.Errorf("failed to tokenize string block: %w", err))
				}
				i.scopeStack = append(i.scopeStack, make(map[string]interface{}))
				err = i.execute(tempTokens)
				i.scopeStack = i.scopeStack[:len(i.scopeStack)-1]
				if err != nil {
					return i.newError(33, token, err)
				}
			} else {
				// Otherwise, push the variable's value onto the stack
				i.push(val)
			}
		} else { // Attempt to parse as a float, boolean, or treat as a string literal
			if token == "true" {
				i.push(true)
			} else if token == "false" {
				i.push(false)
			} else {
				num, err := strconv.ParseFloat(token, 64)
				if err == nil {
					i.push(num)
				} else {
					// If none of the above, it's an unrecognized token
					return i.newError(34, token)
				}
			}
		}
		time.Sleep(time.Millisecond) // Yield to other goroutines
	}
	if commentLevel > 0 {
		return i.newError(37)
	}
	return nil
}

// Eval parses and executes a line of RPN code.
func (i *Interpreter) Eval(line string) error {
	// Drain the interrupted channel before execution
	select {
	case <-i.interrupted:
	default:
	}

	tokens, err := i.tokenize(line) // Use a custom tokenizer
	if err != nil {
		return err
	}

	return i.execute(tokens)
}

// generateSuggestions creates a list of possible completions for the given input.
func (i *Interpreter) generateSuggestions(input string) []string {
	var suggestions []string
	inputLower := strings.ToLower(input)

	// Add opcodes
	for op := range i.opcodes {
		if strings.HasPrefix(op, inputLower) {
			suggestions = append(suggestions, op)
		}
	}

	// Add user-defined words
	for word := range i.words {
		if strings.HasPrefix(word, inputLower) {
			suggestions = append(suggestions, word)
		}
	}

	// Add variables
	for variable := range i.variables {
		if strings.HasPrefix(variable, inputLower) {
			suggestions = append(suggestions, variable)
		}
	}

	sort.Strings(suggestions)
	return suggestions
}

// handleTabCompletion handles the tab key press for autocompletion.
func (i *Interpreter) handleTabCompletion() {
	currentText := i.inputField.GetText()
	parts := strings.Fields(currentText)
	if len(parts) == 0 {
		return
	}

	lastPart := parts[len(parts)-1]

	if i.suggestionIndex == -1 || !strings.HasPrefix(strings.ToLower(lastPart), strings.ToLower(i.suggestions[i.suggestionIndex])) {
		i.suggestions = i.generateSuggestions(lastPart)
		i.suggestionIndex = -1
	}

	if len(i.suggestions) > 0 {
		i.suggestionIndex = (i.suggestionIndex + 1) % len(i.suggestions)
		newText := strings.Join(parts[:len(parts)-1], " ")
		if newText != "" {
			newText += " "
		}
		newText += i.suggestions[i.suggestionIndex]
		i.inputField.SetText(newText)
	}
}

// tokenize splits a line of code into tokens.
func (i *Interpreter) tokenize(line string) ([]string, error) {
	var tokens []string
	inString := false
	inBlock := 0
	current := ""

	for _, r := range line {
		switch {
		case unicode.IsSpace(r) && !inString && inBlock == 0:
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		case r == '"' && inBlock == 0:
			current += string(r)
			if inString {
				tokens = append(tokens, current)
				current = ""
			}
			inString = !inString
		case r == '{' && !inString:
			if current != "" {
				tokens = append(tokens, current)
			}
			tokens = append(tokens, "{")
			current = ""
		case r == '}' && !inString:
			if current != "" {
				tokens = append(tokens, current)
			}
			tokens = append(tokens, "}")
			current = ""
		default:
			current += string(r)
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	if inString {
		return nil, i.newError(36)
	}

	return tokens, nil
}

func formatWord(wordName string, wordDef []string) string {                                                                                      
  var builder strings.Builder                                                                                                                     
  builder.WriteString(fmt.Sprintf(": %s", wordName))                                                                                              
                                                                                                                                                 
  indentLevel := 1                                                                                                                                
  indentUnit := "  "                                                                                                                              
                                                                                                                                                 
  // Group consecutive non-block tokens to print on the same line                                                                                 
  for i := 0; i < len(wordDef); {                                                                                                                 
    // Find the next block token '{' or '}'                                                                                                       
    nextBlockIndex := -1                                                                                                                          
    for j := i; j < len(wordDef); j++ {                                                                                                           
      if wordDef[j] == "{" || wordDef[j] == "}" {                                                                                                 
        nextBlockIndex = j                                                                                                                        
        break                                                                                                                                     
      }                                                                                                                                           
    }                                                                                                                                             
                                                                                                                                                 
    // If there are non-block tokens before the next block token (or at the end)                                                                  
    if nextBlockIndex != i {                                                                                                                      
      end := len(wordDef)                                                                                                                         
      if nextBlockIndex != -1 {                                                                                                                   
        end = nextBlockIndex                                                                                                                      
      }                                                                                                                                           
                                                                                                                                                 
      // Join and print the non-block tokens                                                                                                      
      if i < end {                                                                                                                                
        builder.WriteString("\n" + strings.Repeat(indentUnit, indentLevel))                                                                       
        builder.WriteString(strings.Join(wordDef[i:end], " "))                                                                                    
      }                                                                                                                                           
      i = end                                                                                                                                     
    }                                                                                                                                             

                                                                                                                                                  
    // Handle the block token                                                                                                                     
    if i < len(wordDef) {                                                                                                                         
      tok := wordDef[i]                                                                                                                           
      if tok == "{" {                                                                                                                             
        builder.WriteString("\n" + strings.Repeat(indentUnit, indentLevel))                                                                       
        builder.WriteString("{")                                                                                                                  
        indentLevel++                                                                                                                             
      } else if tok == "}" {                                                                                                                      
        indentLevel--                                                                                                                             
        if indentLevel < 1 {                                                                                                                      
          indentLevel = 1                                                                                                                         
        }                                                                                                                                         
        builder.WriteString("\n" + strings.Repeat(indentUnit, indentLevel))                                                                       
        builder.WriteString("}")                                                                                                                  
      }                                                                                                                                           
      i++                                                                                                                                         
    }                                                                                                                                             
  }                                                                                                                                               
                                                                                                                                                  
  builder.WriteString("\n;")                                                                                                                      
  return builder.String()                                                                                                                         
}

// parseBlock finds a matching '}' for a '{' and returns the inner tokens.
func (i *Interpreter) parseBlock(tokens []string, start int) ([]string, int, error) {
	if tokens[start] != "{" {
		return nil, 0, i.newError(38)
	}
	balance := 1
	for j := start + 1; j < len(tokens); j++ {
		if tokens[j] == "{" {
			balance++
		} else if tokens[j] == "}" {
			balance--
			if balance == 0 {
				return tokens[start+1 : j], j, nil
			}
		}
	}
	return nil, 0, i.newError(39)
}

// updateStackView clears and repopulates the stack table.
func updateStackView(stackTable *tview.Table, stack []interface{}, showType bool) {
	stackTable.SetTitle(fmt.Sprintf("Stack (%d)", len(stack)))
	stackTable.Clear()
	stackTable.SetCell(0, 0, tview.NewTableCell("Index").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	if showType {
		stackTable.SetCell(0, 1, tview.NewTableCell("Type").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	} else {
		stackTable.SetCell(0, 1, tview.NewTableCell("Value").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	}
	for i := len(stack) - 1; i >= 0; i-- {
		item := stack[i]
		stackTable.SetCell(len(stack)-1-i, 0, tview.NewTableCell(fmt.Sprintf("%06d:", len(stack)-1-i)).SetSelectable(false))
		if showType {
			stackTable.SetCell(len(stack)-1-i, 1, tview.NewTableCell(fmt.Sprintf("%T", item)).SetSelectable(false))
		} else {
			stackTable.SetCell(len(stack)-1-i, 1, tview.NewTableCell(fmt.Sprintf("%v", item)).SetSelectable(false))
		}
	}
	stackTable.ScrollToBeginning()
}

func main() {
	var welcome = appName + " v" + version + " - RPN Interpreter written in Go.\n"
	welcome += "Type 'exit', 'quit', 'bye' or press the 'F12' key to exit.\n"
	welcome += "Type 'help' or press the 'F1' key to have a summary of commands.\n\n"

	app := tview.NewApplication()

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Load command history
	loadHistory()

	outputView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(true).SetWordWrap(true)
	outputView.SetBorder(true).SetTitle(appName + " v" + version)
	outputView.SetChangedFunc(func() {
		app.Draw()
		outputView.ScrollToEnd()
	})

	stackTable := tview.NewTable().SetBorders(false).SetSelectable(true, false)
	stackTable.SetBorder(true).SetTitle("Stack")

	angleModeView := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreen)
	angleModeView.SetBorder(true).SetTitle("Mode")

	clockView := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorGreen)
	clockView.SetBorder(true).SetTitle("Time")

	variablesTable := tview.NewTable().SetBorders(false).SetSelectable(true, false)
	variablesTable.SetBorder(true).SetTitle("Variables")

	inputField := tview.NewInputField().SetLabel(myPrompt)
	inputField.SetBorder(true).SetTitle("F2 : Change panel")
	inputField.SetFieldTextColor(tcell.ColorGreen)
	inputField.SetFieldBackgroundColor(tcell.ColorBlack)

	appFlex := tview.NewFlex()

	interpreter := NewInterpreter(app, appFlex, outputView, angleModeView, clockView, variablesTable, stackTable, inputField)

	interpreter.opcodes["exit"] = func(i *Interpreter) error {
		if val, ok := interpreter.variables["_exit_save"].(bool); ok && val { // save to default ?
			interpreter.saveState("default.json")
		}
		app.Stop()
		return nil
	}
	interpreter.opcodes["quit"] = interpreter.opcodes["exit"]
	interpreter.opcodes["bye"] = interpreter.opcodes["exit"]

	interpreter.wordsTable.SetBorder(true)
	interpreter.wordsTable.SetTitle("Words")

	// Initial angle mode display
	updateAngleAndEchoModeView(interpreter)

	// Start goroutine to update clock every second
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			app.QueueUpdateDraw(func() {
				updateClockView(interpreter)
			})
		}
	}()

	// Attempt to load and execute default.json
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(outputView, "Error getting home directory: %v\n", err)
	} else {
		rpnPath := filepath.Join(home, rpnDir)
		defaultRpnFile := filepath.Join(rpnPath, "default.json")
		if _, err := os.Stat(defaultRpnFile); err == nil {
			fmt.Fprintln(outputView, "Loading default.json...")
			if err := interpreter.loadState("default.json"); err != nil {
				fmt.Fprintf(outputView, "Error loading default.json: %v\n", err)
			} else {
				// Check for and execute 'init' variable if it's a code block
				if initVal, ok := interpreter.variables["init"]; ok {

					var initBlock []string
					if blockStr, isStringSlice := initVal.([]string); isStringSlice {
						initBlock = blockStr
					} else if initStr, isString := initVal.(string); isString {
						// If it's a string, try to tokenize it as a block
						if strings.HasPrefix(initStr, "{") && strings.HasSuffix(initStr, "}") {
							innerContent := strings.TrimSpace(initStr[1 : len(initStr)-1])
							tempTokens, err := interpreter.tokenize(innerContent)
							if err == nil {
								initBlock = tempTokens
							} else {
								fmt.Fprintf(outputView, "Warning: 'init' variable is a string that looks like a block, but failed to tokenize: %v\n", err)
							}
						}
					} else if blockIface, isInterfaceSlice := initVal.([]interface{}); isInterfaceSlice {
						convertedBlock := make([]string, len(blockIface))
						for k, v := range blockIface {
							if s, isString := v.(string); isString {
								convertedBlock[k] = s
							} else {
								convertedBlock[k] = fmt.Sprintf("%v", v)
							}
						}
						initBlock = convertedBlock
					}

					if initBlock != nil {
						fmt.Fprintln(outputView, "Executing 'init' variable...")
						if err := interpreter.execute(initBlock); err != nil {
							fmt.Fprintf(outputView, "Error executing 'init' variable: %v\n", err)
						}
					}
				} else {
					// 'init' variable not found.
				}

				// Check for and execute 'init' word if it exists
				if initWord, ok := interpreter.words["init"]; ok {
					fmt.Fprintln(outputView, "Executing 'init' word...")
					if err := interpreter.execute(initWord); err != nil {
						fmt.Fprintf(outputView, "Error executing 'init' word: %v\n", err)
					}
				}
			}
		}
		fmt.Fprintf(outputView, welcome)
	}

	// Initial stack view update
	showStackType := false
	if val, ok := interpreter.variables["_stack_type"].(bool); ok {
		showStackType = val
	}
	updateStackView(stackTable, interpreter.stack, showStackType)
	// Initial variables view update
	showVarsValue := true
	if val, ok := interpreter.variables["_vars_value"].(bool); ok {
		showVarsValue = val
	}
	hideInternalVars := true
	if val, ok := interpreter.variables["_hidden_vars"].(bool); ok {
		hideInternalVars = val
	}
	updateVariablesView(variablesTable, interpreter.variables, showVarsValue, hideInternalVars, outputView)
	updateWordsView(interpreter.wordsTable, interpreter.words)

	// Channel for passing commands from the input field to the interpreter goroutine.
	commandChan := make(chan string, 1)

	// Single goroutine to handle all interpreter execution.
	go func() {
		for text := range commandChan {
			// Execute the command.
			if err := interpreter.Eval(text); err != nil {
				app.QueueUpdateDraw(func() {
					fmt.Fprintln(outputView, err.Error())
				})
			}

			// Update views on the main goroutine after execution.
			app.QueueUpdateDraw(func() {
				showStackType := false
				if val, ok := interpreter.variables["_stack_type"].(bool); ok {
					showStackType = val
				}
				updateStackView(stackTable, interpreter.stack, showStackType)
				showVarsValue := true
				if val, ok := interpreter.variables["_vars_value"].(bool); ok {
					showVarsValue = val
				}
				hideInternalVars := true
				if val, ok := interpreter.variables["_hidden_vars"].(bool); ok {
					hideInternalVars = val
				}
				updateVariablesView(variablesTable, interpreter.variables, showVarsValue, hideInternalVars, outputView)
				updateWordsView(interpreter.wordsTable, interpreter.words)
			})
		}
	}()

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		text := inputField.GetText()

		interpreter.promptMutex.Lock()
		isActive := interpreter.promptActive
		interpreter.promptMutex.Unlock()

		if isActive {
			// A prompt is active, so send the input to the interpreter's
			// dedicated input channel and reset the UI.
			interpreter.inputChan <- text
			inputField.SetLabel(interpreter.originalPrompt)
			inputField.SetText("")
			return
		}

		// This is a regular command.
		if text == "" && len(history) > 0 && historyIndex < len(history) {
			text = history[historyIndex] // Use current history item if input is empty
		}
		if text == "" {
			return
		}

		// Add to history.
		if len(history) == 0 || history[len(history)-1] != text {
			history = append(history, text)
		}
		historyIndex = len(history)
		saveHistory()

		// Echo the command if echo mode is on.
		if val, ok := interpreter.variables["_echo_mode"].(bool); ok && val {
			fmt.Fprintf(outputView, "[yellow]%s%s[white]\n", myPrompt, text)
		}

		// Clear the input field and send the command to the interpreter loop.
		inputField.SetText("")
		commandChan <- text
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if historyIndex > 0 {
				historyIndex--
				inputField.SetText(history[historyIndex])
			}
			return nil
		case tcell.KeyDown:
			if historyIndex < len(history)-1 {
				historyIndex++
				inputField.SetText(history[historyIndex])
			} else if historyIndex == len(history)-1 {
				historyIndex++
				inputField.SetText("")
			}
			return nil
		case tcell.KeyTab:
			interpreter.handleTabCompletion()
			return nil
		case tcell.KeyEsc:
			// Send an interrupt signal to the interpreter
			select {
			case interpreter.interrupted <- struct{}{}:
				fmt.Fprintln(outputView, "[yellow]Execution interrupted.[white]")
			default:
			}
			return nil
		}
		// Reset suggestion index if any other key is pressed
		interpreter.suggestionIndex = -1
		interpreter.suggestions = []string{}
		return event
	})

	// Global key bindings
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			// Execute the "help" command
			commandChan <- "help"
			return nil
		case tcell.KeyF2:
			interpreter.CycleFocus()
			return nil
		case tcell.KeyF12:
			// Execute the "exit" command
			commandChan <- "exit"
			return nil
		case tcell.KeyCtrlC:
			return nil
		}
		return event
	})

	// Layout
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(angleModeView, 3, 0, false).
        AddItem(clockView, 3, 0, false).
        AddItem(stackTable, 0, 1, false).
        AddItem(variablesTable, 0, 1, false).
        AddItem(interpreter.wordsTable, 0, 1, false)

	mainFlex := tview.NewFlex().
		AddItem(outputView, 0, 2, false).
		AddItem(rightPanel, 0, 1, false)

	appFlex.SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, false).
		AddItem(inputField, 3, 0, true)

	app.SetRoot(appFlex, true).SetFocus(inputField)

	// Set initial focusables after all UI elements are created
	interpreter.SetFocusables(inputField, outputView, stackTable, variablesTable, interpreter.wordsTable)

	if err := app.Run(); err != nil {
		panic(err)
	}
	fmt.Println(myPrompt + appName + " v" + version + " - https://github.com/jplozf/polish")
}

// updateVariablesView clears and repopulates the variables table.
func updateVariablesView(variablesTable *tview.Table, variables map[string]interface{}, showValue, hideInternal bool, outputView io.Writer) {
	variablesTable.Clear()
	variablesTable.SetTitle(fmt.Sprintf("Variables (%d)", len(variables)))
	variablesTable.SetCell(0, 0, tview.NewTableCell("Variable").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	if showValue {
		variablesTable.SetCell(0, 1, tview.NewTableCell("Value").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	}

	// Sort keys for consistent order
	keys := make([]string, 0, len(variables))
	for k := range variables {
		if !hideInternal && strings.HasPrefix(k, "_") {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	row := 1
	for _, k := range keys {
		v := variables[k]
		variablesTable.SetCell(row, 0, tview.NewTableCell(k).SetSelectable(false))
		if showValue {
			var displayValue string
			switch val := v.(type) {
			case []string:
				displayValue = "{...}" // Show ellipsis for blocks
			case []interface{}:
				displayValue = "{...}"
			default:
				displayValue = fmt.Sprintf("%v", val)
			}
			variablesTable.SetCell(row, 1, tview.NewTableCell(displayValue).SetSelectable(false))
		}
		row++
	}
	variablesTable.SetTitle(fmt.Sprintf("Variables (%d)", len(keys)))
	variablesTable.ScrollToBeginning()
}

// updateWordsView clears and repopulates the words table.
func updateWordsView(wordsTable *tview.Table, words map[string][]string) {
	wordsTable.Clear()
	wordsTable.SetTitle(fmt.Sprintf("Words (%d)", len(words)))
	wordsTable.SetCell(0, 0, tview.NewTableCell("Word").SetSelectable(false).SetTextColor(tcell.ColorYellow))
	wordsTable.SetCell(0, 1, tview.NewTableCell("Definition").SetSelectable(false).SetTextColor(tcell.ColorYellow))

	// Sort keys for consistent order
	keys := make([]string, 0, len(words))
	for k := range words {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	row := 1
	for _, k := range keys {
		v := words[k]
		wordsTable.SetCell(row, 0, tview.NewTableCell(k).SetSelectable(false))
		// Join the definition tokens into a single string for display
		defStr := strings.Join(v, " ")
		if len(defStr) > 40 { // Truncate long definitions
			defStr = defStr[:37] + "..."
		}
		wordsTable.SetCell(row, 1, tview.NewTableCell(defStr).SetSelectable(false))
		row++
	}
	wordsTable.ScrollToBeginning()
}

// CleanMarkdown takes a markdown string and converts it into a more
// readable plain text format, preserving the document's structure.
func CleanMarkdown(markdown string) string {
	// The order of these replacements is important.
	var result string

	// 1. Remove HTML tags completely
	re := regexp.MustCompile("<[^>]*>")
	result = re.ReplaceAllString(markdown, "")

	// 2. Format headings by adding a newline before them for separation.
	// We'll capture the heading text and just prepend a newline.
	re = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	result = re.ReplaceAllString(result, "\n")

	// 3. Format images to be descriptive in plain text.
	// ![alt text](url) -> Image: alt text (url)
	re = regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	result = re.ReplaceAllString(result, "Image: $1 ($2)")

	// 4. Format links to be descriptive in plain text.
	// [link text](url) -> link text (url)
	re = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	result = re.ReplaceAllString(result, "$1 ($2)")

	// 5. Replace code block fences with a simple line for visual separation.
	re = regexp.MustCompile("(?m)^```.*$")
	result = re.ReplaceAllString(result, "----------------")

	// 6. Normalize list items to have consistent indentation.
	// This makes nested lists look cleaner.
	re = regexp.MustCompile(`(?m)^\s*([*+\-]|[0-9]+\.)\s+`)
	result = re.ReplaceAllString(result, "  $1 ")

	// 7. Clean up consecutive newlines (more than two) into just two.
	re = regexp.MustCompile(`\n{3,}`)
	result = re.ReplaceAllString(result, "\n\n")

	// 8. Trim leading/trailing whitespace.
	result = strings.TrimSpace(result)

	// 9. Add a final newline.
	result += "\n"

	return result
}