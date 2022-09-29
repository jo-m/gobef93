# Golang Befunge93 - now with Unicode support!

(Almost) standard compliant Befunge-93 embeddable runtime.
A CLI frontend is also available.
Read more about Befunge [here](https://github.com/catseye/Befunge-93).

Features:

* No dependencies!
* Unicode support.
* Improved error handling.
* Deterministic randomness.
* Not turing complete. Eat that, [Starlark](https://github.com/bazelbuild/starlark).

Using the default options, this implementation is almost completely identical to the reference implementation.
We lack some option flags but add some new ones, and IO is handled slightly differently.

Get started:

```bash
go install github.com/jo-m/gobef93/cmd/gobef93@latest
gobef93 -help
gobef93 examples/hello_world.bf
gobef93 -allow_unicode examples/hello_wörld.bf
```

## Embedding

Check [main.go](cmd/main.go) for example usage.

## TODOs for later

- [ ] allow to step
- [ ] time travel mode
- [ ] better traceability and debugging tools
