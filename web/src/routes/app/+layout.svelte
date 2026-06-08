<script lang="ts">
  import { goto } from "$app/navigation";
  import { serversClient, userClient, type Server } from "$lib/api";
  import favicon from "$lib/assets/favicon.svg";
  import Navbar from "$lib/components/Navbar.svelte";
  import ServerCard from "$lib/components/ServerCard.svelte";
  import { auth } from "$lib/stores/auth";
  import { onMount, setContext } from "svelte";
  import { page } from "$app/state";
  import CreateServerModal from "$lib/components/CreateServerModal.svelte";
  import { wsClient } from "$lib/websocket";

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

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

{#if loading}
  <div class="flex items-center justify-center min-h-screen bg-base-200">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
{:else}
  <div class="drawer lg:drawer-open bg-secondary/50 h-dvh">
    <input id="my-drawer-4" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content max-h-dvh grid grid-rows-[40px_1fr]">
      <Navbar />
      <div class="rounded-2xl overflow-hidden h-full">
        <main class="bg-base-200 grid grid-cols-[72px_1fr] h-full">
          <div class="flex flex-col items-center gap-2 px-2 py-3">
            {#each servers as server}
              <ServerCard {server} active={server.id === activeServerId} />
            {/each}

            <div class="divider divider-neutral mx-2 my-1"></div>

            <button
              onclick={() => (showCreateModal = true)}
              class="tooltip tooltip-right"
              data-tip="Create server"
              aria-label="Create server"
            >
              <div
                class="w-12 h-12 rounded-2xl flex items-center justify-center text-lg transition-all duration-200 border-2 border-dashed border-base-content/30 text-base-content/30 hover:border-primary hover:text-primary hover:rounded-xl hover:shadow-md"
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
          {@render children()}
        </main>
      </div>
    </div>
  </div>
{/if}

<CreateServerModal
  show={showCreateModal}
  onClose={() => (showCreateModal = false)}
  onCreated={onServerCreated}
/>
