export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number; // Default: 5000ms
  createdAt: Date;
}

// Create a global toasts array using runes
let globalToasts = $state<Toast[]>([]);

export const toastStore = {
  // Get the current toasts array
  get toasts() {
    return globalToasts;
  },
  
  // Add a new toast notification
  add: (toast: Omit<Toast, 'id' | 'createdAt'>) => {
    const id = Math.random().toString(36).substring(2, 9);
    const newToast: Toast = {
      id,
      ...toast,
      createdAt: new Date(),
    };
    
    globalToasts = [...globalToasts, newToast];
    
    // Auto-remove toast after duration or default 5 seconds
    const duration = toast.duration ?? 5000;
    if (duration > 0) {
      setTimeout(() => {
        globalToasts = globalToasts.filter(t => t.id !== id);
      }, duration);
    }
    
    return id;
  },
  
  // Remove a specific toast by ID
  remove: (id: string) => {
    globalToasts = globalToasts.filter(t => t.id !== id);
  },
  
  // Remove all toasts
  clear: () => {
    globalToasts = [];
  },
  
  // Show a success toast
  success: (message: string, duration?: number) => {
    return toastStore.add({
      message,
      type: 'success',
      duration
    });
  },
  
  // Show an error toast
  error: (message: string, duration?: number) => {
    return toastStore.add({
      message,
      type: 'error',
      duration
    });
  },
  
  // Show a warning toast
  warning: (message: string, duration?: number) => {
    return toastStore.add({
      message,
      type: 'warning',
      duration
    });
  },
  
  // Show an info toast
  info: (message: string, duration?: number) => {
    return toastStore.add({
      message,
      type: 'info',
      duration
    });
  }
};