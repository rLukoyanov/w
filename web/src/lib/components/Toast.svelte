<script lang="ts">
  import { notify, type Notification } from "$lib/stores/notifications";
  import { fade, slide } from "svelte/transition";

  let items = $state<Notification[]>([]);

  notify.subscribe((n) => (items = n));

  const typeStyle: Record<string, string> = {
    error: "oklch(0.65 0.18 25)",
    success: "oklch(0.72 0.18 170)",
    warning: "oklch(0.82 0.15 85)",
    info: "oklch(0.62 0.18 285)",
  };
</script>

{#if items.length > 0}
  <div class="toast toast-end toast-bottom z-50 gap-2">
    {#each items as item (item.id)}
      <div
        transition:slide={{ duration: 200 }}
        class="flex items-center gap-2 px-3 py-2 rounded-lg shadow-lg"
        style="background: oklch(0.12 0.006 285); border: 1px solid {typeStyle[item.type] ?? typeStyle.info}40; color: oklch(0.92 0.004 285);"
      >
        <div class="w-1.5 h-1.5 rounded-full shrink-0"
          style="background: {typeStyle[item.type] ?? typeStyle.info};"></div>
        <span class="text-sm">{item.message}</span>
        <button
          onclick={() => notify.remove(item.id)}
          class="btn btn-xs btn-ghost btn-circle shrink-0"
          aria-label="Dismiss"
          style="color: oklch(0.5 0.01 285);"
        >
          ✕
        </button>
      </div>
    {/each}
  </div>
{/if}
