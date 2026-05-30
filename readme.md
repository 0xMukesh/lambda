# λ

a very minimal untyped lambda calculus interpreter

## example

the following is an example which performs `1 + 2` in lambda calculus using [church encoded values](https://en.wikipedia.org/wiki/Church_encoding).

```
λ: (\m.\n.\f.\x.m f (n f x)) (\f.\x.f x) (\f.\x.f (f x))
\f.\x.(f (f (f x)))
λ:
```
