import { Box, NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import { usePageTitle } from '@/utils/hooks';

interface MesheryPageProps {
  /** Short title used for the sidebar/nav (dispatched to Redux). */
  title: string;
  /** Override the base browser-tab text. Defaults to `title`. `| Meshery` is appended automatically unless `noSuffix` is set. */
  headTitle?: string;
  /** Skip appending `| Meshery` to the browser-tab title (e.g. for the 404 page). */
  noSuffix?: boolean;
  children: React.ReactNode;
}

/**
 * Standard page shell: registers the page title in Redux, sets the browser
 * tab title, and wraps children in NoSsr.
 */
export const MesheryPage = ({ title, headTitle, noSuffix, children }: MesheryPageProps) => {
  usePageTitle(title);
  const base = headTitle ?? title;
  const browserTitle = noSuffix ? base : `${base} | Meshery`;

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
