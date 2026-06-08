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
  <h1 class="text-2xl font-bold mb-1">Settings</h1>
  <p class="text-sm text-base-content/60 mb-6">Online users</p>

  {#if loading}
    <span class="loading loading-spinner loading-sm"></span>
  {:else if users.length === 0}
    <p class="text-sm text-base-content/40">No users connected</p>
  {:else}
    <div class="space-y-2">
      {#each users as user}
        <div class="flex items-center gap-3 px-3 py-2 rounded-lg bg-base-100">
          <div
            class="w-2 h-2 rounded-full bg-success shrink-0"
          ></div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <p class="text-sm font-medium truncate">{user.username}</p>
              {#if user.subscribed_channel_name}
                <span class="text-[10px] text-base-content/40">· #{user.subscribed_channel_name}</span>
              {/if}
            </div>
            <p class="text-xs text-base-content/50 truncate">{user.email}</p>
          </div>
          <span class="text-[10px] text-base-content/40 uppercase">online</span>
        </div>
      {/each}
    </div>
  {/if}
</div>
