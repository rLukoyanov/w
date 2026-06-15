<script lang="ts">
  interface Props {
    name?: string;
    size?: number;
  }

  let { name, size = 6 }: Props = $props();

  const avatarBgs = [
    "oklch(0.58 0.2 285 / 0.2)",
    "oklch(0.72 0.18 170 / 0.2)",
    "oklch(0.65 0.18 25 / 0.2)",
    "oklch(0.82 0.15 85 / 0.2)",
    "oklch(0.62 0.18 285 / 0.2)",
  ];

  function getColor(n: string): string {
    const index = n.split("").reduce((a, c) => a + c.charCodeAt(0), 0) % avatarBgs.length;
    return avatarBgs[index];
  }

  function initials(n: string): string {
    return n.slice(0, 2).toUpperCase();
  }
</script>

<div
  class="rounded-full flex items-center justify-center shrink-0 font-semibold font-[family-name:var(--font-family-display)]"
  class:w-6={size === 6}
  class:w-8={size === 8}
  class:w-10={size === 10}
  class:w-12={size === 12}
  style="width: {size * 0.25}rem; height: {size * 0.25}rem; font-size: {size * 0.125}rem; background: {name ? getColor(name) : 'oklch(0.2 0.01 285)'}; color: oklch(0.85 0.01 285);"
>
  {#if name}
    {initials(name)}
  {:else}
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3" style="opacity: 0.5;">
      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  {/if}
</div>
