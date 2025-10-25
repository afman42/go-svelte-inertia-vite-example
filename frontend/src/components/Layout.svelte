<script lang="ts">
  import { Link, page } from '@inertiajs/svelte'
  let { children } = $props()
  const url: string = $page.url
  let { user }: any = $page.props //still search types
</script>

<header class="border-b-4 border-b-purple-600">
  <div
    class="flex items-center py-4 container max-w-[1280px] mx-auto space-x-8 justify-between"
  >
    <div class="flex items-center space-x-8">
      <div class="font-medium">Go + Inertia</div>
      <nav class="flex space-x-4">
        <Link href="/" class={url == '/' ? 'font-semibold' : ''}>Home</Link>
        <Link href="/random" class={url == '/random' ? 'font-semibold' : ''}>
          Random Countries
        </Link>
        <Link href="/all" class={url == '/all' ? 'font-semibold' : ''}>
          All Countries
        </Link>
      </nav>
    </div>

    <div>
      {#if user}
        <div class="flex items-center space-x-4">
          <span class="text-gray-600">Welcome, {user.name}</span>
          <Link href="/profile" class="text-purple-600 hover:underline">
            Profile
          </Link>
          <Link href="/logout" class="text-red-600 hover:underline">
            Logout
          </Link>
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
  </div>
</header>
<main class="container max-w-[1280px] mx-auto">
  {@render children()}
</main>
