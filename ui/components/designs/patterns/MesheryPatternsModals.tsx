/* eslint-disable react/display-name */

import React from 'react';
import {
  importDesignSchema,
  importDesignUiSchema,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  Modal as SistentModal,
} from '@sistent/sistent';
import { RJSFModalWrapper } from '../../shared/Modal/Modal';
import Pattern from '../../../public/static/img/drawer-icons/pattern_svg';

export const ImportDesignModal = React.memo((props) => {
  const { handleClose, handleImportDesign } = props;

  return (
    <>
      <>
        <SistentModal
          open={true}
          closeModal={handleClose}
          headerIcon={
            <Pattern fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
          }
          maxWidth="sm"
          title="Import Design"
          data-testid="import-design-modal"
        >
          <RJSFModalWrapper
            schema={importDesignSchema}
            uiSchema={importDesignUiSchema}
            handleSubmit={handleImportDesign}
            submitBtnText="Import"
            handleClose={handleClose}
          />
        </SistentModal>
      </>
    </>
  );
});

export const PublishModal = React.memo((props) => {
  const { handleClose, handleSubmit, title } = props;

  return (
    <>
      <>
        <SistentModal
          open={true}
          closeModal={handleClose}
          aria-label="catalog publish"
          title={title}
          headerIcon={
            <Pattern fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
          }
          maxWidth="sm"
        >
          <RJSFModalWrapper
            schema={publishCatalogItemSchema}
            uiSchema={publishCatalogItemUiSchema}
            handleSubmit={handleSubmit}
            submitBtnText="Submit for Approval"
            handleClose={handleClose}
            helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
          />
        </SistentModal>
      </>
    </>
  );
});
