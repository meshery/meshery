import _ from 'lodash';
import * as yaml from 'js-yaml';
import { PROMPT_VARIANTS } from '@sistent/sistent';
import { encodeDesignFile, getUnit8ArrayDecodedFile } from '../../../utils/utils';
import { FILE_OPS } from '../../../utils/Enum';
import { EVENT_TYPES } from '../../../lib/event-types';
import downloadContent from '../../../utils/fileDownloader';
import { updateProgress } from '@/store/slices/mesheryUi';
import { buildImportDesignRequestBody } from '../import-design-request';
import { ACTION_TYPES } from './MesheryPatterns.constants';

/**
 * Factory that returns the full set of CRUD / publish / clone / download
 * handlers used by the Designs page. Behavior is identical to the
 * original inline definitions in MesheryPatterns.tsx — the same RTK
 * mutations are called with the same payloads, the same notifications
 * are emitted, and the same modal/notification side-effects fire.
 */
export function createPatternsActions(deps) {
  const {
    // RTK mutations
    clonePattern,
    publishCatalog,
    unpublishCatalog,
    deletePattern,
    deletePatternFile,
    importPattern,
    updatePattern,
    uploadPatternFile,
    deployPatternMutation,
    undeployPatternMutation,
    evaluateRelationships,
    // refs / state
    modalRef,
    meshModels,
    infoModal,
    publishModal,
    user,
    // setters
    setImportModal,
    setPublishModal,
    setSelectedRowData,
    setInfoModal,
    // helpers
    notify,
    sistentInfoModal,
    getPatterns,
  } = deps;

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });

    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
    });
  };

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  const handleDeploy = async ({ design, selectedK8sContexts }) => {
    updateProgress({ showProgress: true });
    await deployPatternMutation({
      patternFile: encodeDesignFile(design),
      patternId: design.id,
      selectedK8sContexts,
    });
    updateProgress({ showProgress: false });
  };

  const handleUndeploy = async ({ design, selectedK8sContexts }) => {
    updateProgress({ showProgress: true });
    await undeployPatternMutation({
      patternFile: encodeDesignFile(design),
      patternId: design.id,
      selectedK8sContexts,
    });
    updateProgress({ showProgress: false });
  };

  const handleUploadImport = () => {
    setImportModal({
      open: true,
    });
  };

  const handleUploadImportClose = () => {
    setImportModal({
      open: false,
    });
  };

  const handleInfoModalClose = () => {
    sistentInfoModal.closeModal();
    setInfoModal({
      open: false,
    });
  };

  const handleInfoModal = (pattern) => {
    sistentInfoModal.openModal({
      title: pattern.name,
    });

    setInfoModal({
      open: true,
      ownerID: pattern.userId,
      selectedResource: pattern,
    });
  };

  const handleUnpublishModal = (ev, pattern) => {
    return async () => {
      let response = await modalRef.current.show({
        title: `Unpublish Catalog item?`,
        subtitle: `Are you sure you want to unpublish ${pattern?.name}?`,
        variant: PROMPT_VARIANTS.DANGER,
        primaryOption: 'UNPUBLISH',
        showInfoIcon:
          "Unpublishing a catalog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want. Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the [documentation](https://docs.meshery.io/concepts/designs).",
      });
      if (response === 'UNPUBLISH') {
        updateProgress({ showProgress: true });
        unpublishCatalog({
          unpublishBody: JSON.stringify({ id: pattern?.id }),
        })
          .unwrap()
          .then(() => {
            updateProgress({ showProgress: false });
            notify({
              message: `Design Unpublished`,
              event_type: EVENT_TYPES.SUCCESS,
            });
          })
          .catch((error) => {
            updateProgress({ showProgress: false });
            handleError(ACTION_TYPES.UNPUBLISH_CATALOG)(error);
          });
      }
    };
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      pattern: {},
      name: '',
    });
  };

  const handlePublish = (formData) => {
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: infoModal.selectedResource?.id,
      catalogData: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishCatalog({
      publishBody: JSON.stringify(payload),
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (user.roleNames.includes('admin')) {
          notify({
            message: `${publishModal?.name} Design Published`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'Design queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch((error) => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.PUBLISH_CATALOG)(error);
      });
  };

  function handleClone(patternID, name) {
    updateProgress({ showProgress: true });
    clonePattern({
      body: JSON.stringify({ name: name + ' (Copy)' }),
      patternID: patternID,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" Design cloned`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `Failed to clone "${name}" Design`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  }

  async function showModal(count, patterns) {
    let response = await modalRef.current.show({
      title: `Delete ${count ? count : ''} Design${count > 1 ? 's' : ''}?`,

      subtitle: `Are you sure you want to delete the ${patterns} design${count > 1 ? 's' : ''}?`,
      variant: PROMPT_VARIANTS.DANGER,
      primaryOption: 'DELETE',
    });
    return response;
  }

  async function handleSubmit({ data, id, name, type, metadata, catalogData }) {
    updateProgress({ showProgress: true });
    if (type === FILE_OPS.DELETE) {
      const response = await showModal(1, name);
      if (response == 'CANCEL') {
        updateProgress({ showProgress: false });
        return;
      }
      deletePatternFile({
        id: id,
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" Design deleted`, event_type: EVENT_TYPES.SUCCESS });
          getPatterns();
          resetSelectedRowData()();
        })
        .catch((error) => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.DELETE_PATTERN)(error);
        });
    }

    if (type === FILE_OPS.UPDATE) {
      const design = yaml.load(data);

      updatePattern({
        updateBody: JSON.stringify({
          id,
          name,
          designFile: design,
          catalogData,
        }),
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" Design updated`, event_type: EVENT_TYPES.SUCCESS });
        })
        .catch((error) => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPDATE_PATTERN)(error);
        });
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body;
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          patternData: {
            name: metadata?.name || name,
            patternFile: getUnit8ArrayDecodedFile(data),
            catalogData,
          },
          save: true,
        });
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({
          url: data,
          save: true,
          name: metadata?.name || name,
          catalogData,
        });
      }
      uploadPatternFile({
        uploadBody: body,
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch((error) => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_PATTERN)(error);
        });
    }
  }

  async function handleImportDesign(data) {
    updateProgress({ showProgress: true });
    const { name } = data;

    const importRequest = await buildImportDesignRequestBody(data);
    if ('errorMessage' in importRequest) {
      updateProgress({ showProgress: false });
      notify({
        message: importRequest.errorMessage,
        event_type: EVENT_TYPES.ERROR,
      });
      return;
    }

    importPattern({
      importBody: importRequest.requestBody,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" design uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        setImportModal((prev) => ({ ...prev, open: false }));
        getPatterns();
      })
      .catch((error) => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.UPLOAD_PATTERN)(error);
      });
  }

  function deletePatterns(patterns) {
    const jsonPatterns = JSON.stringify(patterns);

    updateProgress({ showProgress: true });
    deletePattern({
      deleteBody: jsonPatterns,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        setTimeout(() => {
          notify({
            message: `${patterns.patterns.length} Designs deleted`,
            event_type: EVENT_TYPES.SUCCESS,
          });
          getPatterns();
          resetSelectedRowData()();
        }, 1200);
      })
      .catch((error) => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.DELETE_PATTERN)(error);
      });
  }

  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };

  function handleEvaluateRelationship(pattern) {
    updateProgress({ showProgress: true });
    try {
      const design =
        typeof pattern.patternFile === 'string'
          ? yaml.load(pattern.patternFile)
          : pattern.patternFile;

      evaluateRelationships({
        body: {
          design: { ...design, relationships: [] },
          options: { returnDiffOnly: false, enableTrace: false },
        },
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({
            message: `"${pattern.name}" design evaluated`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        })
        .catch((error) => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.EVALUATE_RELATIONSHIP)(error);
        });
    } catch (error) {
      updateProgress({ showProgress: false });
      handleError(ACTION_TYPES.EVALUATE_RELATIONSHIP)(error);
    }
  }

  return {
    handleError,
    resetSelectedRowData,
    handleDeploy,
    handleUndeploy,
    handleUploadImport,
    handleUploadImportClose,
    handleInfoModalClose,
    handleInfoModal,
    handleUnpublishModal,
    handlePublishModalClose,
    handlePublish,
    handleClone,
    handleSubmit,
    handleImportDesign,
    deletePatterns,
    handleDownload,
    handleEvaluateRelationship,
    showModal,
  };
}
