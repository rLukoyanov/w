<script lang="ts">
  import { goto } from "$app/navigation";

  let {
    serverItems = [],
    activeServer = "dashboard",
  }: {
    serverItems: Array<{
      id: string;
      name: string;
      initials: string;
      color: string;
      avatarUrl: string | null;
    }>;
    activeServer?: string;
  } = $props();

  let tooltip = $state<{ text: string; top: number } | null>(null);

  function showTooltip(e: MouseEvent, text: string) {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
    tooltip = { text, top: rect.top + rect.height / 2 };
  }

  function hideTooltip() {
    tooltip = null;
  }
</script>

<aside
  class="flex w-[72px] shrink-0 flex-col items-center gap-2 overflow-y-auto border-r border-border bg-bg-tertiary py-3"
>
  <div class="flex flex-col items-center gap-2">
    {#each serverItems as server}
      <div class="group relative">
        <button
          onclick={() => goto(`/server/${server.id}`)}
          onmouseenter={(e) => showTooltip(e, server.name)}
          onmouseleave={hideTooltip}
          class="flex h-12 w-12 items-center justify-center overflow-hidden rounded-2xl text-base font-bold text-white transition hover:rounded-xl {!server.avatarUrl
            ? server.color
            : ''}"
        >
          {#if server.avatarUrl}
            <img
              src={server.avatarUrl}
              alt={server.name}
              class="h-full w-full object-cover"
            />
          {:else}
            {server.initials}
          {/if}
        </button>

        {#if activeServer === server.id}
          <div
            class="absolute -left-3 top-1/2 h-10 w-1 -translate-y-1/2 rounded-r-full bg-text-inverse"
          ></div>
        {:else}
          <div
            class="absolute -left-3 top-1/2 h-5 w-1 -translate-y-1/2 rounded-r-full bg-text-inverse opacity-0 transition group-hover:opacity-50"
          ></div>
        {/if}
      </div>
    {/each}
  </div>
</aside>

{#if tooltip}
  <div
    class="pointer-events-none fixed z-50 -translate-y-1/2 rounded-lg bg-surface px-3 py-1.5 text-sm font-medium text-text shadow-lg transition-all duration-150"
    style="left: calc(72px + 12px); top: {tooltip.top}px"
  >
    {tooltip.text}
  </div>
{/if}
