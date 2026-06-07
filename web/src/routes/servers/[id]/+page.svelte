<script lang="ts">
  import {
    type Server,
    type ServerWithChannels,
    type Channel,
    channelsClient,
    serversClient,
  } from "$lib/api/index";
  import { untrack } from "svelte";
  import { ROUTES } from "$lib/routes";
  import { Hash, XCircle } from "lucide-svelte";
  import ChannelItem from "$lib/components/ChannelItem.svelte";
  import CreateChannelForm from "$lib/components/CreateChannelForm.svelte";

  interface Props {
    params: { id: string };
  }

  let { params }: Props = $props();

  // Server data state
  let server = $state<Server | null>(null);
  let channels = $state<Channel[]>([]);
  let isLoading = $state(true);
  let serverError = $state<string | null>(null);

  // Channel creation state
  let isCreating = $state(false);
  let newChannelName = $state("");
  let channelType = $state<"text" | "voice">("text");
  let createError = $state<string | null>(null);

  // Derived values
  const channelCountText = $derived(
    channels.length === 1 ? "1 channel" : `${channels.length} channels`,
  );
  const hasChannels = $derived(channels.length > 0);
  const isFormValid = $derived(newChannelName.trim().length > 0);

  // Load server data with abort controller for cleanup
  $effect(() => {
    const serverId = params.id;
    const controller = new AbortController();

    isLoading = true;
    serverError = null;

    untrack(async () => {
      try {
        const response = await serversClient.get(serverId);
        const data = response as ServerWithChannels;

        if (!controller.signal.aborted) {
          server = {
            id: data.id,
            name: data.name,
            owner_id: data.owner_id,
            created_at: data.created_at,
          };
          channels = data.channels || [];
          isLoading = false;
        }
      } catch (err) {
        if (!controller.signal.aborted) {
          serverError =
            err instanceof Error ? err.message : "Failed to load server";
          isLoading = false;
        }
      }
    });

    return () => controller.abort();
  });

  function resetCreateForm(): void {
    isCreating = false;
    createError = null;
    newChannelName = "";
    channelType = "text";
  }

  async function handleCreateChannel(): Promise<void> {
    if (!isFormValid) {
      createError = "Channel name is required";
      return;
    }

    try {
      const channel = await channelsClient.create(
        params.id,
        newChannelName,
        channelType,
      );
      channels = [...channels, channel];
      resetCreateForm();
    } catch (err) {
      createError =
        err instanceof Error ? err.message : "Failed to create channel";
    }
  }

  async function handleDeleteChannel(
    channelId: string,
    event: MouseEvent,
  ): Promise<void> {
    event.preventDefault();
    event.stopPropagation();

    // Optimistic update: remove from UI immediately
    const deletedChannel = channels.find((ch) => ch.id === channelId);
    if (!deletedChannel) return;

    channels = channels.filter((ch) => ch.id !== channelId);

    try {
      await channelsClient.delete(channelId);
    } catch (err) {
      // Rollback on error: restore the channel
      console.error("Failed to delete channel:", err);
      channels = [...channels, deletedChannel].sort(
        (a, b) =>
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
      );
    }
  }

  function handleKeyDown(event: KeyboardEvent): void {
    if (event.key === "Enter") {
      handleCreateChannel();
    } else if (event.key === "Escape") {
      resetCreateForm();
    } else if (event.key === "Tab") {
      event.preventDefault();
      channelType = channelType === "text" ? "voice" : "text";
    }
  }

  function handleBlur(): void {
    if (!isFormValid) {
      resetCreateForm();
    }
  }
</script>

<main class="min-h-screen bg-base-200">
  <div class="p-4 max-w-2xl">
    {#if isLoading}
      <div class="flex justify-center py-12">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else if serverError}
      <div role="alert" class="alert alert-error">
        <XCircle class="w-6 h-6" />
        <span>{serverError}</span>
      </div>
    {:else if server}
      <!-- Server Header -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold mb-2">{server.name}</h1>
        <p class="text-base-content/60">{channelCountText}</p>
      </header>

      <!-- Channels List -->

      <div class="space-y-1 mb-6">
        {#each channels as channel (channel.id)}
          <ChannelItem
            {channel}
            serverId={params.id}
            channelUrl={ROUTES.SERVER.CHANNEL(params.id, channel.id)}
            onDelete={handleDeleteChannel}
          />
        {/each}
      </div>

      <!-- Create Channel Form -->
      <CreateChannelForm
        {isCreating}
        channelName={newChannelName}
        {channelType}
        error={createError}
        onToggle={() => (isCreating = true)}
        onNameChange={(name) => (newChannelName = name)}
        onTypeChange={(type) => (channelType = type)}
        onKeyDown={handleKeyDown}
        onBlur={handleBlur}
      />
    {/if}
  </div>
</main>
