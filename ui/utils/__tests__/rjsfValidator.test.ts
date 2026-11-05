import { describe, expect, it } from 'vitest';
import customValidator from '../rjsfValidator';

// `rjsfValidator` exports a single AJV-backed validator wired up with
// Kubernetes / OpenAPI custom formats (int32, int64, int-or-string) plus the
// standard `ajv-formats` set (date-time, email, ...). These tests exercise the
// validator end-to-end so a regression in the custom format wiring fails fast.

describe('rjsfValidator (customValidator)', () => {
  it('exposes the @rjsf/validator-ajv8 API surface', () => {
    expect(typeof customValidator).toBe('object');
    expect(typeof customValidator.validateFormData).toBe('function');
    expect(typeof customValidator.isValid).toBe('function');
    expect(typeof customValidator.rawValidation).toBe('function');
  });

  describe('standard formats', () => {
    it('accepts a valid email and rejects an invalid one', () => {
      const schema = { type: 'string', format: 'email' };
      const ok = customValidator.rawValidation(schema, 'devs@meshery.io');
      expect(ok.errors).toBeFalsy();
      const bad = customValidator.rawValidation(schema, 'not-an-email');
      expect(bad.errors?.length).toBeGreaterThan(0);
    });

    it('validates date-time strings', () => {
      const schema = { type: 'string', format: 'date-time' };
      const ok = customValidator.rawValidation(schema, '2025-05-14T12:00:00Z');
      expect(ok.errors).toBeFalsy();
      const bad = customValidator.rawValidation(schema, '14/05/2025');
      expect(bad.errors?.length).toBeGreaterThan(0);
    });
  });

  describe('custom kubernetes/openapi formats', () => {
    it('accepts int32 values in range and rejects out-of-range integers', () => {
      const schema = { type: 'number', format: 'int32' };
      expect(customValidator.rawValidation(schema, 100).errors).toBeFalsy();
      expect(customValidator.rawValidation(schema, -2_147_483_648).errors).toBeFalsy();
      expect(customValidator.rawValidation(schema, 2_147_483_647).errors).toBeFalsy();
      expect(
        customValidator.rawValidation(schema, 2_147_483_648).errors?.length ?? 0,
      ).toBeGreaterThan(0);
      expect(customValidator.rawValidation(schema, 3.14).errors?.length ?? 0).toBeGreaterThan(0);
    });

    it('accepts int64 values within JS safe-integer bounds', () => {
      const schema = { type: 'number', format: 'int64' };
      expect(customValidator.rawValidation(schema, Number.MAX_SAFE_INTEGER).errors).toBeFalsy();
      expect(customValidator.rawValidation(schema, Number.MIN_SAFE_INTEGER).errors).toBeFalsy();
      expect(customValidator.rawValidation(schema, 1.5).errors?.length ?? 0).toBeGreaterThan(0);
    });

    it('accepts int-or-string for both strings and integers', () => {
      // int-or-string is a typeless format - applied via a property with no fixed type
      const schema = {
        type: 'object',
        properties: {
          value: { format: 'int-or-string' },
        },
      };
      expect(customValidator.rawValidation(schema, { value: 'auto' }).errors).toBeFalsy();
      expect(customValidator.rawValidation(schema, { value: 42 }).errors).toBeFalsy();
    });
  });

  describe('object-shape validation via @rjsf API', () => {
    it('reports missing required properties through validateFormData', () => {
      const schema = {
        type: 'object',
        required: ['name'],
        properties: { name: { type: 'string' } },
      };
      const result = customValidator.validateFormData({}, schema);
      expect(result.errors.length).toBeGreaterThan(0);
    });

    it('reports an empty errors list when data satisfies the schema', () => {
      const schema = {
        type: 'object',
        required: ['name'],
        properties: { name: { type: 'string' } },
      };
      const result = customValidator.validateFormData({ name: 'demo' }, schema);
      expect(result.errors).toEqual([]);
    });

    it('isValid returns true/false based on raw validation', () => {
      const schema = { type: 'number', format: 'int32' };
      expect(customValidator.isValid(schema, 1, schema)).toBe(true);
      expect(customValidator.isValid(schema, 5_000_000_000, schema)).toBe(false);
    });
  });
});
