<script lang="ts">
  import { wsConnected } from "$lib/websocket";
</script>

{#if $wsConnected}
  <div class="badge badge-success gap-2">
    <div class="w-2 h-2 rounded-full bg-success-content animate-pulse"></div>
    Connected
  </div>
{:else}
  <div class="badge badge-error gap-2">
    <div class="w-2 h-2 rounded-full bg-error-content"></div>
    Disconnected
  </div>
{/if}
