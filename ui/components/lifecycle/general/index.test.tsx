import { describe, expect, it, vi } from 'vitest';

vi.mock('./empty-state', () => ({
  default: () => null,
}));

vi.mock('./flip-card', () => ({
  default: () => null,
}));

import { EmptyState, FlipCard } from './index';

describe('lifecycle/general index', () => {
  it('re-exports EmptyState and FlipCard', () => {
    expect(EmptyState).toBeDefined();
    expect(FlipCard).toBeDefined();
  });
});
