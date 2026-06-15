<script lang="ts">
  import { fly } from "svelte/transition";
  import type { Server } from "$lib/api";
  import { serversClient } from "$lib/api";
  import { Plus } from "lucide-svelte";
  import { notify } from "$lib/stores/notifications";

  interface Props {
    show: boolean;
    onClose: () => void;
    onCreated: (server: Server) => void;
  }

  let { show, onClose, onCreated }: Props = $props();

  let name = $state("");
  let creating = $state(false);
  let inputEl = $state<HTMLInputElement | null>(null);

  $effect(() => {
    if (show) {
      name = "";
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
      notify.error("Server name is required");
      return;
    }

    creating = true;

    try {
      const server = await serversClient.create(trimmed);
      onCreated(server);
      onClose();
    } catch (err: any) {
      notify.error(err.message ?? "Failed to create server");
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
    <div class="modal-box" transition:fly={{ duration: 150, y: 24 }}
      style="background: oklch(0.105 0.006 285); border: 1px solid oklch(0.2 0.01 285);">
      <form method="dialog">
        <button
          class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
          onclick={onClose}
          aria-label="Close"
          style="color: oklch(0.5 0.01 285);"
        >
          ✕
        </button>
      </form>

      <h3 class="text-lg font-bold mb-1 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">Create Server</h3>
      <p class="text-sm mb-6" style="color: oklch(0.5 0.01 285);">
        Give your new server a name — you can customise it later.
      </p>

      <div class="space-y-4">
        <input
          bind:this={inputEl}
          type="text"
          bind:value={name}
          placeholder="Server name"
          class="input w-full"
          style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
          disabled={creating}
          onkeydown={handleKeydown}
        />
      </div>

      <div class="modal-action">
        <button
          onclick={handleSubmit}
          class="btn"
          style="background: oklch(0.58 0.2 285); color: white; {creating || !name.trim() ? 'opacity: 0.5;' : ''}"
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
        <button onclick={onClose} class="btn btn-ghost" disabled={creating}
          style="color: oklch(0.6 0.01 285);">
          Cancel
        </button>
      </div>
    </div>

    <form method="dialog" class="modal-backdrop">
      <button onclick={onClose}>close</button>
    </form>
  </dialog>
{/key}
