<script lang="ts">
  import { useForm } from '@inertiajs/svelte'
  import Layout from '../../components/Layout.svelte'
  import ErrorBoundary from '../../components/ErrorBoundary.svelte'
  import { toastStore } from '../../stores/toast.svelte'

  let { countries, errors: serverErrors = [] } = $props() as {
    countries: Country[]
    errors?: { field: string; message: string }[]
  }

  // Initialize form with old data if available, otherwise with empty values
  const form = useForm<NewCountryForm>({
    name: '',
    code: ''
  })

  function submit(e: Event) {
    e.preventDefault()
    // Submit the form with options
    $form.post('/countries', {
      preserveState: true,
      onSuccess: () => {
        // Form submission was successful
        $form.reset('name', 'code')
      },
      onError: () => {
        serverErrors.forEach((v: { field: string; message: string }) => {
          toastStore.error(v.message, 2000)
        })
      }
    })
  }
</script>

<Layout>
  <ErrorBoundary>
    <div class="p-4 my-2 shadow bg-white border rounded-lg">
      <h4 class="font-medium mb-4">Add Country</h4>

      <form onsubmit={submit}>
        <div class="mb-3">
          <input
            type="text"
            class="block w-full rounded-sm border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
            placeholder="Name"
            bind:value={$form.name}
          />
        </div>
        <div class="mb-3">
          <input
            type="text"
            class="block w-full rounded-sm border border-gray-300 p-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
            placeholder="Country Code"
            maxlength="2"
            bind:value={$form.code}
          />
        </div>
        <button
          class="block bg-purple-600 text-white text-sm font-medium rounded-sm px-4 py-2 hover:bg-purple-700 disabled:opacity-50"
          type="submit"
          disabled={$form.processing}
        >
          {$form.processing ? 'Submitting...' : 'Submit'}
        </button>
      </form>
    </div>

    <div class="py-6 text-lg">
      <h4 class="font-medium mb-4">All Countries</h4>
      {#if countries && countries.length > 0}
        <ul>
          {#each countries as c, i (i)}
            <li class="mb-2">
              {c.flag}
              {c.name}
            </li>
          {/each}
        </ul>
      {:else}
        <p>No countries available.</p>
      {/if}
    </div>
  </ErrorBoundary>
</Layout>
