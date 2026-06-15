<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Trash2, Server } from "lucide-svelte";

  const brand = "oklch(0.58 0.2 285)";

  let servers = $state<any[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(() => loadServers());

  async function loadServers() {
    loading = true;
    error = null;
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/servers`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("Failed to load servers");
      servers = await res.json();
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function deleteServer(serverId: string, name: string) {
    if (!confirm(`Delete server "${name}"? All channels and messages will be lost.`)) return;
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/servers/${serverId}`, {
        method: "DELETE",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("Failed to delete server");
      servers = servers.filter((s) => s.id !== serverId);
    } catch (e: any) {
      alert(e.message);
    }
  }
</script>

<svelte:head>
  <title>Servers — Admin — W</title>
</svelte:head>

<div class="p-6 max-w-6xl mx-auto">
  <div class="flex items-center gap-3 mb-6">
    <a href="/app/admin" class="text-sm hover:underline" style="color: oklch(0.5 0.01 285);">&larr; Dashboard</a>
    <h1 class="text-xl font-bold font-[family-name:var(--font-family-display)]" style="color: oklch(0.9 0.01 285);">Servers</h1>
  </div>

  {#if loading}
    <div class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg" style="color: {brand};"></span>
    </div>
  {:else if error}
    <div class="text-sm p-4 rounded-xl" style="background: oklch(0.2 0.02 25 / 0.3); color: oklch(0.7 0.18 25); border: 1px solid oklch(0.25 0.04 25 / 0.3);">
      {error}
    </div>
  {:else}
    <div class="overflow-x-auto rounded-xl border" style="border-color: oklch(0.15 0.008 285);">
      <table class="w-full">
        <thead>
          <tr style="border-bottom: 1px solid oklch(0.15 0.008 285);">
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">Server</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">Owner</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider hidden sm:table-cell" style="color: oklch(0.5 0.01 285);">Channels</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider hidden md:table-cell" style="color: oklch(0.5 0.01 285);">Created</th>
            <th class="text-right px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each servers as server}
            <tr class="table-row" style="border-bottom: 1px solid oklch(0.12 0.006 285);">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0"
                    style="background: {brand}15; color: {brand};">
                    <Server class="w-4 h-4" />
                  </div>
                  <span class="text-sm font-medium" style="color: oklch(0.85 0.01 285);">{server.name}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.55 0.01 285);">{server.owner_name || server.owner_id.slice(0, 8)}</td>
              <td class="px-4 py-3 text-sm hidden sm:table-cell" style="color: oklch(0.45 0.01 285);">{server.channels}</td>
              <td class="px-4 py-3 text-sm hidden md:table-cell" style="color: oklch(0.45 0.01 285);">{new Date(server.created_at).toLocaleDateString()}</td>
              <td class="px-4 py-3 text-right">
                <button class="action-btn action-btn-danger"
                  title="Delete server"
                  onclick={() => deleteServer(server.id, server.name)}>
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  :global(.table-row:hover) {
    background: oklch(0.12 0.006 285);
  }
  :global(.action-btn) {
    background: transparent;
    border: none;
    border-radius: 6px;
    padding: 4px 6px;
    color: oklch(0.5 0.01 285);
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }
  :global(.action-btn:hover) {
    background: oklch(0.18 0.01 285);
    color: oklch(0.7 0.01 285);
  }
  :global(.action-btn-danger:hover) {
    background: oklch(0.2 0.04 25 / 0.3);
    color: oklch(0.65 0.18 25);
  }
</style>
