/**
 * DeployConfirmationModal — deploy / undeploy / validate confirmation dialog.
 *
 * Phase 5.b.2 migration target: composes the shared `Modal` primitive so the
 * dialog inherits the project-standard header / chrome instead of using
 * Sistent's `Modal` directly. The tab-based deploy/undeploy/validate flow,
 * Kubernetes context selector, and dry-run integration are preserved 1:1 so
 * the consumer in `MesheryAdapterPlayComponent` keeps working unchanged.
 *
 * The legacy name `ConfirmationMsg` is preserved as the default export to
 * minimise the consumer-side churn. Dead siblings (`SelectDeploymentTarget_`,
 * `SelectDeploymentTarget`) that were never imported elsewhere have been
 * dropped — the lifecycle's actual `SelectDeploymentTarget` lives in
 * `./SelectDeploymentTarget.tsx`.
 */
import React, { FC, useEffect, useState } from 'react';
import {
  Box,
  Checkbox,
  DialogContent,
  DialogContentText,
  DoneAllIcon,
  DoneIcon,
  RemoveDoneIcon,
  SearchIcon as Search,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@sistent/sistent';
import { useDispatch, useSelector } from 'react-redux';
import { useTheme } from '@/theme';
import { Modal } from '@/components/shared/Modal';
import { errorHandlerGenerator, successHandlerGenerator } from '@/utils/helpers/common';
import { useLazyPingKubernetesQuery } from '@/rtk-query/connection';
import { getK8sConfigIdsFromK8sConfig } from '@/utils/multi-ctx';
import { iconMedium, iconSmall } from '../../../css/icons.styles';
import { RoundedTriangleShape } from '@/assets/icons/shapes/RoundedTriangle';
import RedOctagonSvg from '@/assets/icons/shapes/Octagon';
import PatternIcon from '@/assets/icons/Pattern';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { K8sEmptyState } from '../../shared/EmptyState/K8sContextEmptyState';
import { ACTIONS } from '../../../utils/Enum';
import CAN from '../../../utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { TooltipWrappedConnectionChip } from '../../connections/ConnectionChip';
import { setK8sContexts, updateProgress } from '@/store/slices/mesheryUi';
import {
  ActionButton,
  ActionsRow,
  ContextsContainer,
  DialogSubtitle,
  OctagonContainer,
  OctagonText,
  TabLabelWrapper,
  TriangleContainer,
  TriangleNumber,
} from './DeployConfirmationModal.styles';

interface ConfirmationMsgProps {
  open: boolean;
  handleClose: () => void;
  submit: { deploy: () => void; unDeploy: () => void };
  title?: string;
  validationBody?: string;
  componentCount?: number;
  tab: number;
  errors?: { validationError?: number; deploymentError?: number };
  dryRunComponent?: React.ReactNode;
}

const ConfirmationMsg: FC<ConfirmationMsgProps> = (props) => {
  const {
    open,
    handleClose,
    submit,
    title,
    validationBody,
    componentCount,
    tab,
    errors,
    dryRunComponent,
  } = props;

  const [tabVal, setTabVal] = useState(tab);
  const [context, setContexts] = useState<any[]>([]);
  const { notify } = useNotification();
  const [triggerPing] = useLazyPingKubernetesQuery();
  const { selectedK8sContexts, k8sConfig: k8scontext } = useSelector((state: any) => state.ui);

  const isDisabled = !selectedK8sContexts?.length;
  const dispatch = useDispatch();
  useEffect(() => {
    setTabVal(tab);
    setContexts(k8scontext);
  }, [open]);

  const handleTabValChange = (_event: unknown, newVal: number) => {
    setTabVal(newVal);
  };

  const handleKubernetesClick = async (ctxID: string) => {
    updateProgress({ showProgress: true });
    try {
      await triggerPing(ctxID).unwrap();
      updateProgress({ showProgress: false });
      successHandlerGenerator(notify, 'Kubernetes pinged')();
    } catch (err) {
      updateProgress({ showProgress: false });
      errorHandlerGenerator(notify, 'Kubernetes not pinged')(err);
    }
  };

  const handleSubmit = () => {
    if (!selectedK8sContexts?.length) {
      notify({
        message: 'Please select Kubernetes context(s) before proceeding with the operation',
        event_type: EVENT_TYPES.INFO,
      });
      return;
    }

    if (tabVal === ACTIONS.DEPLOY) {
      submit.deploy();
    } else if (tabVal === ACTIONS.UNDEPLOY) {
      submit.unDeploy();
    }
    handleClose();
  };

  const searchContexts = (search: string) => {
    const term = search.toLowerCase();
    if (term === '') {
      setContexts(k8scontext);
      return;
    }
    setContexts(k8scontext.filter((ctx: any) => ctx.name?.toLowerCase().includes(term)));
  };

  const setContextViewer = (id: string) => {
    if (id === 'all') {
      if (selectedK8sContexts?.includes('all')) {
        dispatch(setK8sContexts({ selectedK8sContexts: [] }));
      } else {
        dispatch(setK8sContexts({ selectedK8sContexts: ['all'] }));
      }
      return;
    }

    if (selectedK8sContexts?.includes(id)) {
      const filteredContexts = selectedK8sContexts.filter((cid: string) => cid !== id);
      dispatch(setK8sContexts({ selectedK8sContexts: filteredContexts }));
    } else if (selectedK8sContexts.length > 0 && selectedK8sContexts[0] === 'all') {
      const allContextIds = getK8sConfigIdsFromK8sConfig(k8scontext);
      dispatch(
        setK8sContexts({ selectedK8sContexts: allContextIds.filter((cid: string) => cid !== id) }),
      );
    } else {
      if (selectedK8sContexts.length === k8scontext.length - 1) {
        dispatch(setK8sContexts({ selectedK8sContexts: ['all'] }));
        return;
      }
      dispatch(setK8sContexts({ selectedK8sContexts: [...selectedK8sContexts, id] }));
    }
  };

  const theme = useTheme();
  const isDeployOrUndeploy = tabVal === ACTIONS.DEPLOY || tabVal === ACTIONS.UNDEPLOY;

  return (
    <Modal
      isOpen={open}
      onClose={handleClose}
      title={title || 'Confirmation'}
      headerIcon={
        <PatternIcon style={{ ...iconMedium }} fill={theme.palette.common.white}></PatternIcon>
      }
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
      actions={
        <ActionsRow>
          {isDeployOrUndeploy ? (
            <>
              <ActionButton onClick={handleClose} variant="contained">
                <Typography variant="body2">CANCEL</Typography>
              </ActionButton>

              <ActionButton disabled variant="contained" color="primary" isDisabled={true}>
                <Typography variant="body2">
                  {tabVal === ACTIONS.UNDEPLOY ? 'UNDEPLOY LATER' : 'DEPLOY LATER'}
                </Typography>
              </ActionButton>

              <ActionButton
                onClick={handleSubmit}
                variant="contained"
                color="primary"
                isUndeploy={tabVal === ACTIONS.UNDEPLOY}
                isDisabled={isDisabled}
                disabled={isDisabled}
                data-cy="deploy-btn-confirm"
              >
                <Typography variant="body2">
                  {tabVal === ACTIONS.UNDEPLOY ? 'UNDEPLOY' : 'DEPLOY'}
                </Typography>
              </ActionButton>
            </>
          ) : (
            <ActionButton onClick={handleClose} variant="contained" color="primary">
              <Typography variant="body2">OK</Typography>
            </ActionButton>
          )}
        </ActionsRow>
      }
    >
      <Tabs
        value={validationBody ? tabVal : tabVal === ACTIONS.DEPLOY ? 1 : 0}
        variant="scrollable"
        scrollButtons="auto"
        indicatorColor="primary"
        textColor="primary"
        centered
      >
        {!!validationBody && (
          <Tab
            data-cy="validate-btn-modal"
            onClick={(event) => handleTabValChange(event, ACTIONS.VERIFY)}
            label={
              <div style={{ display: 'flex', alignItems: 'center' }}>
                <DoneIcon
                  style={{
                    margin: '2px',
                    paddingRight: '2px',
                    ...iconSmall,
                  }}
                  fill={theme.palette.icon.default}
                  fontSize="small"
                />
                <TabLabelWrapper>Validate</TabLabelWrapper>
                {!!errors?.validationError && errors.validationError > 0 && (
                  <TriangleContainer>
                    <RoundedTriangleShape color={theme.palette.warning.main}></RoundedTriangleShape>
                    <TriangleNumber style={errors.validationError > 10 ? { left: '25%' } : {}}>
                      {errors.validationError}
                    </TriangleNumber>
                  </TriangleContainer>
                )}
              </div>
            }
            disabled={
              !CAN(
                Keys.CatalogManagementValidateDesign.id,
                Keys.CatalogManagementValidateDesign.function,
              )
            }
          />
        )}
        <Tab
          disabled={
            !CAN(
              Keys.CatalogManagementUndeployDesign.id,
              Keys.CatalogManagementUndeployDesign.function,
            ) ||
            (CAN(
              Keys.CatalogManagementUndeployDesign.id,
              Keys.CatalogManagementUndeployDesign.function,
            ) &&
              isDisabled)
          }
          data-cy="Undeploy-btn-modal"
          onClick={(event) => handleTabValChange(event, ACTIONS.UNDEPLOY)}
          label={
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <div style={{ margin: '2px', paddingRight: '2px' }}>
                {' '}
                <RemoveDoneIcon
                  style={iconSmall}
                  width="20"
                  height="20"
                  fill={theme.palette.icon.default}
                />
              </div>{' '}
              <TabLabelWrapper>Undeploy</TabLabelWrapper>{' '}
            </div>
          }
        />
        <Tab
          disabled={
            !CAN(
              Keys.CatalogManagementDeployDesign.id,
              Keys.CatalogManagementDeployDesign.function,
            ) ||
            (CAN(
              Keys.CatalogManagementDeployDesign.id,
              Keys.CatalogManagementDeployDesign.function,
            ) &&
              isDisabled)
          }
          data-cy="deploy-btn-modal"
          onClick={(event) => handleTabValChange(event, ACTIONS.DEPLOY)}
          label={
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <DoneAllIcon
                style={{ margin: '2px', paddingRight: '2px', ...iconSmall }}
                fill={theme.palette.icon.default}
                fontSize="small"
              />
              <TabLabelWrapper>Deploy</TabLabelWrapper>
              {!!errors?.deploymentError && errors.deploymentError > 0 && (
                <OctagonContainer>
                  <RedOctagonSvg fill={theme.palette.error.main}></RedOctagonSvg>
                  <OctagonText>{errors.deploymentError}</OctagonText>
                </OctagonContainer>
              )}
            </div>
          }
        />
      </Tabs>

      {isDeployOrUndeploy && (
        <DialogSubtitle id="alert-dialog-description">
          <div style={{ height: '100%' }}>{dryRunComponent}</div>
          <div>
            <Typography variant="subtitle1" style={{ marginBottom: '0.8rem' }}>
              {' '}
              {componentCount !== undefined ? (
                <>
                  {' '}
                  {componentCount} component{componentCount > 1 ? 's' : ''}{' '}
                </>
              ) : (
                ''
              )}
            </Typography>
            {k8scontext.length > 0 ? (
              <Typography variant="body1">
                <TextField
                  id="search-ctx"
                  label="Search"
                  size="small"
                  variant="outlined"
                  onChange={(event) => searchContexts(event.target.value)}
                  sx={{
                    width: '100%',
                    backgroundColor: theme.palette.action.disabledBackground,
                    margin: '1px 1px 8px ',
                  }}
                  InputProps={{
                    endAdornment: <Search style={iconMedium} />,
                  }}
                />
                {context.length > 0 ? (
                  <Box sx={{ display: 'table' }}>
                    <Checkbox
                      checked={selectedK8sContexts?.includes('all')}
                      onChange={() => setContextViewer('all')}
                      color="primary"
                    />
                    <span style={{ fontWeight: 'bolder' }}>select all</span>
                  </Box>
                ) : (
                  <Typography variant="subtitle1">No Context found</Typography>
                )}

                <ContextsContainer>
                  {context.map((ctx: any) => (
                    <div
                      key={ctx.id}
                      style={{
                        display: 'flex',
                        justifyContent: 'flex-wrap',
                        alignItems: 'center',
                      }}
                    >
                      <Checkbox
                        checked={
                          selectedK8sContexts?.includes(ctx.id) ||
                          (selectedK8sContexts?.length > 0 && selectedK8sContexts[0] === 'all')
                        }
                        onChange={() => setContextViewer(ctx.id)}
                        color="primary"
                      />
                      <TooltipWrappedConnectionChip
                        title={ctx.name}
                        handlePing={() => handleKubernetesClick(ctx.connectionId)}
                        iconSrc={'/static/img/integrations/kubernetes.svg'}
                      />
                    </div>
                  ))}
                </ContextsContainer>
              </Typography>
            ) : (
              <K8sEmptyState />
            )}
          </div>
        </DialogSubtitle>
      )}
      {tabVal === ACTIONS.VERIFY && (
        <DialogContent>
          <DialogContentText>{validationBody}</DialogContentText>
        </DialogContent>
      )}
    </Modal>
  );
};

export default ConfirmationMsg;
