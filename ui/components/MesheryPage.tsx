import { Box, NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import { usePageTitle } from '@/utils/hooks';

interface MesheryPageProps {
  /** Short title used for the sidebar/nav (dispatched to Redux). */
  title: string;
  /** Full browser-tab title. Defaults to `{title} | Meshery`. */
  headTitle?: string;
  children: React.ReactNode;
}

/**
 * Standard page shell: registers the page title in Redux, sets the browser
 * tab title, and wraps children in NoSsr.
 */
export const MesheryPage = ({ title, headTitle, children }: MesheryPageProps) => {
  usePageTitle(title);
  const browserTitle = headTitle ?? `${title} | Meshery`;

  return (
    <NoSsr>
      <Head>
        <title>{browserTitle}</title>
      </Head>
      {children}
    </NoSsr>
  );
};

/**
 * Centered, overflow-hidden container used by several list/table pages.
 */
export const PageContainer = ({ children }: { children: React.ReactNode }) => (
  <Box sx={{ margin: 'auto', overflow: 'hidden' }}>{children}</Box>
);
