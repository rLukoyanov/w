<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { user } from "$lib/api/users";
	import { pb } from "$lib/api/pocketbase";
	import { page } from "$app/stores";
	import { routes } from "$lib/consts";
	import { invalidateAll } from "$app/navigation";
	import { Settings, LogOut } from "@lucide/svelte";

	let { children } = $props();

	let name = $derived(user.current?.name || "Пользователь");
	let email = $derived(user.current?.email || "");
	let search = $state("");
	let currentRoute = $derived(routes.find((r) => r.path === $page.url.pathname));
	let avatarUrl = $derived(user.current?.avatar
		? pb.files.getUrl(user.current, user.current.avatar, { thumb: "32x32" })
		: null);
	let isAuthPage = $derived($page.url.pathname === "/login" || $page.url.pathname === "/register");
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

	$effect(() => {
		if (dropdownOpen) {
			document.addEventListener("click", onClickOutside);
			return () => document.removeEventListener("click", onClickOutside);
		}
	});

	async function handleLogout() {
		user.logout();
		dropdownOpen = false;
		await invalidateAll();
	}

	function handleSettings() {
		dropdownOpen = false;
	}
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex min-h-screen flex-col bg-bg">
	{#if !isAuthPage}
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
				<span class="text-sm font-medium text-text">{currentRoute.label}</span>
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
				<span class="absolute right-2 top-2 h-2 w-2 rounded-full bg-danger"></span>
			</button>

			<div class="relative ml-2" bind:this={dropdownRef}>
				<button
					onclick={toggleDropdown}
					class="flex items-center gap-2 rounded-lg p-1.5 transition hover:bg-bg-tertiary"
				>
					<div
						class="flex h-7 w-7 items-center justify-center overflow-hidden rounded-full bg-primary-light text-xs font-bold text-primary"
					>
						{#if avatarUrl}
							<img src={avatarUrl} alt={name} class="h-full w-full object-cover" />
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
									<img src={avatarUrl} alt={name} class="h-full w-full object-cover" />
								{:else}
									{name[0]}
								{/if}
							</div>
							<div class="min-w-0">
								<p class="truncate text-sm font-medium text-text">{name}</p>
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
		</nav>
	</header>
	{/if}

	<main class="flex flex-1">
		{@render children()}
	</main>
</div>
