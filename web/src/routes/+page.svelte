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

<main>
	{#if $auth.user}
		<div class="welcome">
			<h1>Welcome, {$auth.user.username}!</h1>
			<p>Email: {$auth.user.email}</p>
			<p>ID: {$auth.user.id}</p>
			<button on:click={handleLogout}>Logout</button>
		</div>
	{:else}
		<div class="auth-container">
			<h1>W</h1>
			<div class="tabs">
				<button
					class:active={mode === 'login'}
					on:click={() => {
						mode = 'login';
						error = '';
					}}
				>
					Login
				</button>
				<button
					class:active={mode === 'register'}
					on:click={() => {
						mode = 'register';
						error = '';
					}}
				>
					Register
				</button>
			</div>

			<form on:submit|preventDefault={handleSubmit}>
				{#if mode === 'register'}
					<div class="form-group">
						<label for="username">Username</label>
						<input
							id="username"
							type="text"
							bind:value={username}
							required
							placeholder="Enter username"
						/>
					</div>
				{/if}

				<div class="form-group">
					<label for="email">Email</label>
					<input id="email" type="email" bind:value={email} required placeholder="Enter email" />
				</div>

				<div class="form-group">
					<label for="password">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						placeholder="Enter password"
					/>
				</div>

				{#if error}
					<div class="error">{error}</div>
				{/if}

				<button type="submit" disabled={loading}>
					{loading ? 'Loading...' : mode === 'login' ? 'Login' : 'Register'}
				</button>
			</form>
		</div>
	{/if}
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell,
			sans-serif;
		background: #1e1e2e;
		color: #cdd6f4;
	}

	main {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		padding: 20px;
	}

	.auth-container {
		background: #181825;
		border-radius: 12px;
		padding: 40px;
		width: 100%;
		max-width: 400px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	}

	.welcome {
		background: #181825;
		border-radius: 12px;
		padding: 40px;
		text-align: center;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	}

	h1 {
		margin: 0 0 30px;
		text-align: center;
		font-size: 32px;
		color: #89b4fa;
	}

	.tabs {
		display: flex;
		gap: 10px;
		margin-bottom: 30px;
	}

	.tabs button {
		flex: 1;
		padding: 12px;
		background: #313244;
		border: none;
		border-radius: 8px;
		color: #cdd6f4;
		cursor: pointer;
		font-size: 14px;
		font-weight: 500;
		transition: all 0.2s;
	}

	.tabs button:hover {
		background: #45475a;
	}

	.tabs button.active {
		background: #89b4fa;
		color: #1e1e2e;
	}

	.form-group {
		margin-bottom: 20px;
	}

	label {
		display: block;
		margin-bottom: 8px;
		font-size: 14px;
		font-weight: 500;
		color: #a6adc8;
	}

	input {
		width: 100%;
		padding: 12px;
		background: #313244;
		border: 2px solid #45475a;
		border-radius: 8px;
		color: #cdd6f4;
		font-size: 14px;
		box-sizing: border-box;
		transition: border-color 0.2s;
	}

	input:focus {
		outline: none;
		border-color: #89b4fa;
	}

	input::placeholder {
		color: #6c7086;
	}

	button[type='submit'] {
		width: 100%;
		padding: 14px;
		background: #89b4fa;
		border: none;
		border-radius: 8px;
		color: #1e1e2e;
		font-size: 16px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	button[type='submit']:hover:not(:disabled) {
		background: #74c7ec;
	}

	button[type='submit']:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.error {
		padding: 12px;
		background: #f38ba8;
		color: #1e1e2e;
		border-radius: 8px;
		margin-bottom: 20px;
		font-size: 14px;
		font-weight: 500;
	}

	.welcome button {
		margin-top: 20px;
		padding: 12px 24px;
		background: #f38ba8;
		border: none;
		border-radius: 8px;
		color: #1e1e2e;
		font-size: 14px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.welcome button:hover {
		background: #eba0ac;
	}

	.welcome p {
		margin: 10px 0;
		color: #a6adc8;
	}
</style>
