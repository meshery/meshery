import { api, mesheryApiPath } from './index';

const TAGS = {
  AI_PROVIDERS: 'aiProviders',
};

export interface AIGenerateRequest {
  prompt: string;
  provider: string;
  model?: string;
  options?: {
    temperature?: number;
    max_tokens?: number;
  };
}

export interface AIGenerateResponse {
  content: string;
  finish_reason?: string;
  usage?: {
    input_tokens: number;
    output_tokens: number;
    total_tokens: number;
  };
}

export interface AIDesignGenerateRequest {
  prompt: string;
  provider: string;
  model?: string;
}

export interface AIDesignGenerateResponse {
  design: object;
  warnings?: string[];
}

export interface AIProvider {
  name: string;
  default?: boolean;
}

export interface AIProvidersResponse {
  providers: AIProvider[];
}

export interface AIHealthResponse {
  status: string;
  providers?: string[];
}

export const aiApi = api
  .enhanceEndpoints({
    addTagTypes: [TAGS.AI_PROVIDERS],
  })
  .injectEndpoints({
    endpoints: (builder) => ({
      getAIProviders: builder.query<AIProvidersResponse, void>({
        query: () => ({
          url: mesheryApiPath('ai/providers'),
          method: 'GET',
        }),
        providesTags: [TAGS.AI_PROVIDERS],
      }),

      getAIHealth: builder.query<AIHealthResponse, void>({
        query: () => ({
          url: mesheryApiPath('ai/health'),
          method: 'GET',
        }),
      }),

      aiGenerate: builder.mutation<AIGenerateResponse, AIGenerateRequest>({
        query: (body) => ({
          url: mesheryApiPath('ai/generate'),
          method: 'POST',
          body,
        }),
      }),

      aiDesignGenerate: builder.mutation<AIDesignGenerateResponse, AIDesignGenerateRequest>({
        query: (body) => ({
          url: mesheryApiPath('ai/design/generate'),
          method: 'POST',
          body,
        }),
      }),
    }),
  });

export const {
  useGetAIProvidersQuery,
  useGetAIHealthQuery,
  useAiGenerateMutation,
  useAiDesignGenerateMutation,
} = aiApi;
