package main

// *****************************************************************************
//  __________      .__  .__       .__
//  \______   \____ |  | |__| _____|  |__
//   |     ___/  _ \|  | |  |/  ___/  |  \
//   |    |  (  <_> )  |_|  |\___ \|   Y  \
//   |____|   \____/|____/__/____  >___|  /
//                               \/     \/
//
//          Polish © jpl@ozf.fr 2024
//
// *****************************************************************************

// *****************************************************************************
// IMPORTS
// *****************************************************************************
import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"polish/color"
	"polish/stack"
	"polish/sto"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// *****************************************************************************
// GLOBALS
// *****************************************************************************
var (
	loopme   = true
	fs       *stack.FStack
	as       *stack.AStack
	appDir   string
	previous string
	varName  string
)

// *****************************************************************************
// init()
// *****************************************************************************
func init() {
	// Get the user folder
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// The application folder is supposed to be into the user folder
	// The application folder will store the serialized stack
	appDir = filepath.Join(userDir, APP_FOLDER)
	if _, err := os.Stat(appDir); errors.Is(err, os.ErrNotExist) {
		// Create this application folder into the user folder if not exists
		err := os.Mkdir(appDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Some blahblah
	greetings()
	// Create the stacks
	fs = stack.NewFStack()
	as = stack.NewAStack()
	// Deserialize the previous stacks, if any
	readStacks()
}

// *****************************************************************************
// main()
// *****************************************************************************
func main() {
	reader := bufio.NewReader(os.Stdin)

	for ok := true; ok; ok = loopme {
		// Display the prompt
		showPrompt()
		// Read the input
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		// Parse the input and execute
		parse(text)
	}
	// Serialize the current stacks
	saveStacks()
	fmt.Printf("\n%s Bye.\n\n", ICON_DISK)
}

// *****************************************************************************
// greetings()
// *****************************************************************************
func greetings() {
	fmt.Printf("%s Welcome to %s\n", ICON_DISK, APP_STRING)
	fmt.Printf("%s %s version %s build %s at %s\n", ICON_DISK, APP_NAME, APP_VERSION, GIT_COMMIT, BUILD_TIME)
	fmt.Printf("%s %s\n\n", ICON_DISK, APP_URL)
}

// *****************************************************************************
// parse()
// *****************************************************************************
func parse(txt string) {
	// https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks
	// words := strings.Fields(txt)
	quoted := false
	words := strings.FieldsFunc(txt, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	for _, w := range words {
		xeq(w)
	}
}

// *****************************************************************************
// xeq()
// *****************************************************************************
func xeq(cmd string) {
	// Is it a string ?
	if cmd[0] == '"' {
		s := ""
		if cmd[len(cmd)-1] == '"' {
			s = cmd[1 : len(cmd)-1]
		} else {
			s = cmd[1:]
		}
		as.Push(s)
	} else {
		// Is it a variable name ?
		if cmd[0] == '$' {
			varName = cmd[1:]
		} else {
			// Is it a number ?
			if isFloat(cmd) {
				// Then push it on the stack
				v, _ := strconv.ParseFloat(cmd, 64)
				fs.Push(v)
			} else {
				// Is it a mathematical function defined into mymath.go
				// under the shape MyFunction ?
				m := My{}
				mName := "My" + strings.Title(strings.ToLower(cmd))
				meth := reflect.ValueOf(m).MethodByName(mName)
				if meth.IsValid() {
					// Yes : call it
					meth.Call(nil)
				} else {
					// Here are special functions, stack handling and alias
					switch cmd {
					case "!!":
						xeq(previous)
					case "exit", "quit", "bye":
						loopme = false
					case "+":
						m.MyAdd()
					case "-":
						m.MySub()
					case "*":
						m.MyMult()
					case "/":
						m.MyDiv()
					case "**":
						m.MyPow()
					case "!":
						m.MyFact()
					case "drop":
						doDrop()
					case "dup":
						doDup()
					case "depth":
						doDepth()
					case "cls", "clr", "clear":
						doClear()
					case ".f":
						showFStack()
					case ".a":
						showAStack()
					case ".v":
						showVars()
					case "swap":
						doSwap()
					case "rot":
						doRot()
					case "adrop":
						doADrop()
					case "adup":
						doADup()
					case "adepth":
						doADepth()
					case "acls", "aclr", "aclear":
						doAClear()
					case "aswap":
						doASwap()
					case "+a", "aadd":
						doAAdd()
					case "+as", "aadds":
						doAAdds()
					case "*a":
						doAMult()
					case "alen":
						doALen()
					case "aright":
						doARight()
					case "aleft":
						doALeft()
					case "amid":
						doAMid()
					case "rcl":
						doRcl()
					case "sto":
						doSto()
					case "asto":
						doASto()
					case "del":
						doDel()
					default:
						raiseError(UNRECOGNIZED_COMMAND)
					}
				}
			}
		}
	}
	if cmd != "!!" {
		previous = cmd
	}
}

// *****************************************************************************
// isFloat()
// *****************************************************************************
func isFloat(c string) bool {
	rc := true
	// We try to convert the entered string to something which looks like a float number
	_, err := strconv.ParseFloat(c, 64)
	if err != nil {
		// It doesn't look like a float number
		rc = false
	}
	return rc
}

// *****************************************************************************
// checkFStack()
// *****************************************************************************
func checkFStack(n int) bool {
	// Do we have enough args on stack to perform the selected operation ?
	if fs.Depth() >= n {
		return true
	} else {
		raiseError(NOT_ENOUGH_ARGS_FLOAT)
		return false
	}
}

// *****************************************************************************
// checkAStack()
// *****************************************************************************
func checkAStack(n int) bool {
	// Do we have enough args on stack to perform the selected operation ?
	if as.Depth() >= n {
		return true
	} else {
		raiseError(NOT_ENOUGH_ARGS_ALPHA)
		return false
	}
}

// *****************************************************************************
// checkStacks()
// *****************************************************************************
func checkStacks(fn int, an int) bool {
	// Do we have enough args on the two stacks to perform the selected operation ?
	return (checkAStack(an) && checkFStack(fn))
}

// *****************************************************************************
// doDrop()
// *****************************************************************************
func doDrop() {
	if checkFStack(1) {
		fs.Pop()
	}
}

// *****************************************************************************
// doDup()
// *****************************************************************************
func doDup() {
	if checkFStack(1) {
		f, _ := fs.Pop()
		fs.Push(f)
		fs.Push(f)
	}
}

// *****************************************************************************
// doClear()
// *****************************************************************************
func doClear() {
	fs.S = nil
}

// *****************************************************************************
// doDepth()
// *****************************************************************************
func doDepth() {
	fs.Push(float64(fs.Depth()))
}

// *****************************************************************************
// showFStack()
// *****************************************************************************
func showFStack() {
	for i, value := range fs.S {
		i = fs.Depth() - 1 - i
		// We use scientific notation if the number of digits is greater than 12
		if math.Log10(math.Abs(value)) > 12 {
			fmt.Printf("[FLOAT] %05d : %21.6E\n", i, value)
		} else {
			fmt.Printf("[FLOAT] %05d : %21.6f\n", i, value)
		}
	}
}

// *****************************************************************************
// showAStack()
// *****************************************************************************
func showAStack() {
	for i, value := range as.S {
		i = as.Depth() - 1 - i
		fmt.Printf("[ALPHA] %05d : %s\n", i, value)
	}
}

// *****************************************************************************
// doADrop()
// *****************************************************************************
func doADrop() {
	if checkAStack(1) {
		as.Pop()
	}
}

// *****************************************************************************
// doADup()
// *****************************************************************************
func doADup() {
	if checkAStack(1) {
		a, _ := as.Pop()
		as.Push(a)
		as.Push(a)
	}
}

// *****************************************************************************
// doAClear()
// *****************************************************************************
func doAClear() {
	as.S = nil
}

// *****************************************************************************
// doADepth()
// *****************************************************************************
func doADepth() {
	fs.Push(float64(as.Depth()))
}

// *****************************************************************************
// doASwap()
// *****************************************************************************
func doASwap() {
	if checkAStack(2) {
		a2, _ := as.Pop()
		a1, _ := as.Pop()
		as.Push(a2)
		as.Push(a1)
	}
}

// *****************************************************************************
// doAAdd()
// *****************************************************************************
func doAAdd() {
	if checkAStack(2) {
		a2, _ := as.Pop()
		a1, _ := as.Pop()
		as.Push(a1 + a2)
	}
}

// *****************************************************************************
// doAAdds()
// *****************************************************************************
func doAAdds() {
	if checkAStack(2) {
		a2, _ := as.Pop()
		a1, _ := as.Pop()
		as.Push(a1 + " " + a2)
	}
}

// *****************************************************************************
// doAMult()
// *****************************************************************************
func doAMult() {
	if checkStacks(1, 1) {
		f, _ := fs.Pop()
		a, _ := as.Pop()
		o := ""
		for i := 1; i <= int(f); i++ {
			o += a
		}
		as.Push(o)
	}
}

// *****************************************************************************
// doALeft()
// *****************************************************************************
func doALeft() {
	if checkStacks(1, 1) {
		f, _ := fs.Pop()
		a, _ := as.Pop()
		if int(f) <= len(a) {
			as.Push(a[0:int(f)])
		} else {
			raiseError(ALPHA_EXTRACT_OUT_OF_BOUNDS)
		}
	}
}

// *****************************************************************************
// doARight()
// *****************************************************************************
func doARight() {
	if checkStacks(1, 1) {
		f, _ := fs.Pop()
		a, _ := as.Pop()
		if int(f) <= len(a) {
			as.Push(a[len(a)-int(f):])
		} else {
			raiseError(ALPHA_EXTRACT_OUT_OF_BOUNDS)
		}
	}
}

// *****************************************************************************
// doAMid()
// *****************************************************************************
func doAMid() {
	if checkStacks(2, 1) {
		/*
			s = "12345"
			s[0:1] = "1"
			s[1:3] = "23"
			s[4:5] = "5"
			A slice is formed by specifying two indices, a low and high bound, separated by a colon:
			a[low : high]
			This selects a half-open range which includes the first element, but excludes the last one.
			The slice is 0-indexed.
		*/
		n, _ := fs.Pop() // Number of extracted chars
		s, _ := fs.Pop() // Start (0-indexed)
		a, _ := as.Pop()
		if int(s) < len(a) && int(s+n) <= len(a) && (s > 0) && (n > 0) {
			as.Push(a[int(s) : int(s)+int(n)])
		} else {
			raiseError(ALPHA_EXTRACT_OUT_OF_BOUNDS)
		}
	}
}

// *****************************************************************************
// doALen()
// *****************************************************************************
func doALen() {
	if checkAStack(1) {
		a, _ := as.Pop()
		fs.Push(float64(len(a)))
	}
}

// *****************************************************************************
// saveStacks()
// *****************************************************************************
func saveStacks() {
	// Serialize the Float64 stack on disk into the application folder
	fsFile, err := os.Create(filepath.Join(appDir, FSTACK_FILE))
	if err == nil {
		fsEncoder := gob.NewEncoder(fsFile)
		fsEncoder.Encode(fs.S)
	}
	fsFile.Close()
	// Serialize the string stack on disk into the application folder
	asFile, err := os.Create(filepath.Join(appDir, ASTACK_FILE))
	if err == nil {
		asEncoder := gob.NewEncoder(asFile)
		asEncoder.Encode(as.S)
	}
	asFile.Close()
	// Serialize the vars on disk into the application folder
	vFile, err := os.Create(filepath.Join(appDir, VARS_FILE))
	if err == nil {
		vEncoder := gob.NewEncoder(vFile)
		vEncoder.Encode(sto.Vars)
	}
	vFile.Close()
}

// *****************************************************************************
// readStacks()
// *****************************************************************************
func readStacks() {
	// Deserialize the previous Float64 stack stored into the application folder, if any
	fsFile, err := os.Open(filepath.Join(appDir, FSTACK_FILE))
	if err == nil {
		fsDecoder := gob.NewDecoder(fsFile)
		fsDecoder.Decode(&fs.S)
	}
	fsFile.Close()
	// Deserialize the previous string stack stored into the application folder, if any
	asFile, err := os.Open(filepath.Join(appDir, ASTACK_FILE))
	if err == nil {
		asDecoder := gob.NewDecoder(asFile)
		asDecoder.Decode(&as.S)
	}
	asFile.Close()
	// Deserialize the previous vars file stored into the application folder, if any
	vFile, err := os.Open(filepath.Join(appDir, VARS_FILE))
	if err == nil {
		vDecoder := gob.NewDecoder(vFile)
		vDecoder.Decode(&sto.Vars)
	}
	vFile.Close()
}

// *****************************************************************************
// doSwap()
// *****************************************************************************
func doSwap() {
	if checkFStack(2) {
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f2)
		fs.Push(f1)
	}
}

// *****************************************************************************
// doRot()
// *****************************************************************************
func doRot() {
	if checkFStack(3) {
		f3, _ := fs.Pop()
		f2, _ := fs.Pop()
		f1, _ := fs.Pop()
		fs.Push(f3)
		fs.Push(f1)
		fs.Push(f2)
	}
}

// *****************************************************************************
// showPrompt()
// *****************************************************************************
func showPrompt() {
	var prompt string
	// Do we have something into the stack to display
	if fs.Depth() > 0 {
		f := fs.S[fs.Depth()-1]
		// We use scientific notation if the number of digits is greater than 12
		if math.Log10(math.Abs(f)) > 12 {
			prompt = fmt.Sprintf("[%05d] %s%21.6E%s %s ", fs.Depth(), color.Green, f, color.Reset, ICON_ARROW)
		} else {
			prompt = fmt.Sprintf("[%05d] %s%21.6f%s %s ", fs.Depth(), color.Green, f, color.Reset, ICON_ARROW)
		}
	} else {
		// Nothing to display
		prompt = fmt.Sprintf("[%05d]           Empty stack %s ", fs.Depth(), ICON_ARROW)
	}
	fmt.Printf("%s", prompt)
}

// *****************************************************************************
// doRcl()
// *****************************************************************************
func doRcl() {
	if varName != "" {
		v := sto.Vars[varName]
		if v != nil {
			switch v.(type) {
			case float64:
				fs.Push(v.(float64))
			case string:
				as.Push(v.(string))
			}
			varName = ""
		} else {
			raiseError(NONEXISTENT_VARIABLE)
		}
	} else {
		raiseError(MISSING_VARIABLE)
	}
}

// *****************************************************************************
// doDel()
// *****************************************************************************
func doDel() {
	if varName != "" {
		v := sto.Vars[varName]
		if v != nil {
			delete(sto.Vars, varName)
			varName = ""
		} else {
			varName = ""
			raiseError(NONEXISTENT_VARIABLE)
		}
	} else {
		raiseError(MISSING_VARIABLE)
	}
}

// *****************************************************************************
// doSto()
// *****************************************************************************
func doSto() {
	if checkFStack(1) {
		if varName != "" {
			v, _ := fs.Pop()
			sto.Vars[varName] = v
			varName = ""
		} else {
			raiseError(MISSING_VARIABLE)
		}
	}
}

// *****************************************************************************
// doASto()
// *****************************************************************************
func doASto() {
	if checkAStack(1) {
		if varName != "" {
			v, _ := as.Pop()
			sto.Vars[varName] = v
			varName = ""
		} else {
			raiseError(MISSING_VARIABLE)
		}
	}
}

// *****************************************************************************
// showVars()
// *****************************************************************************
func showVars() {
	keys := reflect.ValueOf(sto.Vars).MapKeys()
	keysOrder := func(i, j int) bool { return keys[i].Interface().(string) < keys[j].Interface().(string) }
	sort.Slice(keys, keysOrder)

	// process map in key-sorted order
	for _, key := range keys {
		value := sto.Vars[key.Interface().(string)]
		switch value.(type) {
		case float64:
			if math.Log10(math.Abs(value.(float64))) > 12 {
				fmt.Printf("[VAR_F] %16s : %21.6E\n", key, value)
			} else {
				fmt.Printf("[VAR_F] %16s : %21.6f\n", key, value)
			}
		case string:
			fmt.Printf("[VAR_A] %16s : %s\n", key, value)
		}
	}
}
