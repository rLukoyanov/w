<script lang="ts">
  import type { Server } from "$lib/api";
  import { ROUTES } from "$lib/routes";

  interface Props {
    server: Server;
    active: boolean;
  }

  let { server, active }: Props = $props();

  const avatarColors = [
    "bg-primary",
    "bg-secondary",
    "bg-accent",
    "bg-info",
    "bg-success",
    "bg-warning",
    "bg-error",
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
    class="w-12 h-12 rounded-2xl flex items-center justify-center font-bold text-lg transition-all duration-200 {getColor(
      server.name,
    )} text-base-100 {active
      ? 'rounded-xl shadow-md scale-110'
      : 'hover:rounded-xl hover:shadow-md'}"
  >
    {server.name[0]}
  </div>
</a>
