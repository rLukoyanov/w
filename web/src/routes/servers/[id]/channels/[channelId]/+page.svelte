<script lang="ts">
  import type { Channel, Message } from "$lib/api";
  import { channelsClient, messagesClient } from "$lib/api/index";
  import { wsClient, wsConnected } from "$lib/websocket";
  import { onMount, onDestroy } from "svelte";
  import ChannelView from "$lib/components/ChannelView.svelte";

  let { params } = $props();

  let channel = $state<Channel | null>(null);
  let messages = $state<Message[]>([]);
  let loading = $state(true);
  let messagesContainer = $state<HTMLDivElement | null>(null);

  onMount(async () => {
    try {
      channel = await channelsClient.get(params.channelId);
      const msgs = await messagesClient.get(params.channelId);
      messages = msgs;
      loading = false;
      wsClient.on("MESSAGE_CREATE", handleNewMessage);
    } catch {
      messages = [];
      loading = false;
    }
  });

  onDestroy(() => {
    wsClient.off("MESSAGE_CREATE", handleNewMessage);
  });

  function handleNewMessage(data: Message) {
    if (data.channel_id === params.channelId) {
      messages = [...messages, data];
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
  messagesLoading={loading}
  onSend={handleSendMessage}
/>
