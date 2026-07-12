<script lang="ts">
  import { page } from "$app/stores";
  import { routes } from "$lib/consts";
  import UserDropdown from "./UserDropdown.svelte";

  let {
    name = "",
    email = "",
    avatarUrl = null as string | null,
    onlogout = () => {},
    onsettings = () => {},
  } = $props();

  let search = $state("");

  let currentRoute = $derived(
    routes.find((r) => r.path === $page.url.pathname),
  );
</script>

<header
  class="flex items-center gap-4 border-b border-border bg-bg-secondary px-6 py-2"
>
  <a
    href="/dashboard"
    class="flex items-center gap-2 rounded-lg p-1.5 transition hover:bg-bg-tertiary"
  >
    <span
      class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-base font-bold text-text-inverse"
    >
      В
    </span>
    {#if currentRoute}
      <span class="text-sm font-medium text-text"
        >{currentRoute.label}</span
      >
    {/if}
  </a>

  <nav class="ml-auto flex items-center gap-2">
    <div class="relative">
      <svg
        class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        />
      </svg>
      <input
        type="text"
        bind:value={search}
        placeholder="Поиск..."
        class="w-full rounded-lg border border-border bg-bg py-1.5 pl-9 pr-3 text-sm text-text outline-none transition focus:border-primary focus:ring-1 focus:ring-primary"
      />
    </div>
    <button
      title="Уведомления"
      class="relative rounded-lg p-2 text-text-muted transition hover:bg-bg-tertiary hover:text-text"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-5 w-5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
        />
      </svg>
      <span
        class="absolute right-2 top-2 h-2 w-2 rounded-full bg-danger"
      ></span>
    </button>

    <UserDropdown {name} {email} {avatarUrl} {onlogout} {onsettings} />
  </nav>
</header>
