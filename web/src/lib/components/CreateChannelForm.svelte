<script lang="ts">
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
</script>

{#if isCreating}
  <div class="space-y-2">
    <div class="flex items-center gap-3 px-4 py-3 rounded-lg bg-base-100">
      {#if channelType === "voice"}
        <Volume2 class="w-5 h-5 text-primary" />
      {:else}
        <Hash class="w-5 h-5 text-primary" />
      {/if}
      <input
        type="text"
        value={channelName}
        oninput={(e) => onNameChange(e.currentTarget.value)}
        placeholder="channel-name"
        class="flex-1 bg-transparent outline-none text-lg"
        onkeydown={onKeyDown}
        onblur={onBlur}
      />
      <div class="flex gap-1">
        <button
          onclick={() => onTypeChange("text")}
          class="px-2 py-1 rounded text-xs transition-colors {channelType ===
          'text'
            ? 'bg-primary text-primary-content'
            : 'bg-base-200 text-base-content/60 hover:bg-base-300'}"
        >
          Text
        </button>
        <button
          onclick={() => onTypeChange("voice")}
          class="px-2 py-1 rounded text-xs transition-colors {channelType ===
          'voice'
            ? 'bg-primary text-primary-content'
            : 'bg-base-200 text-base-content/60 hover:bg-base-300'}"
        >
          Voice
        </button>
      </div>
    </div>
    <div class="text-sm text-base-content/40 px-4">
      <kbd class="kbd kbd-xs">Tab</kbd> to switch type •
      <kbd class="kbd kbd-xs">Enter</kbd>
      to create • <kbd class="kbd kbd-xs">Esc</kbd> to cancel
    </div>
  </div>
  {#if error}
    <div class="text-error text-sm px-4 py-2">
      {error}
    </div>
  {/if}
{:else}
  <button
    onclick={onToggle}
    class="flex cursor-pointer items-center gap-3 px-4 py-3 rounded-lg hover:bg-base-100 transition-colors group w-full text-left"
  >
    <Plus class="w-5 h-5 text-base-content/40 group-hover:text-primary" />
    <span class="text-lg text-base-content/60 group-hover:text-base-content">
      New Channel
    </span>
  </button>
{/if}
