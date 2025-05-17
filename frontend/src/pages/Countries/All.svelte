<script lang="ts">
import { useForm } from "@inertiajs/svelte";
import Layout from "../../components/Layout.svelte";
let props: EveryFlagCountries = $props();
const form = useForm({
	name: "",
	code: "",
});

function submit(e: any) {
  e.preventDefault();
  $form.post("/countries");
  $form.name = ""
  $form.code = ""
}
</script>

<Layout>
  <div class="p-4 my-2 shadow bg-white border">
    <h4 class="font-medium mb-2">Add Country</h4>
    <form onsubmit={submit}>
      <input
        type="text"
        class="block mb-3 rounded-sm"
        placeholder="Name"
        bind:value={$form.name}
      />
      <input
        type="text"
        class="block mb-3 rounded-sm"
        placeholder="Country Code"
        maxlength="2"
        bind:value={$form.code}
      />
      <button
        class="block bg-purple-600 text-white text-sm font-medium rounded-sm px-4 py-2"
        type="submit"
        disabled={$form.processing}
      >
        Submit
      </button>
    </form>
  </div>
  <div class="py-6 text-lg">
    <h4 class="font-medium">All Countries</h4>
    <ul>
      {#each props.countries as c (c.Name)}
        <li>
          {c.Flag} {c.Name}
        </li>
      {/each}
    </ul>
  </div>
</Layout>
