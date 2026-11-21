import { beforeEach, describe, expect, it } from 'vitest';
import { loadStoredKeys, loadStoredOrganization } from '../session';

describe('loadStoredOrganization', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  it('returns null when nothing is persisted', () => {
    expect(loadStoredOrganization()).toBeNull();
  });

  it('returns null and clears the literal "undefined" string', () => {
    sessionStorage.setItem('currentOrg', 'undefined');
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });

  it('returns null and clears the literal "null" string', () => {
    sessionStorage.setItem('currentOrg', 'null');
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });

  it('returns null and clears malformed JSON without throwing', () => {
    sessionStorage.setItem('currentOrg', '{"id": "123"');
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });

  it('returns the parsed organization when it has a string id', () => {
    const org = { id: 'abc-123', name: 'My Org' };
    sessionStorage.setItem('currentOrg', JSON.stringify(org));
    expect(loadStoredOrganization()).toEqual(org);
    expect(sessionStorage.getItem('currentOrg')).toBe(JSON.stringify(org));
  });

  it('returns null and clears an org with a non-string id', () => {
    sessionStorage.setItem('currentOrg', JSON.stringify({ id: 42, name: 'Bogus' }));
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });

  it('returns null and clears an org with a blank id', () => {
    sessionStorage.setItem('currentOrg', JSON.stringify({ id: '   ', name: 'Bogus' }));
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });

  it('returns null and clears a non-object value', () => {
    sessionStorage.setItem('currentOrg', '"hello"');
    expect(loadStoredOrganization()).toBeNull();
    sessionStorage.setItem('currentOrg', '123');
    expect(loadStoredOrganization()).toBeNull();
    expect(sessionStorage.getItem('currentOrg')).toBeNull();
  });
});

describe('loadStoredKeys', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  it('returns null when nothing is persisted', () => {
    expect(loadStoredKeys()).toBeNull();
  });

  it('returns null and clears the literal "undefined" string', () => {
    sessionStorage.setItem('keys', 'undefined');
    expect(loadStoredKeys()).toBeNull();
    expect(sessionStorage.getItem('keys')).toBeNull();
  });

  it('returns null and clears malformed JSON without throwing', () => {
    sessionStorage.setItem('keys', '["missing-bracket", ');
    expect(loadStoredKeys()).toBeNull();
    expect(sessionStorage.getItem('keys')).toBeNull();
  });

  it('returns the persisted keys array when every entry has a string id', () => {
    const keys = [{ id: 'view', function: 'designs' }];
    sessionStorage.setItem('keys', JSON.stringify(keys));
    expect(loadStoredKeys()).toEqual(keys);
    expect(sessionStorage.getItem('keys')).toBe(JSON.stringify(keys));
  });

  it('returns null and clears an array containing an invalid entry', () => {
    sessionStorage.setItem('keys', JSON.stringify([{ id: 'ok', function: 'designs' }, null]));
    expect(loadStoredKeys()).toBeNull();
    expect(sessionStorage.getItem('keys')).toBeNull();
  });

  it('returns null and clears a non-array value', () => {
    sessionStorage.setItem('keys', JSON.stringify({ id: 'view' }));
    expect(loadStoredKeys()).toBeNull();
    expect(sessionStorage.getItem('keys')).toBeNull();
  });
});
