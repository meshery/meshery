import React, { useState } from "react";
import {
  Box,
  Button,
  Divider,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Typography,
  Tooltip,
} from "@mui/material";
import { styled } from "@mui/material/styles";
import UploadImport from "@/components/UploadImport";
import ViewSwitch from "@/components/ViewSwitch";
import SaveIcon from "@mui/icons-material/Save";
import DeleteIcon from "@mui/icons-material/Delete";
import FullscreenIcon from "@mui/icons-material/Fullscreen";
import FullscreenExitIcon from "@mui/icons-material/FullscreenExit";
import MesheryPatternGrid from "./MesheryPatternGrid";
import MesheryPatternTable from "./MesheryPatternTable";
import AddCircleOutlineRoundedIcon from "@mui/icons-material/AddCircleOutlineRounded";
import EmptyState from "@/components/EmptyStateComponent";
import { useTheme } from "@mui/system";
import PatternConfigurator from "./PatternConfigurator";
function resetSelectedPattern() {
  return { show: false, pattern: null };
}

function Mesherypatterns({ user }) {
  const theme = useTheme();

  const CustomBox = styled(Box)(({ theme }) => ({
    justifySelf: "flex-end",
    marginLeft: "auto",
    paddingLeft: theme.spacing(2),
  }));

  const [patterns, setpatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());
  const [modalOpen, setModalOpen] = useState({
    open: false,
    deploy: false,
    pattern_file: null,
  });
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }
  function handleSubmit({ data, id, name, type }) {
    updateProgress({ showProgress: true });
    if (type === FILE_OPS.DELETE) {
      dataFetch(
        `/api/pattern/${id}`,
        {
          credentials: "include",
          method: "DELETE",
        },
        () => {
          console.log("PatternFile API", `/api/pattern/${id}`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
          resetSelectedRowData()();
        },
        handleError(ACTION_TYPES.DELETE_PATTERN)
      );
    }

    if (type === FILE_OPS.UPDATE) {
      dataFetch(
        `/api/pattern`,
        {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({ pattern_data: { id, pattern_file: data }, save: true }),
        },
        () => {
          console.log("PatternFile API", `/api/pattern`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError(ACTION_TYPES.UPDATE_PATTERN)
      );
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body;
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          pattern_data: {
            name,
            pattern_file: data,
          },
          save: true,
        });
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ url: data, save: true });
      }
      dataFetch(
        `/api/pattern`,
        {
          credentials: "include",
          method: "POST",
          body,
        },
        () => {
          console.log("PatternFile API", `/api/pattern`);
          updateProgress({ showProgress: false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError(ACTION_TYPES.UPLOAD_PATTERN)
      );
    }
  }

  function fetchpatterns(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress: true });

    dataFetch(
      `/api/pattern${query}`,
      { credentials: "include" },
      (result) => {
        console.log("patternFile API", `/api/pattern${query}`);
        updateProgress({ showProgress: false });
        if (result) {
          setpatterns(result.patterns || []);
          setCount(result.total_count || 0);
        }
      },
      // handleError
      handleError(ACTION_TYPES.FETCH_patternS)
    );
  }

  const handleModalClose = () => {
    setModalOpen({
      open: false,
      pattern_file: null,
    });
  };

  const handleModalOpen = (app_file, isDeploy) => {
    setModalOpen({
      open: true,
      deploy: isDeploy,
      pattern_file: app_file,
    });
  };

  function TooltipIcon({ children, onClick, title }) {
    return (
      <Tooltip title={title} placement="top" arrow interactive>
        <IconButton onClick={onClick}>{children}</IconButton>
      </Tooltip>
    );
  }

  function YAMLEditor({ pattern, onClose, onSubmit }) {
    const [yaml, setYaml] = useState("");
    const [fullScreen, setFullScreen] = useState(false);

    const toggleFullScreen = () => {
      setFullScreen(!fullScreen);
    };

    return (
      <Dialog
        onClose={onClose}
        aria-labelledby="pattern-dialog-title"
        open
        maxWidth="md"
        fullScreen={fullScreen}
        fullWidth={!fullScreen}
      >
        <DialogTitle disableTypography id="pattern-dialog-title" sx={{ background: "#fff", color: "#000000" }}>
          <Typography variant="h6">{pattern.name}</Typography>
          <TooltipIcon title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"} onClick={toggleFullScreen}>
            {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
          </TooltipIcon>
          {/* <TooltipIcon title="Exit" onClick={onClose}>
              <CloseIcon />
            </TooltipIcon> */}
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent></DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <Tooltip title="Update Pattern">
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={() =>
                onSubmit({
                  data: yaml,
                  id: pattern.id,
                  name: pattern.name,
                  type: FILE_OPS.UPDATE,
                })
              }
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Pattern">
            <IconButton
              aria-label="Delete"
              color="primary"
              onClick={() =>
                onSubmit({
                  data: yaml,
                  id: pattern.id,
                  name: pattern.name,
                  type: FILE_OPS.DELETE,
                })
              }
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </DialogActions>
      </Dialog>
    );
  }

  return (
    <>
      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()} />
      )}
      {selectedPattern.show && (
        <PatternConfigurator onSubmit={handleSubmit} show={setSelectedPattern} pattern={selectedPattern.pattern} />
      )}
      {!selectedPattern.show && (patterns.length > 0 || viewType === "table") && (
        <Box sx={{ display: "flex" }}>
          <Button
            aria-label="Add Pattern"
            variant="contained"
            color="primary"
            size="large"
            sx={{ marginBottom: theme.spacing(2), marginRight: theme.spacing(2) }}
            onClick={() =>
              setSelectedPattern({
                pattern: { id: null, name: "New Pattern", pattern_file: "name: New Pattern\nservices:" },
                show: true,
              })
            }
          >
            <AddCircleOutlineRoundedIcon sx={{ paddingRight: ".35rem" }} />
            Create Design
          </Button>
          <UploadImport configuration="Designs" />
          <CustomBox>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </CustomBox>
        </Box>
      )}

      {!selectedPattern.show && viewType === "table" && (
        <MesheryPatternTable
          patterns={patterns}
          setSelectedRowData={setSelectedRowData}
          handleModalOpen={handleModalOpen}
          user={user}
        />
      )}
      {!selectedPattern.show && viewType === "grid" && patterns.length === 0 && (
        <EmptyState
          configuration="Designs"
          Button1={
            <Button
              aria-label="Create Design"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              onClick={() =>
                setSelectedPattern({
                  pattern: { id: null, name: "New Pattern", pattern_file: "name: New Pattern\nservices:" },
                  show: true,
                })
              }
              sx={{ marginBottom: theme.spacing(2), marginRight: theme.spacing(2) }}
            >
              <AddCircleOutlineRoundedIcon sx={{ paddingRight: ".35rem" }} />
              Create Design
            </Button>
          }
          Button2={<UploadImport configuration="Designs" />}
        />
      )}
      {!selectedPattern.show && viewType === "grid" && (
        <MesheryPatternGrid
          patterns={patterns}
          //  handleDeploy={handleDeploy}
          //  handleUnDeploy={handleUnDeploy}
          //  handleSubmit={handleSubmit}
          setSelectedPattern={setSelectedPattern}
          selectedPattern={selectedPattern}
          //  pages={Math.ceil(count / pageSize)}
          //  setPage={setPage}
          //  selectedPage={page}
        />
      )}
    </>
  );
}

export default Mesherypatterns;
