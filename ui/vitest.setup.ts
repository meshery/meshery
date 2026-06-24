import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

vi.mock('@/rtk-query/utils', () => ({
  shouldOverrideExisting: false,
  initiateQuery: vi.fn(
    async (
      query: { initiate: (variables: unknown, options?: unknown) => unknown },
      variables?: unknown,
      options?: unknown,
    ) => ({
      data: { stubbed: true, variables, options },
    }),
  ),
}));
vi.mock('./utils', () => ({
  shouldOverrideExisting: false,
  initiateQuery: vi.fn(
    async (
      query: { initiate: (variables: unknown, options?: unknown) => unknown },
      variables?: unknown,
      options?: unknown,
    ) => ({
      data: { stubbed: true, variables, options },
    }),
  ),
}));
