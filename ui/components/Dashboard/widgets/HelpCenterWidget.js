import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { useTheme, PlainCard, DocumentIcon, DesignIcon } from '@sistent/sistent';
import React from 'react';
import { iconMedium } from 'css/icons.styles';

const HelpCenterWidget = (props) => {
  const theme = useTheme();
  const resources = [
    {
      name: 'Cloud Docs',
      link: 'https://docs.layer5.io/cloud',
      external: true,
    },
    {
      name: 'Kanvas Docs',
      link: 'https://docs.layer5.io/kanvas',
      external: true,
    },
    {
      name: 'Support Request',
      link: `${MESHERY_CLOUD_PROD}/support`,
      external: true,
    },
    {
      name: 'Discussion Forum',
      link: 'https://meshery.io/community/#discussion-forums',
      external: true,
    },
    {
      name: 'Slack',
      link: 'https://slack.meshery.io',
      external: true,
    },
  ].map((resource) => ({
    ...resource,
    icon: (
      <DesignIcon
        width="15px"
        height="15px"
        primaryFill="currentColor"
        secondaryFill="currentColor"
      />
    ),
  }));

  return (
    <>
      <PlainCard
        resources={resources}
        icon={
          <DocumentIcon
            {...props.iconsProps}
            fill={theme.palette.icon.default}
            secondaryFill={theme.palette.icon.disabled}
            {...iconMedium}
          />
        }
        title="HELP CENTER"
      />
    </>
  );
};

export default HelpCenterWidget;
