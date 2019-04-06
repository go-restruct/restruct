# Restruct Expr

Expr is a small expression language for Restruct. It allows for simple expressions that can be embedded inside of struct tags. The primary advantage of embedding the language in struct tags is that it keeps all of the serialization and deserialization information inside of the structure itself, allowing us to use nothing more than the raw structure definition to generate fast serialization/deserialization code.

## Motivation

I discovered a library called [Kaitai Struct](https://kaitai.io) not long ago. Kaitai Struct is an amazing library that allows you to use declarative structure definitions written in YAML to define binary structures in a rich manner, not unlike how Restruct uses Go structure definitions. I have been quite impressed with its flexibility and tooling, but it does have some weaknesses; one of the current weaknesses, that I have no doubt will be eventually resolved, is that it only supports deserialization at the moment. I think the world is better if both Kaitai Struct and Restruct live side by side, especially since Restruct offers a pure Go solution. However, I also want to give Restruct some of Kaitai's super powers.

## Implementation

The current implementation uses [goccmack/gocc](https://github.com/goccmack/gocc) to generate both a lexer and parser. This is mostly done for simplicity; a previous prototype used both a handwritten lexer and parser, but it proved very error prone in comparison. (It would be possible, and potentially preferrable, to rewrite the handwritten lexer in the future, since it only needs 1 look-ahead, and the gocc lexer only supports passing a fully-read byte array. However, the parser is likely to stay generated, since handwritting an LALR parser is very cumbersome, and it is unclear what benefit it would bring.)

The AST is handwritten and tries very hard to keep the structure simple. Instead of using enumerations, each possible type of operation is its own kind of AST node. While a little cumbersome to maintain, the hope is that this will help prevent certain kinds of mistakes in the code.

Both the evaluator and AST depend on the expr/value library, which essentially handles the bulk of the language. Most operations are implemented in terms of values, and any constants are encoded as AST nodes containing a value.

The evaluator library implements the base language. The evaluator simply recursively evaluates each expression and returns either a result value or an error.

## Language

The expr language is simple. It follows a grammar that is somewhere between Go and C. A simple expression might be:

```go
len(Field) + 1
```

This document defines the intended behavior of expr, but for now the reference interpreter is the source of truth.

Expr does not define when type checking or compilation may occur. As expr is designed to be compiled into corresponding Go (or perhaps even other language) code, there may be situations where some expressions do not compile when using AOT compilation even though they can run in the interpreter. This is by design, though the interpreter may become more strict over time.

### Identifiers

Expr is focused around structures. As such, there is always a 'context' structure for any given expr context. A structure field can be referred to directly.

Some built-in functions also exist, and these can be referred to directly as well. Since built-in functions are lowercase, they should never conflict with accessible structure fields.

### Built-in functions

Expr has the following built-in functions:

| Definition | Description |
|---|---|
| `int(x)` | Converts ints, uints, and floats into ints. |
| `uint(x)` | Converts ints, uints, and floats into uints.|
| `float(x)` | Converts ints, uints, and floats into floats. |
| `len(x)` | Returns the length of arrays, maps, and strings. |
| `first(x)` | Returns the first element of an array. |
| `last(x)` | Returns the last element of an array. |
| `sum(...)` | Returns sum of all arguments recursively, as int. |
| `usum(...)` | Returns sum of all arguments recursively, as uint. |
| `fsum(...)` | Returns sum of all arguments recursively, as float. |

### Types

Expr has a simplified type system. Because expr is just an expression language, there is no way to refer to types within the language, but there are limitations on what types can exist. The following types exist:

| Type | Description |
|---|---|
| array | An array or slice. Any supported type can be the element type. |
| map | A map. Any supported types can be keys or values. |
| func | A function. Only functions returning exactly 1 argument can be called. |
| struct | A structure reference. |
| bool | A boolean value. |
| int | A signed 64-bit integer. |
| uint | An unsigned 64-bit integer. |
| float | A 64-bit floating point number. |
| string | A UTF-8 string. |

Similar to Go, when an integer is encountered, it is initially not typed. An untyped integer can morph into an int, uint, or float as needed. No other implicit conversions can occur.

### Operators

Like C and Go, expr uses an infix notation for most operators. Most binary operators, such as mathematical operators, require operands to be of the exact same type. Shift operators require uint values for the shift operand. Here is a table of operators, sorted by precedence:

|  Operator |  Name | Precedence | Valid Types |
|---|---|---|---|
| `a[b]` | Index operator. | 7 | array, map |
| `a(...)` | Call operator. | 7 | func |
| `a.b` | Descend operator. | 7 | struct |
| `-a` | Negation operator. | 6 | int, uint, float |
| `!a` | Logical not operator. | 6 | bool |
| `^a` | Bitwise not operator. | 6 | int, uint |
| `a * b` | Multiplication operator. | 5 | int, uint, float |
| `a / b` | Division operator. | 5 | int, uint, float |
| `a % b` | Modulo operator. | 5 | int, uint |
| `a << b` | Bitwise left shift operator. | 5 | int, uint; uint |
| `a >> b` | Bitwise right shift operator. | 5 | int, uint; uint |
| `a & b` | Bitwise and operator. | 5 | int, uint |
| `a &^ b` | Bitwise clear operator. | 5 | int, uint |
| `a + b` | Addition operator. | 4 | int, uint, float |
| `a - b` | Subtraction operator. | 4 | int, uint, float |
| `a | b` | Bitwise or operator. | 4 | int, uint |
| `a ^ b` | Bitwise xor operator. | 4 | int, uint |
| `a == b` | Equality operator. | 3 | int, uint, float, string |
| `a != b` | Inequality operator. | 3 | int, uint, float, string |
| `a > b` | Greater than operator. | 3 | int, uint, float, string |
| `a < b` | Less than operator. | 3 | int, uint, float, string |
| `a >= b` | Greater than or equal operator. | 3 | int, uint, float, string |
| `a <= b` | Less than or equal operator. | 3 | int, uint, float, string |
| `a && b` | Logical and operator. | 2 | bool |
| `a || b` | Logical or operator. | 1 | bool |
| `a ? b : c` | Conditional operator. | 0 | bool; any |

Constant values can be specified using simple specifications.

| Constant | Type |
|---|---|
| 1 | const int |
| 0xff | const int |
| 0777 | const int |
| 0.0 | float |
| 'a' | uint |
| "test" | string |
| true | bool |
| false | bool |

### Ternary

Go omits the C ternary operator, here called the conditional operator. expr includes this operator because it is otherwise impossible to implement branching logic in expr.

It is implemented the same way as C ternary is:

```
   len("test") > 4 ? 1 : 0
```

This expression would return `0`, because `len("test")` is not greater than `4`.

### Errors

If an expr expression does something illegal, it may fail compilation (either at code generation time or compile time) or it may fail at runtime. Such failures can not be handled within the expr language and the error will be propagated up to user code.
