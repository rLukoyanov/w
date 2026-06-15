<script lang="ts">
  import { wsConnected } from "$lib/websocket";
</script>

{#if $wsConnected}
  <div class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-medium"
    style="background: oklch(0.72 0.18 170 / 0.12); color: oklch(0.72 0.18 170);">
    <div class="w-1.5 h-1.5 rounded-full" style="background: oklch(0.72 0.18 170); box-shadow: 0 0 4px oklch(0.72 0.18 170 / 0.5);"></div>
    Connected
  </div>
{:else}
  <div class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-medium"
    style="background: oklch(0.65 0.18 25 / 0.12); color: oklch(0.65 0.18 25);">
    <div class="w-1.5 h-1.5 rounded-full" style="background: oklch(0.65 0.18 25);"></div>
    Disconnected
  </div>
{/if}
