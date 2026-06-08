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
  class="flex items-center gap-1.5 px-2 py-0.5 rounded hover:bg-base-100/50 transition-colors group {active
    ? 'bg-base-100'
    : ''}"
>
  {#if channel.type === "voice"}
    <Volume2 class="w-3.5 h-3.5 text-base-content/40 group-hover:text-primary shrink-0" />
  {:else}
    <Hash class="w-3.5 h-3.5 text-base-content/40 group-hover:text-primary shrink-0" />
  {/if}
  <a href={channelUrl} class="flex-1 text-xs truncate">
    {channel.name}
  </a>
  <button
    onclick={(e) => onDelete(channel.id, e)}
    class="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 hover:text-error shrink-0"
    aria-label="Delete channel {channel.name}"
  >
    <Trash2 class="w-3 h-3" />
  </button>
</div>
