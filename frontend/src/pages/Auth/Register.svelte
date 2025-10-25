<script lang="ts">
  import Layout from '../../components/Layout.svelte'
  import { useForm, router } from '@inertiajs/svelte'

  const form = useForm<RegisterFormData>({
    name: '',
    email: '',
    password: '',
    password_confirmation: '' // Using snake_case to match backend expectations
  })

  function handleSubmit(e: Event) {
    e.preventDefault()

    // Basic validation for password confirmation - access form data function
    if ($form.password !== $form.password_confirmation) {
      // Since we don't have a field-specific error for confirmation, we'll
      // set a general error or handle via backend
      $form.setError('password_confirmation', 'Passwords do not match')
      return
    }

    if ($form.password.length < 6) {
      $form.setError('password', 'Password must be at least 6 characters')
      return
    }

    $form.post('/register', {
      onSuccess: () => {
        // Registration successful
        $form.email = ''
        $form.password = ''
        $form.name = ''
        $form.password_confirmation = ''
        router.visit('/login')
      },
      onError: (errors) => {
        // Errors are automatically handled by Inertia and available in form.errors
        console.log('Registration errors:', errors)
      }
    })
  }
</script>

<svelte:head>
  <title>Register</title>
</svelte:head>

<Layout>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
    <h1 class="text-2xl font-bold mb-6">Create Account</h1>

    {#if $form.errors.name || $form.errors.email || $form.errors.password}
      <div class="mb-4 p-3 bg-red-100 text-red-700 rounded">
        {$form.errors.name || $form.errors.email || $form.errors.password}
      </div>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="mb-4">
        <label for="name" class="block text-gray-700 mb-2">Full Name</label>
        <input
          id="name"
          type="text"
          bind:value={$form.name}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
          required
        />
      </div>

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

      <div class="mb-4">
        <label for="password" class="block text-gray-700 mb-2">Password</label>
        <input
          id="password"
          type="password"
          bind:value={$form.password}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
          required
        />
      </div>

      <div class="mb-6">
        <label for="confirmPassword" class="block text-gray-700 mb-2"
          >Confirm Password</label
        >
        <input
          id="confirmPassword"
          type="password"
          bind:value={$form.password_confirmation}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
          required
        />
      </div>

      <button
        type="submit"
        class="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
        disabled={$form.processing}
      >
        {$form.processing ? 'Creating Account...' : 'Register'}
      </button>
    </form>

    <div class="mt-4 text-center">
      <p class="text-gray-600">
        Already have an account?
        <a href="/login" class="text-purple-600 hover:underline">Login here</a>
      </p>
    </div>
  </div>
</Layout>
