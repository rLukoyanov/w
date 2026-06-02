<script lang="ts">
  import {
    type Channel,
    type Message,
    channelsClient,
    messagesClient,
  } from "$lib/api/index";
  import { wsClient, wsConnected } from "$lib/websocket";
  import { auth } from "$lib/stores/auth";
  import { onMount, onDestroy } from "svelte";

  let { params } = $props();

  let channel = $state<Channel | null>(null);
  let messages = $state<Message[]>([]);
  let loading = $state(true);
  let error = $state("");
  let messageInput = $state("");
  let messagesContainer: HTMLDivElement;

  onMount(async () => {
    try {
      // Load channel info
      channel = await channelsClient.get(params.channelId);

      // Load messages
      const msgs = await messagesClient.get(params.channelId);
      messages = msgs;

      loading = false;

      // Scroll to bottom
      setTimeout(scrollToBottom, 100);

      // Subscribe to new messages
      wsClient.on("MESSAGE_CREATE", handleNewMessage);
    } catch (err: any) {
      error = err.message;
      loading = false;
    }
  });

  onDestroy(() => {
    wsClient.off("MESSAGE_CREATE", handleNewMessage);
  });

  function handleNewMessage(data: Message) {
    if (data.channel_id === params.channelId) {
      messages = [...messages, data];
      setTimeout(scrollToBottom, 10);
    }
  }

  function scrollToBottom() {
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  async function handleSendMessage(e: Event) {
    e.preventDefault();

    if (!messageInput.trim()) return;

    const content = messageInput;
    messageInput = "";

    // Send via WebSocket
    if ($wsConnected) {
      wsClient.sendMessage(params.channelId, content);
    } else {
      // Fallback to HTTP
      try {
        const msg = await messagesClient.create(params.channelId, content);
        messages = [...messages, msg];
        setTimeout(scrollToBottom, 10);
      } catch (err: any) {
        console.error("Failed to send message:", err);
        messageInput = content; // Restore input
      }
    }
  }

  function formatTime(timestamp: string) {
    const date = new Date(timestamp);
    const now = new Date();
    const isToday = date.toDateString() === now.toDateString();

    if (isToday) {
      return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
      });
    }

    return date.toLocaleString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

<main class="h-screen flex flex-col bg-base text-text">
  <!-- Top bar -->
  <div class="border-b border-border bg-surface">
    <div class="px-4 py-4 flex items-center gap-4">
      <a
        href="/servers/{params.id}"
        class="text-sm text-subtext hover:text-text">← Back</a
      >
      {#if channel}
        <h1 class="text-xl font-semibold"># {channel.name}</h1>
      {/if}
      <div class="ml-auto flex items-center gap-2">
        <div
          class="w-2 h-2 rounded-full"
          class:bg-green-500={$wsConnected}
          class:bg-red-500={!$wsConnected}
        ></div>
        <span class="text-xs text-subtext"
          >{$wsConnected ? "Live" : "Offline"}</span
        >
      </div>
    </div>
  </div>

  {#if loading}
    <div class="flex-1 flex items-center justify-center">
      <p class="text-subtext">Loading...</p>
    </div>
  {:else if error}
    <div class="flex-1 flex items-center justify-center">
      <div class="bg-surface border border-red rounded-lg p-6">
        <p class="text-red">{error}</p>
      </div>
    </div>
  {:else}
    <!-- Messages -->
    <div
      bind:this={messagesContainer}
      class="flex-1 overflow-y-auto px-4 py-4 space-y-4"
    >
      {#if messages.length === 0}
        <div class="text-center py-12">
          <p class="text-subtext">No messages yet. Start the conversation!</p>
        </div>
      {:else}
        {#each messages as message}
          <div class="flex gap-3">
            <div
              class="flex-shrink-0 w-10 h-10 rounded-full bg-blue flex items-center justify-center text-sm font-medium"
            >
              {message.author_id.slice(0, 2).toUpperCase()}
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-baseline gap-2 mb-1">
                <span class="text-sm font-medium"
                  >{message.author_id.slice(0, 8)}</span
                >
                <span class="text-xs text-subtext"
                  >{formatTime(message.created_at)}</span
                >
                {#if message.author_id === $auth.user?.id}
                  <span class="text-xs text-blue">You</span>
                {/if}
              </div>
              <p class="text-sm text-text break-words">{message.content}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>

    <!-- Input -->
    <div class="border-t border-border bg-surface p-4">
      <form onsubmit={handleSendMessage} class="flex gap-2">
        <input
          type="text"
          bind:value={messageInput}
          placeholder="Type a message..."
          class="flex-1 px-4 py-2 bg-overlay border border-border rounded text-text placeholder:text-subtext focus:outline-none focus:ring-2 focus:ring-blue"
          disabled={!$wsConnected && loading}
        />
        <button
          type="submit"
          disabled={!messageInput.trim() || (!$wsConnected && loading)}
          class="px-6 py-2 bg-blue text-text text-sm font-medium rounded hover:bg-blue/80 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Send
        </button>
      </form>
      {#if !$wsConnected}
        <p class="text-xs text-red mt-2">
          ⚠️ Disconnected - messages will use HTTP fallback
        </p>
      {/if}
    </div>
  {/if}
</main>
