<script lang="ts">
  import Layout from '../../components/Layout.svelte'
  import { useForm } from '@inertiajs/svelte'

  const form = useForm({
    email: '',
    password: ''
  })

  function handleSubmit(e: Event) {
    e.preventDefault()
    $form.post('/login', {
      onSuccess: () => {
        // Login successful, form will automatically redirect based on server response
        $form.email = ''
        $form.password = ''
      },
      onError: (errors) => {
        // Handle errors if needed (errors will automatically be available in form.errors)
        console.log('Login errors:', errors)
      }
    })
  }
</script>

<svelte:head>
  <title>Login</title>
</svelte:head>

<Layout>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
    <h1 class="text-2xl font-bold mb-6">Login</h1>

    {#if $form.errors.email || $form.errors.password}
      <div class="mb-4 p-3 bg-red-100 text-red-700 rounded">
        {$form.errors.email || $form.errors.password}
      </div>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="mb-4">
        <label for="email" class="block text-gray-700 mb-2">Email</label>
        <input
          id="email"
          type="email"
          bind:value={$form.email}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
          required
        />
      </div>

      <div class="mb-6">
        <label for="password" class="block text-gray-700 mb-2">Password</label>
        <input
          id="password"
          type="password"
          bind:value={$form.password}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
          required
        />
      </div>

      <button
        type="submit"
        class="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
        disabled={$form.processing}
      >
        {$form.processing ? 'Logging in...' : 'Login'}
      </button>
    </form>

    <div class="mt-4 text-center">
      <p class="text-gray-600">
        Don't have an account?
        <a href="/register" class="text-purple-600 hover:underline"
          >Register here</a
        >
      </p>
    </div>
  </div>
</Layout>
