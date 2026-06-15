<script lang="ts">
  import { onMount } from "svelte";
  import { usersClient, type UserWithStatus } from "$lib/api";

  let users = $state<UserWithStatus[]>([]);
  let loading = $state(true);

  onMount(async () => {
    try {
      users = await usersClient.getConnected();
    } catch {
      users = [];
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-6 max-w-lg">
  <h1 class="text-2xl font-bold mb-1 font-[family-name:var(--font-family-display)]"
    style="color: oklch(0.92 0.004 285);">Settings</h1>
  <p class="text-sm mb-6" style="color: oklch(0.5 0.01 285);">Online users</p>

  {#if loading}
    <span class="loading loading-spinner loading-sm" style="color: oklch(0.58 0.2 285);"></span>
  {:else if users.length === 0}
    <p class="text-sm" style="color: oklch(0.4 0.01 285);">No users connected</p>
  {:else}
    <div class="space-y-2">
      {#each users as user}
        <div class="flex items-center gap-3 px-3 py-2.5 rounded-lg"
          style="background: oklch(0.12 0.006 285);">
          <div class="w-2 h-2 rounded-full shrink-0"
            style="background: oklch(0.72 0.18 170); box-shadow: 0 0 6px oklch(0.72 0.18 170 / 0.5);"></div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <p class="text-sm font-medium truncate"
                style="color: oklch(0.88 0.005 285);">{user.username}</p>
              {#if user.subscribed_channel_name}
                <span class="text-[10px]" style="color: oklch(0.4 0.01 285);">· #{user.subscribed_channel_name}</span>
              {/if}
            </div>
            <p class="text-xs truncate" style="color: oklch(0.5 0.01 285);">{user.email}</p>
          </div>
          <span class="text-[10px] uppercase font-medium"
            style="color: oklch(0.72 0.18 170);">online</span>
        </div>
      {/each}
    </div>
  {/if}
</div>
