<script lang="ts">
  import { invitesClient, type JoinResponse } from "$lib/api";
  import { ROUTES } from "$lib/routes";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { Server, CheckCircle, XCircle } from "lucide-svelte";

  let { params } = $props();

  let checking = $state(true);
  let joining = $state(false);
  let done = $state(false);
  let error = $state<string | null>(null);
  let result = $state<JoinResponse | null>(null);

  async function handleJoin() {
    joining = true;
    error = null;

    try {
      result = await invitesClient.join(params.code);
      done = true;
    } catch (err: any) {
      error = err.message ?? "Failed to join server";
    } finally {
      joining = false;
    }
  }

  function goToServer() {
    if (result) {
      goto(ROUTES.SERVER.DETAIL(result.server_id));
    }
  }
</script>

<div class="flex items-center justify-center h-full p-8">
  {#if done && result}
    <div class="text-center max-w-md">
      <CheckCircle class="w-20 h-20 mx-auto mb-6"
        style="color: oklch(0.72 0.18 170);" />
      <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">Joined successfully</h1>
      <p class="mb-8" style="color: oklch(0.5 0.01 285);">
        You are now a member of <strong style="color: oklch(0.88 0.005 285);">{result.server_name}</strong>
      </p>
      <button onclick={goToServer} class="btn btn-lg"
        style="background: oklch(0.58 0.2 285); color: white;">
        Go to Server
      </button>
    </div>
  {:else if error}
    <div class="text-center max-w-md">
      <XCircle class="w-20 h-20 mx-auto mb-6"
        style="color: oklch(0.65 0.18 25);" />
      <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">Cannot join</h1>
      <p class="mb-2" style="color: oklch(0.6 0.01 285);">{error}</p>
      <p class="text-sm mb-8" style="color: oklch(0.4 0.01 285);">
        The invite link may be invalid, expired, or you might already be a member.
      </p>
      <a href={ROUTES.SERVER.INDEX} class="btn"
        style="background: oklch(0.58 0.2 285); color: white;">
        Go to Servers
      </a>
    </div>
  {:else}
    <div class="text-center max-w-md">
      <Server class="w-20 h-20 mx-auto mb-6"
        style="color: oklch(0.2 0.01 285);" />
      <h1 class="text-3xl font-bold mb-2 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">You're invited!</h1>
      <p class="mb-8" style="color: oklch(0.5 0.01 285);">
        Join the server and start chatting with others.
      </p>
      <button
        onclick={handleJoin}
        class="btn btn-lg"
        style="background: oklch(0.58 0.2 285); color: white; {joining ? 'opacity: 0.5;' : ''}"
        disabled={joining}
      >
        {#if joining}
          <span class="loading loading-spinner"></span>
          Joining...
        {:else}
          Join Server
        {/if}
      </button>
    </div>
  {/if}
</div>
