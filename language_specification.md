# Gosling Language Specification

## Table of Contents
1. [Introduction](#introduction)
2. [Lexical Elements](#lexical-elements)
3. [Data Types](#data-types)
4. [Variables and Bindings](#variables-and-bindings)
5. [Expressions](#expressions)
6. [Statements](#statements)
7. [Functions](#functions)
8. [Control Flow](#control-flow)
9. [Operators](#operators)
10. [Comments](#comments)
11. [Grammar](#grammar)
12. [Examples](#examples)

## Introduction

Gosling is a dynamically typed, interpreted programming language with a focus on simplicity and expressiveness. It supports first-class functions, lexical scoping, and basic control flow constructs.
I started off by following along with [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball. and started extending the language and the repl.
I had no intention of naming this after James Gosling (writer of java), and actually was thinking about ducks and the name Gosling was just a random name I thought of with an ending that would be easy to type (.gos).

### Design Goals
- Simple and readable syntax
- Dynamic typing
- First-class functions
- Lexical scoping
- Interactive REPL environment

## Lexical Elements

### Keywords
The following identifiers are reserved keywords and cannot be used as variable names:

```
fn      let     true    false   if      else    return  for
```

### Identifiers
Identifiers are sequences of letters, digits, and underscores, starting with a letter or underscore.

```
identifier = letter { letter | digit | "_" }
letter     = "a" ... "z" | "A" ... "Z" | "_"
digit      = "0" ... "9"
```

**Examples:**
```gosling
x
myVariable
_private
counter1
```

### Literals

#### Integer Literals
Integer literals are sequences of digits representing decimal numbers.

```gosling
42
0
123456
```

#### Boolean Literals
Boolean literals represent truth values.

```gosling
true
false
```

#### String Literals
String literals are sequences of characters enclosed in double quotes.

```gosling
"hello world"
"Gosling is awesome"
""
```

## Data Types

Gosling supports the following built-in data types:

### Integer
64-bit signed integers.

```gosling
let age = 25;
let negative = -42;
```

### Boolean
Truth values: `true` or `false`.

```gosling
let isActive = true;
let isComplete = false;
```

### String
Sequences of Unicode characters.

```gosling
let name = "Alice";
let message = "Hello, World!";
```

### Function
First-class function objects.

```gosling
let add = fn(a, b) { return a + b; };
```

### Null
Can only result from branches that don't produce a value. Users cannot assign null, and any instance of a null being evaluated will result in an error.

```gosling
let a = 5;
let b = 6;
if (a > b) { let c = a + b; } // this branch is a null object and is skipped when assigning value.
```

## Variables and Bindings

Variables are declared using the `let` keyword and are immutable once bound.

### Syntax
```gosling
let identifier = expression;
```

### Examples
```gosling
let x = 42;
let name = "Bob";
let isValid = true;
let calculate = fn(n) { return n * 2; };
```

### Scoping
Gosling uses lexical scoping. Variables are accessible within the scope where they are defined and any nested scopes.

```gosling
let outer = 10;
let func = fn() {
    let inner = 20;
    return outer + inner;  // outer is accessible here
};
```

## Expressions

### Primary Expressions
- Identifiers: `x`, `myVar`
- Literals: `42`, `true`, `"hello"`
- Parenthesized expressions: `(expression)`

### Function Calls
```gosling
functionName(arg1, arg2, ...)
```

### Arithmetic Expressions
```gosling
a + b    // addition
a - b    // subtraction
a * b    // multiplication
a / b    // division
a % b    // modulo
```

### Comparison Expressions
```gosling
a == b   // equality
a != b   // inequality
a < b    // less than
a > b    // greater than
```

### Logical Expressions
```gosling
!a       // logical NOT
```

### Precedence (highest to lowest)
1. Function calls: `f()`
2. Unary operators: `-`, `!`
3. Multiplicative: `*`, `/`, `%`
4. Additive: `+`, `-`
5. Comparison: `<`, `>`, `==`, `!=`

## Statements

### Expression Statements
Any expression can be used as a statement.

```gosling
42;
add(1, 2);
x + y;
```

### Let Statements
Variable declarations and bindings.

```gosling
let x = 10;
let result = calculate(x);
```

### Return Statements
Return a value from a function.

```gosling
return expression;
return 42;
return;  // returns null
```

## Functions

Functions are first-class values in Gosling and support closures.

### Function Literals
```gosling
fn(parameter1, parameter2, ...) {
    // function body
    return expression;
}
```

### Examples
```gosling
// Simple function
let add = fn(a, b) {
    return a + b;
};

// Function with no parameters
let greet = fn() {
    return "Hello!";
};

// Function with closure
let makeCounter = fn() {
    let count = 0;
    return fn() {
        count = count + 1;
        return count;
    };
};
```

### Function Calls
```gosling
add(5, 3);
greet();
let counter = makeCounter();
counter();  // returns 1
counter();  // returns 2
```

## Control Flow

### If Expressions
If expressions evaluate a condition and return a value based on the result.

```gosling
if (condition) {
    // consequence
} else {
    // alternative (optional)
}
```

### Examples
```gosling
let result = if (x > 0) { "positive" } else { "non-positive" };

let abs = fn(n) {
    if (n < 0) {
        return -n;
    } else {
        return n;
    }
};
```

### For Loops
For loops repeatedly execute a block while a condition is true.

```gosling
for (condition) {
    // loop body
}
```

### Example
```gosling
let i = 0;
for (i < 10) {
    // loop body
    i = i + 1;
}
```

## Operators

### Arithmetic Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `5 + 3` → `8` |
| `-` | Subtraction | `5 - 3` → `2` |
| `*` | Multiplication | `5 * 3` → `15` |
| `/` | Division | `6 / 3` → `2` |
| `%` | Modulo | `7 % 3` → `1` |

### Comparison Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `==` | Equal | `5 == 5` → `true` |
| `!=` | Not equal | `5 != 3` → `true` |
| `<` | Less than | `3 < 5` → `true` |
| `>` | Greater than | `5 > 3` → `true` |

### Logical Operators
| Operator | Description | Example |
|----------|-------------|---------|
| `!` | Logical NOT | `!true` → `false` |

### Assignment Operator
| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Assignment | `let x = 5;` |

## Comments

Currently, Gosling does not support comments in the language syntax. This may be added in future versions.

## Grammar

```ebnf
Program = { Statement } .

Statement = LetStatement | ReturnStatement | ExpressionStatement .

LetStatement = "let" identifier "=" Expression ";" .

ReturnStatement = "return" [ Expression ] ";" .

ExpressionStatement = Expression [ ";" ] .

Expression = IfExpression | ForExpression | FunctionLiteral | CallExpression | InfixExpression | PrefixExpression | Primary .

IfExpression = "if" "(" Expression ")" BlockStatement [ "else" BlockStatement ] .

ForExpression = "for" "(" Expression ")" BlockStatement .

FunctionLiteral = "fn" "(" [ ParameterList ] ")" BlockStatement .

ParameterList = identifier { "," identifier } .

CallExpression = Expression "(" [ ArgumentList ] ")" .

ArgumentList = Expression { "," Expression } .

InfixExpression = Expression InfixOperator Expression .

PrefixExpression = PrefixOperator Expression .

Primary = identifier | IntegerLiteral | BooleanLiteral | StringLiteral | "(" Expression ")" .

BlockStatement = "{" { Statement } "}" .

InfixOperator = "+" | "-" | "*" | "/" | "%" | "==" | "!=" | "<" | ">" .

PrefixOperator = "-" | "!" .

IntegerLiteral = digit { digit } .

BooleanLiteral = "true" | "false" .

StringLiteral = '"' { character } '"' .

identifier = letter { letter | digit | "_" } .
```

## Examples

### Basic Arithmetic
```gosling
let a = 10;
let b = 5;
let sum = a + b;        // 15
let product = a * b;    // 50
let remainder = a % 3;  // 1
```

### Functions and Closures
```gosling
// Factorial function
let factorial = fn(n) {
    if (n <= 1) {
        return 1;
    } else {
        return n * factorial(n - 1);
    }
};

// Closure example
let makeAdder = fn(x) {
    return fn(y) {
        return x + y;
    };
};

let add5 = makeAdder(5);
let result = add5(3);  // 8
```

### Control Flow
```gosling
// Conditional logic
let max = fn(a, b) {
    if (a > b) {
        return a;
    } else {
        return b;
    }
};

// Loop example
let sum = 0;
let i = 1;
for (i <= 10) {
    sum = sum + i;
    i = i + 1;
}
// sum is now 55
```

### Higher-Order Functions
```gosling
let apply = fn(f, x) {
    return f(x);
};

let double = fn(n) {
    return n * 2;
};

let result = apply(double, 21);  // 42
```

## Error Handling

Gosling provides runtime error reporting with location information:

- **Division by zero**: `5 / 0`
- **Modulo by zero**: `5 % 0`
- **Unknown identifier**: Using an undefined variable
- **Type errors**: Applying operations to incompatible types
- **Unknown operators**: Using unsupported operator combinations

Errors include file name, line number, and character position when available.

## Future Considerations

The following features may be considered for future versions:

- Comments (`//` and `/* */`)
- Arrays and indexing
- Hash maps/objects
- String interpolation
- More built-in functions
- Import/module system
- Pattern matching
- Exception handling