<script lang="ts">
  import Layout from '../../components/Layout.svelte'
  import { useForm } from '@inertiajs/svelte'
  import { toastStore } from '../../stores/toast.svelte'

  // Get initial data from server (props) using destructuring
  let { errors: serverErrors = [], old = {} } = $props() as {
    errors?: { field: string; message: string }[]
    old?: Record<string, any>
  }

  // Initialize form with old data if available, otherwise with empty values
  const form = useForm({
    email: old?.email || '',
    password: old?.password || ''
  })
  function handleSubmit(e: Event) {
    e.preventDefault()
    // Submit the form with options
    $form.post('/login', {
      preserveState: true,
      onSuccess: () => {
        $form.email = ''
        $form.password = ''
      },
      onError: () => {
        serverErrors.forEach((v: { field: string; message: string }) => {
          toastStore.error(v.message, 2000)
        })
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
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
      </div>

      <div class="mb-6">
        <label for="password" class="block text-gray-700 mb-2">Password</label>
        <input
          id="password"
          type="password"
          bind:value={$form.password}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
        />
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
