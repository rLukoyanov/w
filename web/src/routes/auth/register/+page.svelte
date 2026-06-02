<script lang="ts">
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";

  let email = $state<string>("");
  let password = $state<string>("");
  let username = $state<string>("");
  let error = $state<string>("");
  let loading = $state<boolean>(false);

  async function handleSubmit() {
    error = "";
    loading = true;

    try {
      const response = await userClient.register({
        username: username,
        email: email,
        password: password,
      });
      auth.setUser(response.user);
      goto("/");
    } catch (e) {
      error = e instanceof Error ? e.message : "An error occurred";
    } finally {
      loading = false;
    }
  }
</script>

<div class="mb-8">
  <h1 class="text-4xl font-bold mb-2">Sign up</h1>
  <p class="text-base-content/70">Create your account to get started</p>
</div>

<form
  onsubmit={(e) => {
    e.preventDefault();
    handleSubmit();
  }}
  class="space-y-4"
>
  <div class="form-control">
    <label for="username" class="label">
      <span class="label-text">Username</span>
    </label>
    <input
      id="username"
      type="text"
      bind:value={username}
      required
      placeholder="username"
      class="input input-bordered w-full"
    />
  </div>

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
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="stroke-current shrink-0 h-6 w-6"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span>{error}</span>
    </div>
  {/if}

  <div class="form-control mt-6">
    <button type="submit" disabled={loading} class="btn btn-primary w-full">
      {#if loading}
        <span class="loading loading-spinner"></span>
        Creating account...
      {:else}
        Create account
      {/if}
    </button>
  </div>

  <p class="text-center text-sm mt-4">
    Already have an account?
    <a href="/auth/login" class="link link-primary"> Sign in </a>
  </p>
</form>
