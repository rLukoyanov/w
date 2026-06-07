<script lang="ts">
  import { goto } from "$app/navigation";

  import { onMount } from "svelte";

  onMount(() => {
    goto("/servers");
  });
</script>

<main class="min-h-screen bg-base text-text"></main>
