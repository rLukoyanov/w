<script lang="ts">
  import type { OpenAPISpec, TagGroup } from "$lib/api/docs";
  import { fetchSpec, groupByTag } from "$lib/api/docs";
  import ApiNav from "./ApiNav.svelte";
  import ApiTagGroup from "./ApiTagGroup.svelte";

  let spec = $state<OpenAPISpec | null>(null);
  let error = $state<string | null>(null);
  let groups = $state<TagGroup[]>([]);
  let activeTag = $state("");
  let loading = $state(true);

  $effect(() => {
    fetchSpec()
      .then((data) => {
        spec = data;
        groups = groupByTag(data);
        if (groups.length > 0) activeTag = groups[0].tag;
      })
      .catch((e: Error) => error = e.message)
      .finally(() => loading = false);
  });

  let activeGroup = $derived(groups.find((g) => g.tag === activeTag));
</script>

<div class="flex h-full">
  <!-- Sidebar -->
  <aside class="w-56 shrink-0 overflow-y-auto border-r p-3" style="background: oklch(0.07 0.004 285); border-color: oklch(0.15 0.008 285);">
    <div class="text-xs font-medium mb-3 px-3" style="color: oklch(0.5 0.01 285);">API Endpoints</div>
    {#if groups.length > 0}
      <ApiNav {groups} active={activeTag} onselect={(t) => activeTag = t} />
    {/if}
  </aside>

  <!-- Content -->
  <main class="flex-1 overflow-y-auto">
    {#if loading}
      <div class="flex items-center justify-center h-full">
        <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
      </div>
    {:else if error}
      <div class="flex items-center justify-center h-full">
        <div class="text-center px-6">
          <div class="text-lg font-medium mb-2" style="color: oklch(0.65 0.2 25);">Failed to load API docs</div>
          <div class="text-sm" style="color: oklch(0.5 0.01 285);">{error}</div>
        </div>
      </div>
    {:else if spec && activeGroup}
      <div class="max-w-4xl mx-auto p-6">
        <div class="mb-8">
          <h1 class="text-2xl font-bold mb-1" style="color: oklch(0.9 0.01 285);">
            {spec.info.title}
            <span class="text-base font-mono font-normal" style="color: oklch(0.5 0.01 285);">v{spec.info.version}</span>
          </h1>
          {#if spec.info.description}
            <p class="text-sm" style="color: oklch(0.6 0.01 285);">{spec.info.description}</p>
          {/if}
        </div>

        <div class="mb-6">
          <h2 class="text-lg font-bold mb-3" style="color: oklch(0.72 0.18 285);">{activeGroup.tag}</h2>
          <ApiTagGroup group={activeGroup} />
        </div>
      </div>
    {/if}
  </main>
</div>
