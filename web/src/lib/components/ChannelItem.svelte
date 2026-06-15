<script lang="ts">
  import type { Channel } from "$lib/api/types";
  import { Hash, Volume2, Trash2 } from "lucide-svelte";
  import { slide } from "svelte/transition";

  interface Props {
    channel: Channel;
    serverId: string;
    channelUrl: string;
    onDelete: (id: string, event: MouseEvent) => void;
    active?: boolean;
  }

  let { channel, channelUrl, onDelete, active = false }: Props = $props();
</script>

<div
  transition:slide={{ duration: 200 }}
  class="flex items-center gap-1.5 px-2 py-1 rounded-md transition-colors group"
  style="background: {active ? 'oklch(0.58 0.2 285 / 0.12)' : 'transparent'}; {active ? 'box-shadow: inset 2px 0 0 oklch(0.58 0.2 285);' : ''}"
>
  {#if channel.type === "voice"}
    <Volume2 class="w-3.5 h-3.5 shrink-0" style="color: oklch(0.45 0.01 285);" />
  {:else}
    <Hash class="w-3.5 h-3.5 shrink-0" style="color: {active ? 'oklch(0.58 0.2 285)' : 'oklch(0.45 0.01 285)'};" />
  {/if}
  <a href={channelUrl} class="flex-1 text-xs truncate"
    style="color: {active ? 'oklch(0.92 0.004 285)' : 'oklch(0.6 0.01 285)'};">
    {channel.name}
  </a>
  <button
    onclick={(e) => onDelete(channel.id, e)}
    class="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 shrink-0"
    aria-label="Delete channel {channel.name}"
    style="color: oklch(0.4 0.01 285);"
  >
    <Trash2 class="w-3 h-3" />
  </button>
</div>
