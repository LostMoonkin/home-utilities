import axios from 'axios';

const api = axios.create({
  baseURL: '/api/gateway', // Rewrite handles pointing to backend
});

// Response wrapper based on server/models/response/APIResponse.go
export interface APIResponse<T> {
  biz_code: number;
  message: string;
  data: T;
}

export interface ConfigFileInfo {
  name: string;
  size: number;
  modTime: string; // Time is serialized to string in JSON usually
}

export interface GatewayConfig {
  name: string;
  content: string; // Base64 encoded content
}

export const getConfigs = async () => {
  const response = await api.get<APIResponse<ConfigFileInfo[]>>('/conf/list');
  return response.data;
};

export const getConfig = async (name: string) => {
  // Returns map[string]string (name -> base64 content) wrapped in APIResponse
  const response = await api.get<APIResponse<Record<string, string>>>(
    `/conf?name=${name}`
  );
  return response.data;
};

export const createConfig = async (name: string, content: string) => {
  // content should be base64 encoded
  const response = await api.post<APIResponse<null>>('/conf', { name, content });
  return response.data;
};

export const updateConfig = async (name: string, currentContent: string, expectedContent: string) => {
  const response = await api.put<APIResponse<null>>('/conf', {
    name,
    currentContent,
    expectedContent
  });
  return response.data;
};

export const applyChanges = async () => {
  const response = await api.post<APIResponse<null>>('/apply');
  return response.data;
};
