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

## Performance Differences

```
// BenchmarkSimpleParser benchmarks the simple parser function which resulted in
// the following results on a 2014 MacBook Pro:
//
// BenchmarkSimpleParser-8   1000000     1119 ns/op     384 B/op     22 allocs/op
// BenchmarkParse-8          200000      6360 ns/op     561 B/op     19 allocs/op
```

## Why Lexing and Parsing? The Simple Implementation was Faster!

First, have sympathy for the humans who will need to build, review, vet, ship,
and maintain the solution. In this case, the performance of the lexing/parsing
implementation is probably sufficient for "web scale" traffic. However, the
real value will be recognized over time as requirements change and the 
solution needs to evolve. For example, if we needed to modify the program to
change punctuation, or omit consonants the state functions `stateFn` approach
will prove easy to reason about and maintain over time.

## Where can I learn more about Lexing and Parsing?

The technique was borrowed from one that can be found in the Go standard library.
Specifically, the `text/template` package. Rob Pike explains the advantages of
this design in a talk from 2011 which was recorded and is [available on YouTube](https://www.youtube.com/watch?v=HxaD_trXwRE).

You might also enjoy an article written by Fatih Arslan in 2015 titled [lexers and scanner packages](https://medium.com/@farslan/a-look-at-go-scanner-packages-11710c2655fc).  

Aaron Greenlee has used this approach in his work to build a Token Query Language (TQL)
at PowerChord, Inc. which helps drive the rendering engine of the SaaS platform.







