<script lang="ts">
  import type { Server } from "$lib/api";
  import { ROUTES } from "$lib/routes";

  interface Props {
    server: Server;
    active: boolean;
  }

  let { server, active }: Props = $props();

  const avatarColors = [
    "oklch(0.58 0.2 285)",
    "oklch(0.72 0.18 170)",
    "oklch(0.62 0.18 285)",
    "oklch(0.55 0.15 200)",
    "oklch(0.65 0.15 30)",
    "oklch(0.65 0.18 25)",
    "oklch(0.82 0.15 85)",
  ];

  function getColor(name: string): string {
    const index =
      name.split("").reduce((acc, c) => acc + c.charCodeAt(0), 0) %
      avatarColors.length;
    return avatarColors[index];
  }
</script>

<a
  href={ROUTES.SERVER.DETAIL(server.id)}
  class="tooltip tooltip-right"
  data-tip={server.name}
>
  <div
    class="w-12 h-12 rounded-2xl flex items-center justify-center font-bold text-lg transition-all duration-200 text-white font-[family-name:var(--font-family-display)]"
    style="background: {getColor(server.name)}; border-radius: {active ? '0.75rem' : '1rem'}; box-shadow: {active ? '0 0 16px ' + getColor(server.name) + '/0.4' : 'none'}; transform: {active ? 'scale(1.1)' : 'scale(1)'};"
  >
    {server.name[0]}
  </div>
</a>
