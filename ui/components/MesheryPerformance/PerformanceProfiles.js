import React, {useState, useEffect} from 'react'
import { Box, Button,Typography, } from "@mui/material";
import PerformanceProfileTable from "./PerformanceProfileTable"
import ViewSwitch from "@/components/ViewSwitch";
import AddCircleOutlineRoundedIcon from '@mui/icons-material/AddCircleOutlineRounded';
import  PerformanceProfileGrid from "./PerformanceProfileGrid";
import PerformanceModal from "./PerformanceModal";
import { styled } from "@mui/material/styles";
import CustomModal from "@/components/Modal";
import EmptyState from "@/components/EmptyStateComponent";
import PaperWithoutTitle from "@/components/Paper";
import { useTheme } from "@mui/system";
import fetchPerformanceProfiles from "@/features/performance/graphql/queries/PerformanceProfileQuery";

function PerformanceProfiles() {
  const theme = useTheme();
    const CustomBox = styled(Box)(({theme}) => ({
        justifySelf : "flex-end",
        marginLeft : "auto",
        paddingLeft : theme.spacing(2),
     }))

    const [viewType, setViewType] = useState(
        /**  @type {TypeView} */
        ("grid")
      ); 
      
      const [page, setPage] = useState(0);
      const [search, setSearch] = useState("");
      const [sortOrder, setSortOrder] = useState("");
      const [count, setCount] = useState(0);
      const [testProfiles, setTestProfiles] = useState([]);
      const [profileForModal, setProfileForModal] = useState();
      const [pageSize, setPageSize] = useState(10);
    
        useEffect(() => {
          fetchTestProfiles(page, pageSize, search, sortOrder);
        }, [page, pageSize, search, sortOrder]);

        function fetchTestProfiles(page, pageSize, search, sortOrder) {
          if (!search) search = "";
          if (!sortOrder) sortOrder = "";
    
          fetchPerformanceProfiles({
            selector : {
              pageSize : `${pageSize}`,
              page : `${page}`,
              search : `${encodeURIComponent(search)}`,
              order : `${encodeURIComponent(sortOrder)}`,
            },
          }).subscribe({
            next : (res) => {
              // @ts-ignore
              let result = res?.getPerformanceProfiles;
              if (typeof result !== "undefined") {
                if (result) {
                  setCount(result.total_count || 0);
                  setPageSize(result.page_size || 0);
                  setTestProfiles(result.profiles || []);
                  setPage(result.page || 0);
                }
              }
            },
          });
        }

  return (
    <>
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
              setProfileForModal={setProfileForModal}
              pageSize={pageSize}
              page={page}
              setPage={setPage}
              search={search}
              setSearch={setSearch}
              sortOrder={sortOrder}
              setSortOrder={setSortOrder}
              count={count}
            />
        }
  { viewType==="grid" &&
  <PerformanceProfileGrid
  profiles={testProfiles}
  setProfileForModal={setProfileForModal}
  pages={Math.ceil(count / pageSize)}
              setPage={setPage}
            />
        }

{testProfiles.length == 0 && viewType == "grid" && (
     <EmptyState
     configuration="Performance Profile"
     Button1={
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
  Content= {
    <PaperWithoutTitle sx={{ padding : theme.spacing(10), margin: "0 4rem", }}>
      <PerformanceModal
      profileName={profileForModal?.name}
      // @ts-ignore
      meshName={profileForModal?.service_mesh}
      // @ts-ignore
      url={profileForModal?.endpoints?.[0]}
      // @ts-ignore
      qps={profileForModal?.qps}
      // @ts-ignore
      loadGenerator={profileForModal?.load_generators?.[0]}
      // @ts-ignore
      t={profileForModal?.duration}
      // @ts-ignore
      c={profileForModal?.concurrent_request}
      // @ts-ignore
      reqBody={profileForModal?.request_body}
      // @ts-ignore
      headers={profileForModal?.request_headers}
      // @ts-ignore
      cookies={profileForModal?.request_cookies}
      // @ts-ignore
      contentType={profileForModal?.content_type} />
    </PaperWithoutTitle>
  }
   />
  </>
  )
}

export default PerformanceProfiles