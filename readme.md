# λ

a very minimal untyped lambda calculus interpreter

## example

the following is an example which performs `1 + 2` in lambda calculus using [church encoded values](https://en.wikipedia.org/wiki/Church_encoding).

```bash
λ: add = \m.\n.\f.\x.m f (n f x)
λ: one = \f.\x.f x
λ: two = \f.\x.f (f x)
λ: add one two
\f.\x.(f (f (f x))) # three
λ:
```
