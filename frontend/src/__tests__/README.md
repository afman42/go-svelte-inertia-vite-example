# Frontend Unit Tests

This project includes unit tests for the frontend components using Vitest.

## Test Structure

Due to compatibility issues between Svelte 5's new rune system and Inertia.js with component testing libraries, we've adopted the following testing approach:

### 1. Logic Tests (Recommended)
Logic tests verify business logic and data structures without rendering components:

- Form validation logic
- Data structure validation
- URL handling logic
- State transformations

### 2. Component Tests (Limited/Pending)
Full component rendering tests are currently not working due to compatibility issues between:
- Svelte 5 runes (`$state`, `$props`, etc.)
- Inertia.js runtime
- Testing library's component mounting mechanism

The component test files exist as placeholders but currently skip rendering tests.

## Running Tests

```bash
# Run all tests
pnpm test

# Run tests in watch mode
pnpm test:ui

# Run tests once
pnpm test:run

# Run specific test file
pnpm test:run src/__tests__/Logic.test.ts
```

## Test Files

- `Logic.test.ts` - General business logic tests (fully functional)
- `AuthLogic.test.ts` - Authentication form logic tests (fully functional)
- `Layout.test.ts` - Component tests (placeholder)
- `Home.test.ts` - Component tests (placeholder)
- `Auth.test.ts` - Component tests (placeholder)
- `Countries.test.ts` - Component tests (placeholder)

## Known Limitations

1. Full component rendering tests currently do not work due to Svelte 5 + Inertia integration challenges
2. DOM interaction tests are not available for components that use Inertia features
3. Tests that require `$page` or other Inertia stores are difficult to mock properly

## Current Status

- ✅ Logic tests: Working
- ❌ Component rendering tests: Not working due to Svelte 5/Inertia compatibility issues
- 🔄 Future: Workarounds and updates are being investigated by the community

## Future Improvements

As Svelte 5, Inertia, and testing library ecosystems mature and resolve compatibility issues, we plan to:
- Add more comprehensive component rendering tests
- Implement better mocking strategies for Inertia features
- Create helper utilities for common testing patterns