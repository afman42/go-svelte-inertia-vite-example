# Toast Notifications

A toast notification system built with Svelte 5 runes for displaying temporary messages to users.

## Features

- Four toast types: success, error, warning, info
- Configurable duration (or persistent toasts)
- Auto-dismiss with optional manual dismiss
- Smooth animations and transitions
- Accessible (ARIA attributes, keyboard support)
- Clean, modern styling with Tailwind CSS
- Built with Svelte 5 runes (no external store dependency)

## Components

### ToastContainer
- Container that displays all active toasts
- Should be added to your app layout to appear globally

### ToastNotification
- Individual toast component
- Handles display, dismissal, and animations

## Store

### toastStore
- Centralized store using Svelte 5 runes for managing toast notifications
- Provides methods to add, remove, and clear toasts
- Supports four toast types: success, error, warning, info

## Usage

### 1. Add ToastContainer to your app layout

Add the ToastContainer to your main layout component (e.g., Layout.svelte):

```svelte
<script lang="ts">
  import ToastContainer from './components/ToastContainer.svelte';
  // ... other imports
</script>

<!-- Your layout content -->
<ToastContainer />
```

### 2. Using the toast store in components

```svelte
<script lang="ts">
  import { toastStore } from './stores/toast';
  
  const handleAction = () => {
    // Show a success toast
    toastStore.success('Operation completed successfully!');
    
    // Show an error toast with custom duration (8 seconds)
    toastStore.error('Something went wrong', 8000);
    
    // Show a warning toast
    toastStore.warning('Please check your input');
    
    // Show an info toast
    toastStore.info('New update available');
    
    // Show a persistent toast (won't auto-dismiss)
    const toastId = toastStore.add({
      message: 'This toast requires manual dismissal',
      type: 'info',
      duration: 0  // 0 means persistent
    });
    
    // Manually remove a specific toast
    // toastStore.remove(toastId);
    
    // Clear all toasts
    // toastStore.clear();
  };
</script>

<button on:click={handleAction}>Show Toast</button>
```

### 3. Programmatic usage

You can also use toasts from anywhere in your application:

```ts
import { toastStore } from './stores/toast';

// Direct usage in JS/TS files
toastStore.success('User saved successfully!');

// Or with custom duration
toastStore.error('Failed to save user data', 10000);
```

## Toast Types

| Type | Description | Color |
|------|-------------|-------|
| `success` | Positive feedback | Green |
| `error` | Error messages | Red |
| `warning` | Warning messages | Yellow |
| `info` | Informational messages | Blue |

## API

### toastStore methods

- `toasts`: Getter to access the current toasts array (reactive with runes)
- `add(toast)`: Add a new toast with custom options
- `success(message, duration?)`: Show a success toast
- `error(message, duration?)`: Show an error toast
- `warning(message, duration?)`: Show a warning toast
- `info(message, duration?)`: Show an info toast
- `remove(id)`: Remove a specific toast by ID
- `clear()`: Remove all toasts

### Toast Interface

```ts
interface Toast {
  id: string;           // Auto-generated
  message: string;      // Message to display
  type: 'success' | 'error' | 'warning' | 'info';
  duration?: number;    // Auto-dismiss duration in ms (default: 5000, 0 = persistent)
  createdAt: Date;      // Auto-generated timestamp
}
```

## Implementation Details

The toast system uses Svelte 5 runes for state management:
- `$state` for the reactive toasts array
- No external store dependency
- Automatic reactivity through runes
- Efficient updates without manual subscriptions