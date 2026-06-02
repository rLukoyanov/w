<script lang="ts">
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	onMount(async () => {
		// Check if user is already logged in
		try {
			const user = await api.getMe();
			auth.setUser(user);
			goto('/');
		} catch (e) {
			// Not logged in, that's fine
		}
	});

	async function handleSubmit() {
		error = '';
		loading = true;

		try {
			const response = await api.login({ email, password });
			auth.setUser(response.user);
			goto('/');
		} catch (e) {
			error = e instanceof Error ? e.message : 'An error occurred';
		} finally {
			loading = false;
		}
	}
</script>

<main class="min-h-screen flex items-center justify-center bg-gradient-to-br from-base via-surface to-overlay p-5 relative overflow-hidden">
	<!-- Animated background elements -->
	<div class="absolute top-0 left-0 w-96 h-96 bg-blue/20 rounded-full blur-3xl animate-pulse"></div>
	<div class="absolute bottom-0 right-0 w-96 h-96 bg-sky/20 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
	
	<div class="glass rounded-2xl p-10 w-full max-w-md relative z-10 backdrop-blur-xl">
		<h1 class="text-5xl font-bold bg-gradient-to-r from-blue to-sky bg-clip-text text-transparent text-center mb-4">W</h1>
		<h2 class="text-xl font-medium text-subtext text-center mb-8">Welcome back</h2>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-5"
		>
			<div>
				<label for="email" class="block text-sm font-medium text-subtext mb-2">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					placeholder="Enter email"
					class="glass-input w-full px-4 py-3 rounded-xl text-text placeholder-subtext/50 focus:outline-none transition-all duration-300"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-subtext mb-2">
					Password
				</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					placeholder="Enter password"
					class="glass-input w-full px-4 py-3 rounded-xl text-text placeholder-subtext/50 focus:outline-none transition-all duration-300"
				/>
			</div>

			{#if error}
				<div class="px-4 py-3 bg-gradient-to-r from-red to-pink text-white rounded-xl font-medium shadow-lg">
					{error}
				</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="w-full py-3 bg-gradient-to-r from-blue to-sky text-white font-semibold rounded-xl hover:shadow-lg hover:shadow-blue/50 transition-all duration-300 disabled:opacity-60 disabled:cursor-not-allowed hover:scale-105"
			>
				{loading ? '⏳ Loading...' : '🚀 Login'}
			</button>

			<p class="text-center text-subtext">
				Don't have an account?
				<a href="/register" class="text-blue hover:text-sky font-semibold transition-all duration-300 hover:underline">
					Register →
				</a>
			</p>
		</form>
	</div>
</main>
