<script lang="ts">
  import { goto } from "$app/navigation";
  import { serversClient, userClient, type Server } from "$lib/api";
  import Navbar from "$lib/components/Navbar.svelte";
  import ServerCard from "$lib/components/ServerCard.svelte";
  import { auth } from "$lib/stores/auth";
  import { onMount, setContext } from "svelte";
  import { page } from "$app/state";
  import CreateServerModal from "$lib/components/CreateServerModal.svelte";
  import { wsClient } from "$lib/websocket";
  import { Plus } from "lucide-svelte";

  let { children } = $props();

  let servers = $state<Server[]>([]);
  let loading = $state(true);

  let activeServerId = $derived(page.params.id as string | undefined);
  let showCreateModal = $state(false);

  async function refreshServers(): Promise<void> {
    try {
      servers = await serversClient.getAll();
    } catch (err) {
      console.error("Failed to load servers:", err);
    }
  }

  setContext("servers-context", {
    get servers() {
      return servers;
    },
    get loading() {
      return loading;
    },
    refreshServers,
  });

  function onServerCreated(server: Server) {
    refreshServers();
    goto(`/app/servers/${server.id}`);
  }

  onMount(async () => {
    try {
      const user = await userClient.getMe();
      auth.setUser(user);

      const token = userClient.getToken();
      if (token) {
        wsClient.connect(token);
      }

      await refreshServers();
      loading = false;
    } catch (e) {
      goto("/auth/login");
    }
  });
</script>

{#if loading}
  <div class="flex items-center justify-center min-h-screen"
    style="background: oklch(0.07 0.004 285);">
    <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
  </div>
{:else}
  <div class="h-dvh flex flex-col" style="background: oklch(0.07 0.004 285);">
    <Navbar />
    <div class="flex-1 flex overflow-hidden">
      <div class="flex flex-col items-center gap-2 px-2 py-3 shrink-0"
        style="background: oklch(0.09 0.005 285); width: 72px; border-right: 1px solid oklch(0.16 0.008 285);">
        {#each servers as server}
          <ServerCard {server} active={server.id === activeServerId} />
        {/each}

        <div class="w-8 h-px my-1" style="background: oklch(0.18 0.008 285);"></div>

        <button
          onclick={() => (showCreateModal = true)}
          class="tooltip tooltip-right"
          data-tip="Create server"
          aria-label="Create server"
        >
          <div
            class="w-12 h-12 rounded-2xl flex items-center justify-center text-lg transition-all duration-200 border-2 border-dashed"
            style="border-color: oklch(0.3 0.01 285); color: oklch(0.35 0.01 285);"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="w-5 h-5"
            >
              <path d="M5 12h14" />
              <path d="M12 5v14" />
            </svg>
          </div>
        </button>
      </div>
      <main class="flex-1 overflow-hidden"
        style="background: oklch(0.105 0.006 285);">
        {@render children()}
      </main>
    </div>
  </div>
{/if}

<CreateServerModal
  show={showCreateModal}
  onClose={() => (showCreateModal = false)}
  onCreated={onServerCreated}
/>
