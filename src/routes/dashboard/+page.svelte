<script lang="ts">
  import { user } from "$lib/api/users";

  let name = $derived(user.current?.name || "Пользователь");
</script>

<div class="flex-1 px-8 py-6">
    <h1 class="mb-1 text-2xl font-bold text-text">Дашборд</h1>
    <p class="mb-6 text-sm text-text-muted">Добро пожаловать, {name}!</p>

    <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
      <section class="rounded-xl border border-border bg-surface p-5">
        <h2 class="mb-2 text-sm font-medium text-text-muted">Статистика</h2>
        <p class="text-3xl font-bold text-text">—</p>
      </section>
      <section class="rounded-xl border border-border bg-surface p-5">
        <h2 class="mb-2 text-sm font-medium text-text-muted">Статистика</h2>
        <p class="text-3xl font-bold text-text">—</p>
      </section>
      <section class="rounded-xl border border-border bg-surface p-5">
        <h2 class="mb-2 text-sm font-medium text-text-muted">Статистика</h2>
        <p class="text-3xl font-bold text-text">—</p>
      </section>
    </div>

    <div class="mt-6">
      <a href="/profile" class="text-sm text-primary hover:text-primary-hover"
        >Профиль</a
      >
    </div>
</div>
