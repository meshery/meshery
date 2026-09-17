import React from 'react';
import { ModalBody } from '@sistent/sistent';
import { DoneAll as DoneAllIcon } from '@/assets/icons';
import UndeployIcon from '../../../public/static/img/UndeployIcon';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { UnDeployStepper, DeployStepper } from '../lifecycle/DeployStepper';
import { DryRunDesign } from '../lifecycle/DryRun';
import { ValidateDesign } from '../lifecycle/ValidateDesign';
import { DEPLOYMENT_TYPE } from '../lifecycle/common';
import { parseDesignFile } from '../../../utils/utils';
import { designValidatorCommands } from '../../../machines/validator/designValidator';
import { EVENT_TYPES } from '../../../lib/event-types';

/**
 * Builds the design-lifecycle modal openers (deploy / undeploy /
 * dryrun / validate) and direct (single-click) action executors.
 *
 * Each returned modal opener preserves the exact behavior of the original
 * inline definitions in MesheryPatterns.tsx — same modal title strings,
 * same icons, same stepper / modal-body wrapping, same arg shapes.
 */
export function buildDesignLifecycleHandlers({
  designLifecycleModal,
  designValidationActorRef,
  selectedK8sContexts,
  handleDeploy,
  handleUndeploy,
  notify,
}) {
  const openDeployModal = (e, pattern_file, name) => {
    const design = parseDesignFile(pattern_file);
    e.stopPropagation();
    designLifecycleModal.openModal({
      title: `Deploy design "${name}"`,
      headerIcon: <DoneAllIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <DeployStepper
          handleClose={designLifecycleModal.closeModal}
          validationMachine={designValidationActorRef}
          design={design}
          handleDeploy={handleDeploy}
          deployment_type={DEPLOYMENT_TYPE.DEPLOY}
          selectedK8sContexts={selectedK8sContexts}
        />
      ),
    });
  };

  const openUndeployModal = (e, pattern_file, name) => {
    e.stopPropagation();
    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Undeploy design "${name}"`,
      headerIcon: <UndeployIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <UnDeployStepper
          handleClose={designLifecycleModal.closeModal}
          validationMachine={designValidationActorRef}
          design={design}
          handleUndeploy={handleUndeploy}
          deployment_type={DEPLOYMENT_TYPE.UNDEPLOY}
          selectedK8sContexts={selectedK8sContexts}
        />
      ),
    });
  };

  const openDryRunModal = (e, pattern_file, name) => {
    e.stopPropagation();

    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Dryrun design "${name}"`,
      headerIcon: <DryRunIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <ModalBody style={{ minWidth: '30rem', width: 'auto' }}>
          <DryRunDesign
            handleClose={designLifecycleModal.closeModal}
            validationMachine={designValidationActorRef}
            design={design}
            deployment_type={DEPLOYMENT_TYPE.DEPLOY}
            selectedK8sContexts={selectedK8sContexts}
          />
        </ModalBody>
      ),
    });
  };

  const openValidateModal = (e, pattern_file, name) => {
    e.stopPropagation();

    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Validate design "${name}"`,
      headerIcon: <CheckIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <ModalBody style={{ minWidth: '30rem', width: 'auto' }}>
          <ValidateDesign
            handleClose={designLifecycleModal.closeModal}
            validationMachine={designValidationActorRef}
            design={design}
            deployment_type={DEPLOYMENT_TYPE.DEPLOY}
            selectedK8sContexts={selectedK8sContexts}
          />
        </ModalBody>
      ),
    });
  };

  const directDeploy = async (e, pattern_file, name, id) => {
    e?.stopPropagation?.();
    const design = parseDesignFile(pattern_file);
    if (id) {
      design.id = id;
    }
    await handleDeploy?.({ design, selectedK8sContexts });
    notify?.({
      message: `Deploying design "${name}"`,
      event_type: EVENT_TYPES.INFO,
    });
  };

  const directUndeploy = async (e, pattern_file, name, id) => {
    e?.stopPropagation?.();
    const design = parseDesignFile(pattern_file);
    if (id) {
      design.id = id;
    }
    await handleUndeploy?.({ design, selectedK8sContexts });
    notify?.({
      message: `Undeploying design "${name}"`,
      event_type: EVENT_TYPES.INFO,
    });
  };

  const directDryRun = (e, pattern_file, name) => {
    e?.stopPropagation?.();
    const design = parseDesignFile(pattern_file);
    designValidationActorRef?.send?.(
      designValidatorCommands.dryRunDesignDeployment({
        design,
        k8sContexts: selectedK8sContexts,
      }),
    );
    notify?.({
      message: `Running dry run for design "${name}"`,
      event_type: EVENT_TYPES.INFO,
    });
  };

  const directValidate = (e, pattern_file, name) => {
    e?.stopPropagation?.();
    const design = parseDesignFile(pattern_file);
    designValidationActorRef?.send?.(designValidatorCommands.validateDesignSchema({ design }));
    notify?.({
      message: `Validating design "${name}"`,
      event_type: EVENT_TYPES.INFO,
    });
  };

  return {
    openDeployModal,
    openUndeployModal,
    openDryRunModal,
    openValidateModal,
    directDeploy,
    directUndeploy,
    directDryRun,
    directValidate,
  };
}
