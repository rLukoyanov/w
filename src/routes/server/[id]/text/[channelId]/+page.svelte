<script lang="ts">
  import { channels as channelsApi, type Channel } from '$lib/api/channels';
  import { messages as messagesApi, type Message } from '$lib/api/messages';
  import { tick } from 'svelte';
  import { pb } from '$lib/api/pocketbase';
  import { Hash, ChevronUp, Send, Loader2 } from '@lucide/svelte';

  let { params } = $props();

  let channel = $state<Channel | null>(null);
  let messageList = $state<Message[]>([]);
  let currentPage = $state(1);
  let hasMore = $state(false);
  let loading = $state(false);
  let text = $state('');
  let sending = $state(false);
  let scrollRef = $state<HTMLDivElement | null>(null);
  let keepScroll = $state(false);
  let prevScrollHeight = $state(0);

  let shownMessages = $derived([...messageList].reverse());

  $effect(() => {
    const msgs = messageList;
    if (msgs.length > 0 && scrollRef) {
      tick().then(() => {
        if (keepScroll) {
          keepScroll = false;
          scrollRef.scrollTop = scrollRef.scrollHeight - prevScrollHeight;
        } else {
          scrollRef.scrollTop = scrollRef.scrollHeight;
        }
      });
    }
  });

  async function send() {
    const content = text.trim();
    if (!content || sending) return;

    sending = true;
    try {
      await messagesApi.create({ content, channel_id: params.channelId });
      text = '';
      await loadPage(1);
    } catch (e) {
      console.error('Failed to send message', e);
    } finally {
      sending = false;
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      send();
    }
  }

  async function loadPage(page: number) {
    loading = true;
    if (page > 1) {
      keepScroll = true;
      prevScrollHeight = scrollRef?.scrollHeight ?? 0;
    }
    try {
      const result = await messagesApi.list(params.channelId, page);
      if (page === 1) {
        messageList = result.items;
      } else {
        messageList = [...messageList, ...result.items];
      }
      currentPage = page;
      hasMore = page < result.totalPages;
    } catch (e) {
      console.error('Failed to load messages', e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    const id = params.channelId;
    channel = null;
    messageList = [];
    currentPage = 1;
    hasMore = false;

    channelsApi
      .get(id)
      .then((c) => {
        if (c.type === 'text') channel = c;
      })
      .catch(console.error);

    loadPage(1);
  });

  function userAvatar(msg: Message) {
    const u = msg.expand?.user_id;
    if (!u?.avatar) return null;
    return pb.files.getUrl(u, u.avatar, { thumb: '32x32' });
  }
</script>

<div class="flex min-h-0 flex-1 flex-col">
  {#if channel}
    <div
      class="flex h-12 shrink-0 items-center gap-2 border-b border-border px-4 font-semibold text-text shadow-sm"
    >
      <Hash class="h-5 w-5 text-text-muted" />
      {channel.name}
    </div>

    <div bind:this={scrollRef} class="flex flex-1 flex-col gap-1 overflow-y-auto p-4">
      {#if hasMore}
        <button
          onclick={() => loadPage(currentPage + 1)}
          disabled={loading}
          class="flex w-full cursor-pointer items-center justify-center gap-1 rounded-md py-2 text-sm text-text-muted transition hover:bg-bg-secondary hover:text-text disabled:opacity-50"
        >
          {#if loading}
            <Loader2 class="h-4 w-4 animate-spin" />
            Загрузка...
          {:else}
            <ChevronUp class="h-4 w-4" />
            Загрузить ещё
          {/if}
        </button>
      {/if}

      {#each shownMessages as msg}
        {#if msg.is_deleted}
          <div class="text-sm italic text-text-muted">
            Сообщение удалено
          </div>
        {:else}
          <div class="flex gap-3 rounded-lg px-3 py-2 transition hover:bg-bg-secondary">
            {#if userAvatar(msg)}
              <img
                src={userAvatar(msg)}
                alt=""
                class="mt-0.5 h-8 w-8 shrink-0 rounded-full object-cover"
              />
            {:else}
              <div
                class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary text-sm font-bold text-text-inverse"
              >
                {msg.expand?.user_id?.username?.[0]?.toUpperCase() ?? '?'}
              </div>
            {/if}

            <div class="min-w-0 flex-1">
              <div class="flex items-baseline gap-2">
                <span class="text-sm font-semibold text-text">
                  {msg.expand?.user_id?.username ?? 'Неизвестно'}
                </span>
                <span class="text-xs text-text-muted">
                  {new Date(msg.created).toLocaleString('ru-RU')}
                </span>
              </div>
              <p class="mt-0.5 break-all text-sm text-text-secondary">{msg.content}</p>

              {#if msg.images?.length}
                <div class="mt-2 flex flex-wrap gap-2">
                  {#each msg.images as img}
                    <img
                      src={pb.files.getUrl(msg, img, { thumb: '200x200' })}
                      alt=""
                      class="max-h-48 rounded-lg object-cover"
                    />
                  {/each}
                </div>
              {/if}
            </div>
          </div>
        {/if}
      {/each}

      {#if messageList.length === 0}
        <div class="flex flex-1 items-center justify-center">
          <p class="text-sm text-text-muted">В этом канале пока нет сообщений</p>
        </div>
      {/if}
    </div>

    <div class="flex items-center gap-2 border-t border-border p-4">
      <input
        type="text"
        bind:value={text}
        onkeydown={handleKeydown}
        placeholder="Написать сообщение..."
        disabled={sending}
        class="flex-1 rounded-lg border border-border bg-surface px-4 py-2 text-sm text-text outline-none transition placeholder:text-text-muted focus:border-primary"
      />
      <button
        onclick={send}
        disabled={sending || !text.trim()}
        class="flex h-9 w-9 cursor-pointer items-center justify-center rounded-lg bg-primary text-text-inverse transition hover:bg-primary-hover disabled:cursor-not-allowed disabled:opacity-50"
      >
        <Send class="h-4 w-4" />
      </button>
    </div>
  {:else}
    <div class="flex flex-1 items-center justify-center">
      <p class="text-sm text-text-muted">Загрузка...</p>
    </div>
  {/if}
</div>
