import { describe, expect, it } from 'vitest';
import { URLValidator } from '../URLValidator';

describe('URLValidator', () => {
  it('accepts http URLs with a domain', () => {
    expect(URLValidator('http://example.com')).not.toBeNull();
  });

  it('accepts https URLs with a domain', () => {
    expect(URLValidator('https://example.com')).not.toBeNull();
    expect(URLValidator('https://example.com/path/to/something')).not.toBeNull();
  });

  it('accepts nats and tcp protocols', () => {
    expect(URLValidator('nats://example.com:4222')).not.toBeNull();
    expect(URLValidator('tcp://example.com:9000')).not.toBeNull();
  });

  it('accepts localhost with no domain suffix', () => {
    expect(URLValidator('http://localhost')).not.toBeNull();
    expect(URLValidator('http://localhost:9081')).not.toBeNull();
  });

  it('accepts URLs with an IPv4 address', () => {
    expect(URLValidator('http://127.0.0.1:8080')).not.toBeNull();
    expect(URLValidator('http://192.168.0.1')).not.toBeNull();
  });

  it('accepts URLs with query strings and fragments', () => {
    expect(URLValidator('https://example.com/path?x=1&y=2')).not.toBeNull();
    expect(URLValidator('https://example.com/path#section')).not.toBeNull();
  });

  it('rejects URLs without a protocol', () => {
    expect(URLValidator('example.com')).toBeFalsy();
    expect(URLValidator('www.example.com')).toBeFalsy();
  });

  it('rejects unsupported protocols (ftp, ws, file)', () => {
    expect(URLValidator('ftp://example.com')).toBeFalsy();
    expect(URLValidator('ws://example.com')).toBeFalsy();
    expect(URLValidator('file:///home/user')).toBeFalsy();
  });

  it('rejects obviously malformed strings', () => {
    expect(URLValidator('not a url')).toBeFalsy();
    expect(URLValidator('http:/missing-slash.com')).toBeFalsy();
  });

  it('handles undefined / null gracefully via optional chaining', () => {
    expect(URLValidator(undefined)).toBeFalsy();
    expect(URLValidator(null)).toBeFalsy();
  });

  it('rejects an empty string', () => {
    expect(URLValidator('')).toBeFalsy();
  });
});
