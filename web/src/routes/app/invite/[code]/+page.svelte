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
      <CheckCircle class="w-20 h-20 text-success mx-auto mb-6" />
      <h1 class="text-3xl font-bold mb-2">Joined successfully</h1>
      <p class="text-base-content/60 mb-8">
        You are now a member of <strong>{result.server_name}</strong>
      </p>
      <button onclick={goToServer} class="btn btn-primary btn-lg">
        Go to Server
      </button>
    </div>
  {:else if error}
    <div class="text-center max-w-md">
      <XCircle class="w-20 h-20 text-error mx-auto mb-6" />
      <h1 class="text-3xl font-bold mb-2">Cannot join</h1>
      <p class="text-base-content/60 mb-2">{error}</p>
      <p class="text-sm text-base-content/40 mb-8">
        The invite link may be invalid, expired, or you might already be a member.
      </p>
      <a href={ROUTES.SERVER.INDEX} class="btn btn-primary">
        Go to Servers
      </a>
    </div>
  {:else}
    <div class="text-center max-w-md">
      <Server class="w-20 h-20 text-base-content/20 mx-auto mb-6" />
      <h1 class="text-3xl font-bold mb-2">You're invited!</h1>
      <p class="text-base-content/60 mb-8">
        Join the server and start chatting with others.
      </p>
      <button
        onclick={handleJoin}
        class="btn btn-primary btn-lg"
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
