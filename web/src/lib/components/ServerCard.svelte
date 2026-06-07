<script lang="ts">
  import type { Server } from "$lib/api";
  import { ROUTES } from "$lib/routes";

  let { server } = $props<{ server: Server }>();
</script>

<a
  href={ROUTES.SERVER.DETAIL(server.id)}
  class="card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-1 group"
>
  <div class="w-14 h-14 flex items-center justify-center">
    {server.name[0]}
  </div>
</a>
