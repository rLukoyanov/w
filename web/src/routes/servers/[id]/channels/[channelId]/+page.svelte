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
  import { goto } from "$app/navigation";
  import { Hash, XCircle, MessageCircle, Send, AlertTriangle } from "lucide-svelte";

  let { params } = $props();

  let channel = $state<Channel | null>(null);
  let messages = $state<Message[]>([]);
  let loading = $state(true);
  let error = $state("");
  let messageInput = $state("");
  let messagesContainer = $state<HTMLDivElement | null>(null);

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

  function getInitials(userId: string): string {
    return userId.slice(0, 2).toUpperCase();
  }

  function getUsername(userId: string): string {
    return userId.slice(0, 8);
  }

  const avatarColors = [
    "bg-primary",
    "bg-secondary",
    "bg-accent",
    "bg-info",
    "bg-success",
    "bg-warning",
  ];

  function getUserColor(userId: string): string {
    const index =
      userId.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0) %
      avatarColors.length;
    return avatarColors[index];
  }
</script>

<div class="flex flex-col h-screen bg-base-200">
  <!-- Navbar with breadcrumbs -->
  <div class="navbar bg-base-100 shadow-md">
    <div class="flex-1">
      <div class="breadcrumbs text-sm">
        <ul>
          <li>
            <a href="/app" class="link link-hover">Home</a>
          </li>
          <li>
            <a href="/app/{params.id}" class="link link-hover">Server</a>
          </li>
          {#if channel}
            <li>
              <span class="flex items-center gap-1">
                <Hash class="w-4 h-4" />
                {channel.name}
              </span>
            </li>
          {/if}
        </ul>
      </div>
    </div>
    <div class="flex-none gap-2">
      <!-- WebSocket status badge -->
      {#if $wsConnected}
        <div class="badge badge-success gap-2">
          <div class="w-2 h-2 rounded-full bg-success-content animate-pulse"></div>
          Connected
        </div>
      {:else}
        <div class="badge badge-error gap-2">
          <div class="w-2 h-2 rounded-full bg-error-content"></div>
          Disconnected
        </div>
      {/if}
    </div>
  </div>

  {#if loading}
    <div class="flex-1 flex items-center justify-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error}
    <div class="flex-1 flex items-center justify-center p-4">
      <div class="alert alert-error">
        <XCircle class="w-6 h-6" />
        <span>{error}</span>
      </div>
    </div>
  {:else}
    <!-- Messages container -->
    <div
      bind:this={messagesContainer}
      class="flex-1 overflow-y-auto p-4 space-y-4"
    >
      {#if messages.length === 0}
        <div class="hero min-h-full">
          <div class="hero-content text-center">
            <div class="max-w-md">
              <div class="mb-4">
                <MessageCircle class="w-24 h-24 mx-auto text-base-content/20" />
              </div>
              <h1 class="text-3xl font-bold">No messages yet</h1>
              <p class="py-4 text-base-content/60">
                Start the conversation! Send the first message in this channel.
              </p>
            </div>
          </div>
        </div>
      {:else}
        {#each messages as message}
          <div
            class="chat"
            class:chat-end={message.author_id === $auth.user?.id}
            class:chat-start={message.author_id !== $auth.user?.id}
          >
            <div class="chat-image avatar">
              <div
                class="w-10 rounded-full {getUserColor(message.author_id)} text-base-100 flex items-center justify-center font-semibold"
              >
                {getInitials(message.author_id)}
              </div>
            </div>
            <div class="chat-header">
              {getUsername(message.author_id)}
              <time class="text-xs opacity-50 ml-1">
                {formatTime(message.created_at)}
              </time>
            </div>
            <div
              class="chat-bubble"
              class:chat-bubble-primary={message.author_id === $auth.user?.id}
            >
              {message.content}
            </div>
          </div>
        {/each}
      {/if}
    </div>

    <!-- Message input -->
    <div class="bg-base-100 border-t border-base-300 p-4">
      <form onsubmit={handleSendMessage} class="flex gap-2">
        <input
          type="text"
          bind:value={messageInput}
          placeholder="Type a message..."
          class="input input-bordered flex-1"
          disabled={!$wsConnected && loading}
        />
        <button
          type="submit"
          disabled={!messageInput.trim() || (!$wsConnected && loading)}
          class="btn btn-primary"
        >
          <Send class="w-5 h-5" />
          Send
        </button>
      </form>
      {#if !$wsConnected}
        <div class="alert alert-warning mt-2">
          <AlertTriangle class="w-5 h-5" />
          <span>WebSocket disconnected - using HTTP fallback</span>
        </div>
      {/if}
    </div>
  {/if}
</div>
