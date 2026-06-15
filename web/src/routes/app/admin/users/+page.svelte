<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Shield, User as UserIcon, Trash2 } from "lucide-svelte";

  let users = $state<any[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(() => loadUsers());

  async function loadUsers() {
    loading = true;
    error = null;
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/users`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("Failed to load users");
      users = await res.json();
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function toggleRole(userId: string, currentRole: string) {
    const newRole = currentRole === "admin" ? "user" : "admin";
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/users/${userId}/role`, {
        method: "PATCH",
        headers: token ? { Authorization: `Bearer ${token}`, "Content-Type": "application/json" } : {},
        body: JSON.stringify({ role: newRole }),
      });
      if (!res.ok) throw new Error("Failed to update role");
      users = users.map((u) => (u.id === userId ? { ...u, role: newRole } : u));
    } catch (e: any) {
      alert(e.message);
    }
  }

  async function deleteUser(userId: string, username: string) {
    if (!confirm(`Delete user "${username}"? This cannot be undone.`)) return;
    try {
      const token = baseClient.getToken();
      const res = await fetch(`${BASE_URL}/api/admin/users/${userId}`, {
        method: "DELETE",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) throw new Error("Failed to delete user");
      users = users.filter((u) => u.id !== userId);
    } catch (e: any) {
      alert(e.message);
    }
  }

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString();
  }
</script>

<svelte:head>
  <title>Users — Admin — W</title>
</svelte:head>

<div class="p-6">
  <div class="flex items-center gap-3 mb-6">
    <a href="/app/admin" class="text-sm" style="color: oklch(0.5 0.01 285);">&larr; Dashboard</a>
    <h1 class="text-xl font-bold" style="color: oklch(0.9 0.01 285);">Users</h1>
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
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">User</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Email</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Role</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Joined</th>
            <th class="px-4 py-3 text-xs font-medium uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user}
            <tr class="border-0" style="border-bottom: 1px solid oklch(0.12 0.006 285);">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold"
                    style="background: oklch(0.2 0.01 285); color: oklch(0.7 0.01 285);">
                    {user.username[0].toUpperCase()}
                  </div>
                  <span class="text-sm font-medium" style="color: oklch(0.85 0.01 285);">{user.username}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.6 0.01 285);">{user.email}</td>
              <td class="px-4 py-3">
                <span class="text-xs font-mono px-2 py-0.5 rounded"
                  style="background: {user.role === 'admin' ? 'oklch(0.58 0.2 285 / 0.15)' : 'oklch(0.15 0.006 285)'}; color: {user.role === 'admin' ? 'oklch(0.72 0.18 285)' : 'oklch(0.55 0.01 285)'};">
                  {user.role}
                </span>
              </td>
              <td class="px-4 py-3 text-sm" style="color: oklch(0.5 0.01 285);">{formatDate(user.created_at)}</td>
              <td class="px-4 py-3">
                <div class="flex gap-1">
                  <button
                    class="btn btn-xs btn-ghost"
                    title={user.role === "admin" ? "Demote to user" : "Promote to admin"}
                    onclick={() => toggleRole(user.id, user.role)}>
                    <Shield class="w-3.5 h-3.5" />
                  </button>
                  <button
                    class="btn btn-xs btn-ghost"
                    style="color: oklch(0.65 0.2 25);"
                    title="Delete user"
                    onclick={() => deleteUser(user.id, user.username)}>
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
