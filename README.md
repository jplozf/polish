# Polish RPN Interpreter

Polish is a powerful and interactive Reverse Polish Notation (RPN) interpreter, designed for command-line use. It supports a wide range of operations, including arithmetic, stack manipulation, control flow, variable and word (function) definitions, string operations, and file management.

## Features

*   **Basic Arithmetic & Math Functions:** Perform standard calculations, trigonometry, logarithms, and more.
*   **Stack-based Operations:** Manipulate data directly on the stack with commands like `dup`, `drop`, and `swap`.
*   **Variables & Words (Functions):** Define and manage variables and custom words (functions) for reusable code.
*   **Control Flow:** Implement conditional logic (`if`) and loops (`loop`, `while`).
*   **String Manipulation:** Work with strings, including length, substrings, and case conversion.
*   **File & State Management:** Save and restore interpreter state, import/export RPN scripts, and list available files.
*   **Interactive Editing:** Edit RPN files and code blocks directly within the interpreter's TUI.
*   **History & Persistence:** Command history and interpreter state are saved across sessions.
*   **Configurable Modes:** Toggle angle mode (degrees/radians), echo mode, and variable/stack display options.

## Installation

To build and run Polish, you need Go installed on your system.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/jplozf/polish.git
    cd polish
    ```
2.  **Build the executable:**
    ```bash
    go build -o polish
    ```
3.  **Run the interpreter:**
    ```bash
    ./polish
    ```

## Usage

Upon running `./polish`, you will be presented with an interactive Text User Interface (TUI). Type RPN commands and press Enter to execute them.

### Basic Interaction

*   Enter numbers to push them onto the stack.
*   Enter operators or commands to perform actions on the stack.
*   Enter `exit`, `quit`, `bye` or press the `F12` key to exit the interpreter.
*   Enter `help` or press the `F1` key to show this help.
*   Press the `F2` key to switch the current panel.

## Syntax and Examples

### Numbers and Basic Arithmetic

Numbers are pushed onto the stack. Operators pop arguments, perform calculations, and push the result.

```rpn
10 5 +    ( Result: 15 )
10 5 -    ( Result: 5 )
10 5 *    ( Result: 50 )
10 5 /    ( Result: 2 )
10 3 %    ( Result: 1 )
```

### Stack Manipulation

*   `dup`: Duplicates the top item on the stack.
*   `drop`: Removes the top item from the stack.
*   `swap`: Swaps the top two items on the stack.
*   `clear`: Clears the entire stack.
*   `depth`: Pushes the current stack depth onto the stack.

```rpn
1 2 dup   ( Stack: 1 2 2 )
1 2 drop  ( Stack: 1 )
1 2 swap  ( Stack: 2 1 )
1 2 3 clear ( Stack: empty )
1 2 3 depth ( Stack: 1 2 3 3 )
```

### Variables

Variables can store any data type, including numbers, strings, booleans, and code blocks.

*   `value "name" store`: Stores `value` into `name`.
*   `"name" load`: Loads the value of `name` onto the stack.
*   `delete [var:]<name>`: Deletes a variable.

```rpn
10 "x" store    ( Stores 10 in variable x )
"x" load        ( Pushes 10 onto the stack )
"Hello" "msg" store
"msg" load .    ( Prints "Hello" )
{ 1 2 + } "add_block" store ( Stores a code block )
"add_block" load execute    ( Executes the code block )
delete "x"      ( Deletes variable x )
```

### Words (User-Defined Functions)

Words are named code blocks that can be executed.

*   `: word_name ... ;`: Defines a new word.
*   `word_name`: Executes the defined word.
*   `delete [word:]<name>`: Deletes a word.

```rpn
: greet "Hello, World!" . cr ; ( Defines a word named 'greet' )
greet                         ( Executes 'greet', prints "Hello, World!" )
: square dup * ;              ( Defines a word to square a number )
5 square                      ( Result: 25 )
delete word:greet             ( Deletes the word 'greet' )
```

### Control Flow

*   `condition { then_block } if`: Executes `then_block` if `condition` is true (non-zero for numbers).
*   `condition { then_block } { else_block } if`: Executes `then_block` if `condition` is true, `else_block` otherwise.
*   `count { loop_block } loop`: Executes `loop_block` `count` times. `index` can be used inside the loop.
*   `{ condition_block } { body_block } while`: Executes `body_block` repeatedly as long as `condition_block` evaluates to true.
*   `break`: Exits the current `loop` or `while` block.
*   `continue`: Skips to the next iteration of the current `loop` or `while` block.
*   `index`: Pushes the current loop index onto the stack (0-based).

```rpn
true { "It's true!" . cr } if
false { "This won't run." . cr } { "This will run." . cr } if

3 { "Loop iteration " index . cr } loop
( Output:
  Loop iteration 0
  Loop iteration 1
  Loop iteration 2
)

: countdown 5 ;
{ countdown 0 > } { countdown 1 - "countdown" store countdown . cr } while
( Output:
  4
  3
  2
  1
  0
)
```

### String Manipulation

*   `"string" len`: Pushes the length of the string.
*   `"string" start length mid`: Extracts a substring.
*   `"string" upper`: Converts string to uppercase.
*   `"string" lower`: Converts string to lowercase.
*   `"value" val`: Tries to convert a string to a number or boolean.
*   `value str`: Converts any value to its string representation.
*   `"string" prompt`: Displays the string in the input box and waits for an input.

```rpn
"hello" len       ( Result: 5 )
"world" 1 3 mid   ( Result: "orl" )
"hello" upper     ( Result: "HELLO" )
"WORLD" lower     ( Result: "world" )
"123" val         ( Result: 123 )
true str          ( Result: "true" )
```

### File and State Management

*   `"filename.json" save`: Saves the current interpreter stack, variables, and words to a JSON file.
*   `"filename.json" restore`: Loads interpreter state from a JSON file.
*   `"script.rpn" import`: Executes commands from an RPN file. If a word named `main` is defined in the file, it will be executed automatically after import.
*   `"filename.rpn" "word_name" export`: Exports a defined word to an RPN file.
*   `list`: Lists all `.rpn` files in the interpreter's data directory (`~/.polish`).

```rpn
1 2 3 "mystate.json" save
clear
"mystate.json" restore ( Stack: 1 2 3 )

( Assuming test.rpn contains: : main "Hello from import!" . cr ; )
"test.rpn" import ( Prints "Hello from import!" )
```

### Interactive Editing

*   `editfile "filename.rpn"`: Opens a multi-line editor for an RPN file.
    *   `Ctrl+S`: Save and execute the file content.
    *   `Esc`: Cancel editing.
*   `edit "word_name"` or `edit "var:variable_name"`: Opens the definition of a word or a code block variable in the editor.

```rpn
editfile "my_script.rpn"
edit "my_word"
edit "var:my_block_variable"
```

### Internal Variables (Configuration)

These variables start with `_` and control interpreter behavior. Use `set`, `unset`, or `toggle` to modify them.

*   `_echo_mode`: `true` to echo input commands, `false` otherwise.
*   `_degree_mode`: `true` for degrees in trigonometric functions, `false` for radians.
*   `_vars_value`: `true` to show variable values in the variables view, `false` to show types.
*   `_stack_type`: `true` to show stack item types, `false` to show values.
*   `_hidden_vars`: `true` to show internal variables in the variables view, `false` to hide them.
*   `_exit_save`: `true` to automatically save state to `default.json` on exit.
*   `_last_x`: Stores the last value popped from the stack.

```rpn
"_echo_mode" toggle
"_degree_mode" set
"_vars_value" unset
```

### Comments

Comments start with `(` and end with `)`. They can be nested.

```rpn
( This is a single-line comment )
( This is a multi-line comment
  ( with nested comments )
  that spans multiple lines
)
```

### Error Handling

The interpreter provides specific error messages for various issues, such as stack underflow, type errors, undefined variables, and syntax errors.

*   `_last_error`: Contains the code of the last error. This is a read-only variable.

```rpn
; ( Will now show: error 57: semicolon out of context )
1 + ( Will show: error 1: stack underflow )
"hello" 5 + ( Will show: error 7: type error: '+' requires two numbers or two strings, got string and float64 )
```

## Contributing

Contributions are welcome! Please feel free to open issues or pull requests on the GitHub repository.

## License

This project is licensed under the GNU General Public License - see the [LICENSE.md](LICENSE.md) file for details.
