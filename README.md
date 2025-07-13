# Mars Programming Language

Mars is a modern, statically-typed programming language that compiles to Go. It combines the safety and performance of Go with a more expressive syntax and additional features.

## Features

- **Static Typing**: Strong type system with type inference
- **Immutability by Default**: Variables are immutable unless explicitly marked as mutable
- **Memory Safety**: Garbage collection with optional manual memory management
- **Go Interop**: Seamless integration with Go code and ecosystem
- **Modern Syntax**: Clean and expressive syntax with modern language features

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/mars.git
cd mars

# Build the compiler
go build -o zcc ./cmd/zcc

# Add to your PATH
export PATH=$PATH:$(pwd)
```

### Hello, World!

Create a file named `hello.mars`:

```mars
func main() {
    log("Hello, Mars!");
}
```

Run it:

```bash
zcc run hello.mars
```

## Documentation

- [Tutorial](docs/tutorial.md) - Learn the basics of Mars
- [Grammar](docs/grammar.md) - Formal language specification
- [Architecture](docs/architecture.md) - Compiler design and implementation details

## Project Structure

```
mars/
├── lexer/          # Token definitions and lexical analysis
├── parser/         # Syntax analysis and AST construction
├── analyzer/       # Static analysis and type checking
├── transpiler/     # Go code generation
├── runtime/        # Runtime support and GC
├── cmd/
│   ├── zcc/       # Compiler CLI
│   └── repl/      # Interactive REPL
└── docs/          # Documentation
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Building

```bash
# Build the compiler
go build -o zcc ./cmd/zcc

# Run tests
go test ./...

# Run linter
go vet ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Go Programming Language](https://golang.org/) - For inspiration and the target platform
- [Rust](https://www.rust-lang.org/) - For ideas about memory safety and ownership
- [TypeScript](https://www.typescriptlang.org/) - For type system design inspiration

## Community

- [GitHub Issues](https://github.com/yourusername/mars/issues)
- [Discord Server](https://discord.gg/mars-lang)
- [Blog](https://mars-lang.org/blog)

## Roadmap

- [ ] Package system
- [ ] Standard library
- [ ] Concurrency support
- [ ] Build system integration
- [ ] IDE support
- [ ] Performance optimizations

## Status

Mars is currently in early development. The core language features are implemented, but the ecosystem is still growing. We welcome contributions and feedback! 