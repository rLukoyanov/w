<script lang="ts">
	let name = $state('Пользователь');
	let email = $state('user@example.com');
</script>

<h1>Профиль</h1>

<form>
	<label>
		Имя
		<input type="text" name="name" bind:value={name} />
	</label>
	<label>
		Email
		<input type="email" name="email" bind:value={email} />
	</label>
	<button type="submit">Сохранить</button>
</form>

<nav>
	<a href="/dashboard">На дашборд</a>
</nav>
