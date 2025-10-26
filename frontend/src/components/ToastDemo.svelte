<script lang="ts">
  import { toastStore } from '../stores/toast';
  
  const showSuccessToast = () => {
    toastStore.success('Success! Your action was completed successfully.');
  };
  
  const showErrorToast = () => {
    toastStore.error('Error! Something went wrong. Please try again.');
  };
  
  const showWarningToast = () => {
    toastStore.warning('Warning! Please check your input before proceeding.');
  };
  
  const showInfoToast = () => {
    toastStore.info('Info! A new update is available.');
  };
  
  const showCustomToast = () => {
    toastStore.add({
      message: 'This is a custom toast with 8 second duration',
      type: 'info',
      duration: 8000
    });
  };
  
  const clearAllToasts = () => {
    toastStore.clear();
  };
</script>

<div class="p-6 max-w-md mx-auto bg-white rounded-lg shadow-md">
  <h2 class="text-xl font-semibold mb-4">Toast Notification Demo</h2>
  
  <div class="space-y-3">
    <button 
      class="w-full px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
      on:click={showSuccessToast}
    >
      Show Success Toast
    </button>
    
    <button 
      class="w-full px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
      on:click={showErrorToast}
    >
      Show Error Toast
    </button>
    
    <button 
      class="w-full px-4 py-2 bg-yellow-500 text-gray-800 rounded hover:bg-yellow-600"
      on:click={showWarningToast}
    >
      Show Warning Toast
    </button>
    
    <button 
      class="w-full px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      on:click={showInfoToast}
    >
      Show Info Toast
    </button>
    
    <button 
      class="w-full px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
      on:click={showCustomToast}
    >
      Show Custom Toast (8s)
    </button>
    
    <button 
      class="w-full px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
      on:click={clearAllToasts}
    >
      Clear All Toasts
    </button>
  </div>
</div>