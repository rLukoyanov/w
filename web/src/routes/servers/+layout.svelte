<script lang="ts">
  import { goto } from "$app/navigation";
  import { serversClient, userClient, type Server } from "$lib/api";
  import favicon from "$lib/assets/favicon.svg";
  import Navbar from "$lib/components/Navbar.svelte";
  import ServerCard from "$lib/components/ServerCard.svelte";
  import { auth } from "$lib/stores/auth";
  import { onMount } from "svelte";

  let { children } = $props();

  let servers = $state<Server[]>([]);
  let loading = $state(true);

  onMount(async () => {
    // Check if user is already logged in
    try {
      const user = await userClient.getMe();
      auth.setUser(user);

      // Load servers
      try {
        servers = await serversClient.getAll();
      } catch (err) {
        console.error("Failed to load servers:", err);
      }

      loading = false;
    } catch (e) {
      // Not logged in, redirect to login page
      goto("/auth/login");
    }
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<main>
  <div class="drawer lg:drawer-open bg-secondary/50">
    <input id="my-drawer-4" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">
      <!-- Navbar -->
      <Navbar />
      <!-- Page content here -->
      <div class="rounded-2xl overflow-hidden">
        <main class="min-h-screen bg-base-200 grid grid-cols-[72px_1fr]">
          <div class="flex flex-col gap-4 px-2 py-4">
            {#each servers as server}
              <ServerCard {server} />
            {/each}
          </div>

          {#if loading}
            <div class="flex justify-center py-12 h-full">
              <span class="loading loading-spinner loading-lg"></span>
            </div>
          {/if}
          {@render children()}
        </main>
      </div>
    </div>
  </div>
</main>
