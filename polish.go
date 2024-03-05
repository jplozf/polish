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
	loopme = true
	s      *stack.Stack
	prompt string
	appDir string
)

// *****************************************************************************
// init()
// *****************************************************************************
func init() {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	appDir = filepath.Join(userDir, APP_FOLDER)
	if _, err := os.Stat(appDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(appDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	greetings()
	s = stack.NewStack()
	readStack()
}

// *****************************************************************************
// main()
// *****************************************************************************
func main() {
	reader := bufio.NewReader(os.Stdin)

	for ok := true; ok; ok = loopme {
		if len(s.S) > 0 {
			prompt = fmt.Sprintf("[%05d] %s%20.6f%s ⯈ ", s.Depth(), color.Green, s.S[len(s.S)-1], color.Reset)
		} else {
			prompt = fmt.Sprintf("[%05d]          Empty stack ⯈ ", s.Depth())
		}
		fmt.Printf("%s", prompt)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		parse(text)
	}
	saveStack()
	fmt.Printf("\nBye.\n\n")
}

// *****************************************************************************
// greetings()
// *****************************************************************************
func greetings() {
	fmt.Printf("Welcome to %s\n", APP_STRING)
	fmt.Printf("%s version %s\n", APP_NAME, APP_VERSION)
	fmt.Printf("%s\n\n", APP_URL)
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
	if isFloat(cmd) {
		v, _ := strconv.ParseFloat(cmd, 64)
		s.Push(v)
	} else {
		m := My{}
		mName := "My" + strings.Title(strings.ToLower(cmd))
		meth := reflect.ValueOf(m).MethodByName(mName)
		if meth.IsValid() {
			meth.Call(nil)
		} else {
			switch cmd {
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
}

// *****************************************************************************
// isFloat()
// *****************************************************************************
func isFloat(c string) bool {
	rc := true
	_, err := strconv.ParseFloat(c, 64)
	if err != nil {
		rc = false
	}
	return rc
}

// *****************************************************************************
// checkStack()
// *****************************************************************************
func checkStack(n int) bool {
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
	for k, v := range s.S {
		k = len(s.S) - 1 - k
		fmt.Printf("\t%05d : %20.6f\n", k, v)
	}
}

// *****************************************************************************
// saveStack()
// *****************************************************************************
func saveStack() {
	dataFile, err := os.Create(filepath.Join(appDir, STACK_FILE))
	defer dataFile.Close()

	if err == nil {
		dataEncoder := gob.NewEncoder(dataFile)
		dataEncoder.Encode(s.S)
	}
}

// *****************************************************************************
// readStack()
// *****************************************************************************
func readStack() {
	dataFile, err := os.Open(filepath.Join(appDir, STACK_FILE))
	defer dataFile.Close()

	if err == nil {
		dataDecoder := gob.NewDecoder(dataFile)
		dataDecoder.Decode(&s.S)
	}
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
