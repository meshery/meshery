import React from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/components/shared/LoadingState/LoadingComponentServer', () => ({
  PureHtmlLoadingScreen: ({ id }: { id: string }) => <div id={id} data-testid="pure-html-loader" />,
}));

vi.mock('next/document', () => {
  return {
    default: class MockDocument extends React.Component<any> {},
    Html: ({ children }: any) => <html lang="en">{children}</html>,
    Head: ({ children }: any) => <div>{children}</div>,
    Main: () => <main id="__next" />,
    NextScript: () => <script id="__NEXT_SCRIPT__" />,
  };
});

import MesheryDocument from '../../pages/_document';

describe('MesheryDocument noscript fallback', () => {
  it('renders a noscript fallback with JavaScript required message and links', () => {
    const doc = new MesheryDocument({ emotionStyleTags: [] } as any);
    const htmlString = renderToStaticMarkup(doc.render());

    const noscriptBlocks = htmlString.match(/<noscript>([\s\S]*?)<\/noscript>/g) ?? [];
    const fallback = noscriptBlocks.find((block) => block.includes('JavaScript is required')) ?? '';
    expect(fallback).toContain('JavaScript is required');
    expect(htmlString).toContain('#PRE_REACT_LOADER');
    expect(htmlString).toContain('display: none !important');
    expect(htmlString).toContain('https://docs.meshery.io');
    expect(htmlString).toContain('https://meshery.io');
    expect(htmlString).toContain('https://github.com/meshery/meshery');
  });
});
