<script lang="ts">
  import { fly } from "svelte/transition";
  import Avatar from "./Avatar.svelte";
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { wsClient } from "$lib/websocket";
  import { Server, Settings, BookOpen } from "lucide-svelte";
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
    { name: "Docs", icon: BookOpen, path: ROUTES.DOCS },
    { name: "Settings", icon: Settings, path: ROUTES.SETTINGS },
  ];
</script>

<nav class="flex items-center w-full min-h-10 px-3"
  style="background: oklch(0.07 0.004 285); border-bottom: 1px solid oklch(0.18 0.008 285);">
  <div class="flex w-full justify-between items-center">
    <div class="flex gap-1">
      {#each routes as route}
        <button
          class="btn btn-xs border-0 transition-all duration-200 font-[family-name:var(--font-family-body)]"
          style="background: {isActive(route.path) ? 'oklch(0.58 0.2 285 / 0.15)' : 'transparent'}; color: {isActive(route.path) ? 'oklch(0.72 0.18 285)' : 'oklch(0.6 0.01 285)'}; {isActive(route.path) ? '' : ''}"
          onclick={() => goto(route.path)}
        >
          <route.icon class="w-3.5 h-3.5" />
          {route.name}
        </button>
      {/each}
    </div>

    {#if $auth.user}
      <details class="dropdown dropdown-end dropdown-hover">
        <summary
          class="btn btn-xs border-0 flex gap-2 transition-all duration-200"
          style="background: transparent; color: oklch(0.8 0.01 285);"
        >
          <span class="text-xs font-medium">{$auth.user.username}</span>
          <Avatar name={$auth.user.username} />
        </summary>
        {#key $auth.user.id}
          <ul
            transition:fly={{ duration: 120, y: -4 }}
            class="menu menu-sm dropdown-content rounded-box z-1 mt-3 w-52 p-2 shadow-lg"
            style="background: oklch(0.12 0.006 285); border: 1px solid oklch(0.2 0.01 285);"
          >
            <li class="menu-title">
              <span class="text-xs" style="color: oklch(0.5 0.01 285);">{$auth.user.email}</span>
            </li>
            <li>
              <button
                onclick={handleLogout}
                class="transition-colors duration-150 text-sm"
                style="color: oklch(0.65 0.18 25);"
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
