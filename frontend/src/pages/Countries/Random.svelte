<script lang="ts">
  import { router } from '@inertiajs/svelte'
  import Layout from '../../components/Layout.svelte'
  import ErrorBoundary from '../../components/ErrorBoundary.svelte'
  import LoadingSpinner from '../../components/LoadingSpinner.svelte'

  let { countries } = $props() as { countries: Country[] }
  let isLoading = $state(false)

  function refreshCountries() {
    isLoading = true
    try {
      router.reload()
    } finally {
      isLoading = false
    }
  }
</script>

<Layout>
  <ErrorBoundary>
    <div class="py-6 text-lg">
      {#if isLoading}
        <div class="flex justify-center my-8">
          <LoadingSpinner size="lg" color="purple" centered={true} />
        </div>
      {:else if countries && countries.length > 0}
        <ul>
          {#each countries as c, i (i)}
            <li class="mb-2">
              {c.flag}
              {c.name}
            </li>
          {/each}
        </ul>
        <div class="mt-4">
          <button
            class="block bg-purple-600 text-white text-sm font-medium rounded-sm px-4 py-2 mt-2 hover:bg-purple-700 disabled:opacity-50"
            onclick={refreshCountries}
            disabled={isLoading}
          >
            Refresh
          </button>
        </div>
      {:else}
        <div class="text-center py-8">
          <p>Let's Login first to view countries.</p>
        </div>
      {/if}
    </div>
  </ErrorBoundary>
</Layout>
