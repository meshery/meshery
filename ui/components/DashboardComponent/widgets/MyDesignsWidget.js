import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  CatalogIcon,
  DesignIcon,
  EditIcon,
  GithubIcon,
  useTheme,
  DesignCard,
} from '@layer5/sistent';
import { useState } from 'react';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';

const cardData = [
  {
    id: 1,
    title: 'Create a design from scratch',
    image: <EditIcon width={24} fill={'white'} />,
    redirect: '/extension/meshmap',
  },
  {
    id: 2,
    title: 'Choose a template to start with',
    image: <CatalogIcon width={24} />,
    redirect: '/configuration/catalog',
  },
  {
    id: 3,
    title: 'Import design from github',
    image: <GithubIcon width={24} height={24} />,
    redirect: `${MESHERY_CLOUD_PROD}/connect/github/new`,
  },
];

const MyDesignsWidget = (props) => {
  const [sortOrder, setSortOrder] = useState('updated_at desc');
  const { data: userData } = useGetLoggedInUserQuery();
  const { data: patternsData, isFetching: isPatternsFetching } = useGetUserDesignsQuery({
    expandUser: true,
    page: 0,
    pagesize: 7,
    order: sortOrder,
    user_id: userData?.id,
    metrics: true,
  });
  const theme = useTheme();

  return (
    <>
      <DesignCard
        isPatternsFetching={isPatternsFetching}
        cardData={cardData}
        resources={
          patternsData?.patterns?.map((pattern) => {
            return {
              name: pattern.name,
              timestamp: pattern.updated_at,
              link: `/extension/meshmap?mode=design&design=${pattern.id}`,
              icon: '/static/img/designs.svg',
            };
          }) || []
        }
        icon={
          <DesignIcon
            {...props.iconsProps}
            fill={theme.palette.icon.default}
            primaryFill={theme.palette.background.brand.default}
            secondaryFill={theme.palette.icon.default}
          />
        }
        title="MY DESIGNS"
        actionButton={true}
        href={`${MESHERY_CLOUD_PROD}/catalog/content/my-designs`}
        btnTitle="See All Designs"
        sortOrder={sortOrder}
        setSortOrder={setSortOrder}
      />
    </>
  );
};

export default MyDesignsWidget;
