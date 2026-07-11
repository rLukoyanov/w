<script lang="ts">
	import bg from '$lib/assets/bg.JPG';

	let email = $state('');
	let password = $state('');
</script>

<div class="relative flex flex-1 items-center justify-center overflow-hidden px-4">
	<div
		class="absolute inset-0 bg-cover bg-center bg-no-repeat"
		style="background-image: url({bg})"
	></div>
	<div class="absolute inset-0 bg-bg/70 backdrop-blur-sm"></div>

	<div class="absolute left-6 top-6 flex items-center gap-2">
		<div class="flex h-9 w-9 items-center justify-center rounded-lg bg-primary text-lg font-bold text-text-inverse shadow-lg">
			В
		</div>
		<span class="text-xl font-bold text-text drop-shadow-lg">Вече</span>
	</div>

	<div class="relative flex w-full max-w-3xl overflow-hidden rounded-xl border border-border bg-surface shadow-lg">
		<div class="flex-1 p-10">
			<div class="mb-8">
				<h1 class="text-2xl font-bold text-text">Вход</h1>
				<p class="mt-1 text-sm text-text-muted">Войдите в свой аккаунт</p>
			</div>

			<form method="POST" action="/login" class="space-y-5">
				<div>
					<label for="email" class="mb-1.5 block text-sm font-medium text-text-secondary">Email</label>
					<input
						type="email"
						id="email"
						name="email"
						bind:value={email}
						required
						class="w-full rounded-lg border border-border bg-bg-secondary px-4 py-3 text-text outline-none transition focus:border-primary focus:ring-2 focus:ring-primary-light"
						placeholder="your@email.com"
					/>
				</div>
				<div>
					<label for="password" class="mb-1.5 block text-sm font-medium text-text-secondary">Пароль</label>
					<input
						type="password"
						id="password"
						name="password"
						bind:value={password}
						required
						class="w-full rounded-lg border border-border bg-bg-secondary px-4 py-3 text-text outline-none transition focus:border-primary focus:ring-2 focus:ring-primary-light"
						placeholder="••••••••"
					/>
				</div>
				<button
					type="submit"
					class="w-full rounded-lg bg-primary px-4 py-3 font-medium text-text-inverse transition hover:bg-primary-hover"
				>
					Войти
				</button>
			</form>

			<p class="mt-6 text-center text-sm text-text-muted">
				Нет аккаунта?
				<a href="/register" class="font-medium text-primary hover:text-primary-hover">Зарегистрироваться</a>
			</p>
		</div>

		<div class="hidden w-64 flex-col items-center justify-center gap-3 border-l border-border bg-bg-secondary p-6 md:flex">
			<div class="flex h-44 w-44 items-center justify-center rounded-lg border-2 border-dashed border-border-hover bg-surface">
				<span class="text-center text-xs text-text-muted">QR-код</span>
			</div>
			<p class="text-center text-xs text-text-muted">
				Отсканируйте QR-код<br />для быстрого входа
			</p>
		</div>
	</div>
</div>
