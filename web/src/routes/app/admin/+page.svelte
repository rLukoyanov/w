<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Users, Server, MessageSquare, Hash, Activity, ShieldCheck } from "lucide-svelte";
  import BarChart from "$lib/components/admin/BarChart.svelte";
  import DoughnutChart from "$lib/components/admin/DoughnutChart.svelte";

  const brand = "oklch(0.58 0.2 285)";
  const green = "oklch(0.55 0.18 145)";
  const amber = "oklch(0.6 0.18 85)";

  let stats = $state<{
    total_users: number;
    total_servers: number;
    total_channels: number;
    total_messages: number;
    online_users: number;
    admins: number;
    users_by_role: Record<string, number>;
  } | null>(null);
  let error = $state<string | null>(null);

  let roleData = $derived.by(() => {
    if (!stats) return [];
    return Object.entries(stats.users_by_role || {}).map(([role, count], i) => ({
      label: role,
      value: count as number,
      color: i === 0 ? brand : "oklch(0.55 0.15 25)",
    }));
  });

  let serverData = $derived.by(() => {
    if (!stats) return [];
    return [
      { label: "Servers", value: stats.total_servers, color: brand },
      { label: "Channels", value: stats.total_channels, color: amber },
    ];
  });

  let onlineFraction = $derived(stats ? (stats.online_users / Math.max(stats.total_users, 1)) * 100 : 0);

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

<div class="p-6 max-w-6xl mx-auto">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold font-[family-name:var(--font-family-display)]" style="color: oklch(0.92 0.004 285);">Admin</h1>
      <p class="text-sm mt-0.5" style="color: oklch(0.5 0.01 285);">System overview and management</p>
    </div>
    {#if stats}
      <div class="flex items-center gap-2 px-3 py-1.5 rounded-full text-xs"
        style="background: {stats.online_users > 0 ? green + '18' : 'oklch(0.15 0.008 285)'}; color: {stats.online_users > 0 ? green : 'oklch(0.5 0.01 285)'};">
        <span class="w-1.5 h-1.5 rounded-full"
          style="background: {stats.online_users > 0 ? green : 'oklch(0.4 0.01 285)'}; box-shadow: 0 0 6px {stats.online_users > 0 ? green : 'transparent'};"></span>
        {stats.online_users} online
      </div>
    {/if}
  </div>

  {#if error}
    <div class="text-sm p-4 rounded-xl" style="background: oklch(0.2 0.02 25 / 0.3); color: oklch(0.7 0.18 25); border: 1px solid oklch(0.25 0.04 25 / 0.3);">
      Failed to load stats: {error}
    </div>
  {:else if !stats}
    <div class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg" style="color: {brand};"></span>
    </div>
  {:else}
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mb-8">
      <div class="stat-card">
        <div class="stat-icon" style="background: {green}15; color: {green};">
          <Activity class="w-4 h-4" />
        </div>
        <div class="stat-label">Online</div>
        <div class="stat-value" style="color: {green};">{stats.online_users}</div>
        <div class="stat-trend">
          <div class="trend-bar" style="width: {onlineFraction}%; background: {green};"></div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: {brand}15; color: {brand};">
          <Users class="w-4 h-4" />
        </div>
        <div class="stat-label">Users</div>
        <div class="stat-value" style="color: oklch(0.85 0.01 285);">{stats.total_users}</div>
        <div class="stat-sub">{stats.admins} admin{stats.admins !== 1 ? "s" : ""}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: {amber}15; color: {amber};">
          <Server class="w-4 h-4" />
        </div>
        <div class="stat-label">Servers</div>
        <div class="stat-value" style="color: oklch(0.85 0.01 285);">{stats.total_servers}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: {brand}15; color: {brand};">
          <Hash class="w-4 h-4" />
        </div>
        <div class="stat-label">Channels</div>
        <div class="stat-value" style="color: oklch(0.85 0.01 285);">{stats.total_channels}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: {amber}15; color: {amber};">
          <MessageSquare class="w-4 h-4" />
        </div>
        <div class="stat-label">Messages</div>
        <div class="stat-value" style="color: oklch(0.85 0.01 285);">{stats.total_messages}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: {brand}15; color: {brand};">
          <ShieldCheck class="w-4 h-4" />
        </div>
        <div class="stat-label">Admins</div>
        <div class="stat-value" style="color: oklch(0.85 0.01 285);">{stats.admins}</div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <div class="chart-card">
        <div class="text-sm font-semibold mb-4 font-[family-name:var(--font-family-display)]" style="color: oklch(0.75 0.01 285);">Users by Role</div>
        <div class="flex justify-center py-2">
          <DoughnutChart data={roleData} size={160} />
        </div>
      </div>
      <div class="chart-card">
        <div class="text-sm font-semibold mb-4 font-[family-name:var(--font-family-display)]" style="color: oklch(0.75 0.01 285);">Content Overview</div>
        <div class="py-2">
          <BarChart data={serverData} height={150} />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <a href="/app/admin/users" class="nav-card">
        <div class="text-sm font-semibold font-[family-name:var(--font-family-display)]" style="color: oklch(0.82 0.01 285);">Manage Users</div>
        <div class="text-xs mt-1" style="color: oklch(0.5 0.01 285);">View, promote, demote, or remove users</div>
      </a>
      <a href="/app/admin/servers" class="nav-card">
        <div class="text-sm font-semibold font-[family-name:var(--font-family-display)]" style="color: oklch(0.82 0.01 285);">Manage Servers</div>
        <div class="text-xs mt-1" style="color: oklch(0.5 0.01 285);">Browse and remove servers and their content</div>
      </a>
    </div>
  {/if}
</div>

<style>
  :global(.stat-card) {
    background: oklch(0.1 0.005 285);
    border: 1px solid oklch(0.16 0.008 285);
    border-radius: 12px;
    padding: 14px 16px;
    transition: border-color 0.2s, background 0.2s;
  }
  :global(.stat-card:hover) {
    border-color: oklch(0.22 0.012 285);
    background: oklch(0.115 0.006 285);
  }
  :global(.stat-icon) {
    width: 28px;
    height: 28px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 10px;
  }
  :global(.stat-label) {
    font-size: 10px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: oklch(0.5 0.01 285);
    margin-bottom: 2px;
  }
  :global(.stat-value) {
    font-size: 22px;
    font-weight: 700;
    font-family: var(--font-family-display), system-ui, sans-serif;
    line-height: 1.1;
  }
  :global(.stat-sub) {
    font-size: 10px;
    color: oklch(0.45 0.01 285);
    margin-top: 2px;
  }
  :global(.stat-trend) {
    height: 2px;
    background: oklch(0.15 0.008 285);
    border-radius: 2px;
    margin-top: 8px;
    overflow: hidden;
  }
  :global(.trend-bar) {
    height: 100%;
    border-radius: 2px;
    transition: width 0.6s ease;
  }
  :global(.chart-card) {
    background: oklch(0.1 0.005 285);
    border: 1px solid oklch(0.16 0.008 285);
    border-radius: 12px;
    padding: 20px;
    transition: border-color 0.2s;
  }
  :global(.chart-card:hover) {
    border-color: oklch(0.22 0.012 285);
  }
  :global(.nav-card) {
    display: block;
    background: oklch(0.1 0.005 285);
    border: 1px solid oklch(0.16 0.008 285);
    border-radius: 12px;
    padding: 16px 20px;
    transition: border-color 0.2s, background 0.2s;
  }
  :global(.nav-card:hover) {
    border-color: oklch(0.58 0.2 285 / 0.3);
    background: oklch(0.12 0.008 285);
  }
</style>
