<script>
  import Avatar from "./Avatar.svelte";
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { wsClient } from "$lib/websocket";
  import { Server, Settings } from "lucide-svelte";
  import { page } from "$app/state";
  import { ROUTES } from "$lib/routes";

  function handleLogout() {
    wsClient.disconnect();
    userClient.clearToken();
    auth.logout();
    goto(ROUTES.AUTH.LOGIN);
  }

  const routes = [
    { name: "Servers", icon: Server, path: ROUTES.SERVER.INDEX },
    { name: "Settings", icon: Settings, path: ROUTES.SETTINGS },
  ];
</script>

<nav class="navbar w-full min-h-8 bg-transparent">
  <div class="flex w-full justify-between items-center">
    <div class="flex gap-1">
      {#each routes as route}
        <button
          class="btn btn-neutral text-white btn-ghost gap-2 btn-xs hover:bg-white/20 border-0 {page.url.pathname.includes(
            route.path,
          )
            ? 'bg-white/20'
            : ''}"
          onclick={() => goto(route.path)}
        >
          <route.icon class="w-4 h-4" />
          {route.name}
        </button>
      {/each}
    </div>

    {#if $auth.user}
      <details class="dropdown dropdown-end dropdown-hover">
        <summary
          class="btn btn-xs btn-ghost btn-neutral text-white flex gap-2 hover:bg-white/20 border-0"
        >
          <div>{$auth.user.username}</div>
          <Avatar />
        </summary>
        <ul
          class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
        >
          <li class="menu-title">
            <span class="text-xs text-base-content">{$auth.user.email}</span>
          </li>

          <li>
            <button onclick={handleLogout} class="text-error">Sign out</button>
          </li>
        </ul>
      </details>
    {/if}
  </div>
</nav>
