<script lang="ts">
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import { onMount } from 'svelte';

	let mode: 'login' | 'register' = 'login';
	let email = '';
	let password = '';
	let username = '';
	let error = '';
	let loading = false;

	onMount(async () => {
		// Check if user is already logged in
		try {
			const user = await api.getMe();
			auth.setUser(user);
		} catch (e) {
			// Not logged in, that's fine
		}
	});

	async function handleSubmit() {
		error = '';
		loading = true;

		try {
			if (mode === 'login') {
				const response = await api.login({ email, password });
				auth.setUser(response.user);
			} else {
				const response = await api.register({ username, email, password });
				auth.setUser(response.user);
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'An error occurred';
		} finally {
			loading = false;
		}
	}

	function handleLogout() {
		api.clearToken();
		auth.logout();
		email = '';
		password = '';
		username = '';
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
	{:else}
		<div class="bg-surface rounded-xl p-10 w-full max-w-md shadow-2xl">
			<h1 class="text-4xl font-bold text-blue text-center mb-8">W</h1>

			<div class="flex gap-2 mb-8">
				<button
					class="flex-1 py-3 px-4 rounded-lg font-medium transition-all {mode === 'login'
						? 'bg-blue text-base'
						: 'bg-overlay text-text hover:bg-overlay/80'}"
					onclick={() => {
						mode = 'login';
						error = '';
					}}
				>
					Login
				</button>
				<button
					class="flex-1 py-3 px-4 rounded-lg font-medium transition-all {mode === 'register'
						? 'bg-blue text-base'
						: 'bg-overlay text-text hover:bg-overlay/80'}"
					onclick={() => {
						mode = 'register';
						error = '';
					}}
				>
					Register
				</button>
			</div>

			<form on:submit|preventDefault={handleSubmit} class="space-y-5">
				{#if mode === 'register'}
					<div>
						<label for="username" class="block text-sm font-medium text-subtext mb-2">
							Username
						</label>
						<input
							id="username"
							type="text"
							bind:value={username}
							required
							placeholder="Enter username"
							class="w-full px-4 py-3 bg-overlay border-2 border-overlay rounded-lg text-text placeholder-subtext/50 focus:outline-none focus:border-blue transition-colors"
						/>
					</div>
				{/if}

				<div>
					<label for="email" class="block text-sm font-medium text-subtext mb-2">Email</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						required
						placeholder="Enter email"
						class="w-full px-4 py-3 bg-overlay border-2 border-overlay rounded-lg text-text placeholder-subtext/50 focus:outline-none focus:border-blue transition-colors"
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
						class="w-full px-4 py-3 bg-overlay border-2 border-overlay rounded-lg text-text placeholder-subtext/50 focus:outline-none focus:border-blue transition-colors"
					/>
				</div>

				{#if error}
					<div class="px-4 py-3 bg-red text-base rounded-lg font-medium">
						{error}
					</div>
				{/if}

				<button
					type="submit"
					disabled={loading}
					class="w-full py-3 bg-blue text-base font-semibold rounded-lg hover:bg-sky transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
				>
					{loading ? 'Loading...' : mode === 'login' ? 'Login' : 'Register'}
				</button>
			</form>
		</div>
	{/if}
</main>

