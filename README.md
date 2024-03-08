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
| abs | ( n1 -- n2 ) | n2 is the absolute value of n1 |
| acos | ( n1 -- n2 ) | n2 is the principal radian angle whose cosine is n1 |
| acosh | ( n1 -- n2 ) | n2 is the value whose hyperbolic cosine is n1 |
| asin | ( n1 -- n2 ) | n2 is the principal radian angle whose sine is n1 |
| asinh | ( n1 -- n2 ) | n2 is the value whose hyperbolic sine is n1 |
| atan | ( n1 -- n2 ) | n2 is the principal radian angle whose tangent is n1 |
| atanh | ( n1 -- n2 ) | n2 is the value whose hyperbolic tangent is n1 |
| atan2 | ( n1 n2 -- n3 ) | n3 is the principal radian angle (between -&pi; and &pi;) whose tangent is (n1/n2) |