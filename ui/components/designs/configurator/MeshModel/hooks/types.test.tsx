import { describe, expect, it } from 'vitest';
import { Types } from './types';

describe('configurator/MeshModel/hooks/types', () => {
  it('exports the Types JSDoc-only namespace as an empty object', () => {
    expect(Types).toEqual({});
  });
});
