<script lang="ts">
  import type { Channel } from "$lib/api/types";
  import { Hash, Volume2, Trash2 } from "lucide-svelte";
  import { slide } from "svelte/transition";

  interface Props {
    channel: Channel;
    serverId: string;
    channelUrl: string;
    onDelete: (id: string, event: MouseEvent) => void;
  }

  let { channel, channelUrl, onDelete }: Props = $props();
</script>

<div
  transition:slide={{ duration: 200 }}
  class="flex items-center gap-3 px-4 py-3 rounded-lg hover:bg-base-100 transition-colors group"
>
  {#if channel.type === "voice"}
    <Volume2 class="w-5 h-5 text-base-content/40 group-hover:text-primary" />
  {:else}
    <Hash class="w-5 h-5 text-base-content/40 group-hover:text-primary" />
  {/if}
  <a href={channelUrl} class="flex-1 text-lg">
    {channel.name}
  </a>
  <button
    onclick={(e) => onDelete(channel.id, e)}
    class="opacity-0 group-hover:opacity-100 transition-opacity p-1 hover:text-error"
    aria-label="Delete channel {channel.name}"
  >
    <Trash2 class="w-4 h-4" />
  </button>
</div>
