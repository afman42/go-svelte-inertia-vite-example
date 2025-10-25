<script lang="ts">
  import Layout from '../components/Layout.svelte'
  let props: HomePageProps = $props();
</script>

<svelte:head>
  <title>Welcome</title>
</svelte:head>

<Layout>
  <div class="py-6">
    <!-- <Head title="Welcome" /> -->
    <h1>Welcome</h1>
    <p>Hello, welcome to your first Inertia app!</p>
  </div>
</Layout>
