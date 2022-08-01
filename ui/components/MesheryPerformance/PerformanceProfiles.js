import React, {useState} from 'react'
import { Avatar, Box, Button, Divider, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, Paper, TableCell, Typography, Tooltip, TableSortLabel } from "@mui/material";
import PerformanceProfileTable from "./PerformanceProfileTable"
import ViewSwitch from "@/components/ViewSwitch";
import AddCircleOutlineRoundedIcon from '@mui/icons-material/AddCircleOutlineRounded';
import  PerformanceProfileGrid from "./PerformanceProfileGrid";
import PerformanceModal from "./PerformanceModal";
import { styled } from "@mui/material/styles";
import CustomModal from "@/components/Modal";
import EmptyConfigurationList from "@/components/emptylist";

function PerformanceProfiles() {
    
    const CustomBox = styled(Box)(({theme}) => ({
        justifySelf : "flex-end",
        marginLeft : "auto",
        paddingLeft : theme.spacing(2),
     }))

    const [viewType, setViewType] = useState(
        /**  @type {TypeView} */
        ("grid")
      ); 

      const [count, setCount] = useState(0);
      const [testProfiles, setTestProfiles] = useState([]);
      const [profileForModal, setProfileForModal] = useState();
      
      console.log("profileForModal", profileForModal)

     const profiles123 = [
        {
        id: "27d7425d-f8f2-41ae-8f5f-443bad744b55",
        name: "ebd",
        user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
        load_generators: [
        "nighthawk"
        ],
        endpoints: [
        "https://layer5.io/"
        ],
        service_mesh: "linkerd",
        concurrent_request: 4,
        qps: 4,
        duration: "30s",
        created_at: "2022-07-04T06:37:07.233393Z",
        updated_at: "2022-07-04T06:37:07.233403Z"
        },
        {
        id: "54df3a92-dec6-4566-8274-0cc2b8bd62ab",
        name: "sdv",
        user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
        load_generators: [
        "nighthawk"
        ],
        endpoints: [
        "https://layer5.io/blog"
        ],
        service_mesh: "app mesh",
        duration: "30s",
        created_at: "2022-07-04T06:35:23.043739Z",
        updated_at: "2022-07-04T06:35:23.043748Z"
        },
        {
        id: "acbff055-721d-4e57-a363-cdc51f9d43cc",
        name: "ebd",
        user_id: "f714c166-5113-4f52-844c-38f0672b5e60",
        load_generators: [
        "fortio"
        ],
        endpoints: [
        "https://youtu.be/JNoL5CLrY68"
        ],
        service_mesh: "app mesh",
        duration: "30s",
        last_run: "2022-04-01T18:18:17.449935Z",
        total_results: 2,
        created_at: "2022-03-19T08:02:16.125012Z",
        updated_at: "2022-03-19T08:02:16.125022Z"
        }
        ];
     
        const handleClick = () => setTestProfiles(profiles123);
  return (
    <>
    <Button onClick={handleClick}>QWE</Button>
    {(testProfiles.length > 0 || viewType === "table") &&
    <Box sx={{display: "flex", paddingLeft : "1rem"}}>
          <Button aria-label="Add Performance Profile" variant="contained"
          color="primary"
          size="large"
          onClick={() => setProfileForModal({})}
            >
           <AddCircleOutlineRoundedIcon sx={{ paddingRight : ".35rem" }}  />
           Add Performance Profile
        </Button> 
    <CustomBox >
    <ViewSwitch view={viewType} changeView={setViewType} />
  </CustomBox>
  </Box>
}
  { viewType==="table" &&
  <PerformanceProfileTable
              testProfiles={testProfiles}
            />
        }
  { viewType==="grid" &&
  <PerformanceProfileGrid
  profiles={testProfiles}
            />
        }

{testProfiles.length == 0 && viewType == "grid" && (
     <EmptyConfigurationList
     configuration="Performance Profile"
     NewButton={
      <Button
                aria-label="Add Performance Profile"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={() => setProfileForModal({})}
              >
                <AddCircleOutlineRoundedIcon sx={{ paddingRight : ".35rem" }}  />
                <Typography className="addIcon">Add Performance Profile</Typography>
              </Button>
     } />
        )}     
   <CustomModal
   open={profileForModal}
   handleClose={() => {
    setProfileForModal(undefined);
  }}
   Content={
    <div>
      <PerformanceModal />
    </div>
   }
   />
  </>
  )
}

export default PerformanceProfiles