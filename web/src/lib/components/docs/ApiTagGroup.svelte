<script lang="ts">
  import type { TagGroup } from "$lib/api/docs";
  import ApiEndpoint from "./ApiEndpoint.svelte";

  let { group }: { group: TagGroup } = $props();
</script>

<div class="flex flex-col gap-3">
  {#each group.endpoints as endpoint}
    <ApiEndpoint {endpoint} />
  {/each}
</div>
