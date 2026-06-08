<script lang="ts">
  import { onMount } from "svelte";
  import { serversClient, invitesClient, type Server, type Invite, type ServerWithChannels } from "$lib/api";
  import { notify } from "$lib/stores/notifications";
  import { ROUTES } from "$lib/routes";
  import { auth } from "$lib/stores/auth";
  import { page } from "$app/state";
  import { Copy, Trash2, Plus, Link, Check } from "lucide-svelte";
  import CreateInviteModal from "$lib/components/CreateInviteModal.svelte";

  let { params } = $props();

  let server = $state<Server | null>(null);
  let invites = $state<Invite[]>([]);
  let loading = $state(true);
  let showCreateModal = $state(false);
  let copiedId = $state<string | null>(null);

  let isOwner = $derived(server?.owner_id === $auth.user?.id);

  onMount(async () => {
    try {
      const data = await serversClient.get(params.id) as ServerWithChannels;
      server = { id: data.id, name: data.name, owner_id: data.owner_id, created_at: data.created_at };
      invites = await invitesClient.list(params.id);
    } catch {
      notify.error("Failed to load server settings");
    } finally {
      loading = false;
    }
  });

  async function loadInvites() {
    try {
      invites = await invitesClient.list(params.id);
    } catch {
      notify.error("Failed to load invites");
    }
  }

  async function handleDelete(inviteId: string) {
    try {
      await invitesClient.delete(inviteId);
      invites = invites.filter((i) => i.id !== inviteId);
      notify.success("Invite deleted");
    } catch {
      notify.error("Failed to delete invite");
    }
  }

  async function copyLink(invite: Invite) {
    const url = `${window.location.origin}${ROUTES.INVITE.JOIN(invite.code)}`;
    try {
      await navigator.clipboard.writeText(url);
      copiedId = invite.id;
      setTimeout(() => (copiedId = null), 2000);
    } catch {
      notify.error("Failed to copy link");
    }
  }

  function formatExpiry(invite: Invite): string {
    if (!invite.expires_at) return "Never";
    const date = new Date(invite.expires_at);
    return date.toLocaleDateString();
  }

  function usesLeft(invite: Invite): string {
    if (!invite.max_uses) return "Unlimited";
    const left = invite.max_uses - invite.use_count;
    return `${left} left`;
  }
</script>

<div class="p-6 max-w-2xl">
  {#if loading}
    <span class="loading loading-spinner loading-sm"></span>
  {:else if server}
    <h1 class="text-2xl font-bold mb-1">{server.name}</h1>
    <p class="text-sm text-base-content/60 mb-6">Server Settings</p>

    <section class="mb-8">
      <div class="flex items-center justify-between mb-3">
        <h2 class="text-lg font-semibold">Invites</h2>
        {#if isOwner}
          <button
            onclick={() => (showCreateModal = true)}
            class="btn btn-primary btn-sm gap-1"
          >
            <Plus class="w-4 h-4" />
            New Invite
          </button>
        {/if}
      </div>

      {#if invites.length === 0}
        <p class="text-sm text-base-content/40">No invites yet</p>
      {:else}
        <div class="space-y-2">
          {#each invites as invite}
            <div class="flex items-center gap-3 px-4 py-3 rounded-lg bg-base-100">
              <Link class="w-4 h-4 text-base-content/40 shrink-0" />
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <code class="text-xs font-mono">{invite.code}</code>
                  <span class="text-[10px] text-base-content/40">
                    · {usesLeft(invite)} · Expires {formatExpiry(invite)}
                  </span>
                </div>
              </div>
              <div class="flex gap-1">
                <button
                  onclick={() => copyLink(invite)}
                  class="btn btn-ghost btn-xs btn-square"
                  aria-label="Copy invite link"
                >
                  {#if copiedId === invite.id}
                    <Check class="w-3.5 h-3.5 text-success" />
                  {:else}
                    <Copy class="w-3.5 h-3.5" />
                  {/if}
                </button>
                {#if isOwner}
                  <button
                    onclick={() => handleDelete(invite.id)}
                    class="btn btn-ghost btn-xs btn-square hover:text-error"
                    aria-label="Delete invite"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </section>
  {/if}
</div>

<CreateInviteModal
  show={showCreateModal}
  serverId={params.id}
  onClose={() => (showCreateModal = false)}
  onCreated={loadInvites}
/>
