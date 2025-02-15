import React, { useEffect, useState } from 'react';
import { useTheme, PlainCard, BellIcon } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

const LatestBlogs = (props) => {
  const theme = useTheme();
  const [resources, setResources] = useState([]);
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
            icon: '/static/img/designs.svg',
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
  }, [resources]);

  return (
    <UsesSistent>
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
    </UsesSistent>
  );
};

export default LatestBlogs;
