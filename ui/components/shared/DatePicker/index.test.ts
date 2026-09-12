import { describe, expect, it, vi } from 'vitest';

vi.mock('@mui/x-date-pickers/LocalizationProvider', () => ({
  LocalizationProvider: () => null,
}));

vi.mock('@mui/x-date-pickers/AdapterMoment', () => ({
  AdapterMoment: class {},
}));

vi.mock('./MesheryDateTimePicker', () => ({
  default: () => null,
}));

describe('DatePicker index', () => {
  it('re-exports the date picker primitives', async () => {
    const mod = await import('./index');
    expect(mod.LocalizationProvider).toBeDefined();
    expect(mod.AdapterMoment).toBeDefined();
    expect(mod.MesheryDateTimePicker).toBeDefined();
  });
});
