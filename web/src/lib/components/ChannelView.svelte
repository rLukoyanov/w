<script lang="ts">
  import type { Channel, Message } from "$lib/api";
  import { Hash, Send, MessageCircle } from "lucide-svelte";

  interface Props {
    channel: Channel | null;
    messages: Message[];
    messagesLoading: boolean;
    onSend: (content: string) => Promise<void>;
  }

  let { channel, messages, messagesLoading, onSend }: Props = $props();
  let input = $state("");

  async function handleSend() {
    const content = input.trim();
    if (!content) return;
    await onSend(content);
    input = "";
  }

  function formatTime(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }
</script>

<div class="flex flex-col overflow-hidden h-full">
  {#if channel && channel.type === "text"}
    <header
      class="px-4 py-2 border-b border-base-300 flex items-center gap-2 shrink-0"
    >
      <Hash class="w-4 h-4 text-base-content/60" />
      <h2 class="font-semibold text-sm">{channel.name}</h2>
    </header>

    <div class="flex-1 overflow-y-auto p-4 space-y-3">
      {#if messagesLoading}
        <div class="flex justify-center py-8">
          <span class="loading loading-spinner loading-sm"></span>
        </div>
      {:else if messages.length === 0}
        <div
          class="flex flex-col items-center justify-center h-full text-center"
        >
          <MessageCircle class="w-16 h-16 text-base-content/20 mb-4" />
          <p class="text-base-content/60 text-sm">No messages yet</p>
        </div>
      {:else}
        {#each messages as msg}
          <div class="flex gap-2">
            <div
              class="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-xs font-bold text-primary shrink-0"
            >
              {msg.author_id.slice(0, 2).toUpperCase()}
            </div>
            <div class="min-w-0">
              <div class="flex items-baseline gap-2">
                <span class="text-xs font-semibold"
                  >{msg.author_id.slice(0, 8)}</span
                >
                <span class="text-[10px] text-base-content/40"
                  >{formatTime(msg.created_at)}</span
                >
              </div>
              <p class="text-sm wrap-break-word">{msg.content}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>

    <div class="shrink-0 border-t border-base-300 p-3">
      <form
        onsubmit={(e) => {
          e.preventDefault();
          handleSend();
        }}
        class="flex gap-2"
      >
        <input
          type="text"
          bind:value={input}
          placeholder="Type a message..."
          class="input input-bordered input-sm flex-1"
        />
        <button
          type="submit"
          disabled={!input.trim()}
          class="btn btn-primary btn-sm"
        >
          <Send class="w-4 h-4" />
        </button>
      </form>
    </div>
  {:else if channel && channel.type === "voice"}
    <div
      class="flex items-center justify-center h-full text-base-content/40 text-sm"
    >
      Voice channel selected
    </div>
  {:else}
    <div
      class="flex items-center justify-center h-full text-base-content/40 text-sm"
    >
      Select a channel to start chatting
    </div>
  {/if}
</div>
