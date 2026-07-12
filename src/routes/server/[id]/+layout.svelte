<script lang="ts">
  import { channels as channelsApi, type Channel } from '$lib/api/channels';
  import { servers as serversApi, type Server } from '$lib/api/servers';
  import { goto } from '$app/navigation';
  import { Hash, Volume2 } from '@lucide/svelte';

  let { params, children } = $props();

  let server = $state<Server | null>(null);
  let channelList = $state<Channel[]>([]);
  let textChannels = $derived(channelList.filter((c) => c.type === 'text'));
  let voiceChannels = $derived(channelList.filter((c) => c.type === 'voice'));

  $effect(() => {
    const id = params.id;
    server = null;
    channelList = [];

    serversApi.get(id).then((s) => (server = s)).catch(console.error);
    channelsApi.list(id).then((c) => (channelList = c)).catch(console.error);
  });
</script>

<aside
  class="flex w-60 shrink-0 flex-col border-r border-border bg-bg-secondary"
>
  <div
    class="flex h-12 shrink-0 items-center border-b border-border px-4 font-semibold text-text shadow-sm"
  >
    {server?.name ?? 'Загрузка...'}
  </div>

  <div class="flex-1 overflow-y-auto py-2">
    {#if textChannels.length > 0}
      <div class="mb-1 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-text-muted">
        Текстовые каналы
      </div>
      {#each textChannels as ch}
        <button
          onclick={() => goto(`/server/${params.id}/text/${ch.id}`)}
          class="flex w-full cursor-pointer items-center gap-2 rounded-md px-3 py-1.5 text-sm text-text-secondary transition hover:bg-bg-tertiary hover:text-text"
        >
          <Hash class="h-4 w-4 shrink-0 text-text-muted" />
          <span class="truncate">{ch.name}</span>
        </button>
      {/each}
    {/if}

    {#if voiceChannels.length > 0}
      <div class="mb-1 mt-3 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-text-muted">
        Голосовые каналы
      </div>
      {#each voiceChannels as ch}
        <button
          class="flex w-full cursor-pointer items-center gap-2 rounded-md px-3 py-1.5 text-sm text-text-secondary transition hover:bg-bg-tertiary hover:text-text"
        >
          <Volume2 class="h-4 w-4 shrink-0 text-text-muted" />
          <span class="truncate">{ch.name}</span>
        </button>
      {/each}
    {/if}

    {#if channelList.length === 0}
      <p class="px-4 py-8 text-center text-sm text-text-muted">
        Каналы не найдены
      </p>
    {/if}
  </div>
</aside>

<div class="flex min-h-0 flex-1">
  {@render children()}
</div>
