# Golang Befunge93 - With Unicode support!

(Almost) standard compliant Befunge-93 embeddable runtime.
A CLI frontend is also available.
Read more about Befunge [here](https://github.com/catseye/Befunge-93).

Features:

* No dependencies but Go stdlib.
* Unicode support.
* Improved error handling.
* Deterministic randomness.
* Not turing complete. Eat that, [Starlark](https://github.com/bazelbuild/starlark).

Using the default options, this implementation is almost completely identical to the reference implementation.

Get started:

```bash
go install jo-m.ch/go/gobef93/cmd/gobef93@latest
gobef93 -help
gobef93 examples/hello_world.bf
gobef93 -allow_unicode examples/hello_wörld.bf
```

## Embedding

Check [main.go](cmd/gobef93/main.go) for example usage.

## TODOs and ideas

- [ ] Allow to step
- [ ] Time travel mode
- [ ] Better traceability and debugging tools
- [ ] Implement remaining options (see TODOs in `pkg/bef93/prog.go`)
