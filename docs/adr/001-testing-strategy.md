# ADR-001: Testing Strategy

**Date:** 2025-06-30
**Status:** Accepted
**Author:** Nick Kirby

## Context

The ts-copy project currently has no automated testing infrastructure. The codebase consists of a single `main.go` file with file discovery, concurrent processing, and Tailscale integration logic. As the project grows and adds features, we need a testing strategy to:

- Ensure reliability of file transfer operations
- Prevent regressions during refactoring
- Enable confident feature development
- Support CI/CD pipeline requirements

The current development workflow relies on manual testing using `just dry-run` in the test directory, which is insufficient for comprehensive validation.

## Decision

We will start with **unit tests** as the initial testing strategy, with integration tests to be added in a future phase:

### Unit Tests

- **Target**: Individual functions and components in isolation
- **Framework**: Go's built-in `testing` package with `testify/assert` for assertions
- **Coverage**:
  - File discovery logic (`filepath.Walk` operations)
  - Extension filtering functionality
  - Configuration parsing and validation
  - Error handling paths
- **Mocking**: Mock `os/exec` calls to avoid actual Tailscale operations

### Testing Infrastructure

- **Organization**: `*_test.go` files alongside source code
- **Test Data**: Structured test fixtures in `testdata/` directory
- **CI Integration**: Maintain existing `go test ./... -v` in GitHub Actions
- **Coverage Target**: Aim for 100% code coverage on core logic

### Future Considerations

- **Integration Tests**: Will be evaluated in a future ADR once unit testing is established
- **End-to-end Testing**: May be considered if the project grows in complexity

## Status

Accepted

## Consequences

### Positive

- **Reliability**: Catch bugs before they reach users
- **Refactoring Confidence**: Safe to restructure code with test coverage
- **Documentation**: Tests serve as executable documentation of expected behavior
- **CI/CD Integration**: Tests already integrated in release workflow
- **Maintainability**: Easier to validate changes and onboard contributors

### Negative

- **Development Overhead**: Initial time investment to write tests
- **Maintenance Burden**: Tests must be updated when functionality changes
- **Mocking Complexity**: Tailscale integration requires careful mocking to avoid side effects
- **Limited Real-world Validation**: Mocked tests may miss integration issues with actual Tailscale

### Implementation Considerations

- **Refactoring Required**: Current monolithic `main.go` needs function extraction for testability
- **Dependency Management**: Add `testify` for better assertions and mocking capabilities
- **Test Execution**: Ensure tests can run in CI environment without Tailscale dependencies
- **Documentation**: Update README.md and CLAUDE.md with testing instructions

### Alternatives Considered

- **Integration Tests First**: Rejected due to setup complexity and external dependencies
- **Hybrid Unit/Integration**: Deferred to future ADR once unit testing foundation is established
- **Manual testing only**: Current approach, insufficient for growing codebase
