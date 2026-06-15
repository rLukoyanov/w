<script lang="ts">
  import { onMount } from "svelte";
  import { baseClient, BASE_URL } from "$lib/api/client";
  import { Shield, Trash2 } from "lucide-svelte";

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
    <a href="/app/admin" class="text-sm hover:underline" style="color: #94a3b8;">&larr; Dashboard</a>
    <h1 class="text-xl font-bold" style="color: #f1f5f9;">Users</h1>
  </div>

  {#if loading}
    <div class="flex justify-center py-16">
      <span class="loading loading-spinner loading-lg" style="color: #818cf8;"></span>
    </div>
  {:else if error}
    <div class="text-sm p-4 rounded-xl" style="background: #450a0a; color: #fca5a5; border: 1px solid #7f1d1d;">
      {error}
    </div>
  {:else}
    <div class="overflow-x-auto rounded-xl border" style="border-color: #334155;">
      <table class="w-full">
        <thead>
          <tr style="border-bottom: 1px solid #334155;">
            <th class="text-left px-4 py-3 text-[11px] font-semibold uppercase tracking-wider" style="color: #64748b;">User</th>
            <th class="text-left px-4 py-3 text-[11px] font-semibold uppercase tracking-wider hidden sm:table-cell" style="color: #64748b;">Email</th>
            <th class="text-left px-4 py-3 text-[11px] font-semibold uppercase tracking-wider" style="color: #64748b;">Role</th>
            <th class="text-left px-4 py-3 text-[11px] font-semibold uppercase tracking-wider hidden md:table-cell" style="color: #64748b;">Joined</th>
            <th class="text-right px-4 py-3 text-[11px] font-semibold uppercase tracking-wider" style="color: #64748b;">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user}
            <tr class="table-row" style="border-bottom: 1px solid #1e293b;">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2.5">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
                    style="background: #334155; color: #94a3b8;">
                    {user.username[0].toUpperCase()}
                  </div>
                  <div>
                    <div class="text-sm font-medium" style="color: #e2e8f0;">{user.username}</div>
                    <div class="text-xs sm:hidden" style="color: #64748b;">{user.email}</div>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3 text-sm hidden sm:table-cell" style="color: #94a3b8;">{user.email}</td>
              <td class="px-4 py-3">
                <span class="inline-block text-[11px] font-semibold px-2 py-0.5 rounded-md"
                  style="background: {user.role === 'admin' ? '#818cf822' : '#1e293b'}; color: {user.role === 'admin' ? '#818cf8' : '#94a3b8'};">
                  {user.role}
                </span>
              </td>
              <td class="px-4 py-3 text-sm hidden md:table-cell" style="color: #64748b;">{new Date(user.created_at).toLocaleDateString()}</td>
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
    background: #1e293b;
  }
  :global(.action-btn) {
    background: transparent;
    border: none;
    border-radius: 6px;
    padding: 4px 6px;
    color: #64748b;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }
  :global(.action-btn:hover) {
    background: #334155;
    color: #cbd5e1;
  }
  :global(.action-btn-danger:hover) {
    background: #450a0a;
    color: #fca5a5;
  }
</style>