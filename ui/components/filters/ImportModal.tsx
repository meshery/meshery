import React from 'react';
import { Modal as SistentModal, importFilterSchema, importFilterUiSchema } from '@sistent/sistent';
import { RJSFModalWrapper } from '../shared/Modal/Modal';
import Filter from '../../public/static/img/drawer-icons/filter_svg';
import type { ImportModalProps } from './Filters.types';

const ImportModal = React.memo(({ handleClose, handleImportFilter }: ImportModalProps) => {
  return (
    <>
      <SistentModal
        open={true}
        closeModal={handleClose}
        headerIcon={
          <Filter fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
        }
        title="Import Design"
        maxWidth="sm"
      >
        <RJSFModalWrapper
          schema={importFilterSchema}
          uiSchema={importFilterUiSchema}
          handleSubmit={handleImportFilter}
          submitBtnText="Import"
          handleClose={handleClose}
        />
      </SistentModal>
    </>
  );
});

ImportModal.displayName = 'ImportModal';

export default ImportModal;
