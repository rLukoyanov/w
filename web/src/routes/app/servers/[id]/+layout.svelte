<script lang="ts">
  import { untrack } from "svelte";
  import { serversClient, channelsClient, type Server, type Channel, type ServerWithChannels } from "$lib/api";
  import { ROUTES } from "$lib/routes";
  import { page } from "$app/state";
  import ChannelItem from "$lib/components/ChannelItem.svelte";
  import CreateChannelForm from "$lib/components/CreateChannelForm.svelte";
  import { notify } from "$lib/stores/notifications";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { Settings, Plus } from "lucide-svelte";

  interface Props {
    params: { id: string };
    children: import("svelte").Snippet;
  }

  let { params, children }: Props = $props();

  let server = $state<Server | null>(null);
  let channels = $state<Channel[]>([]);
  let isLoading = $state(true);
  let activeChannelId = $derived(page.params.channelId as string | undefined);

  let isCreating = $state(false);
  let newChannelName = $state("");
  let channelType = $state<"text" | "voice">("text");
  let isFormValid = $derived(newChannelName.trim().length > 0);

  $effect(() => {
    const serverId = params.id;
    const controller = new AbortController();
    isLoading = true;

    untrack(async () => {
      try {
        const data = await serversClient.get(serverId, controller.signal) as ServerWithChannels;
        if (!controller.signal.aborted) {
          server = { id: data.id, name: data.name, owner_id: data.owner_id, created_at: data.created_at };
          channels = data.channels || [];
        }
      } catch (err) {
        if (!controller.signal.aborted) {
          notify.error(err instanceof Error ? err.message : "Failed to load server");
        }
      } finally {
        if (!controller.signal.aborted) {
          isLoading = false;
        }
      }
    });

    return () => controller.abort();
  });

  function resetForm() {
    isCreating = false;
    newChannelName = "";
    channelType = "text";
  }

  async function handleCreateChannel() {
    if (!isFormValid) {
      notify.error("Channel name is required");
      return;
    }

    try {
      const channel = await channelsClient.create(params.id, newChannelName, channelType);
      channels = [...channels, channel];
      resetForm();
      goto(ROUTES.SERVER.CHANNEL(params.id, channel.id));
    } catch (err) {
      notify.error(err instanceof Error ? err.message : "Failed to create channel");
    }
  }

  async function handleDeleteChannel(channelId: string, event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();

    const deleted = channels.find((ch) => ch.id === channelId);
    if (!deleted) return;

    channels = channels.filter((ch) => ch.id !== channelId);

    if (activeChannelId === channelId) {
      const next = channels.find((c) => c.type === "text");
      goto(next ? ROUTES.SERVER.CHANNEL(params.id, next.id) : ROUTES.SERVER.DETAIL(params.id));
    }

    try {
      await channelsClient.delete(channelId);
    } catch (err) {
      console.error("Failed to delete channel:", err);
      channels = [...channels, deleted].sort(
        (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
      );
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter") handleCreateChannel();
    else if (event.key === "Escape") resetForm();
    else if (event.key === "Tab") {
      event.preventDefault();
      channelType = channelType === "text" ? "voice" : "text";
    }
  }
</script>

{#if isLoading}
  <div class="flex items-center justify-center h-full">
    <span class="loading loading-spinner loading-lg" style="color: oklch(0.58 0.2 285);"></span>
  </div>
{:else if server}
  <div class="flex h-full">
    <div class="flex flex-col overflow-y-auto shrink-0"
      style="width: 280px; border-right: 1px solid oklch(0.16 0.008 285);">
      <header class="px-3 py-2.5 shrink-0"
        style="border-bottom: 1px solid oklch(0.16 0.008 285);">
        <div class="flex items-center justify-between">
          <div class="min-w-0">
            <h1 class="text-sm font-bold truncate font-[family-name:var(--font-family-display)]"
              style="color: oklch(0.92 0.004 285);">{server.name}</h1>
            <p class="text-[10px]" style="color: oklch(0.45 0.01 285);">
              {channels.length} channel{channels.length !== 1 ? "s" : ""}
            </p>
          </div>
          <div class="flex items-center gap-0.5">
            <button
              onclick={() => (isCreating = !isCreating)}
              class="btn btn-ghost btn-sm btn-square shrink-0 transition-all duration-200"
              aria-label="Create channel"
              style="color: {isCreating ? 'oklch(0.62 0.18 285)' : 'oklch(0.5 0.01 285)'};"
            >
              <Plus class="w-3.5 h-3.5" />
            </button>
            {#if server?.owner_id === $auth.user?.id}
              <a
                href={ROUTES.SERVER.SETTINGS(params.id)}
                class="btn btn-ghost btn-sm btn-square shrink-0 transition-all duration-200"
                aria-label="Server settings"
                style="color: oklch(0.5 0.01 285); hover:color: oklch(0.62 0.18 285);"
              >
                <Settings class="w-3.5 h-3.5" />
              </a>
            {/if}
          </div>
        </div>
      </header>

      {#if isCreating}
        <div class="px-2 py-1.5 shrink-0" style="background: oklch(0.09 0.005 285); border-bottom: 1px solid oklch(0.16 0.008 285);">
          <CreateChannelForm
            {isCreating}
            channelName={newChannelName}
            {channelType}
            error={null}
            onToggle={() => (isCreating = true)}
            onNameChange={(name) => (newChannelName = name)}
            onTypeChange={(type) => (channelType = type)}
            onKeyDown={handleKeyDown}
          />
        </div>
      {/if}

      <div class="flex-1 p-1.5 space-y-px">
        {#each channels as channel (channel.id)}
          <a
            href={ROUTES.SERVER.CHANNEL(params.id, channel.id)}
            class="block cursor-pointer"
          >
            <ChannelItem
              {channel}
              serverId={params.id}
              channelUrl={ROUTES.SERVER.CHANNEL(params.id, channel.id)}
              onDelete={handleDeleteChannel}
              active={channel.id === activeChannelId}
            />
          </a>
        {/each}
      </div>

    </div>

    <div class="flex-1 overflow-hidden">
      {@render children()}
    </div>
  </div>
{/if}
