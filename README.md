# `go-typemeta`

This module is intended for three things:

- Easily reading information about types.
- Adding special information to types to be used later when reading information about types. For instance, type names, enums for primitive types, struct field descriptions and default values, and so on.
- Converting values between types.

It is most useful for working with schemas such as GraphQL or Postman, but also for reducing the amount of reflect calls by storing type information in-memory.

> :warning: **This package is a work in progress, and you may encounter serious bugs.**

## Developers

- Ludvig Aldén [@ludvigalden](https://github.com/ludvigalden)
