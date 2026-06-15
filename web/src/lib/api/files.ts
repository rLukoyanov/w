import { baseClient } from './client';
import type { Attachment } from './types';

class FilesClient {
  async upload(channelId: string, file: File): Promise<Attachment> {
    const formData = new FormData();
    formData.append("file", file);

    const token = baseClient.getToken();
    const res = await fetch(`http://localhost:8090/api/channels/${channelId}/attachments`, {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      body: formData,
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: "Upload failed" }));
      throw new Error(err.error || `HTTP ${res.status}`);
    }

    return res.json();
  }
}

export const filesClient = new FilesClient();
