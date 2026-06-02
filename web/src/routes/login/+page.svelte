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

<main class="min-h-screen flex items-center justify-center bg-base p-5">
	<div class="bg-surface border border-border rounded-lg p-8 w-full max-w-sm">
		<div class="text-center mb-8">
			<h1 class="text-2xl font-semibold text-text mb-2">Sign in</h1>
			<p class="text-sm text-subtext">Enter your credentials to continue</p>
		</div>

		<form
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="space-y-5"
		>
			<div>
				<label for="email" class="block text-xs font-medium text-subtext mb-1.5">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					placeholder="your@email.com"
					class="w-full px-3 py-2 bg-overlay border border-border rounded text-sm text-text placeholder-subtext/50 focus:outline-none focus:border-blue transition-colors"
				/>
			</div>

			<div>
				<label for="password" class="block text-xs font-medium text-subtext mb-1.5">
					Password
				</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					placeholder="••••••••"
					class="w-full px-3 py-2 bg-overlay border border-border rounded text-sm text-text placeholder-subtext/50 focus:outline-none focus:border-blue transition-colors"
				/>
			</div>

			{#if error}
				<div class="px-3 py-2 bg-red/10 border border-red/20 rounded text-sm text-red">
					{error}
				</div>
			{/if}

			<button
				type="submit"
				disabled={loading}
				class="w-full px-4 py-2 bg-blue text-white text-sm font-medium rounded hover:bg-blue/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
			>
				{loading ? 'Signing in...' : 'Sign in'}
			</button>

			<p class="text-center text-sm text-subtext">
				Don't have an account?
				<a href="/register" class="text-blue hover:underline">
					Sign up
				</a>
			</p>
		</form>
	</div>
</main>
