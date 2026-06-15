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
      "oklch(0.58 0.2 285 / 0.2)",
      "oklch(0.72 0.18 170 / 0.2)",
      "oklch(0.65 0.18 25 / 0.2)",
      "oklch(0.82 0.15 85 / 0.2)",
      "oklch(0.62 0.18 285 / 0.2)",
      "oklch(0.55 0.15 200 / 0.2)",
    ];
    const seed = msg.author?.username ?? msg.author_id;
    const index =
      seed.split("").reduce((a, c) => a + c.charCodeAt(0), 0) % colors.length;
    return colors[index];
  }
</script>

<div class="flex flex-col overflow-hidden h-full">
  {#if channel && channel.type === "text"}
    <header class="px-4 py-2.5 border-b shrink-0"
      style="border-color: oklch(0.2 0.01 285);">
      <div class="flex items-center gap-2">
        <Hash class="w-4 h-4" style="color: oklch(0.58 0.2 285);" />
        <h2 class="text-sm font-semibold font-[family-name:var(--font-family-display)]"
          style="color: oklch(0.92 0.004 285);">{channel.name}</h2>
      </div>
    </header>

    <div bind:this={scrollContainer} class="flex-1 overflow-y-auto p-4 space-y-2">
      {#if messagesLoading}
        <div class="flex justify-center py-8">
          <span class="loading loading-spinner loading-sm"></span>
        </div>
      {:else if messages.length === 0}
        <div class="flex flex-col items-center justify-center h-full text-center">
          <MessageCircle class="w-16 h-16 mb-4" style="color: oklch(0.25 0.01 285);" />
          <p class="text-sm" style="color: oklch(0.5 0.01 285);">No messages yet</p>
          <p class="text-xs mt-1" style="color: oklch(0.4 0.01 285);">Start the conversation</p>
        </div>
      {:else}
        <div bind:this={sentinel} class="h-px"></div>

        {#if loadingMore}
          <div class="flex justify-center py-2">
            <span class="loading loading-spinner loading-xs"></span>
          </div>
        {/if}

        {#each messages as msg}
          <div class="flex gap-2.5 group px-2 py-1.5 rounded-lg transition-colors hover:oklch(0.15 0.008 285 / 0.3)">
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shrink-0 font-[family-name:var(--font-family-display)]"
              style="background: {avatarColor(msg)};">
              {initials(msg)}
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-baseline gap-2">
                <span class="text-xs font-semibold font-[family-name:var(--font-family-display)]"
                  style="color: oklch(0.85 0.01 285);"
                  >{msg.author?.username ?? msg.author_id.slice(0, 8)}</span>
                <span class="text-[10px]" style="color: oklch(0.4 0.01 285);"
                  >{formatTime(msg.created_at)}</span>
              </div>
              <p class="text-sm leading-relaxed mt-0.5 wrap-break-word"
                style="color: oklch(0.88 0.005 285);">{msg.content}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>

    <div class="shrink-0 border-t p-3" style="border-color: oklch(0.2 0.01 285);">
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
          class="input input-sm flex-1"
          style="background: oklch(0.12 0.006 285); border-color: oklch(0.22 0.01 285); color: oklch(0.92 0.004 285);"
        />
        <button
          type="submit"
          disabled={!input.trim()}
          class="btn btn-sm"
          style="background: oklch(0.58 0.2 285); color: white; {!input.trim() ? 'opacity: 0.5;' : ''}"
        >
          <Send class="w-4 h-4" />
        </button>
      </form>
    </div>
  {:else if channel && channel.type === "voice"}
    <div class="flex items-center justify-center h-full text-sm"
      style="color: oklch(0.45 0.01 285);">
      Voice channel selected
    </div>
  {:else}
    <div class="flex items-center justify-center h-full text-sm"
      style="color: oklch(0.45 0.01 285);">
      Select a channel to start chatting
    </div>
  {/if}
</div>
