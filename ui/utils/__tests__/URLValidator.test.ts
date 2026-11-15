import { describe, expect, it } from 'vitest';
import { URLValidator } from '../URLValidator';

describe('URLValidator', () => {
  it('accepts http URLs with a domain', () => {
    expect(URLValidator('http://example.com')).toBe(true);
  });

  it('accepts https URLs with a domain', () => {
    expect(URLValidator('https://example.com')).toBe(true);
    expect(URLValidator('https://example.com/path/to/something')).toBe(true);
  });

  it('accepts nats and tcp protocols', () => {
    expect(URLValidator('nats://example.com:4222')).toBe(true);
    expect(URLValidator('tcp://example.com:9000')).toBe(true);
  });

  it('accepts localhost with no domain suffix', () => {
    expect(URLValidator('http://localhost')).toBe(true);
    expect(URLValidator('http://localhost:9081')).toBe(true);
  });

  it('accepts URLs with an IPv4 address', () => {
    expect(URLValidator('http://127.0.0.1:8080')).toBe(true);
    expect(URLValidator('http://192.168.0.1')).toBe(true);
  });

  it('accepts URLs with query strings and fragments', () => {
    expect(URLValidator('https://example.com/path?x=1&y=2')).toBe(true);
    expect(URLValidator('https://example.com/path#section')).toBe(true);
  });

  it('rejects URLs without a protocol', () => {
    expect(URLValidator('example.com')).toBe(false);
    expect(URLValidator('www.example.com')).toBe(false);
  });

  it('rejects unsupported protocols (ftp, ws, file)', () => {
    expect(URLValidator('ftp://example.com')).toBe(false);
    expect(URLValidator('ws://example.com')).toBe(false);
    expect(URLValidator('file:///home/user')).toBe(false);
  });

  it('rejects numeric dotted hosts, port 0, overlong ports, and double dots', () => {
    expect(URLValidator('http://256.256.256.256')).toBe(false);
    expect(URLValidator('http://example.com:0')).toBe(false);
    expect(URLValidator('http://example.com:65536')).toBe(false);
    expect(URLValidator('http://example..com')).toBe(false);
  });

  it('rejects obviously malformed strings', () => {
    expect(URLValidator('not a url')).toBe(false);
    expect(URLValidator('http:/missing-slash.com')).toBe(false);
  });

  it('handles undefined / null gracefully via optional chaining', () => {
    expect(URLValidator(undefined)).toBe(false);
    expect(URLValidator(null)).toBe(false);
  });

  it('rejects an empty string', () => {
    expect(URLValidator('')).toBe(false);
  });

  it.each([24, 40, 100, 1000])('settles quickly on a %i-character run without ReDoS', (n) => {
    const hostile = `http://${'a'.repeat(n)}!`;
    const started = performance.now();
    expect(URLValidator(hostile)).toBe(false);
    expect(performance.now() - started).toBeLessThan(100);
  });
});
