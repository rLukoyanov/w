<script lang="ts">
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { XCircle } from "lucide-svelte";
  import { ROUTES } from "$lib/routes";

  let email = $state<string>("");
  let password = $state<string>("");
  let error = $state<string>("");
  let loading = $state<boolean>(false);

  async function handleSubmit() {
    error = "";
    loading = true;

    try {
      const response = await userClient.login({ email, password });
      auth.setUser(response.user);
      goto(ROUTES.SERVER.INDEX);
    } catch (e) {
      error = e instanceof Error ? e.message : "An error occurred";
    } finally {
      loading = false;
    }
  }
</script>

<div class="mb-8">
  <h1 class="text-4xl font-bold mb-2">Sign in</h1>
  <p class="text-base-content/70">Enter your credentials to continue</p>
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
      <span class="label-text">Email</span>
    </label>
    <input
      id="email"
      type="email"
      bind:value={email}
      required
      placeholder="your@email.com"
      class="input input-bordered w-full"
    />
  </div>

  <div class="form-control">
    <label for="password" class="label">
      <span class="label-text">Password</span>
    </label>
    <input
      id="password"
      type="password"
      bind:value={password}
      required
      placeholder="••••••••"
      class="input input-bordered w-full"
    />
  </div>

  {#if error}
    <div role="alert" class="alert alert-error">
      <XCircle class="w-6 h-6" />
      <span>{error}</span>
    </div>
  {/if}

  <div class="form-control mt-6">
    <button type="submit" disabled={loading} class="btn btn-primary w-full">
      {#if loading}
        <span class="loading loading-spinner"></span>
        Signing in...
      {:else}
        Sign in
      {/if}
    </button>
  </div>

  <p class="text-center text-sm mt-4">
    Don't have an account?
    <a href="/auth/register" class="link link-primary"> Sign up </a>
  </p>
</form>
