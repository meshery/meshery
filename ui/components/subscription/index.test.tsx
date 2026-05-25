import { describe, expect, it } from 'vitest';
import * as index from './index';

describe('subscription/index', () => {
  it('module loads without throwing', () => {
    expect(index).toBeDefined();
  });
});
