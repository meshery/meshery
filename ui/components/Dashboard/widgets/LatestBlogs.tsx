import React, { useEffect, useState } from 'react';
import { useTheme, PlainCard, BellIcon, DesignIcon } from '@sistent/sistent';

type LatestBlogsProps = {
  iconsProps?: object;
};

type Resource = { name: string; link?: string; external?: boolean; icon?: React.ReactNode };

const LatestBlogs = (props: LatestBlogsProps) => {
  const theme = useTheme();
  const [resources, setResources] = useState<Resource[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchBlogFeed = async () => {
      try {
        const response = await fetch('https://meshery.io/feed.xml');
        const text = await response.text();
        const xmlDoc = new DOMParser().parseFromString(text, 'application/xml');
        const items = Array.from(xmlDoc.getElementsByTagName('entry'));

        const newResources = items.slice(0, 5).map((item) => {
          const title = item.getElementsByTagName('title')[0]?.textContent || 'No Title';
          const link = item.getElementsByTagName('content')[0]?.getAttribute('xml:base') || '#';
          return {
            name: title,
            link: link,
            external: true,
            icon: (
              <DesignIcon
                width="15px"
                height="15px"
                primaryFill="currentColor"
                secondaryFill="currentColor"
              />
            ),
          };
        });

        setResources(newResources);
      } catch (error) {
        console.error('Error fetching latest blogs:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchBlogFeed();
  }, []);

  return (
    <PlainCard
      resources={loading ? [{ name: 'Loading...' }] : resources}
      icon={
        <BellIcon
          {...props.iconsProps}
          fill={theme.palette.icon.default}
          secondaryFill={theme.palette.background.brand.default}
        />
      }
      title="LATEST BLOGS"
    />
  );
};

export default LatestBlogs;
