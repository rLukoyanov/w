<script lang="ts">
  import { fly } from "svelte/transition";
  import { invitesClient } from "$lib/api";
  import { notify } from "$lib/stores/notifications";
  import { Link, Plus } from "lucide-svelte";

  interface Props {
    show: boolean;
    serverId: string;
    onClose: () => void;
    onCreated: () => void;
  }

  let { show, serverId, onClose, onCreated }: Props = $props();

  let maxUses = $state<string>("");
  let expiresHours = $state<string>("");
  let creating = $state(false);
  let createdLink = $state<string | null>(null);

  $effect(() => {
    if (show) {
      maxUses = "";
      expiresHours = "";
      creating = false;
      createdLink = null;
    }
  });

  async function handleSubmit() {
    creating = true;

    try {
      const invite = await invitesClient.create(
        serverId,
        maxUses ? parseInt(maxUses) : null,
        expiresHours ? parseInt(expiresHours) : null,
      );
      createdLink = `${window.location.origin}/app/invite/${invite.code}`;
      onCreated();
    } catch (err: any) {
      notify.error(err.message ?? "Failed to create invite");
    } finally {
      creating = false;
    }
  }

  async function copyLink() {
    if (!createdLink) return;
    try {
      await navigator.clipboard.writeText(createdLink);
      notify.success("Link copied");
    } catch {
      notify.error("Failed to copy");
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") onClose();
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

      {#if createdLink}
        <h3 class="text-lg font-bold mb-1 font-[family-name:var(--font-family-display)]"
          style="color: oklch(0.92 0.004 285);">Invite created</h3>
        <p class="text-sm mb-6" style="color: oklch(0.5 0.01 285);">
          Share this link with others to let them join the server.
        </p>

        <div class="space-y-4">
          <div class="flex items-center gap-2 px-3 py-2 rounded-lg"
            style="background: oklch(0.14 0.007 285);">
            <code class="text-xs flex-1 truncate" style="color: oklch(0.72 0.18 170);">{createdLink}</code>
            <button onclick={copyLink} class="btn btn-xs btn-ghost btn-square shrink-0" aria-label="Copy link"
              style="color: oklch(0.5 0.01 285);">
              <Link class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        <div class="modal-action">
          <button onclick={onClose} class="btn"
            style="background: oklch(0.58 0.2 285); color: white;">Done</button>
        </div>
      {:else}
        <h3 class="text-lg font-bold mb-1 font-[family-name:var(--font-family-display)]"
          style="color: oklch(0.92 0.004 285);">Create Invite</h3>
        <p class="text-sm mb-6" style="color: oklch(0.5 0.01 285);">
          Set limits for your invite link.
        </p>

        <div class="space-y-4" role="none" onkeydown={handleKeydown}>
          <div class="form-control">
            <label class="label py-1" for="modal-max-uses">
              <span class="label-text text-sm" style="color: oklch(0.7 0.01 285);">Max uses</span>
            </label>
            <select
              id="modal-max-uses"
              bind:value={maxUses}
              class="select w-full"
              style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
              disabled={creating}
            >
              <option value="">Unlimited</option>
              <option value="1">1 use</option>
              <option value="5">5 uses</option>
              <option value="10">10 uses</option>
              <option value="25">25 uses</option>
              <option value="50">50 uses</option>
              <option value="100">100 uses</option>
            </select>
          </div>

          <div class="form-control">
            <label class="label py-1" for="modal-expires-hours">
              <span class="label-text text-sm" style="color: oklch(0.7 0.01 285);">Expires</span>
            </label>
            <select
              id="modal-expires-hours"
              bind:value={expiresHours}
              class="select w-full"
              style="background: oklch(0.14 0.007 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
              disabled={creating}
            >
              <option value="">Never</option>
              <option value="1">1 hour</option>
              <option value="6">6 hours</option>
              <option value="12">12 hours</option>
              <option value="24">24 hours</option>
              <option value="48">2 days</option>
              <option value="72">3 days</option>
              <option value="168">7 days</option>
            </select>
          </div>
        </div>

        <div class="modal-action">
          <button
            onclick={handleSubmit}
            class="btn"
            style="background: oklch(0.58 0.2 285); color: white; {creating ? 'opacity: 0.5;' : ''}"
            disabled={creating}
          >
            {#if creating}
              <span class="loading loading-spinner loading-xs"></span>
              Creating...
            {:else}
              <Plus class="w-4 h-4" />
              Create Invite
            {/if}
          </button>
          <button onclick={onClose} class="btn btn-ghost" disabled={creating}
            style="color: oklch(0.6 0.01 285);">
            Cancel
          </button>
        </div>
      {/if}
    </div>

    <form method="dialog" class="modal-backdrop">
      <button onclick={onClose}>close</button>
    </form>
  </dialog>
{/key}
