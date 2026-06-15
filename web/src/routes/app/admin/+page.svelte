<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";

  let stats = $state<{
    total_users: number;
    total_servers: number;
    total_channels: number;
    total_messages: number;
    online_users: number;
  } | null>(null);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/stats`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("Failed to load stats");
      stats = await res.json();
    } catch (e: any) {
      error = e.message;
    }
  });
</script>

<svelte:head>
  <title>Admin — W</title>
</svelte:head>

<div class="p-6">
  <h1 class="text-2xl font-bold mb-1" style="color: oklch(0.9 0.01 285);">Admin Dashboard</h1>
  <p class="text-sm mb-6" style="color: oklch(0.5 0.01 285);">System overview and management</p>

  {#if error}
    <div class="text-sm p-4 rounded-lg" style="background: oklch(0.2 0.02 25 / 0.3); color: oklch(0.7 0.18 25);">
      Failed to load stats: {error}
    </div>
  {:else if !stats}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
      <div class="rounded-xl p-4 border" style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);">
        <div class="text-xs font-medium mb-1" style="color: oklch(0.5 0.01 285);">Online Users</div>
        <div class="text-3xl font-bold" style="color: oklch(0.6 0.2 145);">{stats.online_users}</div>
      </div>
      <div class="rounded-xl p-4 border" style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);">
        <div class="text-xs font-medium mb-1" style="color: oklch(0.5 0.01 285);">Total Users</div>
        <div class="text-3xl font-bold" style="color: oklch(0.72 0.18 285);">{stats.total_users}</div>
      </div>
      <div class="rounded-xl p-4 border" style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);">
        <div class="text-xs font-medium mb-1" style="color: oklch(0.5 0.01 285);">Servers</div>
        <div class="text-3xl font-bold" style="color: oklch(0.65 0.2 85);">{stats.total_servers}</div>
      </div>
      <div class="rounded-xl p-4 border" style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);">
        <div class="text-xs font-medium mb-1" style="color: oklch(0.5 0.01 285);">Channels</div>
        <div class="text-3xl font-bold" style="color: oklch(0.6 0.2 245);">{stats.total_channels}</div>
      </div>
      <div class="rounded-xl p-4 border" style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);">
        <div class="text-xs font-medium mb-1" style="color: oklch(0.5 0.01 285);">Messages</div>
        <div class="text-3xl font-bold" style="color: oklch(0.82 0.15 85);">{stats.total_messages}</div>
      </div>
    </div>

    <div class="mt-8 flex gap-4">
      <a href="/app/admin/users"
        class="flex-1 rounded-xl p-4 border transition-colors admin-card"
        style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);"
        >
        <div class="text-sm font-medium" style="color: oklch(0.8 0.01 285);">Manage Users</div>
        <div class="text-xs mt-1" style="color: oklch(0.5 0.01 285);">View all users, change roles, remove accounts</div>
      </a>
      <a href="/app/admin/servers"
        class="flex-1 rounded-xl p-4 border transition-colors admin-card"
        style="background: oklch(0.1 0.005 285); border-color: oklch(0.18 0.008 285);"
        >
        <div class="text-sm font-medium" style="color: oklch(0.8 0.01 285);">Manage Servers</div>
        <div class="text-xs mt-1" style="color: oklch(0.5 0.01 285);">Browse all servers, remove inappropriate content</div>
      </a>
    </div>
  {/if}
</div>

<style>
  .admin-card:hover {
    border-color: oklch(0.3 0.01 285) !important;
  }
</style>
