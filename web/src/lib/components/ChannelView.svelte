<script lang="ts">
  import type { Channel, Message } from "$lib/api";
  import { Hash, Send, MessageCircle } from "lucide-svelte";
  import { onMount, tick } from "svelte";

  interface Props {
    channel: Channel | null;
    messages: Message[];
    messagesLoading: boolean;
    loadingMore?: boolean;
    hasMore?: boolean;
    onSend: (content: string) => Promise<void>;
    onLoadMore?: () => void;
  }

  let {
    channel,
    messages,
    messagesLoading,
    loadingMore = false,
    hasMore = false,
    onSend,
    onLoadMore,
  }: Props = $props();
  let input = $state("");
  let scrollContainer = $state<HTMLDivElement | null>(null);
  let sentinel = $state<HTMLDivElement | null>(null);

  onMount(() => {
    scrollToBottom();

    if (!scrollContainer || !sentinel || !onLoadMore) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loadingMore) {
          onLoadMore();
        }
      },
      { root: scrollContainer, threshold: 0 },
    );

    observer.observe(sentinel);
    return () => observer.disconnect();
  });

  // Scroll to bottom when new messages arrive and user is near bottom
  $effect(() => {
    if (messages.length > 0 && onLoadMore === undefined) {
      // Initial load — also handled by onMount
    }
  });

  // Auto-scroll on new messages
  let prevLength = $state(0);
  $effect(() => {
    if (messages.length > prevLength && prevLength > 0) {
      tick().then(() => {
        if (scrollContainer) {
          const threshold = 100;
          const atBottom =
            scrollContainer.scrollHeight - scrollContainer.scrollTop - scrollContainer.clientHeight < threshold;
          if (atBottom) scrollToBottom();
        }
      });
    }
    prevLength = messages.length;
  });

  function scrollToBottom() {
    tick().then(() => {
      if (scrollContainer) {
        scrollContainer.scrollTop = scrollContainer.scrollHeight;
      }
    });
  }

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

  function initials(msg: Message): string {
    const name = msg.author?.username;
    if (name) return name[0].toUpperCase();
    return msg.author_id.slice(0, 1).toUpperCase();
  }

  function avatarColor(msg: Message): string {
    const colors = [
      "bg-primary/20 text-primary",
      "bg-secondary/20 text-secondary",
      "bg-accent/20 text-accent",
      "bg-info/20 text-info",
      "bg-success/20 text-success",
      "bg-warning/20 text-warning",
    ];
    const seed = msg.author?.username ?? msg.author_id;
    const index =
      seed.split("").reduce((a, c) => a + c.charCodeAt(0), 0) % colors.length;
    return colors[index];
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

    <div bind:this={scrollContainer} class="flex-1 overflow-y-auto p-4 space-y-3">
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
        <div bind:this={sentinel} class="h-px"></div>

        {#if loadingMore}
          <div class="flex justify-center py-2">
            <span class="loading loading-spinner loading-xs"></span>
          </div>
        {/if}

        {#each messages as msg}
          <div class="flex gap-2">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0 {avatarColor(msg)}"
            >
              {initials(msg)}
            </div>
            <div class="min-w-0">
              <div class="flex items-baseline gap-2">
                <span class="text-xs font-semibold"
                  >{msg.author?.username ?? msg.author_id.slice(0, 8)}</span
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
