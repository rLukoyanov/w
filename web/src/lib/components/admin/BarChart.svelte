<script lang="ts">
  let { data, height = 150 }: {
    data: { label: string; value: number; color?: string }[];
    height?: number;
  } = $props();

  const barColor = "oklch(0.58 0.2 285)";
  let maxVal = $derived(Math.max(...data.map((d) => d.value), 1));
</script>

<div class="flex items-end justify-center gap-3 w-full px-2" style="height: {height}px;">
  {#each data as item}
    <div class="flex flex-col items-center justify-end flex-1 max-w-[60px] h-full">
      <span class="text-[11px] font-medium mb-1.5 tabular-nums" style="color: oklch(0.65 0.01 285);">{item.value}</span>
      <div class="w-full rounded-full transition-all duration-700 ease-out"
        style="height: {(item.value / maxVal) * (height - 36)}px; min-height: 3px;
               background: linear-gradient(to top, {item.color || barColor}, {item.color || barColor}cc);
               box-shadow: 0 0 8px {item.color || barColor}44;">
      </div>
      <span class="text-[10px] mt-1.5 font-medium tracking-wide uppercase" style="color: oklch(0.5 0.01 285);">{item.label}</span>
    </div>
  {/each}
</div>
