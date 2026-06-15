<script lang="ts">
  import type { TagGroup } from "$lib/api/docs";

  let { groups, active, onselect }: {
    groups: TagGroup[];
    active: string;
    onselect: (tag: string) => void;
  } = $props();
</script>

<nav class="flex flex-col gap-1">
  {#each groups as group}
    <button
      class="text-left px-3 py-2 rounded-lg text-sm transition-all duration-150 font-[family-name:var(--font-family-body)]"
      style="background: {active === group.tag ? 'oklch(0.58 0.2 285 / 0.15)' : 'transparent'}; color: {active === group.tag ? 'oklch(0.72 0.18 285)' : 'oklch(0.6 0.01 285)'};"
      onclick={() => onselect(group.tag)}
    >
      <div class="font-medium">{group.tag}</div>
      <div class="text-xs mt-0.5" style="color: oklch(0.45 0.01 285);">
        {group.endpoints.length} {group.endpoints.length === 1 ? "endpoint" : "endpoints"}
      </div>
    </button>
  {/each}
</nav>
