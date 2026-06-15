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

      goto(`/app/servers/${server.id}`);
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
      <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
    </div>
  {:else if ctx.servers.length === 0}
    <div class="flex flex-col items-center justify-center h-full text-center">
      <ServerIcon class="w-24 h-24 mb-6" style="color: oklch(0.2 0.01 285);" />
      <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">Welcome to W</h1>
      <p class="mb-8 max-w-md" style="color: oklch(0.5 0.01 285);">
        You don't have any servers yet. Create your first server to get started!
      </p>

      {#if showCreateServer}
        <div class="w-full max-w-md space-y-4">
          <input
            type="text"
            bind:value={newServerName}
            placeholder="Server name"
            class="input w-full"
            style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
            disabled={isCreating}
          />
          <div class="flex gap-2">
            <button
              onclick={handleCreateServer}
              class="btn flex-1"
              style="background: oklch(0.58 0.2 285); color: white; {isCreating || !newServerName.trim() ? 'opacity: 0.5;' : ''}"
              disabled={isCreating || !newServerName.trim()}
            >
              {#if isCreating}
                <span class="loading loading-spinner"></span>
                Creating...
              {:else}
                Create Server
              {/if}
            </button>
            <button onclick={resetForm} class="btn btn-ghost"
              style="color: oklch(0.6 0.01 285);">Cancel</button>
          </div>
        </div>
      {:else}
        <button
          onclick={() => (showCreateServer = true)}
          class="btn btn-lg gap-2"
          style="background: oklch(0.58 0.2 285); color: white;"
        >
          <Plus class="w-5 h-5"></Plus>
          Create your first server
        </button>
      {/if}
    </div>
  {:else}
    <div class="flex flex-col items-center justify-center h-full text-center">
      <ServerIcon class="w-24 h-24 mb-6" style="color: oklch(0.2 0.01 285);" />
      <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">Select a server</h1>
      <p class="mb-8" style="color: oklch(0.5 0.01 285);">
        Choose a server from the sidebar to start chatting
      </p>

      {#if showCreateServer}
        <div class="w-full max-w-md space-y-4">
          <input
            type="text"
            bind:value={newServerName}
            placeholder="Server name"
            class="input w-full"
            style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
            disabled={isCreating}
          />
          <div class="flex gap-2">
            <button
              onclick={handleCreateServer}
              class="btn flex-1"
              style="background: oklch(0.58 0.2 285); color: white; {isCreating || !newServerName.trim() ? 'opacity: 0.5;' : ''}"
              disabled={isCreating || !newServerName.trim()}
            >
              {#if isCreating}
                <span class="loading loading-spinner"></span>
                Creating...
              {:else}
                Create Server
              {/if}
            </button>
            <button onclick={resetForm} class="btn btn-ghost"
              style="color: oklch(0.6 0.01 285);">Cancel</button>
          </div>
        </div>
      {:else}
        <button
          onclick={() => (showCreateServer = true)}
          class="btn btn-outline btn-lg gap-2"
          style="border-color: oklch(0.3 0.01 285); color: oklch(0.6 0.01 285);"
        >
          <Plus class="w-5 h-5" />
          New Server
        </button>
      {/if}
    </div>
  {/if}
</main>
