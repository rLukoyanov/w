<script lang="ts">
  let { data, size = 140 }: {
    data: { label: string; value: number; color: string }[];
    size?: number;
  } = $props();

  let total = $derived(data.reduce((s, d) => s + d.value, 0) || 1);
  let cx = $derived(size / 2);
  let cy = $derived(size / 2);
  let r = $derived(size / 2 - 8);
  let circumference = $derived(2 * Math.PI * r);

  let segments = $derived.by(() => {
    let offset = 0;
    return data.map((d) => {
      const pct = d.value / total;
      const len = pct * circumference;
      const seg = { ...d, offset, length: len, pct };
      offset += len;
      return seg;
    });
  });
</script>

<div class="flex flex-col items-center gap-3">
  <svg width={size} height={size} viewBox="0 0 {size} {size}">
    <circle cx={cx} cy={cy} r={r} fill="none" stroke="oklch(0.15 0.008 285)" stroke-width="14" />
    {#each segments as seg}
      {#if seg.value > 0}
        <circle cx={cx} cy={cy} r={r} fill="none"
          stroke={seg.color}
          stroke-width="14"
          stroke-dasharray="{seg.length} {circumference - seg.length}"
          stroke-dashoffset={-seg.offset}
          transform="rotate(-90 {cx} {cy})"
          style="transition: stroke-dasharray 0.5s;"
        />
      {/if}
    {/each}
    <text x={cx} y={cy - 4} text-anchor="middle" class="text-lg font-bold" fill="oklch(0.85 0.01 285)">{total}</text>
    <text x={cx} y={cy + 12} text-anchor="middle" class="text-[10px]" fill="oklch(0.5 0.01 285)">total</text>
  </svg>

  <div class="flex flex-wrap gap-3 justify-center">
    {#each segments as seg}
      <div class="flex items-center gap-1.5 text-xs">
        <div class="w-2.5 h-2.5 rounded-sm" style="background: {seg.color};"></div>
        <span style="color: oklch(0.6 0.01 285);">{seg.label}</span>
        <span class="font-medium" style="color: oklch(0.8 0.01 285);">{(seg.pct * 100).toFixed(0)}%</span>
      </div>
    {/each}
  </div>
</div>
