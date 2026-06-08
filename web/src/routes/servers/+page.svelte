<script lang="ts">
  import { getContext } from "svelte";
  import { goto } from "$app/navigation";
  import { serversClient, type Server } from "$lib/api";
  import { Plus, Server as ServerIcon } from "lucide-svelte";
  import { notify } from "$lib/stores/notifications";

  let ctx = getContext<{
    servers: Server[];
    loading: boolean;
    refreshServers: () => Promise<void>;
  }>("servers-context");

  let showCreateServer = $state(false);
  let newServerName = $state("");
  let isCreating = $state(false);

  async function handleCreateServer() {
    if (!newServerName.trim()) {
      notify.error("Server name is required");
      return;
    }

    isCreating = true;

    try {
      const server = await serversClient.create(newServerName);
      newServerName = "";
      showCreateServer = false;
      isCreating = false;

      goto(`/servers/${server.id}`);
    } catch (err: any) {
      notify.error(err.message ?? "Failed to create server");
      isCreating = false;
    }
  }

  function resetForm() {
    showCreateServer = false;
    newServerName = "";
  }
</script>

<main class="h-full p-8">
  {#if ctx.loading}
    <div class="flex items-center justify-center h-full">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if ctx.servers.length === 0}
    <!-- Empty state: no servers yet -->
    <div class="flex flex-col items-center justify-center h-full text-center">
      <ServerIcon class="w-24 h-24 text-base-content/20 mb-6" />
      <h1 class="text-3xl font-bold mb-2">Welcome to Weche</h1>
      <p class="text-base-content/60 mb-8 max-w-md">
        You don't have any servers yet. Create your first server to get started!
      </p>

      {#if showCreateServer}
        <div class="w-full max-w-md space-y-4">
          <input
            type="text"
            bind:value={newServerName}
            placeholder="Server name"
            class="input input-bordered w-full"
            disabled={isCreating}
          />
          <div class="flex gap-2">
            <button
              onclick={handleCreateServer}
              class="btn btn-primary flex-1"
              disabled={isCreating || !newServerName.trim()}
            >
              {#if isCreating}
                <span class="loading loading-spinner"></span>
                Creating...
              {:else}
                Create Server
              {/if}
            </button>
            <button onclick={resetForm} class="btn btn-ghost">Cancel</button>
          </div>
        </div>
      {:else}
        <button
          onclick={() => (showCreateServer = true)}
          class="btn btn-primary btn-lg gap-2"
        >
          <Plus class="w-5 h-5"></Plus>
          Create your first server
        </button>
      {/if}
    </div>
  {:else}
    <!-- Servers exist: show selection prompt -->
    <div class="flex flex-col items-center justify-center h-full text-center">
      <ServerIcon class="w-24 h-24 text-base-content/20 mb-6" />
      <h1 class="text-3xl font-bold mb-2">Select a server</h1>
      <p class="text-base-content/60 mb-8">
        Choose a server from the sidebar to start chatting
      </p>

      {#if showCreateServer}
        <div class="w-full max-w-md space-y-4">
          <input
            type="text"
            bind:value={newServerName}
            placeholder="Server name"
            class="input input-bordered w-full"
            disabled={isCreating}
          />
          <div class="flex gap-2">
            <button
              onclick={handleCreateServer}
              class="btn btn-primary flex-1"
              disabled={isCreating || !newServerName.trim()}
            >
              {#if isCreating}
                <span class="loading loading-spinner"></span>
                Creating...
              {:else}
                Create Server
              {/if}
            </button>
            <button onclick={resetForm} class="btn btn-ghost">Cancel</button>
          </div>
        </div>
      {:else}
        <button
          onclick={() => (showCreateServer = true)}
          class="btn btn-outline btn-lg gap-2"
        >
          <Plus class="w-5 h-5" />
          New Server
        </button>
      {/if}
    </div>
  {/if}
</main>
