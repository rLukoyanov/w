<script lang="ts">
  import { type Server, serversClient, userClient } from "$lib/api/index";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import ServerCard from "$lib/components/ServerCard.svelte";
  import { Plus, XCircle, Server as ServerIcon } from "lucide-svelte";

  let servers = $state<Server[]>([]);
  let loading = $state(true);
  let showCreateServer = $state(false);
  let newServerName = $state("");
  let createError = $state("");

  async function handleCreateServer() {
    if (!newServerName.trim()) {
      createError = "Server name is required";
      return;
    }

    try {
      const server = await serversClient.create(newServerName);
      servers.push(server);
      newServerName = "";
      showCreateServer = false;
      createError = "";

      // Navigate to the new server
      goto(`/servers/${server.id}`);
    } catch (err: any) {
      createError = err.message;
    }
  }
</script>

<main class="min-h-screen">
</main>
