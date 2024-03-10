# Polish | *a console RPN calculator written in Go*

## Summary

- Polish is a RPN calulator for text console environment.
- Polish is cross platform, tested on Linux, Windows and Android (Termux), should run on macOS.
- Polish is written in Go language.

## Authors

- Main developer : [**J.-P. Liguori**](https://github.com/jplozf/polish)
- Help and advices wanted...

## License

- This project is licensed under the GNU General Public License - see the [LICENSE.md](LICENSE.md) file for details.

## Features

* Algebraic functions
    - add, substract, multiply and divide
    - power, logarithms, exponential and factorial
* Trigonometric functions
    - sin, cos, tan, atan
    - degrees to/from radians conversions
    - degrees decimal to/from degrees, minutes, seconds conversion
* Date and Time functions
* Alpha (strings) stack and related functions
* Stack manipulations
* Stay tuned, more to come...

### [Reverse Polish Notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation) *(taken from Wikipedia)*

Reverse Polish notation (RPN), also known as reverse Łukasiewicz notation, Polish postfix notation or simply postfix notation, is a mathematical notation in which operators follow their operands, in contrast to prefix or Polish notation (PN), in which operators precede their operands. The notation does not need any parentheses for as long as each operator has a fixed number of operands.

### Numerical stack

The term postfix notation describes the general scheme in mathematics and computer sciences, whereas the term reverse Polish notation typically refers specifically to the method used to enter calculations into hardware or software calculators, which often have additional side effects and implications depending on the actual implementation involving a stack. The description "Polish" refers to the nationality of logician **Jan Łukasiewicz**, who invented Polish notation in 1924.

* In reverse Polish notation, the operators follow their operands. For example, to add **3** and **4** together, the expression is :
```
    ⯈ 3 4 + 
```
rather than **3 + 4**.

* The conventional notation expression **3 − 4 + 5** becomes :
```
    ⯈ 3 4 − 5 + 
```
in reverse Polish notation: **4** is first subtracted from **3**, then **5** is added to it.

The concept of a stack, a last-in/first-out construct, is integral to the left-to-right evaluation of RPN. In the example **3 4 −**, first the **3** is put onto the stack, then the **4**; the **4** is now on top and the **3** below it. The subtraction operator removes the top two items from the stack, performs **3 − 4**, and puts the result of **−1** onto the stack.

The common terminology is that added items are pushed on the stack and removed items are popped.

The advantage of reverse Polish notation is that it removes the need for order of operations and parentheses that are required by infix notation and can be evaluated linearly, left-to-right.

* For example, the infix expression **(3 × 4) + (5 × 6)** becomes :
```
    ⯈ 3 4 × 5 6 × +
```
in reverse Polish notation.

### Alphabetical stack

In addition to the float numbers stack, there is also an alphabetical strings stack.

* Alpha strings are entered beween double quotes
``` 
    ⯈ "hello" "world" +as
``` 
* If only one alpha string is entered on command line, the last quote could be omitted
``` 
    ⯈ "bye
``` 
  
### Variables

Finally, data, whether numeric or alphabetical, can also be stored in named variables.

* Variables names are introduced by the prefix $
``` 
    ⯈ $foo 5 sto
    ⯈ $foo rcl
``` 

### Data permanence

The floating point stack, alpha string stack, and variables are serialized to disk when the program exits. This data will be automatically restored the next time you open the interface.

## Algebraic Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| + | ( f1 f2 -- f3 ) | Add f1 to f2, giving the sum f3 |
| - | ( f1 f2 -- f3 ) | Substract f2 from f1, giving f3 |
| * | ( f1 f2 -- f3 ) | Multiply f1 by f2, giving f3 |
| / | ( f1 f2 -- f3 ) | Divide f1 by f2, giving the quotient f3 |
| ** pow | ( f1 f2 -- f3 ) | Raise f1 to the power of f2, giving f3 |
| abs | ( f1 -- f2 ) | f2 is the absolute value of f1 |
| ! fact | ( f1 -- f2 ) | f2 is the factorial of f1 |
| cbrt | ( f1 -- f2 ) | f2 is the cube root of f1 |
| ceil | ( f1 -- f2 ) | f2 is the least integer value greater than or equal to f1 |
| exp | ( f1 -- f2 ) | f2 is the base-e exponential of f1 |
| exp2 | ( f1 -- f2 ) | f2 is the base-2 exponential of f1 |
| floor | ( f1 -- f2 ) | f2 is the greatest integer value less than or equal to f1 |
| gamma | ( f1 -- f2 ) | f2 is the Gamma function of f1 |
| inv | ( f1 -- f2 ) | f2 is the value of 1/f1 |
| log | ( f1 -- f2 ) | f2 is the natural logarithm of f1 |
| log10 | ( f1 -- f2 ) | f2 is the decimal logarithm of f1 |
| log2 | ( f1 -- f2 ) | f2 is the binary logarithm of f1 |
| rnd | ( -- f ) | f is a pseudo-random number in the half-open interval [0.0,1.0) |
| round | ( f1 -- f2 ) | f2 is the nearest integer of f1, rounding half away from zero |
| sqr
| sqrt
| trunc
| frac
| root
| mod

## Trigonometric Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| acos | ( f1 -- f2 ) | f2 is the principal radian angle whose cosine is f1 |
| acosh | ( f1 -- f2 ) | f2 is the value whose hyperbolic cosine is f1 |
| asin | ( f1 -- f2 ) | f2 is the principal radian angle whose sine is f1 |
| asinh | ( f1 -- f2 ) | f2 is the value whose hyperbolic sine is f1 |
| atan | ( f1 -- f2 ) | f2 is the principal radian angle whose tangent is f1 |
| atanh | ( f1 -- f2 ) | f2 is the value whose hyperbolic tangent is f1 |
| atan2 | ( f1 f2 -- f3 ) | f3 is the principal radian angle (between -&pi; and &pi;) whose tangent is (f1/f2) |
| cos | ( f1 -- f2 ) | f2 is the cosine of the radian angle f1 |
| cosh | ( f1 -- f2 ) | f2 is the hyperbolic cosine of f1 |
| sin | ( f1 -- f2 ) | f2 is the sine of the radian angle f1 |
| sinh
| tan
| tanh
| hypot | ( f1 f2 -- f3 ) | f3 is the Sqrt(f1\*f1 + f2\*f2) |
| torad
| todeg
| todms
| todec

## Date and Time Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| now
| day
| month
| year
| hour
| minute
| second
| date
| ddiff
| d2000

## Alpha Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| +a aadd | ( a1 a2 -- a3 ) | Concatenate the top two alpha stack items, giving a3 |
| +as aadds | ( a1 a2 -- a3 ) | Concatenate the top two alpha stack items with a space between them, giving a3 |
| *a | ( a1 f -- a2 ) | a2 is the a1 string replicated f times |
| alen | ( a -- f ) | f is the length of the alpha string a |
| aright | ( a1 f -- a2 ) | a2 is the alpha string of the f right most characters of a1 |
| aleft | ( a1 f -- a2 ) | a2 is the alpha string of the f left most characters of a1 |
| amid | ( a1 f1 f2 -- a2 ) | a2 is the a1 substring starting from f1 (0-indexed) with f2 characters |
| ftoa | ( f -- a ) | Convert the float value f to an alpha string a |

## Misc Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| exit quit bye | ( -- ) | Exit the Polish program, stacks and variables are saved |
| .a | ( -- ) | Display the alpha stack |
| .f | ( -- ) | Display the float stack |
| .v | ( -- ) | Display all the variables stored |

## Stacks and Variables Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| drop | ( f -- ) | Remove f from the stack |
| dup | ( f -- f f ) | Duplicate the number f on the stack |
| depth | ( -- f ) | f is the number of values contained on the stack |
| swap | ( f1 f2 -- f2 f1 ) | Exchange the top two stack items |
| rot | ( f1 f2 f3 -- f2 f3 f1 ) | Rotate the top three stack entries |
| clear cls clr | ( -- ) | Clear the stack |
| rcl | ( $ -- a ) or ( $ -- f ) | Push the value of the $ variable |
| sto | ( $ f -- ) | Store the value of the float value f into the $ variable |
| del | ( $ -- ) | Delete the $ variable |
| adrop | ( a -- ) | Remove a from the alpha stack |
| adup | ( a -- a a ) | Duplicate the string a on the alpha stack | 
| aclear acls aclr | ( -- ) | Clear the alpha stack |
| adepth | ( -- f ) | f is the number of values contained on the alpha stack |
| aswap | ( a1 a2 -- a2 a1 ) | Exchange the top two alpha stack items |
| asto | ( $ a -- ) | Store the value of the alpha string a into the $ variable |

## Constants
| Constant | Value | Feature |
|:--------:|-------|---------|
| phi | 1.6180339887498948482045 | &phi; = Golden ratio |
| e | 2.7182818284590452353602 | e = Euler's constant |
| pi | 3.1415926535897932384626 | &pi;  = Ratio of a circle's circumference to its diameter |