<script lang="ts">
  import Layout from '../../components/Layout.svelte'
  import { useForm } from '@inertiajs/svelte'
  import FormError from '../../components/FormError.svelte'

  // Get initial data from server (props) using destructuring
  let { errors: serverErrors = {}, old = {} } = $props() as {
    errors?: Record<string, string[]>
    old?: Record<string, any>
  }

  // Initialize form with old data if available, otherwise with empty values
  const form = useForm({
    email: (old.email as string) || '',
    password: (old.password as string) || ''
  })

  function handleSubmit(e: Event) {
    e.preventDefault()
    // Submit the form with options
    $form.post('/login', {
      preserveState: true,
      onSuccess: () => {
        // Login successful, form will automatically redirect based on server response
        $form.email = ''
        $form.password = ''
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

    <form onsubmit={handleSubmit}>
      <div class="mb-4">
        <label for="email" class="block text-gray-700 mb-2">Email</label>
        <input
          id="email"
          type="email"
          bind:value={$form.email}
          class="w-full px-3 py-2 border {serverErrors.hasOwnProperty(
            'email'
          ) && serverErrors.email.length > 0
            ? 'border-red-500'
            : 'border-gray-300'} rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
        {#if serverErrors.hasOwnProperty('email') && serverErrors.email.length > 0}
          <FormError errors={serverErrors.email} />
        {/if}
      </div>

      <div class="mb-6">
        <label for="password" class="block text-gray-700 mb-2">Password</label>
        <input
          id="password"
          type="password"
          bind:value={$form.password}
          class="w-full px-3 py-2 border {serverErrors.hasOwnProperty(
            'password'
          ) && serverErrors.password.length > 0
            ? 'border-red-500'
            : 'border-gray-300'} rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
        {#if serverErrors.hasOwnProperty('password') && serverErrors.password.length > 0}
          <FormError errors={serverErrors.password} />
        {/if}
      </div>

      <button
        type="submit"
        class="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50"
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
