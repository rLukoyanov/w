<script lang="ts">
  import type { Channel, Message, Attachment } from "$lib/api";
  import { Hash, Send, Paperclip, File, Image, X, Loader, MessageCircle } from "lucide-svelte";
  import { onMount, tick } from "svelte";
  import { BASE_URL } from "$lib/api/client";

  interface Props {
    channel: Channel | null;
    messages: Message[];
    messagesLoading: boolean;
    loadingMore?: boolean;
    hasMore?: boolean;
    onSend: (content: string) => Promise<void>;
    onLoadMore?: () => void;
    onTyping?: () => void;
    onStopTyping?: () => void;
    typingUsers?: string[];
  }

  let {
    channel,
    messages,
    messagesLoading,
    loadingMore = false,
    hasMore = false,
    onSend,
    onLoadMore,
    onTyping,
    onStopTyping,
    typingUsers = [],
  }: Props = $props();
  let input = $state("");
  let scrollContainer = $state<HTMLDivElement | null>(null);
  let sentinel = $state<HTMLDivElement | null>(null);
  let fileUploading = $state(false);
  let fileInput = $state<HTMLInputElement | null>(null);

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

  function handleInput() {
    if (onTyping) onTyping();
  }

  function handleBlur() {
    if (onStopTyping) onStopTyping();
  }

  async function handleSend() {
    const content = input.trim();
    if (!content) return;
    await onSend(content);
    input = "";
    if (onStopTyping) onStopTyping();
  }

  async function handleFileSelect(e: Event) {
    const target = e.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file || !channel) return;

    fileUploading = true;
    try {
      const formData = new FormData();
      formData.append("file", file);

      const token = localStorage.getItem("token");
      const res = await fetch(`${BASE_URL}/api/channels/${channel.id}/attachments`, {
        method: "POST",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: formData,
      });

      if (!res.ok) throw new Error("Upload failed");

      const attachment: Attachment = await res.json();

      await onSend(attachment.url);

      scrollToBottom();
    } catch (err) {
      console.error("Upload error:", err);
    } finally {
      fileUploading = false;
      if (fileInput) fileInput.value = "";
    }
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
    const index = seed.split("").reduce((a, c) => a + c.charCodeAt(0), 0) % colors.length;
    return colors[index];
  }

  function formatSize(bytes: number): string {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
    return (bytes / (1024 * 1024)).toFixed(1) + " MB";
  }

  function isImage(filename: string): boolean {
    return /\.(jpg|jpeg|png|gif|webp|svg|bmp)$/i.test(filename);
  }

  function typingText(): string {
    if (typingUsers.length === 0) return "";
    if (typingUsers.length === 1) return `${typingUsers[0]} is typing...`;
    if (typingUsers.length === 2) return `${typingUsers[0]} and ${typingUsers[1]} are typing...`;
    return `${typingUsers[0]} and ${typingUsers.length - 1} others are typing...`;
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
          <div class="flex gap-2.5 group px-2 py-1.5 rounded-lg transition-colors">
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
              {#if msg.attachments}
                {#each msg.attachments as att}
                  <a href={att.url} target="_blank" rel="noopener noreferrer"
                    class="inline-flex items-center gap-2 mt-1.5 px-3 py-2 rounded-lg text-xs transition-colors att-link"
                    style="background: oklch(0.12 0.006 285); color: oklch(0.72 0.18 285);">
                    {#if isImage(att.filename)}
                      <Image class="w-3.5 h-3.5 shrink-0" />
                    {:else}
                      <File class="w-3.5 h-3.5 shrink-0" />
                    {/if}
                    <span class="truncate max-w-[200px]">{att.filename}</span>
                    <span style="color: oklch(0.45 0.01 285);">({formatSize(att.size)})</span>
                  </a>
                {/each}
              {/if}
            </div>
          </div>
        {/each}

        {#if typingText()}
          <div class="flex items-center gap-2 px-2 py-1 text-xs italic"
            style="color: oklch(0.5 0.01 285);">
            <span class="loading loading-dots loading-xs"></span>
            {typingText()}
          </div>
        {/if}
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
        <button
          type="button"
          disabled={fileUploading}
          class="btn btn-sm btn-ghost px-2"
          style="color: oklch(0.5 0.01 285);"
          onclick={() => fileInput?.click()}
        >
          {#if fileUploading}
            <Loader class="w-4 h-4 animate-spin" />
          {:else}
            <Paperclip class="w-4 h-4" />
          {/if}
        </button>
        <input
          type="file"
          bind:this={fileInput}
          class="hidden"
          onchange={handleFileSelect}
        />
        <input
          type="text"
          bind:value={input}
          oninput={handleInput}
          onblur={handleBlur}
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

<style>
  .att-link:hover {
    background: oklch(0.15 0.008 285) !important;
  }
</style>
