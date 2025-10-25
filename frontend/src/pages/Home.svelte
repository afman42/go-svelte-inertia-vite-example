<script lang="ts">
  import Layout from '../components/Layout.svelte'
  import { page } from '@inertiajs/svelte'
  let { user }: any = $page.props
</script>

<svelte:head>
  <title>Welcome</title>
</svelte:head>

<Layout>
  <div class="py-6">
    <!-- <Head title="Welcome" /> -->
    <h1>Welcome</h1>
    <p>Hello, welcome to your first Inertia app!</p>
    {#if user}
      Email: {user.email} | Name: {user.name}
    {/if}
  </div>
</Layout>
