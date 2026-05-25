import React from 'react';
import {
  Modal as SistentModal,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
} from '@sistent/sistent';
import { RJSFModalWrapper } from '../shared/Modal/Modal';
import Filter from '../../public/static/img/drawer-icons/filter_svg';
import type { PublishModalProps } from './Filters.types';

const PublishModal = React.memo(({ handleClose, handleSubmit, title }: PublishModalProps) => {
  return (
    <>
      <SistentModal
        open={true}
        headerIcon={
          <Filter fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
        }
        closeModal={handleClose}
        aria-label="catalog publish"
        title={title}
        maxWidth="sm"
      >
        <RJSFModalWrapper
          schema={publishCatalogItemSchema}
          uiSchema={publishCatalogItemUiSchema}
          submitBtnText="Submit for Approval"
          handleSubmit={handleSubmit}
          helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
          handleClose={handleClose}
        />
      </SistentModal>
    </>
  );
});

PublishModal.displayName = 'PublishModal';

export default PublishModal;
