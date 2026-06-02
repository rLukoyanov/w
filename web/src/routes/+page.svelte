<script lang="ts">
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
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
			<h2 class="text-2xl font-semibold text-text text-center mb-8">Login</h2>

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
					{loading ? 'Loading...' : 'Login'}
				</button>

				<p class="text-center text-subtext">
					Don't have an account?
					<a href="/register" class="text-blue hover:text-sky font-medium transition-colors">
						Register
					</a>
				</p>
			</form>
		</div>
	{/if}
</main>

