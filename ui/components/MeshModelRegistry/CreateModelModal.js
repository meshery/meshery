import React from 'react';
import { Modal } from '@layer5/sistent';
import UrlStepper from './Stepper/UrlStepper';
import { updateProgress } from 'lib/store';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';

const CreateModelModal = ({ isCreateModalOpen, setIsCreateModalOpen }) => {
  const [importModelReq] = useImportMeshModelMutation();
  const handleGenerateModal = async (data) => {
    const { url, model, register } = data;
    const requestBody = {
      importBody: {
        url: url,
        model: model,
      },
      uploadType: 'url',
      register: register,
    };

    updateProgress({ showProgress: true });
    await importModelReq({ importBody: requestBody });
    updateProgress({ showProgress: false });
  };

  return (
    <Modal
      maxWidth="sm"
      open={isCreateModalOpen}
      closeModal={() => setIsCreateModalOpen(false)}
      title="Create Model"
      style={{
        zIndex: 1500,
      }}
    >
      <UrlStepper
        handleGenerateModal={handleGenerateModal}
        handleClose={() => setIsCreateModalOpen(false)}
      />
    </Modal>
  );
};

export default CreateModelModal;
