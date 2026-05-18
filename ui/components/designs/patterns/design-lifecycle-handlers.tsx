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

/**
 * Builds the four design-lifecycle modal openers (deploy / undeploy /
 * dryrun / validate).
 *
 * Each returned function preserves the exact behavior of the original
 * inline definitions in MesheryPatterns.tsx — same modal title strings,
 * same icons, same stepper / modal-body wrapping, same arg shapes.
 */
export function buildDesignLifecycleHandlers({
  designLifecycleModal,
  designValidationActorRef,
  selectedK8sContexts,
  handleDeploy,
  handleUndeploy,
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

  return { openDeployModal, openUndeployModal, openDryRunModal, openValidateModal };
}
