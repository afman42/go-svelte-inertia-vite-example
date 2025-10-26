<script lang="ts">
  let { children } = $props();
  let hasError = $state(false);
  let errorInfo = $state<string | null>(null);

  function resetError() {
    hasError = false;
    errorInfo = null;
  }
</script>

{#if hasError}
  <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
    <strong class="font-bold">Error! </strong>
    <span class="block sm:inline">An unexpected error occurred.</span>
    <button 
      onclick={resetError} 
      class="absolute top-0 bottom-0 right-0 px-4 py-3"
    >
      Reset
    </button>
  </div>
{:else}
  {@render children()}
{/if}