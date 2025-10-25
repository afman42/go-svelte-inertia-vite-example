<script lang="ts">
import { Link, page } from "@inertiajs/svelte"
let props = $props();
const url: string = $page.url;
</script>
<header class="border-b-4 border-b-purple-600">
  <div class="flex items-center py-4 container  max-w-[1280px] mx-auto space-x-8">
    <div class="font-medium">Go + Inertia</div>
    <nav class="flex space-x-4">
      <Link href="/" className={url == "/" ? "font-semibold" : ""}>
        Home
      </Link>
      <Link
        href="/random"
        class={url == "/random" ? "font-semibold" : ""}
      >
        Random Countries
      </Link>
      <Link href="/all" class={url == "/all" ? "font-semibold" : ""}>
        All Countries
      </Link>
    </nav>
  </div>
</header>
<main class="container max-w-[1280px] mx-auto">{@render props.children()}</main>
