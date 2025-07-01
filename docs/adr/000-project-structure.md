# ADR-000: Project Structure Reorganization

**Date:** 2025-07-01
**Status:** Accepted
**Author:** Nick Kirby

## Context

The ts-copy project currently uses a monolithic structure with all logic contained in a single `main.go` file. This presents several challenges:

- **Testability**: All logic is embedded in the main function, making it difficult to unit test individual components
- **Maintainability**: A single 125-line file contains CLI parsing, file discovery, concurrent processing, and external command execution
- **Go Conventions**: The current structure doesn't follow standard Go project layout conventions
- **Separation of Concerns**: Different responsibilities (CLI, business logic, I/O operations) are tightly coupled

As the project grows and we implement ADR-001's testing strategy, we need a more structured approach that enables proper unit testing and follows Go best practices.

## Decision

We will restructure the project to follow standard Go project layout conventions with clear separation of concerns:

### Proposed Structure

```
ts-copy/
├── cmd/
│   └── tscp/
│       └── main.go             # CLI entrypoint and coordination
├── internal/
│   ├── config/
│   │   └── args.go             # Argument parsing and validation
│   ├── discovery/
│   │   └── files.go            # File discovery and filtering logic
│   ├── transfer/
│   │   └── tailscale.go        # Tailscale copy operations
│   └── worker/
│       └── pool.go             # Concurrent worker pool management
├── testdata/                   # Test fixtures for unit tests
├── go.mod
├── go.sum
└── ...
```

### Package Responsibilities

- **cmd/tscp**: Minimal main function that coordinates between packages
- **internal/config**: CLI argument parsing, validation, and configuration structure
- **internal/discovery**: File system traversal and extension filtering
- **internal/transfer**: Tailscale command execution and error handling
- **internal/worker**: Concurrent processing logic and worker pool management

### Migration Approach

1. Create new package structure
2. Extract functions from main.go into appropriate packages
3. Update main.go to use new packages
4. Ensure existing functionality remains unchanged
5. Update build and release configurations if needed

## Status

Proposed

## Consequences

### Positive

- **Testability**: Each package can be unit tested in isolation
- **Go Conventions**: Follows standard Go project layout (cmd/, internal/)
- **Separation of Concerns**: Clear boundaries between CLI, business logic, and I/O
- **Maintainability**: Easier to locate and modify specific functionality
- **ADR-001 Enablement**: Provides the foundation for comprehensive unit testing
- **Future Growth**: Better structure for adding new features and packages

### Negative

- **Initial Complexity**: More files and directories to manage
- **Import Dependencies**: Need to manage internal package imports
- **Refactoring Effort**: Significant restructuring of existing code
- **Learning Curve**: Team members need to understand new package structure

### Implementation Considerations

- **Backward Compatibility**: Ensure CLI interface remains unchanged
- **Build Process**: Verify GoReleaser configuration works with new structure
- **Import Paths**: Use internal/ to prevent external package usage
- **Error Handling**: Maintain consistent error handling across packages
- **Documentation**: Update CLAUDE.md and README.md to reflect new structure

### Alternatives Considered

- **Keep Current Structure**: Rejected due to testability and maintainability concerns
- **Pkg/ Instead of Internal/**: Rejected because packages are not intended for external use
- **Flat Package Structure**: Rejected due to lack of clear separation of concerns
- **Domain-Driven Structure**: Deferred as current scope doesn't warrant domain complexity