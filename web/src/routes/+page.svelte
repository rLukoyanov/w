<script lang="ts">
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(async () => {
		// Check if user is already logged in
		try {
			const user = await api.getMe();
			auth.setUser(user);
		} catch (e) {
			// Not logged in, redirect to login page
			goto('/login');
		}
	});

	function handleLogout() {
		api.clearToken();
		auth.logout();
		goto('/login');
	}
</script>

<main class="min-h-screen flex items-center justify-center bg-base p-5">
	{#if $auth.user}
		<div class="bg-surface border border-border rounded-lg p-8 max-w-md w-full">
			<h1 class="text-2xl font-semibold text-text mb-6">Welcome back</h1>
			<div class="bg-overlay border border-border rounded p-4 mb-6 space-y-2">
				<div class="flex items-center justify-between">
					<span class="text-sm text-subtext">Username</span>
					<span class="text-sm text-text font-medium">{$auth.user.username}</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-subtext">Email</span>
					<span class="text-sm text-text">{$auth.user.email}</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm text-subtext">ID</span>
					<span class="text-xs text-subtext font-mono">{$auth.user.id}</span>
				</div>
			</div>
			<button
				onclick={handleLogout}
				class="w-full px-4 py-2 bg-overlay border border-border text-text text-sm font-medium rounded hover:bg-border/50 transition-colors"
			>
				Sign out
			</button>
		</div>
	{/if}
</main>

