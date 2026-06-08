import { writable } from "svelte/store";

export type NotificationType = "error" | "success" | "warning" | "info";

export interface Notification {
  id: string;
  type: NotificationType;
  message: string;
}

const DEFAULT_DURATION = 4000;

function createNotificationStore() {
  const notifications = writable<Notification[]>([]);

  function push(type: NotificationType, message: string, duration = DEFAULT_DURATION) {
    const id = crypto.randomUUID();
    notifications.update((n) => [...n, { id, type, message }]);

    if (duration > 0) {
      setTimeout(() => remove(id), duration);
    }

    return id;
  }

  function remove(id: string) {
    notifications.update((n) => n.filter((item) => item.id !== id));
  }

  function error(message: string) {
    return push("error", message);
  }

  function success(message: string) {
    return push("success", message);
  }

  function warning(message: string) {
    return push("warning", message);
  }

  function info(message: string) {
    return push("info", message);
  }

  return {
    subscribe: notifications.subscribe,
    push,
    remove,
    error,
    success,
    warning,
    info,
  };
}

export const notify = createNotificationStore();
