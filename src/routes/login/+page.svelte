<script lang="ts">
	let email = $state('');
	let password = $state('');
</script>

<h1>Вход</h1>

<form method="POST" action="/login">
	<label>
		Email
		<input type="email" name="email" bind:value={email} required />
	</label>
	<label>
		Пароль
		<input type="password" name="password" bind:value={password} required />
	</label>
	<button type="submit">Войти</button>
</form>

<p>
	Нет аккаунта? <a href="/register">Зарегистрироваться</a>
</p>
