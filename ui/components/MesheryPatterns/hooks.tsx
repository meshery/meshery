import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { updateProgress } from '@/store/slices/mesheryUi';
import { useRouter } from 'next/router';
import yaml from 'js-yaml';
import downloadContent from '../../utils/fileDownloader';

export const usePatternActions = (
  modalRef,
  getPatterns,
  clonePattern,
  deletePattern,
  deployPattern,
  undeployPattern,
  unpublishPattern,
) => {
  const { notify } = useNotification();
  const router = useRouter();

  const handleClone = (id, name) => {
    clonePattern({ id })
      .unwrap()
      .then(() => {
        notify({
          message: `"${name}" design cloned successfully`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        getPatterns();
      })
      .catch((err) => {
        notify({
          message: `Failed to clone design: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  const handleOpenInConfigurator = (id) => {
    router.push(`/configuration/designs/configurator?design_id=${id}`);
  };

  const deletePatterns = async (patterns, resetSelectedRowData) => {
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
          resetSelectedRowData && resetSelectedRowData();
        }, 1200);
      })
      .catch((err) => {
        updateProgress({ showProgress: false });
        notify({
          message: `Failed to delete designs: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  const handleUnpublish = async (rowData, setUnpublishModal) => {
    unpublishPattern({ id: rowData.id })
      .unwrap()
      .then(() => {
        notify({
          message: `"${rowData.name}" design unpublished`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        setUnpublishModal && setUnpublishModal({ open: false });
        getPatterns();
      })
      .catch((err) => {
        notify({
          message: `Failed to unpublish design: ${err.error}`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  const handleDesignDownload = (e, design, source_type, params) => {
    e?.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      const id = design.id;
      const name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (err) {
      console.error(err);
      updateProgress({ showProgress: false });
    }
  };

  const handleSubmit = async (
    { data, id, name, type, catalog_data },
    updatePattern,
    uploadPatternFile,
    refetchPatterns,
  ) => {
    updateProgress({ showProgress: true });
    if (type === 'update') {
      const design = yaml.load(data);
      updatePattern({
        updateBody: JSON.stringify({
          id,
          name,
          design_file: design,
          catalog_data,
        }),
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" Design updated`, event_type: EVENT_TYPES.SUCCESS });
          refetchPatterns();
        })
        .catch((err) => {
          updateProgress({ showProgress: false });
          notify({
            message: `Failed to update design: ${err.error}`,
            event_type: EVENT_TYPES.ERROR,
          });
        });
    }
  };

  // Suppress unused param warnings for hooks passed but not yet used in this hook
  void deployPattern;
  void undeployPattern;
  void modalRef;

  return {
    handleClone,
    handleOpenInConfigurator,
    deletePatterns,
    handleUnpublish,
    handleDesignDownload,
    handleSubmit,
  };
};
