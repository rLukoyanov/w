<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";
  import { user } from "$lib/api/users";
  import { pb } from "$lib/api/pocketbase";
  import { page } from "$app/stores";
  import { invalidateAll } from "$app/navigation";
  import { servers as serversApi } from "$lib/api/servers";
  import { onMount } from "svelte";
  import Header from "$lib/components/Header.svelte";
  import ServerNav from "$lib/components/ServerNav.svelte";

  let { children } = $props();

  let name = $derived(user.current?.username || "Пользователь");
  let email = $derived(user.current?.email || "");
  let avatarUrl = $derived(
    user.current?.avatar
      ? pb.files.getUrl(user.current, user.current.avatar, { thumb: "32x32" })
      : null,
  );
  let isAuthPage = $derived(
    $page.url.pathname === "/login" || $page.url.pathname === "/register",
  );
  let isServerPage = $derived($page.url.pathname.startsWith("/server/"));

  const serverColors = [
    "bg-blue-500",
    "bg-purple-500",
    "bg-green-500",
    "bg-amber-500",
    "bg-rose-500",
    "bg-cyan-500",
    "bg-indigo-500",
    "bg-emerald-500",
  ];

  let serverItems = $state<
    Array<{
      id: string;
      name: string;
      initials: string;
      color: string;
      avatarUrl: string | null;
    }>
  >([]);

  onMount(async () => {
    try {
      const records = await serversApi.list();
      serverItems = records.map((s, i) => ({
        id: s.id,
        name: s.name,
        initials: s.name[0]?.toUpperCase() || "?",
        color: serverColors[i % serverColors.length],
        avatarUrl: s.avatar
          ? pb.files.getUrl(s, s.avatar, { thumb: "48x48" })
          : null,
      }));
    } catch (e) {
      console.error("Failed to load servers", e);
    }
  });

  let activeServer = $derived(
    $page.url.pathname.startsWith("/server/")
      ? $page.url.pathname.replace("/server/", "")
      : "",
  );

  async function handleLogout() {
    user.logout();
    await invalidateAll();
  }
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen flex-col bg-bg">
  {#if !isAuthPage}
    <div class="flex flex-1 overflow-hidden">
      <div class="flex flex-1 flex-col">
        {#if !isServerPage}
          <Header {name} {email} {avatarUrl} onlogout={handleLogout} />
        {/if}
        <div class="flex min-h-0 flex-1">
          <ServerNav {serverItems} {activeServer} />
          <main class="flex min-h-0 flex-1">
            {@render children()}
          </main>
        </div>
      </div>
    </div>
  {:else}
    <main class="flex flex-1">
      {@render children()}
    </main>
  {/if}
</div>
