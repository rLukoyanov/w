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
		<div class="bg-surface rounded-xl p-10 text-center shadow-2xl">
			<h1 class="text-4xl font-bold text-blue mb-8">Welcome, {$auth.user.username}!</h1>
			<p class="text-subtext mb-2">Email: {$auth.user.email}</p>
			<p class="text-subtext mb-6">ID: {$auth.user.id}</p>
			<button
				onclick={handleLogout}
				class="px-6 py-3 bg-red text-base font-semibold rounded-lg hover:bg-pink transition-colors"
			>
				Logout
			</button>
		</div>
	{/if}
</main>

