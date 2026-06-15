<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Users, Server, MessageSquare, Hash, Activity, ShieldCheck } from "lucide-svelte";
  import BarChart from "$lib/components/admin/BarChart.svelte";
  import DoughnutChart from "$lib/components/admin/DoughnutChart.svelte";

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
      color: i === 0 ? "#818cf8" : "#f87171",
    }));
  });

  let serverData = $derived.by(() => {
    if (!stats) return [];
    return [
      { label: "Servers", value: stats.total_servers, color: "#818cf8" },
      { label: "Channels", value: stats.total_channels, color: "#fbbf24" },
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
      <h1 class="text-2xl font-bold" style="color: #f1f5f9;">Admin</h1>
      <p class="text-sm mt-0.5" style="color: #94a3b8;">System overview and management</p>
    </div>
    {#if stats}
      <div class="flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium"
        style="background: #{'22c55e'}18; color: #22c55e;">
        <span class="w-1.5 h-1.5 rounded-full" style="background: #22c55e;"></span>
        {stats.online_users} online
      </div>
    {/if}
  </div>

  {#if error}
    <div class="text-sm p-4 rounded-xl" style="background: #450a0a; color: #fca5a5; border: 1px solid #7f1d1d;">
      Failed to load stats: {error}
    </div>
  {:else if !stats}
    <div class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg" style="color: #818cf8;"></span>
    </div>
  {:else}
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mb-8">
      <div class="stat-card">
        <div class="stat-icon" style="background: #22c55e22; color: #22c55e;">
          <Activity class="w-4 h-4" />
        </div>
        <div class="stat-label">Online</div>
        <div class="stat-value" style="color: #22c55e;">{stats.online_users}</div>
        <div class="stat-bar">
          <div class="stat-bar-fill" style="width: {onlineFraction}%; background: #22c55e;"></div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #818cf822; color: #818cf8;">
          <Users class="w-4 h-4" />
        </div>
        <div class="stat-label">Users</div>
        <div class="stat-value" style="color: #f1f5f9;">{stats.total_users}</div>
        <div class="stat-sub">{stats.admins} admin{stats.admins !== 1 ? "s" : ""}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #fbbf2422; color: #fbbf24;">
          <Server class="w-4 h-4" />
        </div>
        <div class="stat-label">Servers</div>
        <div class="stat-value" style="color: #f1f5f9;">{stats.total_servers}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #818cf822; color: #818cf8;">
          <Hash class="w-4 h-4" />
        </div>
        <div class="stat-label">Channels</div>
        <div class="stat-value" style="color: #f1f5f9;">{stats.total_channels}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #fbbf2422; color: #fbbf24;">
          <MessageSquare class="w-4 h-4" />
        </div>
        <div class="stat-label">Messages</div>
        <div class="stat-value" style="color: #f1f5f9;">{stats.total_messages}</div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #818cf822; color: #818cf8;">
          <ShieldCheck class="w-4 h-4" />
        </div>
        <div class="stat-label">Admins</div>
        <div class="stat-value" style="color: #f1f5f9;">{stats.admins}</div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
      <div class="chart-card">
        <div class="text-sm font-semibold mb-4" style="color: #e2e8f0;">Users by Role</div>
        <div class="flex justify-center py-2">
          <DoughnutChart data={roleData} size={160} />
        </div>
      </div>
      <div class="chart-card">
        <div class="text-sm font-semibold mb-4" style="color: #e2e8f0;">Content Overview</div>
        <div class="py-2">
          <BarChart data={serverData} height={150} />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <a href="/app/admin/users" class="nav-card">
        <div class="text-sm font-semibold" style="color: #e2e8f0;">Manage Users</div>
        <div class="text-xs mt-1" style="color: #94a3b8;">View, promote, demote, or remove users</div>
      </a>
      <a href="/app/admin/servers" class="nav-card">
        <div class="text-sm font-semibold" style="color: #e2e8f0;">Manage Servers</div>
        <div class="text-xs mt-1" style="color: #94a3b8;">Browse and remove servers and their content</div>
      </a>
    </div>
  {/if}
</div>

<style>
  :global(.stat-card) {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 12px;
    padding: 14px 16px;
    transition: border-color 0.2s, background 0.2s;
  }
  :global(.stat-card:hover) {
    border-color: #475569;
    background: #1e293b;
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
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: #94a3b8;
    margin-bottom: 2px;
  }
  :global(.stat-value) {
    font-size: 24px;
    font-weight: 800;
    line-height: 1.1;
  }
  :global(.stat-sub) {
    font-size: 11px;
    color: #64748b;
    margin-top: 2px;
  }
  :global(.stat-bar) {
    height: 3px;
    background: #334155;
    border-radius: 3px;
    margin-top: 8px;
    overflow: hidden;
  }
  :global(.stat-bar-fill) {
    height: 100%;
    border-radius: 3px;
    transition: width 0.6s ease;
  }
  :global(.chart-card) {
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 12px;
    padding: 20px;
    transition: border-color 0.2s;
  }
  :global(.chart-card:hover) {
    border-color: #475569;
  }
  :global(.nav-card) {
    display: block;
    background: #1e293b;
    border: 1px solid #334155;
    border-radius: 12px;
    padding: 16px 20px;
    transition: border-color 0.2s, background 0.2s;
  }
  :global(.nav-card:hover) {
    border-color: #6366f1;
    background: #1e293b;
  }
</style>
