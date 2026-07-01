import React, { useEffect, useMemo, useState } from 'react';
import { useTheme, PlainCard, BellIcon, DesignIcon } from '@sistent/sistent';
import WidgetErrorFallback from './WidgetErrorFallback';

type DashboardIconProps = {
  fill?: string;
  primaryFill?: string;
  secondaryFill?: string;
  width?: number | string;
  height?: number | string;
};

type LatestBlogsProps = {
  iconsProps?: DashboardIconProps;
};

type Resource = { name: string; link: string; external: true };

const BLOG_FEED_URL = 'https://meshery.io/feed.xml';
const LOADING_RESOURCES = [{ name: 'Loading...' }];
const EMPTY_RESOURCES = [{ name: 'No blog posts found.' }];

const LatestBlogs = ({ iconsProps }: LatestBlogsProps) => {
  const theme = useTheme();
  const [resources, setResources] = useState<Resource[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasError, setHasError] = useState(false);

  useEffect(() => {
    let isActive = true;

    const fetchBlogFeed = async () => {
      try {
        const response = await fetch(BLOG_FEED_URL);
        if (!response.ok) {
          throw new Error(`Unexpected response status: ${response.status}`);
        }
        const text = await response.text();
        const xmlDoc = new DOMParser().parseFromString(text, 'application/xml');
        const items = Array.from(xmlDoc.getElementsByTagName('entry'));

        const nextResources = items.slice(0, 5).map((item) => {
          const title = item.getElementsByTagName('title')[0]?.textContent || 'No Title';
          const link = item.getElementsByTagName('content')[0]?.getAttribute('xml:base') || '#';

          return {
            name: title,
            link,
            external: true as const,
          };
        });

        if (isActive) {
          setResources(nextResources);
        }
      } catch (error) {
        if (isActive) {
          console.error('Error fetching latest blogs:', error);
          setHasError(true);
        }
      } finally {
        if (isActive) {
          setLoading(false);
        }
      }
    };

    fetchBlogFeed();

    return () => {
      isActive = false;
    };
  }, []);

  const cardResources = useMemo(
    () =>
      resources.map((resource) => ({
        ...resource,
        icon: (
          <DesignIcon
            width="15px"
            height="15px"
            primaryFill="currentColor"
            secondaryFill="currentColor"
          />
        ),
      })),
    [resources],
  );

  if (!loading && hasError) {
    return (
      <WidgetErrorFallback
        widgetTitle="Latest Blogs"
        message="Unable to load the latest blog posts. Please try again later."
      />
    );
  }

  return (
    <PlainCard
      resources={
        loading ? LOADING_RESOURCES : cardResources.length > 0 ? cardResources : EMPTY_RESOURCES
      }
      icon={
        <BellIcon
          {...iconsProps}
          fill={theme.palette.icon.default}
          secondaryFill={theme.palette.background.brand.default}
        />
      }
      title="LATEST BLOGS"
    />
  );
};

export default LatestBlogs;
