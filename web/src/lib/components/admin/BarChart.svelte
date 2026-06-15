<script lang="ts">
  let { data, height = 160, barGap = 4 }: {
    data: { label: string; value: number; color?: string }[];
    height?: number;
    barGap?: number;
  } = $props();

  let maxVal = $derived(Math.max(...data.map((d) => d.value), 1));
</script>

<div class="flex items-end gap-[{barGap}px] w-full" style="height: {height}px;">
  {#each data as item}
    <div class="flex flex-col items-center justify-end flex-1 h-full">
      <span class="text-[10px] font-medium mb-1" style="color: oklch(0.6 0.01 285);">{item.value}</span>
      <div
        class="w-full rounded-t transition-all duration-500"
        style="height: {(item.value / maxVal) * (height - 24)}px; min-height: 2px; background: {item.color || 'oklch(0.58 0.2 285)'}; opacity: 0.85;">
      </div>
      <span class="text-[10px] mt-1 truncate max-w-full text-center" style="color: oklch(0.45 0.01 285);">{item.label}</span>
    </div>
  {/each}
</div>
