import api from './api';

export interface ChatResponse {
  reply: string;
  suggestions?: Array<{
    id: number;
    name: string;
    brand: string;
    slug: string;
    thumbnail_url?: string;
    price?: number;
    reason?: string;
  }>;
  follow_up_hint?: string;
}

export async function chat(message: string, sessionId: string, history?: string[]): Promise<ChatResponse> {
  const { data } = await api.post('/ai/chat', {
    message,
    session_id: sessionId,
    history,
  });
  return data as ChatResponse;
}

export interface RecommendResponse {
  suggestions: Array<{
    id: number;
    name: string;
    brand: string;
    slug: string;
    thumbnail_url?: string;
    price?: number;
    reason?: string;
  }>;
}

export async function recommend(query: string): Promise<RecommendResponse> {
  const { data } = await api.post('/ai/recommend', { query });
  return data as RecommendResponse;
}
