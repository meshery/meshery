import React, { useEffect, useState } from 'react';
import Head from 'next/head';
import { Grid2 } from '@sistent/sistent';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import { Adapters } from '../../components/extensions';
import DefaultError from '@/components/General/error-404';
import { EVENT_TYPES } from '../../lib/event-types';
import { EXTENSION_NAMES } from '../../utils/Enum';
import { useNotification } from '../../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useDispatch, useSelector } from 'react-redux';
import { toggleCatalogContent, updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';
import MeshMapSignUpCard from './cards/MeshMapSignUpCard';
import MeshMapSnapshotCard from './cards/MeshMapSnapshotCard';
import PerformanceActionCard from './cards/PerformanceActionCard';
import DockerExtensionCard from './cards/DockerExtensionCard';
import HelmKanvasCard from './cards/HelmKanvasCard';
import DesignEmbedCard from './cards/DesignEmbedCard';
import CloudExtensionCard from './cards/CloudExtensionCard';
import KubectlPluginCard from './cards/KubectlPluginCard';
import KubectlMeshSyncCard from './cards/KubectlMeshSyncCard';
import CatalogCard from './cards/CatalogCard';

const Extensions = () => {
  const [catalogContent, setCatalogContent] = useState(true);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [hasAccessToMeshMap, setHasAccessToMeshMap] = useState(false);
  const { notify } = useNotification();
  const [updateUserPref] = useUpdateUserPrefMutation();
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Extensions' }));
  }, []);

  const { capabilitiesRegistry } = useSelector((state) => state.ui);

  const {
    data: userData,
    isSuccess: userDataFetched,
    isError: isUserError,
    error: userError,
  } = useGetUserPrefQuery();

  const handleToggle = () => {
    dispatch(toggleCatalogContent({ catalogVisibility: !catalogContent }));
    setCatalogContent(!catalogContent);
    handleCatalogPreference(!catalogContent);
  };

  const fetchUser = () => {
    if (userDataFetched && userData) {
      setExtensionPreferences(userData?.usersExtensionPreferences);
      setCatalogContent(userData?.usersExtensionPreferences?.catalogContent);
    } else if (isUserError) {
      console.log(userError);
    }
  };

  useEffect(() => {
    fetchUser();
  }, [userData]);

  const handleCatalogPreference = (catalogPref) => {
    let body = Object.assign({}, extensionPreferences);
    body['catalogContent'] = catalogPref;
    updateUserPref({ usersExtensionPreferences: body })
      .unwrap()
      .then(() => {
        notify({
          message: `Catalog Content was ${catalogPref ? 'enab' : 'disab'}led`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((error) => {
        console.log(error);
      });
  };

  useEffect(() => {
    const meshMapExtensionExists = capabilitiesRegistry?.extensions?.navigator?.filter(
      (val) => val.title.toLowerCase() === EXTENSION_NAMES.KANVAS,
    );
    if (typeof meshMapExtensionExists === 'object' && meshMapExtensionExists.length)
      setHasAccessToMeshMap(true);
  }, []);

  return (
    <>
      <React.Fragment>
        <Head>
          <title>Extensions | Meshery</title>
        </Head>
        {CAN(keys.VIEW_EXTENSIONS.action, keys.VIEW_EXTENSIONS.subject) ? (
          <Grid2 container spacing={2} size="grow">
            <MeshMapSnapshotCard githubActionEnabled={false} />
            <PerformanceActionCard githubActionEnabled={false} />
            <MeshMapSignUpCard hasAccessToMeshMap={hasAccessToMeshMap} />
            <HelmKanvasCard />
            <DockerExtensionCard />
            <DesignEmbedCard />
            <CloudExtensionCard />
            <KubectlPluginCard />
            <KubectlMeshSyncCard />
            <CatalogCard catalogContent={catalogContent} handleToggle={handleToggle} />
            <Adapters />
          </Grid2>
        ) : (
          <DefaultError />
        )}
      </React.Fragment>
    </>
  );
};

export default Extensions;
