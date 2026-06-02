<script lang="ts">
	import { api, type Server } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { wsClient } from '$lib/websocket';

	let servers = $state<Server[]>([]);
	let loading = $state(true);
	let showCreateServer = $state(false);
	let newServerName = $state('');
	let createError = $state('');

	onMount(async () => {
		// Check if user is already logged in
		try {
			const user = await api.getMe();
			auth.setUser(user);
			
			// Connect to WebSocket
			const token = api.getToken();
			if (token) {
				wsClient.connect(token);
			}
			
			// Load servers (when we add the API endpoint)
			// For now, show empty state
			loading = false;
		} catch (e) {
			// Not logged in, redirect to login page
			goto('/login');
		}
	});

	async function handleCreateServer() {
		if (!newServerName.trim()) {
			createError = 'Server name is required';
			return;
		}

		try {
			const server = await api.createServer(newServerName);
			servers.push(server);
			newServerName = '';
			showCreateServer = false;
			createError = '';
			
			// Navigate to the new server
			goto(`/servers/${server.id}`);
		} catch (err: any) {
			createError = err.message;
		}
	}

	function handleLogout() {
		wsClient.disconnect();
		api.clearToken();
		auth.logout();
		goto('/login');
	}
</script>

<main class="min-h-screen bg-base text-text">
	{#if $auth.user}
		<!-- Top bar -->
		<div class="border-b border-border bg-surface">
			<div class="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
				<div class="flex items-center gap-3">
					<h1 class="text-xl font-semibold">W</h1>
					<div class="flex items-center gap-2">
						<div class="w-2 h-2 rounded-full" class:bg-green-500={$wsClient.connected} class:bg-red-500={!$wsClient.connected}></div>
						<span class="text-xs text-subtext">{$wsClient.connected ? 'Connected' : 'Disconnected'}</span>
					</div>
				</div>
				<div class="flex items-center gap-3">
					<span class="text-sm text-subtext">{$auth.user.username}</span>
					<button
						onclick={handleLogout}
						class="px-3 py-1.5 bg-overlay border border-border text-text text-xs font-medium rounded hover:bg-border/50 transition-colors"
					>
						Sign out
					</button>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="max-w-7xl mx-auto px-4 py-8">
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-2xl font-semibold">Your Servers</h2>
				<button
					onclick={() => showCreateServer = !showCreateServer}
					class="px-4 py-2 bg-blue text-text text-sm font-medium rounded hover:bg-blue/80 transition-colors"
				>
					Create Server
				</button>
			</div>

			{#if showCreateServer}
				<div class="bg-surface border border-border rounded-lg p-6 mb-6">
					<h3 class="text-lg font-medium mb-4">Create a new server</h3>
					<input
						type="text"
						bind:value={newServerName}
						placeholder="Server name"
						class="w-full px-4 py-2 bg-overlay border border-border rounded text-text placeholder:text-subtext focus:outline-none focus:ring-2 focus:ring-blue mb-3"
					/>
					{#if createError}
						<p class="text-sm text-red mb-3">{createError}</p>
					{/if}
					<div class="flex gap-2">
						<button
							onclick={handleCreateServer}
							class="px-4 py-2 bg-blue text-text text-sm font-medium rounded hover:bg-blue/80 transition-colors"
						>
							Create
						</button>
						<button
							onclick={() => { showCreateServer = false; createError = ''; newServerName = ''; }}
							class="px-4 py-2 bg-overlay border border-border text-text text-sm font-medium rounded hover:bg-border/50 transition-colors"
						>
							Cancel
						</button>
					</div>
				</div>
			{/if}

			{#if loading}
				<div class="text-center py-12">
					<p class="text-subtext">Loading...</p>
				</div>
			{:else if servers.length === 0}
				<div class="bg-surface border border-border rounded-lg p-12 text-center">
					<p class="text-subtext mb-4">You don't have any servers yet.</p>
					<p class="text-sm text-subtext">Create your first server to get started!</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each servers as server}
						<a
							href="/servers/{server.id}"
							class="bg-surface border border-border rounded-lg p-6 hover:bg-overlay transition-colors"
						>
							<h3 class="text-lg font-medium mb-2">{server.name}</h3>
							<p class="text-xs text-subtext">Created {new Date(server.created_at).toLocaleDateString()}</p>
						</a>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</main>

