<script lang="ts">
  import { slide } from "svelte/transition";
  import { Hash, Volume2, Plus } from "lucide-svelte";

  interface Props {
    isCreating: boolean;
    channelName: string;
    channelType: "text" | "voice";
    error: string | null;
    onToggle: () => void;
    onNameChange: (name: string) => void;
    onTypeChange: (type: "text" | "voice") => void;
    onKeyDown: (event: KeyboardEvent) => void;
    onBlur: () => void;
  }

  let {
    isCreating,
    channelName,
    channelType,
    error,
    onToggle,
    onNameChange,
    onTypeChange,
    onKeyDown,
    onBlur,
  }: Props = $props();

  let inputEl = $state<HTMLInputElement | null>(null);

  $effect(() => {
    if (isCreating && inputEl) {
      inputEl.focus();
    }
  });
</script>

{#if isCreating}
  <div transition:slide={{ duration: 150 }}>
    <div class="flex items-center gap-1.5 px-2 py-1 rounded bg-base-100">
      {#if channelType === "voice"}
        <Volume2 class="w-3.5 h-3.5 text-base-content/60 shrink-0" />
      {:else}
        <Hash class="w-3.5 h-3.5 text-base-content/60 shrink-0" />
      {/if}
      <input
        bind:this={inputEl}
        type="text"
        value={channelName}
        oninput={(e) => onNameChange(e.currentTarget.value)}
        placeholder="channel-name"
        class="flex-1 bg-transparent outline-none text-xs"
        onkeydown={onKeyDown}
        onblur={onBlur}
      />
      <div class="flex gap-0.5">
        <button
          onclick={() => onTypeChange("text")}
          class="px-1.5 py-0.5 rounded text-[10px] transition-colors {channelType ===
          'text'
            ? 'bg-primary text-primary-content'
            : 'bg-base-200 text-base-content/60 hover:bg-base-300'}"
        >
          Text
        </button>
        <button
          onclick={() => onTypeChange("voice")}
          class="px-1.5 py-0.5 rounded text-[10px] transition-colors {channelType ===
          'voice'
            ? 'bg-primary text-primary-content'
            : 'bg-base-200 text-base-content/60 hover:bg-base-300'}"
        >
          Voice
        </button>
      </div>
    </div>
    <div class="text-[10px] text-base-content/40 px-2 py-0.5">
      <kbd class="kbd kbd-xs">Tab</kbd> type •
      <kbd class="kbd kbd-xs">Enter</kbd> create •
      <kbd class="kbd kbd-xs">Esc</kbd> cancel
    </div>
  </div>
{:else}
  <button
    onclick={onToggle}
    class="flex cursor-pointer items-center gap-1.5 px-2 py-1 rounded hover:bg-base-100/50 transition-colors group w-full text-left"
  >
    <Plus class="w-3.5 h-3.5 text-base-content/40 group-hover:text-primary shrink-0" />
    <span class="text-xs text-base-content/60 group-hover:text-base-content">
      New Channel
    </span>
  </button>
{/if}
