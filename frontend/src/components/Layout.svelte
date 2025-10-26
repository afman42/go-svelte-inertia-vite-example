<script lang="ts">
  import { Link, page, router } from '@inertiajs/svelte'
  import ConfirmModal from './ConfirmModal.svelte'
  import ToastContainer from './ToastContainer.svelte'
  let { children } = $props()
  const url: string = $page.url
  let { user }: any = $page.props //still search types

  // State for logout confirmation modal
  let showLogoutModal = $state(false)
  let mobileMenuOpen = $state(false)

  function handleLogout() {
    // Perform logout request
    // Must POST Method
    router.get('/logout', {
      preserveScroll: true
    })
  }

  function toggleMobileMenu() {
    mobileMenuOpen = !mobileMenuOpen
  }
</script>

<header class="border-b-4 border-b-purple-600">
  <div
    class="flex items-center py-4 container max-w-[1280px] mx-auto justify-between"
  >
    <div class="flex items-center">
      <div class="font-medium text-lg">Go + Inertia</div>
    </div>

    <!-- Desktop Navigation -->
    <nav class="hidden md:flex items-center space-x-8">
      <div class="flex space-x-6">
        <Link
          href="/"
          class={url == '/'
            ? 'font-semibold text-purple-600'
            : 'text-gray-600 hover:text-purple-600'}
        >
          Home
        </Link>
        <Link
          href="/random"
          class={url == '/random'
            ? 'font-semibold text-purple-600'
            : 'text-gray-600 hover:text-purple-600'}
        >
          Random Countries
        </Link>
        <Link
          href="/all"
          class={url == '/all'
            ? 'font-semibold text-purple-600'
            : 'text-gray-600 hover:text-purple-600'}
        >
          All Countries
        </Link>
      </div>

      <div>
        {#if user}
          <div class="flex items-center space-x-6">
            <span class="text-gray-600">Welcome, {user.name}</span>
            <Link href="/profile" class="text-purple-600 hover:underline">
              Profile
            </Link>
            <button
              onclick={() => (showLogoutModal = true)}
              class="text-red-600 hover:underline"
            >
              Logout
            </button>
          </div>
        {:else}
          <div class="flex space-x-4">
            <Link href="/login" class="text-purple-600 hover:underline">
              Login
            </Link>
            <Link href="/register" class="text-purple-600 hover:underline">
              Register
            </Link>
          </div>
        {/if}
      </div>
    </nav>

    <!-- Mobile menu button -->
    <div class="md:hidden flex items-center">
      <button
        onclick={toggleMobileMenu}
        class="text-gray-600 hover:text-purple-600 focus:outline-none"
        aria-label="Toggle menu"
      >
        {#if !mobileMenuOpen}
          <svg
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
        {:else}
          <svg
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        {/if}
      </button>
    </div>
  </div>

  <!-- Mobile Navigation -->
  {#if mobileMenuOpen}
    <div class="md:hidden container max-w-[1280px] mx-auto pb-4">
      <div class="flex flex-col space-y-3 px-4 pt-2 pb-3">
        <Link
          href="/"
          class="block py-2 px-3 rounded-md {url == '/'
            ? 'bg-purple-100 text-purple-600'
            : 'text-gray-600'}"
          onclick={() => (mobileMenuOpen = false)}
        >
          Home
        </Link>
        <Link
          href="/random"
          class="block py-2 px-3 rounded-md {url == '/random'
            ? 'bg-purple-100 text-purple-600'
            : 'text-gray-600'}"
          onclick={() => (mobileMenuOpen = false)}
        >
          Random Countries
        </Link>
        <Link
          href="/all"
          class="block py-2 px-3 rounded-md {url == '/all'
            ? 'bg-purple-100 text-purple-600'
            : 'text-gray-600'}"
          onclick={() => (mobileMenuOpen = false)}
        >
          All Countries
        </Link>

        {#if user}
          <div class="pt-4 border-t border-gray-200">
            <p class="text-gray-600 mb-2">Welcome, {user.name}</p>
            <Link
              href="/profile"
              class="block py-2 px-3 text-purple-600 hover:bg-purple-50 rounded-md"
              onclick={() => (mobileMenuOpen = false)}
            >
              Profile
            </Link>
            <button
              onclick={() => {
                showLogoutModal = true
                mobileMenuOpen = false
              }}
              class="block w-full text-left py-2 px-3 text-red-600 hover:bg-red-50 rounded-md"
            >
              Logout
            </button>
          </div>
        {:else}
          <div class="pt-4 border-t border-gray-200 space-y-2">
            <Link
              href="/login"
              class="block py-2 px-3 text-purple-600 hover:bg-purple-50 rounded-md"
              onclick={() => (mobileMenuOpen = false)}
            >
              Login
            </Link>
            <Link
              href="/register"
              class="block py-2 px-3 text-purple-600 hover:bg-purple-50 rounded-md"
              onclick={() => (mobileMenuOpen = false)}
            >
              Register
            </Link>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</header>
<main class="container max-w-[1280px] mx-auto py-6">
  {@render children()}
</main>

<ToastContainer />
<ConfirmModal
  show={showLogoutModal}
  message="Are you sure you want to log out?"
  onConfirm={handleLogout}
  onClose={() => (showLogoutModal = false)}
/>
