<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Trash2, Server } from "lucide-svelte";

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

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString();
  }
</script>

<svelte:head>
  <title>Servers — Admin — W</title>
</svelte:head>

<div class="p-6">
  <div class="flex items-center gap-3 mb-6">
    <a href="/app/admin" class="text-sm" style="color: oklch(0.5 0.01 285);">&larr; Dashboard</a>
    <h1 class="text-xl font-bold" style="color: oklch(0.9 0.01 285);">Servers</h1>
  </div>

  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
    </div>
  {:else if error}
    <div class="text-sm p-4 rounded-lg" style="background: oklch(0.2 0.02 25 / 0.3); color: oklch(0.7 0.18 25);">
      {error}
    </div>
  {:else}
    <div class="overflow-x-auto rounded-xl border" style="border-color: oklch(0.15 0.008 285);">
      <table class="table table-sm w-full">
        <thead>
          <tr style="color: oklch(0.5 0.01 285);">
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Server</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Owner</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Channels</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Created</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each servers as server}
            <tr class="border-0" style="border-bottom: 1px solid oklch(0.12 0.006 285);">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <Server class="w-4 h-4" style="color: oklch(0.58 0.2 285);" />
                  <span class="text-sm font-medium" style="color: oklch(0.85 0.01 285);">{server.name}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.6 0.01 285);">{server.owner_name || server.owner_id.slice(0, 8)}</td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.5 0.01 285);">{server.channels}</td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.5 0.01 285);">{formatDate(server.created_at)}</td>
              <td class="px-4 py-3">
                <button
                  class="btn btn-xs btn-ghost"
                  style="color: oklch(0.65 0.2 25);"
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
