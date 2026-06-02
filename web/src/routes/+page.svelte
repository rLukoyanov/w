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

<main class="min-h-screen flex items-center justify-center bg-gradient-to-br from-base via-surface to-overlay p-5 relative overflow-hidden">
	<!-- Animated background elements -->
	<div class="absolute top-0 left-0 w-96 h-96 bg-blue/20 rounded-full blur-3xl animate-pulse"></div>
	<div class="absolute bottom-0 right-0 w-96 h-96 bg-pink/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
	
	{#if $auth.user}
		<div class="glass rounded-2xl p-10 text-center relative z-10 backdrop-blur-xl">
			<h1 class="text-5xl font-bold bg-gradient-to-r from-blue to-sky bg-clip-text text-transparent mb-8">Welcome, {$auth.user.username}!</h1>
			<div class="glass rounded-xl p-6 mb-6 border border-blue/20">
				<p class="text-subtext mb-2 text-lg">✉️ {$auth.user.email}</p>
				<p class="text-subtext/60 text-sm font-mono">ID: {$auth.user.id}</p>
			</div>
			<button
				onclick={handleLogout}
				class="px-8 py-3 bg-gradient-to-r from-red to-pink text-white font-semibold rounded-xl hover:shadow-lg hover:shadow-red/50 transition-all duration-300 hover:scale-105"
			>
				Logout
			</button>
		</div>
	{/if}
</main>

