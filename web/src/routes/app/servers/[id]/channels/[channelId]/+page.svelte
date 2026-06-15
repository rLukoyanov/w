<script lang="ts">
  import type { Channel, Message } from "$lib/api/index";
  import { channelsClient, messagesClient } from "$lib/api/index";
  import { wsClient, wsConnected } from "$lib/websocket";
  import { onMount, onDestroy } from "svelte";
  import ChannelView from "$lib/components/ChannelView.svelte";

  let { params } = $props();

  let channel = $state<Channel | null>(null);
  let messages = $state<Message[]>([]);
  let loading = $state(true);
  let loadingMore = $state(false);
  let hasMore = $state(true);
  let typingUsers = $state<{ user_id: string; username: string }[]>([]);

  onMount(async () => {
    try {
      channel = await channelsClient.get(params.channelId);
      const msgs = await messagesClient.get(params.channelId);
      messages = msgs.reverse();
      hasMore = msgs.length >= 50;
      loading = false;

      wsClient.on("MESSAGE_CREATE", handleNewMessage);
      wsClient.on("TYPING_START", handleTypingStart);
      wsClient.on("TYPING_STOP", handleTypingStop);
      wsClient.subscribe(params.channelId);
    } catch {
      messages = [];
      loading = false;
    }
  });

  onDestroy(() => {
    wsClient.off("MESSAGE_CREATE", handleNewMessage);
    wsClient.off("TYPING_START", handleTypingStart);
    wsClient.off("TYPING_STOP", handleTypingStop);
    wsClient.unsubscribe();
  });

  async function handleLoadMore() {
    if (loadingMore || messages.length === 0) return;
    loadingMore = true;

    try {
      const oldest = messages[0];
      const msgs = await messagesClient.get(params.channelId, 50, oldest.created_at);
      if (msgs.length < 50) hasMore = false;
      messages = [...msgs.reverse(), ...messages];
    } catch {
      // ignore
    } finally {
      loadingMore = false;
    }
  }

  function handleNewMessage(data: any) {
    if (data.channel_id === params.channelId) {
      const msg: Message = data.author
        ? data
        : {
            ...data,
            author: {
              id: data.author_id,
              username: data.author_username ?? data.author_id.slice(0, 8),
              email: "",
            },
          };
      messages = [...messages, msg];
    }
  }

  const typingTimeouts = new Map<string, number>();

  function handleTypingStart(data: any) {
    if (data.channel_id !== params.channelId) return;

    const existing = typingUsers.find((u) => u.user_id === data.user_id);
    if (existing) {
      if (typingTimeouts.has(data.user_id)) {
        clearTimeout(typingTimeouts.get(data.user_id)!);
      }
    } else {
      typingUsers = [...typingUsers, { user_id: data.user_id, username: data.username || data.user_id.slice(0, 8) }];
    }

    const timeout = setTimeout(() => {
      handleTypingStop(data);
    }, 4000) as unknown as number;
    typingTimeouts.set(data.user_id, timeout);
  }

  function handleTypingStop(data: any) {
    if (data.channel_id !== params.channelId) return;
    typingUsers = typingUsers.filter((u) => u.user_id !== data.user_id);
    if (typingTimeouts.has(data.user_id)) {
      clearTimeout(typingTimeouts.get(data.user_id)!);
      typingTimeouts.delete(data.user_id);
    }
  }

  async function handleSendMessage(content: string) {
    if ($wsConnected) {
      wsClient.sendMessage(params.channelId, content);
    } else {
      const msg = await messagesClient.create(params.channelId, content);
      messages = [...messages, msg];
    }
  }

  function handleTyping() {
    wsClient.sendTyping(params.channelId);
  }

  function handleStopTyping() {
    wsClient.sendTypingStop(params.channelId);
  }
</script>

<ChannelView
  {channel}
  {messages}
  {loadingMore}
  {hasMore}
  messagesLoading={loading}
  onSend={handleSendMessage}
  onLoadMore={handleLoadMore}
  onTyping={handleTyping}
  onStopTyping={handleStopTyping}
  typingUsers={typingUsers.map((u) => u.username)}
/>
