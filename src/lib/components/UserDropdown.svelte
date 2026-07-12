<script lang="ts">
  import { Settings, LogOut } from "@lucide/svelte";

  let {
    name = "",
    email = "",
    avatarUrl = null as string | null,
    onlogout = () => {},
    onsettings = () => {},
  } = $props();

  let dropdownOpen = $state(false);
  let dropdownRef = $state<HTMLDivElement>();

  function toggleDropdown() {
    dropdownOpen = !dropdownOpen;
  }

  function onClickOutside(e: MouseEvent) {
    if (dropdownRef && !dropdownRef.contains(e.target as Node)) {
      dropdownOpen = false;
    }
  }

  function handleLogout() {
    onlogout();
    dropdownOpen = false;
  }

  function handleSettings() {
    onsettings();
    dropdownOpen = false;
  }

  $effect(() => {
    if (dropdownOpen) {
      document.addEventListener("click", onClickOutside);
      return () => document.removeEventListener("click", onClickOutside);
    }
  });
</script>

<div class="relative ml-2" bind:this={dropdownRef}>
  <button
    onclick={toggleDropdown}
    class="flex items-center gap-2 rounded-lg p-1.5 transition hover:bg-bg-tertiary"
  >
    <div
      class="flex h-7 w-7 items-center justify-center overflow-hidden rounded-full bg-primary-light text-xs font-bold text-primary"
    >
      {#if avatarUrl}
        <img
          src={avatarUrl}
          alt={name}
          class="h-full w-full object-cover"
        />
      {:else}
        {name[0]}
      {/if}
    </div>
  </button>

  {#if dropdownOpen}
    <div
      class="absolute right-0 top-full mt-2 w-60 rounded-xl border border-border bg-surface p-2 shadow-lg"
    >
      <div class="flex items-center gap-3 px-3 py-3">
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-full bg-primary-light text-sm font-bold text-primary"
        >
          {#if avatarUrl}
            <img
              src={avatarUrl}
              alt={name}
              class="h-full w-full object-cover"
            />
          {:else}
            {name[0]}
          {/if}
        </div>
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-text">
            {name}
          </p>
          <p class="truncate text-xs text-text-muted">{email}</p>
        </div>
      </div>

      <hr class="mx-2 border-border" />

      <a
        href="/settings"
        onclick={handleSettings}
        class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-text transition hover:bg-bg-tertiary"
      >
        <Settings class="h-4 w-4" />
        Настройки
      </a>

      <hr class="mx-2 border-border" />

      <button
        onclick={handleLogout}
        class="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-danger transition hover:bg-bg-tertiary"
      >
        <LogOut class="h-4 w-4" />
        Выйти
      </button>
    </div>
  {/if}
</div>
