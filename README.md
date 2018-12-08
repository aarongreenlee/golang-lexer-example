# golang-lexer-example

This repository contains two solutions to reverse words in a string while
preserving whitespace between words:

**Lexing and Parsing**

```
result, err := Parse("a string      value")
// returns: "a gnirts      eulav"
```

**Simple Function**

```
result := SimpleParser("a string      value")
// returns: "a gnirts      eulav"
```

## Why Lexing and Parsing?

The lexing and parsing solution is superior to the simple function and will
demonstrate it's value as requirements continue to change. For example, if we
needed to modify the program to change punctuation, or omit consonants the
state functions `stateFn` approach will proove easy to reason about and 
maintain over time.
 
The technique was borrowed from one that can be found in the Go standard library.
Specifically, the `text/template` package. Rob Pike explains the advantages of
this design in a talk from 2011 which was recorded and is [available on YouTube](https://www.youtube.com/watch?v=HxaD_trXwRE).

Aaron Greenlee has used this approach in his work to build a Token Query Language (TQL)
at PowerChord, Inc. which helps drive the rendering engine of the SaaS platform.


