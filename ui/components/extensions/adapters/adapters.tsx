import React, { useEffect, useState } from 'react';
import { isUndefined } from 'lodash';
import { CardContainer, FrontSideDescription, ImageWrapper } from '../../../css/icons.styles';
import { ADAPTER_STATUS, adaptersList } from './constants';
import type { AdaptersListType } from './constants';
import changeAdapterState from '../../graphql/mutations/AdapterStatusMutation';
import { LARGE_6_MED_12_GRID_STYLE } from '../../../css/grid.style';
import { promisifiedDataFetch } from '../../../lib/data-fetch';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { Grid2, Switch, Typography, useTheme } from '@sistent/sistent';
import { updateProgress } from '@/store/slices/mesheryUi';

interface AdapterPayload {
  status: string;
  adapter: string;
  targetPort: string;
}

interface MeshAdapter {
  name: string;
  adapter_location: string;
}

interface SyncResult {
  meshAdapters: MeshAdapter[];
}

const Adapters: React.FC = () => {
  // States.
  const [availableAdapters, setAvailableAdapters] = useState<AdaptersListType>(adaptersList);

  // Hooks.
  const { notify } = useNotification();

  // useEffects.
  useEffect(() => {
    handleAdapterSync();
  }, []);

  const theme = useTheme();

  // Handlers.
  const handleAdapterSync = async (showLoader = true): Promise<void> => {
    showLoader && updateProgress({ showProgress: true });

    promisifiedDataFetch('/api/system/sync', {
      method: 'GET',
      credentials: 'include',
    })
      .then((result: SyncResult) => {
        showLoader && updateProgress({ showProgress: false });

        if (!isUndefined(result)) {
          // Deep copying to avoid mutability.
          const currentAdaptersList = structuredClone(adaptersList);

          result.meshAdapters.forEach((element) => {
            const adapterId = element.name;
            if (adapterId && currentAdaptersList[adapterId]) {
              currentAdaptersList[adapterId].enabled = true;
              currentAdaptersList[adapterId].url = element.adapter_location;
            }
          });
          setAvailableAdapters(currentAdaptersList);
        }
      })
      .catch(() => handleError('Unable to fetch list of adapters.'));
  };

  const handleAdapterDeployment = (
    payload: AdapterPayload,
    msg: string,
    selectedAdapter: AdaptersListType[string],
    _adapterId: string,
  ): void => {
    updateProgress({ showProgress: true });

    changeAdapterState((_response: unknown) => {
      updateProgress({ showProgress: false });

      const actionText = payload.status.toLowerCase();
      notify({
        message: `${selectedAdapter.name} adapter ${actionText}`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    }, payload);
  };

  const handleError =
    (msg: string) =>
      (error: Error | null = null): void => {
        updateProgress({ showProgress: false });
        if (error) {
          notify({
            message: `${msg}: ${error.toString()}`,
            event_type: EVENT_TYPES.ERROR,
          });
        } else {
          notify({
            message: msg,
            event_type: EVENT_TYPES.ERROR,
          });
        }
      };

  const handleToggle = (selectedAdapter: AdaptersListType[string], adapterId: string): void => {
    setAvailableAdapters({
      ...availableAdapters,
      [adapterId]: { ...selectedAdapter, enabled: !selectedAdapter.enabled },
    });
    let payload: AdapterPayload;
    let msg: string;
    if (!selectedAdapter.enabled) {
      payload = {
        status: ADAPTER_STATUS.ENABLED,
        adapter: selectedAdapter.label,
        targetPort: String(selectedAdapter.defaultPort ?? ''),
      };
      msg = 'Unable to deploy adapter';
    } else {
      payload = {
        status: ADAPTER_STATUS.DISABLED,
        adapter: selectedAdapter.label,
        targetPort: String(selectedAdapter.defaultPort ?? ''),
      };
      msg = 'Unable to undeploy adapter';
    }
    handleAdapterDeployment(payload, msg, selectedAdapter, adapterId);
  };

  // Render.
  return (
    <>
      {Object.entries(availableAdapters).map(([adapterId, adapter]) => (
        <Grid2 size={LARGE_6_MED_12_GRID_STYLE} key={adapterId}>
          <CardContainer>
            <Typography variant="h5" component="div">
              Meshery Adapter for {adapter.name}
            </Typography>

            <FrontSideDescription variant="body1">
              <ImageWrapper src={adapter.imageSrc} />
              <div
                style={{
                  display: 'inline',
                  position: 'relative',
                }}
              >
                {adapter.description}
              </div>
            </FrontSideDescription>

            <Grid2
              container
              spacing={2}
              direction="row"
              size="grow"
              justifyContent="space-between"
              alignItems="baseline"
              style={{
                paddingRight: '3rem',
                paddingLeft: '.5rem',
                bottom: '1.5rem',
              }}
            >
              <Typography variant="subtitle2" style={{ fontStyle: 'italic' }}>
                <a
                  href="https://docs.meshery.io/concepts/architecture/adapters"
                  target="_blank"
                  rel="noreferrer"
                  style={{
                    textDecoration: 'none',
                    color: theme.palette.text.brand,
                  }}
                  data-testid={`adapter-docs-${String(adapter.name).toLowerCase()}`}
                >
                  Open Adapter docs
                </a>
              </Typography>

              <div style={{ textAlign: 'right' }}>
                <Switch
                  checked={adapter.enabled}
                  onChange={() => handleToggle(adapter, adapterId)}
                  name="OperatorSwitch"
                  color="primary"
                />
              </div>
            </Grid2>
          </CardContainer>
        </Grid2>
      ))}
    </>
  );
};

export default Adapters;
