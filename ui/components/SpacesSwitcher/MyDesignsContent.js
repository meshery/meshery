//@ts-check
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import React, { useState } from 'react';
import MainDesignsContent from './MainDesignsContent';
import { VISIBILITY } from '@/utils/Enum';

const MyDesignsContent = () => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE, VISIBILITY.PUBLISHED];

  const [visibility, setVisibility] = useState(visibilityItems);
  const [searchQuery, setSearchQuery] = useState('');
  const [modified, setModified] = useState('updated_at desc');
  const handleModifiedChange = (event) => {
    setPage(0);
    setModified(event.target.value);
  };
  const handleVisibilityChange = (event) => {
    const value = event.target.value;
    setVisibility(typeof value === 'string' ? value.split(',') : value);
    setPage(0);
  };
  const onSearchChange = (e) => {
    setPage(0);
    setSearchQuery(e.target.value);
  };

  const [page, setPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: page,
      pagesize: 10,
      order: 'updated_at desc',
      user_id: currentUser?.id,
      metrics: true,
    },
    {
      skip: !currentUser?.id,
    },
  );
  const hasMore = designsData?.total_count > designsData?.page_size * (designsData?.page + 1);
  const total_count = designsData?.total_count || 0;
  return (
    <>
      <MainDesignsContent
        designs={designsData?.patterns}
        isFetching={isFetching}
        isLoading={isLoading}
        setPage={setPage}
        hasMore={hasMore}
        total_count={total_count}
      />
    </>
  );
};

export default MyDesignsContent;
