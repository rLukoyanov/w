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

  onMount(async () => {
    try {
      channel = await channelsClient.get(params.channelId);
      const msgs = await messagesClient.get(params.channelId);
      messages = msgs.reverse();
      hasMore = msgs.length >= 50;
      loading = false;

      wsClient.on("MESSAGE_CREATE", handleNewMessage);
      wsClient.subscribe(params.channelId);
    } catch {
      messages = [];
      loading = false;
    }
  });

  onDestroy(() => {
    wsClient.off("MESSAGE_CREATE", handleNewMessage);
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

  async function handleSendMessage(content: string) {
    if ($wsConnected) {
      wsClient.sendMessage(params.channelId, content);
    } else {
      const msg = await messagesClient.create(params.channelId, content);
      messages = [...messages, msg];
    }
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
/>
