import type React from 'react';
import _ from 'lodash';
import { FILE_OPS } from '../../utils/Enum';
import { EVENT_TYPES } from '../../lib/event-types';
import { getUnit8ArrayDecodedFile } from '../../utils/utils';
import { trueRandom } from '../../lib/trueRandom';
import downloadContent from '../../utils/fileDownloader';
import { updateProgress } from '@/store/slices/mesheryUi';
import { ACTION_TYPES } from './Filters.constants';

type Notify = (_args: { message: string; event_type: string; details?: string }) => void;

type ErrorHandler = (_action: { error_msg: string }) => (_error: Error | string) => void;

type SubmitArgs = {
  data: any;
  name: string;
  id?: string;
  type: string;
  metadata?: any;
  catalogData?: any;
};

type CreateSubmitHandlerArgs = {
  notify: Notify;
  handleError: ErrorHandler;
  showmodal: (_count: number, _name?: string) => Promise<string | undefined>;
  resetSelectedRowData: () => () => void;
  deleteFilterFile: (_args: { id: string }) => { unwrap: () => Promise<any> };
  uploadFilterFile: (_args: { uploadBody: any }) => { unwrap: () => Promise<any> };
  updateFilterFile: (_args: { updateBody: any }) => { unwrap: () => Promise<any> };
};

export function createHandleSubmit({
  notify,
  handleError,
  showmodal,
  resetSelectedRowData,
  deleteFilterFile,
  uploadFilterFile,
  updateFilterFile,
}: CreateSubmitHandlerArgs) {
  return async function handleSubmit({ data, name, id, type, metadata, catalogData }: SubmitArgs) {
    updateProgress({ showProgress: true });
    if (type === FILE_OPS.DELETE) {
      const response = await showmodal(1, name);
      if (response !== 'Delete') {
        updateProgress({ showProgress: false });
        return;
      }
      if (!id) {
        updateProgress({ showProgress: false });
        return;
      }
      deleteFilterFile({ id })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" filter deleted`, event_type: EVENT_TYPES.SUCCESS });
          resetSelectedRowData()();
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.DELETE_FILTERS);
        });
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      // todo: remove this
      let body: any = { save: true };
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          ...body,
          filterData: { filterFile: data, name: metadata.name },
          config: metadata.config,
        });
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url: data, name: metadata.name, config: metadata.config });
      }
      uploadFilterFile({ uploadBody: body })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_FILTERS);
        });
    }

    if (type === FILE_OPS.UPDATE) {
      updateFilterFile({
        updateBody: JSON.stringify({
          filterData: { id, name: name, catalogData },
          config: data,
          save: true,
        }),
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_FILTERS);
        });
    }
  };
}

type CreateHandleDownloadArgs = {
  notify: Notify;
};

export function createHandleDownload({ notify }: CreateHandleDownloadArgs) {
  return (e: React.MouseEvent, id: string, name: string) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      downloadContent({ id, name, type: 'filter' });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" filter downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };
}

type CreateUploadHandlerArgs = {
  handleSubmit: (_args: SubmitArgs) => Promise<void>;
};

export function createUploadHandler({ handleSubmit }: CreateUploadHandlerArgs) {
  return function uploadHandler(
    ev: React.ChangeEvent<HTMLInputElement>,
    _: unknown,
    metadata: any,
  ) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener('load', (event) => {
      let uint8 = new Uint8Array(event.target!.result as ArrayBuffer);
      handleSubmit({
        data: Array.from(uint8) as any,
        name: file?.name || 'meshery_' + Math.floor(trueRandom() * 100),
        type: FILE_OPS.FILE_UPLOAD,
        metadata: metadata,
      });
    });
    reader.readAsArrayBuffer(file);
  };
}

type CreateHandleImportFilterArgs = {
  notify: Notify;
  handleError: ErrorHandler;
  updateFilterFile: (_args: { updateBody: any }) => { unwrap: () => Promise<any> };
  getFilters: () => void;
};

export function createHandleImportFilter({
  notify,
  handleError,
  updateFilterFile,
  getFilters,
}: CreateHandleImportFilterArgs) {
  return function handleImportFilter(data: {
    uploadType: string;
    config: string;
    name: string;
    url: string;
    file: ArrayBuffer | string;
  }) {
    updateProgress({ showProgress: true });
    const { uploadType, name, config, url, file } = data;
    let requestBody: any = null;
    switch (uploadType) {
      case 'File Upload':
        requestBody = JSON.stringify({
          config,
          save: true,
          filter_data: {
            name,
            filter_file: getUnit8ArrayDecodedFile(file),
          },
        });
        break;
      case 'URL Upload':
        requestBody = JSON.stringify({
          config,
          save: true,
          url,
          filter_data: {
            name,
          },
        });
        break;
    }

    updateFilterFile({ updateBody: requestBody })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" filter uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        getFilters();
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.UPLOAD_FILTERS);
      });
  };
}

type CreateHandleUnpublishArgs = {
  canPublishFilter: boolean;
  modalRef: React.MutableRefObject<{ show: (_args: any) => Promise<string> } | null>;
  unpublishFilter: (_args: { unpublishBody: string }) => { unwrap: () => Promise<any> };
  notify: Notify;
  handleError: ErrorHandler;
};

export function createHandleUnpublishModal({
  canPublishFilter,
  modalRef,
  unpublishFilter,
  notify,
  handleError,
}: CreateHandleUnpublishArgs) {
  return (ev: React.MouseEvent, filter: any) => {
    if (canPublishFilter) {
      return async () => {
        let response = await modalRef.current?.show({
          title: `Unpublish Catalog item?`,
          subtitle: `Are you sure that you want to unpublish "${filter?.name}"?`,
          primaryOption: 'Yes',
          showInfoIcon: `Unpublishing a catolog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want.  Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the documentation https://docs.meshery.io/concepts/designs.`,
        });
        if (response === 'Yes') {
          updateProgress({ showProgress: true });
          unpublishFilter({ unpublishBody: JSON.stringify({ id: filter?.id }) })
            .unwrap()
            .then(() => {
              updateProgress({ showProgress: false });
              notify({
                message: `"${filter?.name}" filter unpublished`,
                event_type: EVENT_TYPES.SUCCESS,
              });
            })
            .catch(() => {
              updateProgress({ showProgress: false });
              handleError(ACTION_TYPES.UNPUBLISH_CATALOG);
            });
        }
      };
    }
  };
}

type CreateHandlePublishArgs = {
  meshModels: any[];
  publishModal: { filter: any };
  user: any;
  publishFilter: (_args: { publishBody: string }) => { unwrap: () => Promise<any> };
  notify: Notify;
  handleError: ErrorHandler;
};

export function createHandlePublish({
  meshModels,
  publishModal,
  user,
  publishFilter,
  notify,
  handleError,
}: CreateHandlePublishArgs) {
  return (formData: any) => {
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: publishModal.filter?.id,
      catalogData: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishFilter({ publishBody: JSON.stringify(payload) })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (user.roleNames.includes('admin')) {
          notify({
            message: `${publishModal.filter?.name} filter published to Meshery Catalog`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'filters queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.PUBLISH_CATALOG);
      });
  };
}

type CreateHandleCloneArgs = {
  cloneFilter: (_args: { body: string; filterID: string }) => { unwrap: () => Promise<any> };
  notify: Notify;
  handleError: ErrorHandler;
};

export function createHandleClone({ cloneFilter, notify, handleError }: CreateHandleCloneArgs) {
  return function handleClone(filterID: string, name: string) {
    updateProgress({ showProgress: true });
    cloneFilter({
      body: JSON.stringify({ name: name + ' (Copy)' }),
      filterID: filterID,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" filter cloned`, event_type: EVENT_TYPES.SUCCESS });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.CLONE_FILTERS);
      });
  };
}

type CreateDeleteFilterArgs = {
  deleteFilterFile: (_args: { id: string }) => { unwrap: () => Promise<any> };
  notify: Notify;
  handleError: ErrorHandler;
};

export function createDeleteFilter({
  deleteFilterFile,
  notify,
  handleError,
}: CreateDeleteFilterArgs) {
  return function deleteFilter(id: string) {
    deleteFilterFile({ id: id })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({ message: `Filter deleted`, event_type: EVENT_TYPES.SUCCESS });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.DELETE_FILTERS);
      });
  };
}
