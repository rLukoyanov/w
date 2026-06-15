<script lang="ts">
  let { data, size = 150 }: {
    data: { label: string; value: number; color: string }[];
    size?: number;
  } = $props();

  const bgColor = "#334155";
  let total = $derived(data.reduce((s, d) => s + d.value, 0) || 1);
  let cx = $derived(size / 2);
  let cy = $derived(size / 2);
  let r = $derived(size / 2 - 10);
  let circ = $derived(2 * Math.PI * r);

  let segments = $derived.by(() => {
    let offset = 0;
    return data.map((d) => {
      const pct = d.value / total;
      const len = pct * circ;
      const seg = { ...d, offset, length: len, pct };
      offset -= len;
      return seg;
    });
  });
</script>

<div class="flex flex-col items-center gap-3">
  <svg width={size} height={size} viewBox="0 0 {size} {size}">
    <circle cx={cx} cy={cy} r={r} fill="none" stroke={bgColor} stroke-width="16" />
    {#each segments as seg}
      {#if seg.value > 0}
        <circle cx={cx} cy={cy} r={r} fill="none"
          stroke={seg.color}
          stroke-width="16"
          stroke-linecap="round"
          stroke-dasharray="{seg.length} {circ - Math.max(seg.length, 1)}"
          stroke-dashoffset={seg.offset}
          transform="rotate(-90 {cx} {cy})"
          style="transition: stroke-dasharray 0.6s ease, stroke-dashoffset 0.6s ease;" />
      {/if}
    {/each}
    <text x={cx} y={cy - 3} text-anchor="middle"
      font-size="22" font-weight="700" fill="#f1f5f9">{total}</text>
    <text x={cx} y={cy + 13} text-anchor="middle"
      font-size="10" fill="#94a3b8">total</text>
  </svg>

  <div class="flex flex-wrap justify-center gap-x-4 gap-y-1.5">
    {#each segments as seg}
      <div class="flex items-center gap-1.5 text-xs">
        <div class="w-2 h-2 rounded-full" style="background: {seg.color}; box-shadow: 0 0 4px {seg.color}66;"></div>
        <span style="color: #94a3b8;">{seg.label}</span>
        <span class="font-semibold tabular-nums" style="color: #e2e8f0;">{(seg.pct * 100).toFixed(0)}%</span>
      </div>
    {/each}
  </div>
</div>