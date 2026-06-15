<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Shield, Trash2 } from "lucide-svelte";

  const brand = "oklch(0.58 0.2 285)";

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
</script>

<svelte:head>
  <title>Users — Admin — W</title>
</svelte:head>

<div class="p-6 max-w-6xl mx-auto">
  <div class="flex items-center gap-3 mb-6">
    <a href="/app/admin" class="text-sm hover:underline" style="color: oklch(0.5 0.01 285);">&larr; Dashboard</a>
    <h1 class="text-xl font-bold font-[family-name:var(--font-family-display)]" style="color: oklch(0.9 0.01 285);">Users</h1>
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
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">User</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider hidden sm:table-cell" style="color: oklch(0.5 0.01 285);">Email</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">Role</th>
            <th class="text-left px-4 py-3 text-[10px] font-semibold uppercase tracking-wider hidden md:table-cell" style="color: oklch(0.5 0.01 285);">Joined</th>
            <th class="text-right px-4 py-3 text-[10px] font-semibold uppercase tracking-wider" style="color: oklch(0.5 0.01 285);">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user}
            <tr class="table-row" style="border-bottom: 1px solid oklch(0.12 0.006 285);">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                    style="background: oklch(0.18 0.01 285); color: oklch(0.7 0.01 285);">
                    {user.username[0].toUpperCase()}
                  </div>
                  <div>
                    <div class="text-sm font-medium" style="color: oklch(0.85 0.01 285);">{user.username}</div>
                    <div class="text-xs sm:hidden" style="color: oklch(0.5 0.01 285);">{user.email}</div>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3 text-sm hidden sm:table-cell" style="color: oklch(0.55 0.01 285);">{user.email}</td>
              <td class="px-4 py-3">
                <span class="inline-block text-[11px] font-semibold px-2 py-0.5 rounded-md"
                  style="background: {user.role === 'admin' ? brand + '18' : 'oklch(0.14 0.006 285)'}; color: {user.role === 'admin' ? brand : 'oklch(0.5 0.01 285)'};">
                  {user.role}
                </span>
              </td>
              <td class="px-4 py-3 text-sm hidden md:table-cell" style="color: oklch(0.45 0.01 285);">{new Date(user.created_at).toLocaleDateString()}</td>
              <td class="px-4 py-3 text-right">
                <div class="flex gap-1 justify-end">
                  <button class="action-btn"
                    title={user.role === "admin" ? "Demote to user" : "Promote to admin"}
                    onclick={() => toggleRole(user.id, user.role)}>
                    <Shield class="w-3.5 h-3.5" />
                  </button>
                  <button class="action-btn action-btn-danger"
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
