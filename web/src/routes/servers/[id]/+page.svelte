<script lang="ts">
	import { api, type Server, type Channel } from '$lib/api';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let { params } = $props();
	
	let server = $state<Server | null>(null);
	let channels = $state<Channel[]>([]);
	let loading = $state(true);
	let error = $state('');
	let showCreateChannel = $state(false);
	let newChannelName = $state('');
	let createError = $state('');

	onMount(async () => {
		try {
			server = await api.getServer(params.id);
			// Load channels when we add the API endpoint
			// For now, show empty state
			loading = false;
		} catch (err: any) {
			error = err.message;
			loading = false;
		}
	});

	async function handleCreateChannel() {
		if (!newChannelName.trim()) {
			createError = 'Channel name is required';
			return;
		}

		try {
			const channel = await api.createChannel(params.id, newChannelName);
			channels.push(channel);
			newChannelName = '';
			showCreateChannel = false;
			createError = '';
			
			// Navigate to the new channel
			goto(`/servers/${params.id}/channels/${channel.id}`);
		} catch (err: any) {
			createError = err.message;
		}
	}
</script>

<main class="min-h-screen bg-base text-text">
	<!-- Top bar -->
	<div class="border-b border-border bg-surface">
		<div class="max-w-7xl mx-auto px-4 py-4 flex items-center gap-4">
			<a href="/" class="text-sm text-subtext hover:text-text">← Back to servers</a>
			{#if server}
				<h1 class="text-xl font-semibold">{server.name}</h1>
			{/if}
		</div>
	</div>

	<!-- Content -->
	<div class="max-w-7xl mx-auto px-4 py-8">
		{#if loading}
			<div class="text-center py-12">
				<p class="text-subtext">Loading...</p>
			</div>
		{:else if error}
			<div class="bg-surface border border-red rounded-lg p-6">
				<p class="text-red">{error}</p>
			</div>
		{:else if server}
			<div class="flex items-center justify-between mb-6">
				<h2 class="text-2xl font-semibold">Channels</h2>
				<button
					onclick={() => showCreateChannel = !showCreateChannel}
					class="px-4 py-2 bg-blue text-text text-sm font-medium rounded hover:bg-blue/80 transition-colors"
				>
					Create Channel
				</button>
			</div>

			{#if showCreateChannel}
				<div class="bg-surface border border-border rounded-lg p-6 mb-6">
					<h3 class="text-lg font-medium mb-4">Create a new channel</h3>
					<input
						type="text"
						bind:value={newChannelName}
						placeholder="Channel name"
						class="w-full px-4 py-2 bg-overlay border border-border rounded text-text placeholder:text-subtext focus:outline-none focus:ring-2 focus:ring-blue mb-3"
					/>
					{#if createError}
						<p class="text-sm text-red mb-3">{createError}</p>
					{/if}
					<div class="flex gap-2">
						<button
							onclick={handleCreateChannel}
							class="px-4 py-2 bg-blue text-text text-sm font-medium rounded hover:bg-blue/80 transition-colors"
						>
							Create
						</button>
						<button
							onclick={() => { showCreateChannel = false; createError = ''; newChannelName = ''; }}
							class="px-4 py-2 bg-overlay border border-border text-text text-sm font-medium rounded hover:bg-border/50 transition-colors"
						>
							Cancel
						</button>
					</div>
				</div>
			{/if}

			{#if channels.length === 0}
				<div class="bg-surface border border-border rounded-lg p-12 text-center">
					<p class="text-subtext mb-4">No channels yet.</p>
					<p class="text-sm text-subtext">Create your first channel to start chatting!</p>
				</div>
			{:else}
				<div class="space-y-2">
					{#each channels as channel}
						<a
							href="/servers/{params.id}/channels/{channel.id}"
							class="block bg-surface border border-border rounded-lg p-4 hover:bg-overlay transition-colors"
						>
							<h3 class="text-lg font-medium"># {channel.name}</h3>
							<p class="text-xs text-subtext mt-1">{channel.type} channel</p>
						</a>
					{/each}
				</div>
			{/if}
		{/if}
	</div>
</main>
