<script lang="ts">
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
</script>

<h1>Регистрация</h1>

<form method="POST" action="/register">
	<label>
		Email
		<input type="email" name="email" bind:value={email} required />
	</label>
	<label>
		Пароль
		<input type="password" name="password" bind:value={password} required />
	</label>
	<label>
		Подтвердите пароль
		<input type="password" name="confirmPassword" bind:value={confirmPassword} required />
	</label>
	<button type="submit">Зарегистрироваться</button>
</form>

<p>
	Уже есть аккаунт? <a href="/login">Войти</a>
</p>
