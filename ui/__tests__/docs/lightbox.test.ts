// @vitest-environment node
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { JSDOM } from 'jsdom';
import { describe, expect, it } from 'vitest';

const script = readFileSync(resolve(process.cwd(), '../docs/static/js/lightbox.js'), 'utf8');

describe('documentation lightbox links', () => {
  it.each([
    ['/images/full.PNG?version=2#preview', false, true],
    ['/images/full.png', true, false],
    ['/search?file=diagram.png', false, false],
    ['/guide#diagram.png', false, false],
  ])('handles %s with download=%s', (href, download, opensLightbox) => {
    const dom = new JSDOM('<main><a><img width="100" height="100" src="/thumb.png"></a></main>', {
      url: 'https://docs.example.test/',
      runScripts: 'outside-only',
    });
    const link = dom.window.document.querySelector('a')!;
    link.href = href;
    if (download) link.setAttribute('download', '');
    dom.window.eval(script);
    const event = new dom.window.MouseEvent('click', {
      bubbles: true,
      cancelable: true,
      button: 0,
    });
    link.querySelector('img')!.dispatchEvent(event);

    expect(event.defaultPrevented).toBe(opensLightbox);
    expect(dom.window.document.getElementById('lightbox')!.classList.contains('active')).toBe(
      opensLightbox,
    );
    if (opensLightbox) {
      expect(dom.window.document.getElementById('lightbox-img')!.getAttribute('src')).toBe(
        link.href,
      );
    }
    dom.window.close();
  });
});
