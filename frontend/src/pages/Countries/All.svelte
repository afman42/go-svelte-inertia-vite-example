<script lang="ts">
  import { useForm } from '@inertiajs/svelte'
  import Layout from '../../components/Layout.svelte'
  let { countries } = $props() as { countries: Country[] }
  const form = useForm<NewCountryForm>({
    name: '',
    code: ''
  })

  function submit(e: Event) {
    e.preventDefault()
    if ($form.name.length <= 0) {
      $form.setError('name', 'Name Still Empty')
      return
    }
    if ($form.code.length <= 0) {
      $form.setError('code', 'Code Still Empty')
      return
    }
    $form.post('/countries', {
      onSuccess: () => {
        // Form submission was successful
        $form.name = ''
        $form.code = ''
      },
      onError: (errors) => {
        // Handle errors, maybe display them to the user
        console.error('Form submission errors:', errors)
      }
    })
  }
</script>

<Layout>
  <div class="p-4 my-2 shadow bg-white border">
    <h4 class="font-medium mb-2">Add Country</h4>
    <form onsubmit={submit}>
      <input
        type="text"
        class="block mb-3 rounded-sm border border-gray-300 p-2"
        placeholder="Name"
        bind:value={$form.name}
      />
      {#if $form.errors.name}
        <div class="text-red-500 text-sm mb-3">{$form.errors.name}</div>
      {/if}
      <input
        type="text"
        class="block mb-3 rounded-sm border border-gray-300 p-2"
        placeholder="Country Code"
        maxlength="2"
        bind:value={$form.code}
      />
      {#if $form.errors.code}
        <div class="text-red-500 text-sm mb-3">{$form.errors.code}</div>
      {/if}
      <button
        class="block bg-purple-600 text-white text-sm font-medium rounded-sm px-4 py-2"
        type="submit"
        disabled={$form.processing}
      >
        {$form.processing ? 'Submitting...' : 'Submit'}
      </button>
    </form>
  </div>
  <div class="py-6 text-lg">
    <h4 class="font-medium">All Countries</h4>
    <ul>
      {#each countries as c (c.flag)}
        <li>
          {c.flag}
          {c.name}
        </li>
      {/each}
    </ul>
  </div>
</Layout>
