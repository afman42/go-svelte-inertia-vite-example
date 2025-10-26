<script lang="ts">
  import Layout from '../../components/Layout.svelte'
  import ErrorBoundary from '../../components/ErrorBoundary.svelte'

  // Extract user data from props
  let { user }: { user: User } = $props()
</script>

<svelte:head>
  <title>User Profile</title>
</svelte:head>

<Layout>
  <ErrorBoundary>
    <div class="max-w-2xl mx-auto py-6">
      <h1 class="text-2xl font-bold mb-6">User Profile</h1>

      <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="mb-4">
          <div class="block text-gray-700 font-medium mb-1">Name</div>
          <p class="text-gray-900">{user.name}</p>
        </div>

        <div class="mb-4">
          <div class="block text-gray-700 font-medium mb-1">Email</div>
          <p class="text-gray-900">{user.email}</p>
        </div>

        <div class="mb-4">
          <div class="block text-gray-700 font-medium mb-1">User ID</div>
          <p class="text-gray-900">{user.id}</p>
        </div>

        <div class="mt-8">
          <p class="text-gray-600 mb-4">
            Use the logout button in the header to log out safely.
          </p>
        </div>
      </div>
    </div>
  </ErrorBoundary>
</Layout>
