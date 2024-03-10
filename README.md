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

* Algebrical functions
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

## Functions
| Function | Stack | Feature |
|:--------:|:-----:|---------|
| + | ( n1 n2 -- n3 ) | Add n1 to n2, giving the sum n3 |
| - | ( n1 n2 -- n3 ) | Substract n2 from n1, giving n3 |
| * | ( n1 n2 -- n3 ) | Multiply n1 by n2, giving n3 |
| / | ( n1 n2 -- n3 ) | Divide n1 by n2, giving the quotient n3 |
| ** pow | ( n1 n2 -- n3 ) | Raise n1 to the power of n2, giving n3 |
| abs | ( n1 -- n2 ) | n2 is the absolute value of n1 |
| acos | ( n1 -- n2 ) | n2 is the principal radian angle whose cosine is n1 |
| acosh | ( n1 -- n2 ) | n2 is the value whose hyperbolic cosine is n1 |
| asin | ( n1 -- n2 ) | n2 is the principal radian angle whose sine is n1 |
| asinh | ( n1 -- n2 ) | n2 is the value whose hyperbolic sine is n1 |
| atan | ( n1 -- n2 ) | n2 is the principal radian angle whose tangent is n1 |
| atanh | ( n1 -- n2 ) | n2 is the value whose hyperbolic tangent is n1 |
| atan2 | ( n1 n2 -- n3 ) | n3 is the principal radian angle (between -&pi; and &pi;) whose tangent is (n1/n2) |
| ! fact | ( n1 -- n2 ) | n2 is the factorial of n1 |
| drop | ( n -- ) | Remove n from the stack |
| exit quit bye | ( -- ) | Exit the Polish program, stacks and variables are saved |
| dup | ( n -- n n ) | Duplicate the number n on the stack |
| depth | ( -- n ) | n is the number of values contained on the stack |
| swap | ( n1 n2 -- n2 n1 ) | Exchange the top two stack items |
| rot | ( n1 n2 n3 -- n2 n3 n1 ) | Rotate the top three stack entries |
| clear cls clr | ( -- ) | Clear the stack |
| adrop | ( a -- ) | Remove a from the alpha stack |
| adup | ( a -- a a ) | Duplicate the string a on the alpha stack | 
| aclear acls aclr | ( -- ) | Clear the alpha stack |
| adepth | ( -- n ) | n is the number of values contained on the alpha stack |
| aswap | ( a1 a2 -- a2 a1 ) | Exchange the top two alpha stack items |
| +a aadd | ( a1 a2 -- a3 ) | Concatenate the top two alpha stack items, giving a3 |
| +as aadds | ( a1 a2 -- a3 ) | Concatenate the top two alpha stack items with a space between them, giving a3 |
| *a | ( a1 n -- a2 ) | a2 is the a1 string replicated n times |
| alen
| aright
| aleft
| amid
| rcl
| sto
| asto
| del 
| .a
| .f
| .v
| cbrt
| ceil
| cos | ( n1 -- n2 ) | n2 is the cosine of the radian angle n1 |
| cosh | ( n1 -- n2 ) | n2 is the hyperbolic cosine of n1 |
| exp
| exp2
| floor | ( n1 -- n2 ) | n2 is the greatest integer value less than or equal to n1 |
| gamma | ( n1 -- n2 ) | n2 is the Gamma function of n1 |
| hypot | ( n1 n2 -- n3 ) | n3 is the Sqrt(n1\*n1 + n2\*n2) |
| inv | ( n1 -- n2 ) | n2 is the value of 1/n1 |
| log | ( n1 -- n2 ) | n2 is the natural logarithm of n1 |
| log10 | ( n1 -- n2 ) | n2 is the decimal logarithm of n1 |
| rnd
| round
| sin
| sinh
| sqr
| sqrt
| tan
| tanh
| trunc
| frac
| pi
| e
| phi
| root
| torad
| todeg
| now
| todms
| todec
| day
| month
| year
| hour
| minute
| second
| date
| ddiff
| mod
| d2000
