import React, {useState, useEffect} from 'react'
import { Box, Button,} from "@mui/material";
import { styled } from "@mui/material/styles";
import UploadImport from "@/components/UploadImport";
import ViewSwitch from "@/components/ViewSwitch";
import MesheryApplicationGrid from "./ApplicationsGrid"
import ApplicationTable from "./ApplicationTable";
import { useTheme } from "@mui/system";
import YAMLEditor from  "@/components/YamlDialog";
import EmptyState from "@/components/EmptyStateComponent";
import  fetchApplicationsQuery from "@/features/applications/ApplicationsQuery";

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
    
 const [page, setPage] = useState(0);
 const [count, setCount] = useState(0);
 const [search] = useState("");
 const [sortOrder] = useState("");
 const [pageSize, setPageSize] = useState(10)
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
   
    useEffect(() => {
      fetchApplications(page, pageSize, search, sortOrder);
    }, [page, pageSize, search, sortOrder]);

    function fetchApplications(page, pageSize, search, sortOrder) {
      if (!search) search = "";
      if (!sortOrder) sortOrder = "";
       
      fetchApplicationsQuery({
        selector : {
          pageSize : `${pageSize}`,
          page : `${page}`,
          search : `${encodeURIComponent(search)}`,
          order : `${encodeURIComponent(sortOrder)}`,
        },
      }).subscribe({
        next : (res) => {
         let result = res?.fetchApplications;
         if (typeof result !== "undefined") {
          if (result) {
            setApplications(result.applications || []);
            setCount(result.total_count || 0);
            setPage(result.page || 0);
            setPageSize(result.page_size || 0);
          }
        }
    }
  })
}

    
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


  return (
    <div> 
                    {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()}  />
          )}  
      {!selectedApplication.show  &&  (applications.length > 0 || viewType === "table") &&
      <Box sx={{display: "flex"}}>  
       <UploadImport configuration="Applications" />
       <CustomBox >
            <ViewSwitch view={viewType} changeView={setViewType} />
          </CustomBox>
       </Box> 
}
       
    {!selectedApplication.show &&  viewType==="table" &&   
         <ApplicationTable
         applications = {applications}
         setSelectedRowData = {setSelectedRowData}
         user ={user}
         count ={count}
         pageSize={pageSize}
         page={page}
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
             pages={Math.ceil(count / pageSize)}
             setPage={setPage}
             selectedPage={page}
           />
}
{!selectedApplication.show && viewType === "grid" && applications.length === 0 &&
            <EmptyState configuration="Applications" 
            Button1={
              <UploadImport configuration="Applications" />
            }
            />
        } 
  
  </div>
  )
}

export default MesheryApplications