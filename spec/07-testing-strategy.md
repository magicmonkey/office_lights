# Testing Strategy

## Overview
Define testing approach for the office lights control system.

## Tasks

### 1. Unit Tests for LED Strip Driver
- Test `NewLEDStrip()` constructor
- Test `SetColor()` with valid values
- Test `SetColor()` with invalid values (out of range)
- Test JSON message formatting
- Test state retrieval methods
- Mock MQTT client for testing

### 2. Unit Tests for LED Bar Driver
- Test `NewLEDBar()` constructor
- Test `SetRGBW()` with valid and invalid inputs
- Test `SetWhite()` with valid and invalid inputs
- Test message formatting (verify 77 values)
- Test convenience methods (TurnOffAll, etc.)
- Verify correct value ordering in CSV output
- Mock MQTT client for testing

### 3. Unit Tests for Video Light Driver
- Test `NewVideoLight()` constructor
- Test `SetState()` with valid and invalid inputs
- Test brightness validation (0-100)
- Test message formatting
- Test convenience methods
- Mock MQTT client for testing

### 4. Integration Tests
- Test MQTT message publishing end-to-end
- Test with actual MQTT broker (local test broker)
- Verify messages are published to correct topics
- Verify message payloads are correct

### 5. Create Mock MQTT Client
- Implement mock for testing without actual MQTT broker
- Capture published messages for verification
- Allow testing of error conditions

### 6. Test Utilities
- Create helper functions for common test setups
- Create test data generators
- Create assertion helpers for message validation

### 7. Manual Testing Checklist
- Connect to actual lights (when available)
- Verify LED strip color changes
- Verify LED bar segments light correctly
- Verify video lights turn on/off at correct brightness
- Test edge cases (maximum values, minimum values)

### 8. Documentation Tests
- Ensure code examples in documentation work
- Verify all exported functions have proper comments
- Check that README matches actual implementation

## Test Coverage Goals
- Aim for >80% code coverage on driver packages
- 100% coverage on message formatting logic
- All error paths should be tested

## Continuous Integration
- Set up GitHub Actions or similar CI
- Run tests on every commit
- Run linting (golangci-lint)
- Check code formatting (gofmt)

## Success Criteria
- All unit tests pass
- Integration tests pass with test MQTT broker
- Test coverage meets goals
- CI pipeline runs successfully
