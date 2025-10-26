<script lang="ts">
  import type { Toast, ToastType } from '../stores/toast.svelte'

  let {
    toast,
    remove
  }: {
    toast: Toast
    remove: (id: string) => void
  } = $props()

  // Define type-specific styling
  const getTypeClasses = (type: ToastType) => {
    const baseClasses =
      'p-4 rounded-md shadow-lg transform transition-all duration-300 ease-in-out'

    switch (type) {
      case 'success':
        return `${baseClasses} bg-green-500 text-white`
      case 'error':
        return `${baseClasses} bg-red-500 text-white`
      case 'warning':
        return `${baseClasses} bg-yellow-500 text-gray-800`
      case 'info':
        return `${baseClasses} bg-blue-500 text-white`
      default:
        return `${baseClasses} bg-gray-700 text-white`
    }
  }

  // Auto-dismiss toast on hover out after a delay
  let dismissTimeout: ReturnType<typeof setTimeout> | null = null

  const startDismissTimer = () => {
    if (dismissTimeout) {
      clearTimeout(dismissTimeout)
    }

    dismissTimeout = setTimeout(() => {
      remove(toast.id)
    }, 3000)
  }

  // If toast has no duration (persistent), start the dismiss timer on mouse leave
  if (toast.duration === 0) {
    startDismissTimer()
  }

  $effect(() => {
    return () => {
      if (dismissTimeout) {
        clearTimeout(dismissTimeout)
      }
    }
  })
</script>

<div
  class="toast-container toast {getTypeClasses(toast.type)}"
  data-type={toast.type}
  role="alert"
  aria-live="assertive"
  onmouseenter={() => {
    if (dismissTimeout) {
      clearTimeout(dismissTimeout)
    }
  }}
  onmouseleave={() => {
    if (toast.duration === 0) {
      startDismissTimer()
    }
  }}
>
  <div class="flex items-start">
    {#if toast.type === 'success'}
      <div class="mr-3">✓</div>
    {:else if toast.type === 'error'}
      <div class="mr-3">✗</div>
    {:else if toast.type === 'warning'}
      <div class="mr-3">⚠</div>
    {:else if toast.type === 'info'}
      <div class="mr-3">ℹ</div>
    {/if}

    <div class="flex-1 min-w-0">
      <p class="text-sm font-medium">{toast.message}</p>
    </div>

    <button
      type="button"
      class="ml-4 text-white hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-white rounded-full"
      aria-label="Dismiss notification"
      onclick={() => remove(toast.id)}
    >
      <svg
        class="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M6 18L18 6M6 6l12 12"
        />
      </svg>
    </button>
  </div>
</div>
