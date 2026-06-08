<script lang="ts">
  import { fly } from "svelte/transition";
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

  let isActive = (path: string) => page.url.pathname.includes(path);

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
          class="btn btn-neutral text-white btn-ghost gap-2 btn-xs border-0 transition-all duration-200 hover:scale-105 {isActive(
            route.path,
          )
            ? 'bg-white/20 scale-105'
            : 'hover:bg-white/10'}"
          onclick={() => goto(route.path)}
        >
          <route.icon
            class="w-4 h-4 transition-all duration-200 {isActive(
              route.path,
            )
              ? 'scale-110'
              : ''} {route.name === 'Settings' && isActive(route.path)
              ? 'icon-spin'
              : ''}"
          />
          {route.name}
        </button>
      {/each}
    </div>

    {#if $auth.user}
      <details class="dropdown dropdown-end dropdown-hover">
        <summary
          class="btn btn-xs btn-ghost btn-neutral text-white flex gap-2 hover:bg-white/20 border-0 transition-all duration-200"
        >
          <div>{$auth.user.username}</div>
          <Avatar />
        </summary>
        {#key $auth.user.id}
          <ul
            transition:fly={{ duration: 120, y: -4 }}
            class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
          >
            <li class="menu-title">
              <span class="text-xs text-base-content">{$auth.user.email}</span>
            </li>

            <li>
              <button
                onclick={handleLogout}
                class="text-error transition-colors duration-150 hover:bg-error/10"
              >
                Sign out
              </button>
            </li>
          </ul>
        {/key}
      </details>
    {/if}
  </div>
</nav>

<style>
  :global(.icon-spin) {
    animation: spin 3s linear infinite;
  }

  @keyframes spin {
    from { transform: rotate(0deg) scale(1.1); }
    to { transform: rotate(360deg) scale(1.1); }
  }
</style>
