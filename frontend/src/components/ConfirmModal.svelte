<script lang="ts">
  let { show = false, message = 'Are you sure?', confirmText = 'Confirm', cancelText = 'Cancel', onConfirm, onClose }: { 
    show?: boolean, 
    message?: string, 
    confirmText?: string, 
    cancelText?: string,
    onConfirm?: () => void,
    onClose?: () => void
  } = $props();
  
  let localShow = $state(show);
  
  $effect(() => {
    localShow = show;
  });
  
  const close = () => {
    localShow = false;
    onClose?.();
  }

  const confirm = () => {
    onConfirm?.();
    close();
  }
</script>

{#if localShow}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg shadow-xl max-w-md w-full">
      <h3 class="text-lg font-medium mb-4">Confirm Action</h3>
      <p class="mb-6">{message}</p>
      <div class="flex justify-end space-x-3">
        <button 
          onclick={close} 
          class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded"
        >
          {cancelText}
        </button>
        <button 
          onclick={confirm} 
          class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        >
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}