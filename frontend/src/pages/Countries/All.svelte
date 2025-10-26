<script lang="ts">
  import { useForm } from '@inertiajs/svelte'
  import Layout from '../../components/Layout.svelte'
  import ErrorBoundary from '../../components/ErrorBoundary.svelte'
  import FormError from '../../components/FormError.svelte'

  let { countries, errors: serverErrors = {} } = $props() as {
    countries: Country[]
    errors?: Record<string, string[]>
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
            class="block w-full rounded-sm border {serverErrors.hasOwnProperty(
              'name'
            ) && serverErrors.name.length > 0
              ? 'border-red-500'
              : 'border-gray-300'} p-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
            placeholder="Name"
            bind:value={$form.name}
          />
          {#if serverErrors.hasOwnProperty('name') && serverErrors.name.length > 0}
            <FormError errors={serverErrors.name} />
          {/if}
        </div>
        <div class="mb-3">
          <input
            type="text"
            class="block w-full rounded-sm border {serverErrors.hasOwnProperty(
              'code'
            ) && serverErrors.code.length > 0
              ? 'border-red-500'
              : 'border-gray-300'} p-2 focus:outline-none focus:ring-2 focus:ring-purple-500"
            placeholder="Country Code"
            maxlength="2"
            bind:value={$form.code}
          />
          {#if serverErrors.hasOwnProperty('code') && serverErrors.code.length > 0}
            <FormError errors={serverErrors.code} />
          {/if}
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
