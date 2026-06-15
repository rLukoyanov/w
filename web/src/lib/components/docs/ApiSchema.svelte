<script lang="ts">
  import type { SchemaObject } from "$lib/api/docs";
  import ApiSchema from "./ApiSchema.svelte";

  let { schema, name }: { schema: SchemaObject; name?: string } = $props();
</script>

{#if !schema.properties && !schema.items && !schema.enum}
  <span class="inline-flex items-center gap-1.5">
    {#if name}
      <span class="font-medium text-xs" style="color: oklch(0.7 0.15 285);">{name}:</span>
    {/if}
    <span class="font-mono text-xs px-1.5 py-0.5 rounded"
      style="background: oklch(0.15 0.006 285); color: oklch(0.75 0.12 185);">
      {schema.type || "any"}{schema.format ? `<${schema.format}>` : ""}
    </span>
    {#if schema.description}
      <span class="text-xs" style="color: oklch(0.5 0.01 285);">// {schema.description}</span>
    {/if}
  </span>
{:else if schema.enum}
  <div class="flex flex-wrap items-center gap-1">
    {#if name}
      <span class="font-medium text-xs" style="color: oklch(0.7 0.15 285);">{name}:</span>
    {/if}
    <span class="font-mono text-xs px-1.5 py-0.5 rounded" style="background: oklch(0.15 0.006 285); color: oklch(0.75 0.12 285);">enum</span>
    <span class="font-mono text-xs" style="color: oklch(0.55 0.01 285);">[ {schema.enum.map(v => `"${v}"`).join(", ")} ]</span>
  </div>
{:else if schema.properties}
  <div class="flex flex-col gap-1">
    {#if name}
      <div class="text-xs font-medium" style="color: oklch(0.7 0.15 285);">{name}:</div>
    {/if}
    <div class="flex flex-col gap-1.5 pl-3 border-l-2" style="border-color: oklch(0.2 0.01 285);">
      {#each Object.entries(schema.properties) as [propName, propSchema]}
        <div class="flex items-start gap-2">
          <span class="font-mono text-xs font-medium shrink-0" style="color: oklch(0.8 0.01 285);">{propName}</span>
          {#if schema.required?.includes(propName)}
            <span class="text-xs font-mono shrink-0" style="color: oklch(0.65 0.2 25);">required</span>
          {/if}
          <ApiSchema schema={propSchema} />
        </div>
      {/each}
    </div>
  </div>
{:else if schema.items}
  <div class="flex flex-col gap-1">
    {#if name}
      <span class="font-medium text-xs" style="color: oklch(0.7 0.15 285);">{name}:</span>
    {/if}
    <div class="flex items-center gap-1">
      <span class="font-mono text-xs" style="color: oklch(0.55 0.01 285);">array&lt;</span>
      <ApiSchema schema={schema.items} />
      <span class="font-mono text-xs" style="color: oklch(0.55 0.01 285);">&gt;</span>
    </div>
  </div>
{/if}
