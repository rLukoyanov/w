<script lang="ts">
  import type { EndpointInfo, ParameterObject } from "$lib/api/docs";
  import ApiSchema from "./ApiSchema.svelte";

  let { endpoint }: { endpoint: EndpointInfo } = $props();

  const methodColors: Record<string, string> = {
    get: "oklch(0.6 0.2 145)",
    post: "oklch(0.6 0.2 245)",
    patch: "oklch(0.65 0.2 85)",
    delete: "oklch(0.6 0.2 25)",
  };

  let methodColor = $derived(methodColors[endpoint.method] || "oklch(0.6 0.01 285)");

  let op = $derived(endpoint.operation);

  let params = $derived((op.parameters || []).filter((p): p is ParameterObject => "name" in p));
  let hasParams = $derived(params.length > 0);
  let hasBody = $derived(!!op.requestBody);
  let hasResponses = $derived(Object.keys(op.responses).length > 0);
  let expanded = $state(false);
</script>

<div class="rounded-xl border transition-all duration-200"
  style="background: oklch(0.1 0.005 285); border-color: {expanded ? 'oklch(0.25 0.015 285)' : 'oklch(0.15 0.008 285)'};">
  <button
    class="w-full flex items-center gap-3 px-4 py-3 text-left"
    onclick={() => expanded = !expanded}
  >
    <span class="font-mono text-xs font-bold uppercase tracking-wider px-2 py-1 rounded min-w-[4.5rem] text-center"
      style="background: {methodColor}22; color: {methodColor};">
      {endpoint.method}
    </span>
    <span class="font-mono text-sm flex-1" style="color: oklch(0.8 0.01 285);">
      {endpoint.path}
    </span>
    {#if op.summary}
      <span class="text-sm hidden lg:inline" style="color: oklch(0.55 0.01 285);">{op.summary}</span>
    {/if}
    <span class="text-xs transition-transform duration-200" style="color: oklch(0.45 0.01 285); transform: rotate({expanded ? 180 : 0}deg);">
      ▼
    </span>
  </button>

  {#if expanded}
    <div class="px-4 pb-4 flex flex-col gap-4 border-t" style="border-color: oklch(0.15 0.008 285);">
      {#if op.description}
        <div class="pt-3 text-sm" style="color: oklch(0.65 0.01 285);">{op.description}</div>
      {/if}

      {#if op.summary && op.summary !== op.description}
        <div class="pt-3 text-sm" style="color: oklch(0.65 0.01 285);">{op.summary}</div>
      {/if}

      {#if hasParams}
        <div>
          <div class="text-xs font-medium mb-2" style="color: oklch(0.6 0.01 285);">Parameters</div>
          <div class="flex flex-col gap-1">
            {#each params as param}
              <div class="flex items-start gap-3 py-1.5 px-3 rounded-lg" style="background: oklch(0.12 0.006 285);">
                <span class="font-mono text-sm font-medium min-w-[8rem]" style="color: oklch(0.8 0.01 285);">{param.name}</span>
                <span class="font-mono text-xs px-1.5 py-0.5 rounded" style="background: oklch(0.18 0.008 285); color: oklch(0.55 0.01 285);">{param.in}</span>
                {#if param.required}
                  <span class="text-xs font-mono" style="color: oklch(0.65 0.2 25);">required</span>
                {/if}
                <span class="font-mono text-xs" style="color: oklch(0.75 0.12 185);">{param.schema?.type || "string"}</span>
                {#if param.description}
                  <span class="text-xs flex-1" style="color: oklch(0.5 0.01 285);">{param.description}</span>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if hasBody}
        <div>
          <div class="text-xs font-medium mb-2" style="color: oklch(0.6 0.01 285);">
            Request Body
            {#if op.requestBody!.required}
              <span class="font-mono" style="color: oklch(0.65 0.2 25);">required</span>
            {/if}
          </div>
          <div class="p-3 rounded-lg" style="background: oklch(0.12 0.006 285);">
            {#each Object.entries(op.requestBody!.content) as [contentType, media]}
              <div class="flex flex-col gap-2">
                <span class="font-mono text-xs" style="color: oklch(0.5 0.01 285);">{contentType}</span>
                <ApiSchema schema={media.schema} />
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if hasResponses}
        <div>
          <div class="text-xs font-medium mb-2" style="color: oklch(0.6 0.01 285);">Responses</div>
          <div class="flex flex-col gap-2">
            {#each Object.entries(op.responses) as [status, resp]}
              <div class="p-3 rounded-lg" style="background: oklch(0.12 0.006 285);">
                <div class="flex items-center gap-2 mb-2">
                  <span class="font-mono text-sm font-bold" style="color: {Number(status) < 300 ? 'oklch(0.6 0.2 145)' : Number(status) < 400 ? 'oklch(0.65 0.2 85)' : 'oklch(0.6 0.2 25)'};">{status}</span>
                  <span class="text-sm" style="color: oklch(0.6 0.01 285);">{resp.description}</span>
                </div>
                {#if resp.content}
                  {#each Object.entries(resp.content) as [contentType, media]}
                    <div class="flex flex-col gap-1">
                      <span class="font-mono text-xs" style="color: oklch(0.5 0.01 285);">{contentType}</span>
                      <ApiSchema schema={media.schema} />
                    </div>
                  {/each}
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
