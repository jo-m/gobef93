# Befunge93 interpreter in Go

This package implements the Befunge 93 standard interpreter
(read more about Befunge [here](https://github.com/catseye/Befunge-93)).

No dependencies!

Can be used as library to embed Befunge into your own application, and also offers a CLI frontend.
Using the default options, this implementation is almost completely identical to the reference implementation.
The CLI frontend differs slightly in its behavior from the reference:
It lacks some option flags, adds some new ones, and IO is handled slightly differently.

For usage of the library, check [main.go](cmd/main.go).

CLI usage: TODO

## TODOs for later

- [ ] allow to step
- [ ] time travel mode
- [ ] better traceability and debugging tools
- [ ] handle TODOs littered in the code
- [ ] installation & CLI usage
