<script lang="ts">
  import { fly } from "svelte/transition";
  import type { Server } from "$lib/api";
  import { serversClient } from "$lib/api";
  import { XCircle, Plus } from "lucide-svelte";

  interface Props {
    show: boolean;
    onClose: () => void;
    onCreated: (server: Server) => void;
  }

  let { show, onClose, onCreated }: Props = $props();

  let name = $state("");
  let error = $state<string | null>(null);
  let creating = $state(false);
  let inputEl = $state<HTMLInputElement | null>(null);

  $effect(() => {
    if (show) {
      name = "";
      error = null;
      creating = false;
    }
  });

  $effect(() => {
    if (show && inputEl) {
      inputEl.focus();
    }
  });

  async function handleSubmit() {
    const trimmed = name.trim();
    if (!trimmed) {
      error = "Server name is required";
      return;
    }

    creating = true;
    error = null;

    try {
      const server = await serversClient.create(trimmed);
      onCreated(server);
      onClose();
    } catch (err: any) {
      error = err.message ?? "Failed to create server";
    } finally {
      creating = false;
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && !creating) {
      handleSubmit();
    } else if (e.key === "Escape") {
      onClose();
    }
  }
</script>

{#key show}
  <dialog class="modal {show ? 'modal-open' : ''}">
    <div class="modal-box" transition:fly={{ duration: 150, y: 24 }}>
      <form method="dialog">
        <button
          class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          onclick={onClose}
          aria-label="Close"
        >
          ✕
        </button>
      </form>

      <h3 class="text-lg font-bold mb-1">Create Server</h3>
      <p class="text-sm text-base-content/60 mb-6">
        Give your new server a name — you can customise it later.
      </p>

      <div class="space-y-4">
        <input
          bind:this={inputEl}
          type="text"
          bind:value={name}
          placeholder="Server name"
          class="input input-bordered w-full"
          disabled={creating}
          onkeydown={handleKeydown}
        />

        {#if error}
          <div role="alert" class="alert alert-error py-2">
            <XCircle class="w-5 h-5 shrink-0" />
            <span class="text-sm">{error}</span>
          </div>
        {/if}
      </div>

      <div class="modal-action">
        <button
          onclick={handleSubmit}
          class="btn btn-primary"
          disabled={creating || !name.trim()}
        >
          {#if creating}
            <span class="loading loading-spinner loading-xs"></span>
            Creating...
          {:else}
            <Plus class="w-4 h-4" />
            Create Server
          {/if}
        </button>
        <button onclick={onClose} class="btn btn-ghost" disabled={creating}>
          Cancel
        </button>
      </div>
    </div>

    <form method="dialog" class="modal-backdrop">
      <button onclick={onClose}>close</button>
    </form>
  </dialog>
{/key}
