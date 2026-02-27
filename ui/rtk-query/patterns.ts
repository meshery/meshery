import { api } from './index';
import { ctxUrl } from '../utils/multi-ctx';

const PATTERN_ENDPOINT = '/api/pattern';

export const patternsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    deployPattern: builder.mutation({
      query: ({ patternFileData, contexts, verify = false, dryRun = false }) => {
        const endpoint = ctxUrl(PATTERN_ENDPOINT + '/deploy', contexts);
        const params = new URLSearchParams();
        if (verify) params.append('verify', 'true');
        if (dryRun) params.append('dryRun', 'true');
        const queryString = params.toString();
        const url = queryString ? `${endpoint}&${queryString}` : endpoint;

        return {
          url,
          method: 'POST',
          body: patternFileData,
        };
      },
    }),
  }),
  overrideExisting: false,
});

export const { useDeployPatternMutation } = patternsApi;

// Export helper functions for backward compatibility
export async function deployPatternWithData(patternFileData, contexts, options) {
  const { verify = false, dryRun = false } = options;
  const endpoint = `${ctxUrl(PATTERN_ENDPOINT + '/deploy', contexts)}${
    verify ? '&verify=true' : ''
  }${dryRun ? '&dryRun=true' : ''}`;

  // For now, we'll use fetch directly for backward compatibility
  const response = await fetch(endpoint, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: patternFileData,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}

export async function dryRunPattern(patternFileData, contexts) {
  return deployPatternWithData(patternFileData, contexts, {
    verify: false,
    dryRun: true,
  });
}
