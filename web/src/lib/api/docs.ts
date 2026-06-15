import { BASE_URL } from "./client";

export interface OpenAPISpec {
  openapi: string;
  info: { title: string; description?: string; version: string };
  servers?: { url: string; description?: string }[];
  paths: Record<string, PathItem>;
  components?: {
    schemas?: Record<string, SchemaObject>;
  };
}

export interface PathItem {
  get?: Operation;
  post?: Operation;
  patch?: Operation;
  delete?: Operation;
}

export interface Operation {
  tags?: string[];
  summary?: string;
  description?: string;
  operationId?: string;
  parameters?: ParameterObject[];
  requestBody?: RequestBodyObject;
  responses: Record<string, ResponseObject>;
  security?: Record<string, string[]>[];
}

export interface ParameterObject {
  name: string;
  in: "path" | "query" | "header" | "cookie";
  required?: boolean;
  description?: string;
  schema: SchemaObject;
}

export interface SchemaObject {
  type?: string;
  format?: string;
  properties?: Record<string, SchemaObject>;
  items?: SchemaObject;
  required?: string[];
  enum?: string[];
  description?: string;
  nullable?: boolean;
  example?: unknown;
}

export interface RequestBodyObject {
  description?: string;
  required?: boolean;
  content: Record<string, { schema: SchemaObject }>;
}

export interface ResponseObject {
  description: string;
  content?: Record<string, { schema: SchemaObject }>;
}

export interface EndpointInfo {
  path: string;
  method: string;
  operation: Operation;
}

export interface TagGroup {
  tag: string;
  endpoints: EndpointInfo[];
}

export function groupByTag(spec: OpenAPISpec): TagGroup[] {
  const groups = new Map<string, EndpointInfo[]>();

  for (const [path, item] of Object.entries(spec.paths)) {
    for (const method of ["get", "post", "patch", "delete"] as const) {
      const op = item[method];
      if (!op) continue;

      const tag = (op.tags && op.tags[0]) || "default";
      if (!groups.has(tag)) groups.set(tag, []);
      groups.get(tag)!.push({ path, method, operation: op });
    }
  }

  const order = ["Auth", "Servers", "Channels", "Messages", "Users", "Invites", "WebSocket"];

  return Array.from(groups.entries()).sort((a, b) => {
    const ia = order.indexOf(a[0]);
    const ib = order.indexOf(b[0]);
    if (ia !== -1 && ib !== -1) return ia - ib;
    if (ia !== -1) return -1;
    if (ib !== -1) return 1;
    return a[0].localeCompare(b[0]);
  }).map(([tag, endpoints]) => ({ tag, endpoints }));
}

export function resolveRef(ref: string): string {
  return ref.replace("#/components/schemas/", "");
}

export async function fetchSpec(): Promise<OpenAPISpec> {
  const headers: Record<string, string> = {};
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE_URL}/api/openapi.json`, { headers });
  if (!res.ok) throw new Error(`Failed to load spec: ${res.statusText}`);
  return res.json();
}
