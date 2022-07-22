import React, {useState} from 'react'
import { AppBar, Box, Button, Grid, Dialog, DialogTitle, Divider, DialogContent, DialogActions, Tooltip, Hidden, IconButton, Toolbar, Typography, TableCell, TableSortLabel } from "@mui/material";
import MUIDataTable from "mui-datatables";
// import Moment from "react-moment";
import { styled } from "@mui/material/styles";
import UploadImport from "@/components/UploadImport";
import ViewSwitch from "@/components/ViewSwitch";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import UndeployIcon from "../../public/static/img/UndeployIcon";
import MesheryApplicationGrid from "./ApplicationsGrid"
import { useTheme } from "@mui/system";
import YAMLEditor from  "@/components/YamlDialog";

function resetSelectedApplication() {
  return { show : false, application : null };
}


function MesheryApplications({user}) {
  const theme = useTheme();  
  const CustomBox = styled(Box)(({theme}) => ({
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : theme.spacing(2),
 }))

    const [applications, setApplications] = useState([]);
    const [selectedRowData, setSelectedRowData] = useState(null);
    const [selectedApplication, setSelectedApplication] = useState(resetSelectedApplication());
    const [modalOpen, setModalOpen] = useState({
      open : false,
      deploy : false,
      application_file : null
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

    function fetchApplications(page, pageSize, search, sortOrder) {
      if (!search) search = "";
      if (!sortOrder) sortOrder = "";
  
      const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
        sortOrder
      )}`;
  
      updateProgress({ showProgress : true });
  
      dataFetch(
        `/api/application${query}`,
        { credentials : "include", },
        (result) => {
          console.log("ApplicationFile API", `/api/application${query}`);
          updateProgress({ showProgress : false });
          if (result) {
            setApplications(result.applications || []);
            setCount(result.total_count || 0);
          }
        },
        // handleError
        handleError(ACTION_TYPES.FETCH_APPLICATIONS)
      );
    }

    const applications1 = [
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

   
    const handleClick = () => setApplications(applications1);;
    
    const handleModalClose = () => {
      setModalOpen({
        open : false,
        application_file : null
      });
    }
  
    const handleModalOpen = (app_file, isDeploy) => {
      setModalOpen({
        open : true,
        deploy : isDeploy,
        application_file : app_file
      });
    }
    

    const columns = [
        {
          name : "name",
          label : "Application Name",
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
            // customBodyRender : function CustomBody(value) {
            //   return <Moment format="LLLL">{value}</Moment>;
            // },
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
              const rowData = applications[tableMeta.rowIndex];
              return (
                <>
                  <IconButton
                    title="Deploy"
                    onClick={() => handleModalOpen(rowData.application_file, true)}
                  >
                    <DoneAllIcon data-cy="deploy-button" />
                  </IconButton>
                  <IconButton
                    title="Undeploy"
                    onClick={() => handleModalOpen(rowData.application_file, false)}
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
            text : "application(s) selected"
          }
        },
        onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(applications[meta.rowIndex]),
    }

  return (
    <div> 
                    {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()}  />
          )}
      <Button onClick={handleClick}>QWE</Button>     
      <Box sx={{display: "flex"}}>  

       <UploadImport configuration="Applications" />
       <CustomBox >
            <ViewSwitch view={viewType} changeView={setViewType} />
          </CustomBox>
       </Box> 
       {console.log(applications1.name +"qwert")}
    {!selectedApplication.show &&  viewType==="table" &&   
         <MUIDataTable
    title={<div>Applications</div>}
    data={applications}
    columns={columns}
    options={options}    
  />
}
    {!selectedApplication.show &&  viewType==="grid" &&   
             <MesheryApplicationGrid
            applications={applications}
            //  handleDeploy={handleDeploy}
            //  handleUnDeploy={handleUnDeploy}
            //  handleSubmit={handleSubmit}
             setSelectedApplication={setSelectedApplication}
             selectedApplication={selectedApplication}
            //  pages={Math.ceil(count / pageSize)}
            //  setPage={setPage}
            //  selectedPage={page}
           />
}
  
  </div>
  )
}

export default MesheryApplications