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
	"reflect"
	"strconv"
	"strings"
)

// *****************************************************************************
// GLOBALS
// *****************************************************************************
var (
	loopme   = true
	s        *stack.Stack
	appDir   string
	previous string
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
	// Create the stack
	s = stack.NewStack()
	// Deserialize the previous stack, if any
	readStack()
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
	// Serialize the current stack
	saveStack()
	fmt.Printf("\n%s Bye.\n\n", ICON_DISK)
}

// *****************************************************************************
// greetings()
// *****************************************************************************
func greetings() {
	fmt.Printf("%s Welcome to %s\n", ICON_DISK, APP_STRING)
	fmt.Printf("%s %s version %s\n", ICON_DISK, APP_NAME, APP_VERSION)
	fmt.Printf("%s %s\n\n", ICON_DISK, APP_URL)
}

// *****************************************************************************
// parse()
// *****************************************************************************
func parse(txt string) {
	words := strings.Fields(txt)
	for _, w := range words {
		xeq(w)
	}
}

// *****************************************************************************
// xeq()
// *****************************************************************************
func xeq(cmd string) {
	// Is it a number ?
	if isFloat(cmd) {
		// Then push it on the stack
		v, _ := strconv.ParseFloat(cmd, 64)
		s.Push(v)
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
			case "show", ".s":
				showStack()
			case "swap":
				doSwap()
			case "rot":
				doRot()
			default:
				fmt.Printf("\t"+color.Red+"Unrecognized command '%s'\n"+color.Reset, cmd)
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
// checkStack()
// *****************************************************************************
func checkStack(n int) bool {
	// Do we have enough args on stack to perform the selected operation ?
	if s.Depth() >= n {
		return true
	} else {
		fmt.Println("\t" + color.Red + "Not enough arguments on stack" + color.Reset)
		return false
	}
}

// *****************************************************************************
// doDrop()
// *****************************************************************************
func doDrop() {
	if checkStack(1) {
		s.Pop()
	}
}

// *****************************************************************************
// doDup()
// *****************************************************************************
func doDup() {
	if checkStack(1) {
		f, _ := s.Pop()
		s.Push(f)
		s.Push(f)
	}
}

// *****************************************************************************
// doClear()
// *****************************************************************************
func doClear() {
	s.S = nil
}

// *****************************************************************************
// doDepth()
// *****************************************************************************
func doDepth() {
	s.Push(float64(len(s.S)))
}

// *****************************************************************************
// showStack()
// *****************************************************************************
func showStack() {
	for i, value := range s.S {
		i = len(s.S) - 1 - i
		if math.Log10(math.Abs(value)) > 12 {
			fmt.Printf("\t%05d : %21.6E\n", i, value)
		} else {
			fmt.Printf("\t%05d : %21.6f\n", i, value)
		}
	}
}

// *****************************************************************************
// saveStack()
// *****************************************************************************
func saveStack() {
	// Serialize the stack on disk into the application folder
	dataFile, err := os.Create(filepath.Join(appDir, STACK_FILE))
	if err == nil {
		dataEncoder := gob.NewEncoder(dataFile)
		dataEncoder.Encode(s.S)
	}
	dataFile.Close()
}

// *****************************************************************************
// readStack()
// *****************************************************************************
func readStack() {
	// Deserialize the previous stack stored into the application folder, if any
	dataFile, err := os.Open(filepath.Join(appDir, STACK_FILE))
	if err == nil {
		dataDecoder := gob.NewDecoder(dataFile)
		dataDecoder.Decode(&s.S)
	}
	dataFile.Close()
}

// *****************************************************************************
// doSwap()
// *****************************************************************************
func doSwap() {
	if checkStack(2) {
		f2, _ := s.Pop()
		f1, _ := s.Pop()
		s.Push(f2)
		s.Push(f1)
	}
}

// *****************************************************************************
// doRot()
// *****************************************************************************
func doRot() {
	if checkStack(3) {
		f3, _ := s.Pop()
		f2, _ := s.Pop()
		f1, _ := s.Pop()
		s.Push(f3)
		s.Push(f1)
		s.Push(f2)
	}
}

// *****************************************************************************
// showPrompt()
// *****************************************************************************
func showPrompt() {
	var prompt string
	// Do we have something into the stack to display
	if len(s.S) > 0 {
		f := s.S[len(s.S)-1]
		// We use scientific notation if the number of digits is greater than 12
		if math.Log10(math.Abs(f)) > 12 {
			prompt = fmt.Sprintf("[%05d] %s%21.6E%s %s ", s.Depth(), color.Green, f, color.Reset, ICON_ARROW)
		} else {
			prompt = fmt.Sprintf("[%05d] %s%21.6f%s %s ", s.Depth(), color.Green, f, color.Reset, ICON_ARROW)
		}
	} else {
		// Nothing to display
		prompt = fmt.Sprintf("[%05d]           Empty stack %s ", s.Depth(), ICON_ARROW)
	}
	fmt.Printf("%s", prompt)
}
