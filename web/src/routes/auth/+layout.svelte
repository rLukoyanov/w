<script lang="ts">
  import { onMount } from "svelte";
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";

  let { children } = $props();

  onMount(async () => {
    // Check if user is already logged in
    try {
      const user = await userClient.getMe();
      auth.setUser(user);
      goto("/");
    } catch (e) {
      // Not logged in, that's fine
    }
  });
</script>

<main class="grid grid-cols-2 min-h-screen">
  <!-- Left side - Auth Form -->
  <div class="flex items-center justify-center bg-base-100 p-8">
    <div class="w-full max-w-md">
      {@render children()}
    </div>
  </div>

  <!-- Right side - Empty for now -->
  <div class="bg-base-200">
    <!-- Empty space -->
  </div>
</main>
