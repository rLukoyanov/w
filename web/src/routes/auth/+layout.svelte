<script lang="ts">
  import { onMount } from "svelte";
  import { userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { Server, Shield, Zap, Lock } from "lucide-svelte";

  let { children } = $props();

  onMount(async () => {
    try {
      const user = await userClient.getMe();
      auth.setUser(user);
      goto("/app/servers");
    } catch (e) {
      // Not logged in, that's fine
    }
  });

  const features = [
    { icon: Server, text: "Team chat that belongs to you" },
    { icon: Zap, text: "One binary. Zero compromises." },
    { icon: Shield, text: "Self-hosted. Your data stays yours." },
    { icon: Lock, text: "End-to-end encrypted voice & video" },
  ];
</script>

<main class="grid grid-cols-2 min-h-screen">
  <div class="flex items-center justify-center p-8">
    <div class="w-full max-w-md">
      {@render children()}
    </div>
  </div>

  <div class="relative flex items-center justify-center p-12 overflow-hidden"
    style="background: linear-gradient(135deg, oklch(0.12 0.01 285) 0%, oklch(0.08 0.02 285) 50%, oklch(0.12 0.01 285) 100%);"
  >
    <div class="absolute inset-0 opacity-[0.03]"
      style="background-image: radial-gradient(circle at 25% 50%, oklch(0.58 0.2 285) 0%, transparent 50%), radial-gradient(circle at 75% 50%, oklch(0.72 0.18 170) 0%, transparent 50%);"
    ></div>

    <div class="relative z-10 max-w-md text-center">
      <div class="w-16 h-16 mx-auto mb-8 rounded-2xl flex items-center justify-center"
        style="background: linear-gradient(135deg, oklch(0.58 0.2 285), oklch(0.62 0.18 285)); box-shadow: 0 0 40px oklch(0.58 0.2 285 / 0.3);"
      >
        <span class="text-3xl font-bold font-[family-name:var(--font-family-display)] text-white">W</span>
      </div>

      <h1 class="text-4xl font-bold mb-3 font-[family-name:var(--font-family-display)]"
        style="color: oklch(0.92 0.004 285);">
        Own your team chat
      </h1>
      <p class="text-base mb-10"
        style="color: oklch(0.6 0.01 285);">
        W is an open-source communication platform packed into a single binary. No Docker. No Kubernetes. No cloud dependency.
      </p>

      <div class="space-y-4 text-left">
        {#each features as feature}
          <div class="flex items-center gap-3 px-4 py-3 rounded-xl"
            style="background: oklch(0.16 0.008 285); border: 1px solid oklch(0.22 0.01 285);">
            <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0"
              style="background: oklch(0.58 0.2 285 / 0.15);">
              <feature.icon class="w-4 h-4" style="color: oklch(0.62 0.18 285);" />
            </div>
            <span class="text-sm" style="color: oklch(0.85 0.005 285);">{feature.text}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>
</main>
