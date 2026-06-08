<script lang="ts">
  import { notify, type Notification } from "$lib/stores/notifications";
  import { fade, slide } from "svelte/transition";

  let items = $state<Notification[]>([]);

  notify.subscribe((n) => (items = n));

  const typeClass: Record<string, string> = {
    error: "alert-error",
    success: "alert-success",
    warning: "alert-warning",
    info: "alert-info",
  };
</script>

{#if items.length > 0}
  <div class="toast toast-end toast-bottom z-50 gap-2">
    {#each items as item (item.id)}
      <div
        transition:slide={{ duration: 200 }}
        class="alert {typeClass[item.type] ?? 'alert-info'} shadow-lg flex items-center gap-2 py-2 pr-2"
      >
        <span class="text-sm">{item.message}</span>
        <button
          onclick={() => notify.remove(item.id)}
          class="btn btn-xs btn-ghost btn-circle shrink-0"
          aria-label="Dismiss"
        >
          ✕
        </button>
      </div>
    {/each}
  </div>
{/if}
