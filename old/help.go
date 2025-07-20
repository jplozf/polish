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

type outHelp struct {
	stack    string
	helpText string
}

var (
	help map[string][]outHelp
)

func init() {
/*
	help["+"] = append(help["+"], outHelp{"( f1 f2 -- f3 )", "Add f1 to f2, giving the sum f3"})
	help["-"] = append(help["-"], outHelp{"( f1 f2 -- f3 )", "Substract f2 from f1, giving f3"})
	help["*"] = append(help["*"], outHelp{"( f1 f2 -- f3 )", "Multiply f1 by f2, giving f3"})
	help["/"] = append(help["/"], outHelp{"( f1 f2 -- f3 )", "Divide f1 by f2, giving the quotient f3"})
	help["**"] = append(help["**"], outHelp{"( f1 f2 -- f3 )", "Raise f1 to the power of f2, giving f3"})
	help["pow"] = append(help["pow"], outHelp{"( f1 f2 -- f3 )", "Raise f1 to the power of f2, giving f3"})
	help["abs"] = append(help["abs"], outHelp{"( f1 -- f2 )", "f2 is the absolute value of f1"})
*/
}

/*
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
*/
