<script lang="ts">
  import Layout from '../components/Layout.svelte'
  import ErrorBoundary from '../components/ErrorBoundary.svelte'
  import { page } from '@inertiajs/svelte'
  let { user }: any = $page.props
</script>

<svelte:head>
  <title>Welcome</title>
</svelte:head>

<Layout>
  <ErrorBoundary>
    <div class="py-6">
      <!-- <Head title="Welcome" /> -->
      <h1 class="text-2xl font-bold mb-4">Welcome</h1>
      <p class="mb-4">Hello, welcome to your first Inertia app!</p>
      {#if user}
        <div class="bg-gray-100 p-4 rounded">
          <p><span class="font-semibold">Email:</span> {user.email}</p>
          <p><span class="font-semibold">Name:</span> {user.name}</p>
        </div>
      {/if}
    </div>
  </ErrorBoundary>
</Layout>
