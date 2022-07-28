import React, {useState} from 'react'
import { Avatar, Box, Button, Divider, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, TableCell, Typography, Tooltip, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { styled } from "@mui/material/styles";
import UploadImport from "@/components/UploadImport";
import ViewSwitch from "@/components/ViewSwitch";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import UndeployIcon from "../../public/static/img/UndeployIcon";
import SaveIcon from '@mui/icons-material/Save';
import DeleteIcon from '@mui/icons-material/Delete';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import MesheryPatternGrid from "./MesheryPatternGrid"
import AddCircleOutlineRoundedIcon from '@mui/icons-material/AddCircleOutlineRounded';
import EmptyConfigurationList from "@/components/EmptyConfigurationList";
import { useTheme } from "@mui/system";

function resetSelectedPattern() {
  return { show : false, pattern : null };
}


function Mesherypatterns({user}) {
  const theme = useTheme();  

  const CustomBox = styled(Box)(({theme}) => ({
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : theme.spacing(2),
 }))

    const [patterns, setpatterns] = useState([]);
    const [selectedRowData, setSelectedRowData] = useState(null);
    const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());
    const [modalOpen, setModalOpen] = useState({
      open : false,
      deploy : false,
      pattern_file : null
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

    function fetchpatterns(page, pageSize, search, sortOrder) {
      if (!search) search = "";
      if (!sortOrder) sortOrder = "";

      const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
        sortOrder
      )}`;

      updateProgress({ showProgress : true });

      dataFetch(
        `/api/pattern${query}`,
        { credentials : "include", },
        (result) => {
          console.log("patternFile API", `/api/pattern${query}`);
          updateProgress({ showProgress : false });
          if (result) {
            setpatterns(result.patterns || []);
            setCount(result.total_count || 0);
          }
        },
        // handleError
        handleError(ACTION_TYPES.FETCH_patternS)
      );
    }

    const patterns1 = [
      {
        id: "e7ccec75-bec6-4b28-b450-272aefa8a182",
        name: "IstioFilterPattern.yaml",
        user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
        created_at: "2022-07-14T13:50:46.375252Z",
updated_at: "2022-07-14T13:53:35.479756Z"
        },
        {
        id: "fb43eb24-de45-481a-955c-3916440276eb",
        name: "IstioFilterPattern (1).yaml",
        user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
        created_at: "2022-07-14T13:50:46.375252Z",
updated_at: "2022-07-14T13:53:35.479756Z"
  }
    ];


    const handleClick = () => setpatterns(patterns1);;

    const handleModalClose = () => {
      setModalOpen({
        open : false,
        pattern_file : null
      });
    }

    const handleModalOpen = (app_file, isDeploy) => {
      setModalOpen({
        open : true,
        deploy : isDeploy,
        pattern_file : app_file
      });
    }

    function TooltipIcon({ children, onClick, title }) {
      return (
        <Tooltip title={title} placement="top" arrow interactive >
          <IconButton onClick={onClick}>
            {children}
          </IconButton>
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
        <Dialog onClose={onClose} aria-labelledby="pattern-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
          <DialogTitle disableTypography id="pattern-dialog-title" sx={{background: "#fff", color: "#000000" }}>
             <Typography variant="h6" >
              {pattern.name}
            </Typography>
            <TooltipIcon
              title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"}
              onClick={toggleFullScreen}>
              {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
            </TooltipIcon>
            {/* <TooltipIcon title="Exit" onClick={onClose}>
              <CloseIcon />
            </TooltipIcon> */}
          </DialogTitle>
          <Divider variant="fullWidth" light />
          <DialogContent>

          </DialogContent>
          <Divider variant="fullWidth" light />
          <DialogActions>
            <Tooltip title="Update Pattern">
              <IconButton
                aria-label="Update"
                color="primary"
                onClick={() => onSubmit({
                  data : yaml, id : pattern.id, name : pattern.name, type : FILE_OPS.UPDATE
                })}
              >
                <SaveIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Delete Pattern">
              <IconButton
                aria-label="Delete"
                color="primary"
                onClick={() => onSubmit({
                  data : yaml,
                  id : pattern.id,
                  name : pattern.name,
                  type : FILE_OPS.DELETE
                })}
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
          </DialogActions>
        </Dialog>
      );
    }
    
    const columns = [
        {
          name : "name",
          label : "Name",
          options : {
            filter : false,
            sort : true,
            searchable : true,
            customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
              return (
                <TableCell key={index} onClick={() => sortColumn(index)}>
                  <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                    <b>{column.label}</b>
                  </TableSortLabel>
                </TableCell>
              );
            },
          },
        },
        {
          name : "created_at",
          label : "Upload Timestamp",
          options : {
            filter : false,
            sort : true,
            searchable : true,
            customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
              return (
                <TableCell key={index} onClick={() => sortColumn(index)}>
                  <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                    <b>{column.label}</b>
                  </TableSortLabel>
                </TableCell>
              );
            },
            customBodyRender : function CustomBody(value) {
              return <Moment format="LLLL">{value}</Moment>;
            },
          },
        },
        {
          name : "updated_at",
          label : "Update Timestamp",
          options : {
            filter : false,
            sort : true,
            searchable : true,
            customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
              return (
                <TableCell key={index} onClick={() => sortColumn(index)}>
                  <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                    <b>{column.label}</b>
                  </TableSortLabel>
                </TableCell>
              );
            },
            customBodyRender : function CustomBody(value) {
              return <Moment format="LLLL">{value}</Moment>;
            },
          },
        },
        {
          name : "Actions",
          options : {
            filter : false,
            sort : false,
            searchable : false,
            customHeadRender : function CustomHead({ index, ...column }) {
              return (
                <TableCell key={index}>
                  <b>{column.label}</b>
                </TableCell>
              );
            },
            customBodyRender : function CustomBody(_, tableMeta) {
              const rowData = patterns[tableMeta.rowIndex];
              return (
                <>
              <IconButton onClick={() => setSelectedPattern({ pattern : patterns[tableMeta.rowIndex], show : true })}>
                <Avatar src="/static/img/pattwhite.svg" sx={{ width: 27, height: 27, marginRight: 1 }} />
              </IconButton>
                  <IconButton
                    title="Deploy"
                    onClick={() => handleModalOpen(rowData.pattern_file, true)}
                  >
                    <DoneAllIcon data-cy="deploy-button" />
                  </IconButton>
                  <IconButton
                    title="Undeploy"
                    onClick={() => handleModalOpen(rowData.pattern_file, false)}
                  >
                    <UndeployIcon fill="rgba(0, 0, 0, 0.54)" data-cy="undeploy-button" />
                  </IconButton>
                </>
              );
            },
          },
        },
      ];

      const options = {
        filter : false,
        sort : !(user && user.user_id === "meshery"),
        search : !(user && user.user_id === "meshery"),
        filterType : "textField",
        // responsive : "scrollFullHeight",
        resizableColumns : true,
        serverSide : true,
        rowsPerPageOptions : [10, 20, 25],
        fixedHeader : true,
        print : false,
        download : false,
        textLabels : {
          selectedRows : {
            text : "pattern(s) selected"
          }
        },
        onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(patterns[meta.rowIndex]),
    }

  return (
    <> 
            {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()}  />
        )}
          <Button onClick={handleClick}>QWE</Button>  
             
      {!selectedPattern.show  &&  (patterns.length > 0 || viewType === "table") &&
       <Box sx={{display: "flex"}}>
         
      <Button aria-label="Add Pattern" variant="contained"
          color="primary"
          size="large"
           sx={{marginBottom: theme.spacing(2) , marginRight: theme.spacing(2)}} >
           <AddCircleOutlineRoundedIcon sx={{ paddingRight : ".35rem" }} />
          Create Design
        </Button>    
       <UploadImport configuration="Designs" />
       <CustomBox >
            <ViewSwitch view={viewType} changeView={setViewType} />
          </CustomBox>
       </Box>
      }
    
    {!selectedPattern.show &&  viewType==="table" &&   
         <MUIDataTable
    title={<div>Designs</div>}
    data={patterns}
    columns={columns}
    options={options}    
  />
}
{!selectedPattern.show && viewType === "grid" && patterns.length === 0 &&
            <EmptyConfigurationList configuration="Designs" 
            NewButton={
              <Button
              aria-label="Create Design"
              variant="contained"
              color="primary"
              size="large"
              // @ts-ignore
              sx={{marginBottom: theme.spacing(2) , marginRight: theme.spacing(2)}} >  
              <AddCircleOutlineRoundedIcon  sx={{ paddingRight : ".35rem" }} />
              Create Design
            </Button>
            }
            />
        } 
    {!selectedPattern.show &&  viewType==="grid" &&   
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
}
  </>
  )
}

export default Mesherypatterns