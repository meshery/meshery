import { Box, NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import type { ComponentProps } from 'react';
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
 *
 * Head is intentionally rendered outside NoSsr so the <title> tag is
 * included in the server-side HTML response (better SEO, no hydration flicker).
 */
export const MesheryPage = ({ title, headTitle, noSuffix, children }: MesheryPageProps) => {
  usePageTitle(title);
  const base = headTitle ?? title;
  const browserTitle = noSuffix ? base : `${base} | Meshery`;

  return (
    <>
      <Head>
        <title>{browserTitle}</title>
      </Head>
      <NoSsr>{children}</NoSsr>
    </>
  );
};

MesheryPage.displayName = 'MesheryPage';

/**
 * Centered, overflow-hidden container used by several list/table pages.
 * Pass sx to extend or override the default layout styles.
 */
export const PageContainer = ({
  children,
  sx,
}: {
  children: React.ReactNode;
  sx?: ComponentProps<typeof Box>['sx'];
}) => <Box sx={{ margin: 'auto', overflow: 'hidden', ...sx }}>{children}</Box>;

PageContainer.displayName = 'PageContainer';
