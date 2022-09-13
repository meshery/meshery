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
import MesheryFilterGrid from "./MesheryFilterGrid";
import MesheryFilterTable from "./MesheryFilterTable";
import AddCircleOutlineRoundedIcon from "@mui/icons-material/AddCircleOutlineRounded";
import EmptyState from "@/components/EmptyStateComponent";
import { useTheme } from "@mui/system";
function resetSelectedFilter() {
  return { show: false, filter: null };
}

const ViewBox = styled(Box)(({ theme }) => ({
  justifySelf: "flex-end",
  marginLeft: "auto",
  paddingLeft: theme.spacing(2),
}));

function MesheryFilters({ user }) {
  const theme = useTheme();

  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);

  const [filters, setFilters] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [modalOpen, setModalOpen] = useState({
    open: false,
    deploy: false,
    filter_file: null,
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

  const handleModalClose = () => {
    setModalOpen({
      open: false,
      filter_file: null,
    });
  };

  const handleModalOpen = (app_file, isDeploy) => {
    setModalOpen({
      open: true,
      deploy: isDeploy,
      filter_file: app_file,
    });
  };

  function TooltipIcon({ children, onClick, title }) {
    return (
      <Tooltip title={title} placement="top" arrow interactive>
        <IconButton onClick={onClick}>{children}</IconButton>
      </Tooltip>
    );
  }

  function YAMLEditor({ filter, onClose, onSubmit }) {
    const [yaml, setYaml] = useState("");
    const [fullScreen, setFullScreen] = useState(false);

    const toggleFullScreen = () => {
      setFullScreen(!fullScreen);
    };

    return (
      <Dialog
        onClose={onClose}
        aria-labelledby="filter-dialog-title"
        open
        maxWidth="md"
        fullScreen={fullScreen}
        fullWidth={!fullScreen}
      >
        <DialogTitle disableTypography id="filter-dialog-title" sx={{ background: "#fff", color: "#000000" }}>
          <Typography variant="h6">{filter.name}</Typography>
          <TooltipIcon title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"} onClick={toggleFullScreen}>
            {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
          </TooltipIcon>
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent></DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <Tooltip title="Update Filter">
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={() =>
                onSubmit({
                  data: yaml,
                  id: filter.id,
                  name: filter.name,
                  type: FILE_OPS.UPDATE,
                })
              }
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Filter">
            <IconButton
              aria-label="Delete"
              color="primary"
              onClick={() =>
                onSubmit({
                  data: yaml,
                  id: filter.id,
                  name: filter.name,
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
        <YAMLEditor filter={selectedRowData} onClose={resetSelectedRowData()} />
      )}

      {!selectedFilter.show && (filters.length > 0 || viewType === "table") && (
        <Box sx={{ display: "flex" }}>
          <Button
            aria-label="Add Filter"
            variant="contained"
            color="primary"
            size="large"
            sx={{ marginBottom: theme.spacing(2), marginRight: theme.spacing(2) }}
          >
            <AddCircleOutlineRoundedIcon sx={{ paddingRight: ".35rem" }} />
            Create Filter
          </Button>
          <UploadImport configuration="Filters" />
          <ViewBox>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </ViewBox>
        </Box>
      )}

      {!selectedFilter.show && viewType === "table" && (
        <MesheryFilterTable
          filters={filters}
          setSelectedRowData={setSelectedRowData}
          handleModalOpen={handleModalOpen}
          user={user}
        />
      )}
      {!selectedFilter.show && viewType === "grid" && filters.length === 0 && (
        <EmptyState
          configuration="Filters"
          Button1={
            <Button
              aria-label="Create Filter"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              sx={{ marginBottom: theme.spacing(2), marginRight: theme.spacing(2) }}
            >
              <AddCircleOutlineRoundedIcon sx={{ paddingRight: ".35rem" }} />
              Create Filters
            </Button>
          }
          Button2={<UploadImport configuration="Filters" />}
        />
      )}
      {!selectedFilter.show && viewType === "grid" && (
        <MesheryFilterGrid
          filters={filters}
          //  handleDeploy={handleDeploy}
          //  handleUnDeploy={handleUnDeploy}
          //  handleSubmit={handleSubmit}
          setSelectedFilter={setSelectedFilter}
          selectedFilter={selectedFilter}
          // pages={Math.ceil(count / pageSize)}
          // setPage={setPage}
          // selectedPage={page}
        />
      )}
    </>
  );
}

export default MesheryFilters;
