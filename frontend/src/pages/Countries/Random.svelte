<script lang="ts">
  import { router } from '@inertiajs/svelte'
  import Layout from '../../components/Layout.svelte'
  let { countries } = $props() as { countries: Country[] }
</script>

<Layout>
  <div class="py-6 text-lg">
    <ul>
      {#each countries as c (c.name)}
        <li>
          {c.flag}
          {c.name}
        </li>
      {/each}
    </ul>
    <div>
      {#if countries.length > 0}
        <button
          class="block bg-purple-600 text-white text-sm font-medium rounded-sm px-4 py-2 mt-2"
          onclick={() => router.reload()}
        >
          Refresh
        </button>
      {:else}
        Let's Login first
      {/if}
    </div>
  </div>
</Layout>
