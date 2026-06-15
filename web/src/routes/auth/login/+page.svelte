<script lang="ts">
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { ROUTES } from "$lib/routes";
  import { notify } from "$lib/stores/notifications";

  let email = $state<string>("");
  let password = $state<string>("");
  let loading = $state<boolean>(false);

  async function handleSubmit() {
    loading = true;

    try {
      const response = await userClient.login({ email, password });
      auth.setUser(response.user);
      goto(ROUTES.SERVER.INDEX);
    } catch (e) {
      notify.error(e instanceof Error ? e.message : "An error occurred");
    } finally {
      loading = false;
    }
  }
</script>

<div class="mb-8">
  <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
    style="color: oklch(0.92 0.004 285);">Sign in</h1>
  <p class="text-sm" style="color: oklch(0.5 0.01 285);">Enter your credentials to continue</p>
</div>

<form
  onsubmit={(e) => {
    e.preventDefault();
    handleSubmit();
  }}
  class="space-y-4"
>
  <div class="form-control">
    <label for="email" class="label">
      <span class="label-text text-sm" style="color: oklch(0.7 0.01 285);">Email</span>
    </label>
    <input
      id="email"
      type="email"
      bind:value={email}
      required
      placeholder="your@email.com"
      class="input w-full"
      style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
    />
  </div>

  <div class="form-control">
    <label for="password" class="label">
      <span class="label-text text-sm" style="color: oklch(0.7 0.01 285);">Password</span>
    </label>
    <input
      id="password"
      type="password"
      bind:value={password}
      required
      placeholder="••••••••"
      class="input w-full"
      style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
    />
  </div>

  <div class="form-control mt-6">
    <button type="submit" disabled={loading} class="btn w-full"
      style="background: oklch(0.58 0.2 285); color: white; {loading ? 'opacity: 0.5;' : ''}">
      {#if loading}
        <span class="loading loading-spinner"></span>
        Signing in...
      {:else}
        Sign in
      {/if}
    </button>
  </div>

  <p class="text-center text-sm mt-4" style="color: oklch(0.5 0.01 285);">
    Don't have an account?
    <a href="/auth/register" style="color: oklch(0.62 0.18 285);"> Sign up </a>
  </p>
</form>
