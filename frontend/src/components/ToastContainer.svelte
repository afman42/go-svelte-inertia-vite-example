<script lang="ts">
  import { toastStore } from '../stores/toast.svelte'
  import ToastNotification from './ToastNotification.svelte'

  // Function to remove a toast
  const removeToast = (id: string) => {
    toastStore.remove(id)
  }
</script>

{#if toastStore.toasts.length > 0}
  <div
    class="toast-container fixed top-4 right-4 z-50 space-y-2 w-full max-w-sm"
  >
    {#each toastStore.toasts as toast (toast.id)}
      <ToastNotification {toast} remove={removeToast} />
    {/each}
  </div>
{/if}
