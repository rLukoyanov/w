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
    <div class="flex items-center gap-1.5 px-2 py-1 rounded-md"
      style="background: oklch(0.14 0.007 285);">
      {#if channelType === "voice"}
        <Volume2 class="w-3.5 h-3.5 shrink-0" style="color: oklch(0.5 0.01 285);" />
      {:else}
        <Hash class="w-3.5 h-3.5 shrink-0" style="color: oklch(0.5 0.01 285);" />
      {/if}
      <input
        bind:this={inputEl}
        type="text"
        value={channelName}
        oninput={(e) => onNameChange(e.currentTarget.value)}
        placeholder="channel-name"
        class="flex-1 bg-transparent outline-none text-xs"
        style="color: oklch(0.92 0.004 285);"
        onkeydown={onKeyDown}
      />
      <div class="flex gap-0.5">
        <button
          onclick={() => onTypeChange("text")}
          class="px-1.5 py-0.5 rounded text-[10px] transition-colors cursor-pointer"
          style="background: {channelType === 'text' ? 'oklch(0.58 0.2 285)' : 'oklch(0.18 0.008 285)'}; color: {channelType === 'text' ? 'white' : 'oklch(0.5 0.01 285)'};"
        >
          Text
        </button>
        <button
          onclick={() => onTypeChange("voice")}
          class="px-1.5 py-0.5 rounded text-[10px] transition-colors cursor-pointer"
          style="background: {channelType === 'voice' ? 'oklch(0.58 0.2 285)' : 'oklch(0.18 0.008 285)'}; color: {channelType === 'voice' ? 'white' : 'oklch(0.5 0.01 285)'};"
        >
          Voice
        </button>
      </div>
    </div>
    <div class="text-[10px] px-2 py-0.5" style="color: oklch(0.4 0.01 285);">
      <kbd class="kbd kbd-xs">Tab</kbd> type •
      <kbd class="kbd kbd-xs">Enter</kbd> create •
      <kbd class="kbd kbd-xs">Esc</kbd> cancel
    </div>
  </div>
{:else}
  <button
    onclick={onToggle}
    class="flex cursor-pointer items-center gap-1.5 px-2 py-1 rounded-md transition-colors group w-full text-left"
    style="color: oklch(0.5 0.01 285);"
  >
    <Plus class="w-3.5 h-3.5 shrink-0" style="color: oklch(0.4 0.01 285);" />
    <span class="text-xs">New Channel</span>
  </button>
{/if}
