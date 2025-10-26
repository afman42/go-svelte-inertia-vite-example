<script lang="ts">
  let { errors }: { errors: string[] } = $props()
</script>

{#if errors && errors.length > 0}
  <div class="mt-1">
    {#each errors as error}
      <p class="text-sm text-red-600">{error}</p>
    {/each}
  </div>
{/if}
