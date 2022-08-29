import React , {useState} from "react";
import { AppBar, Paper, Tooltip, IconButton, Button, Tabs, Tab } from "@mui/material";
import { styled } from "@mui/material/styles";
import CloudIcon from '@mui/icons-material/Cloud';
import InsertChartIcon from '@mui/icons-material/InsertChart';
import MeshConfigComponent from "@/components/SettingsComponents/MeshConfigComponent"
import MeshAdapterConfigComponent from "@/components/SettingsComponents/MeshAdapterConfigComponent"
import PrometheusComponent from "@/components/SettingsComponents/PrometheusComponent"
import GrafanaComponent from "@/components/SettingsComponents/GrafanaComponent"

export default function Settings () {
    
    const [tabVal, setTabVal] = useState(0);
    const [subTabVal, setSubTabVal] = useState(0);
  
    const CustomTabIcon = styled("img")(({ theme }) => ({
      display : 'inline',
      verticalAlign : 'text-top',
      width : theme.spacing(1.75),
      marginLeft : theme.spacing(0.5),
  }));

  const  handleChange = (val) => {
        return (event, newVal) => {
        if (val === 'tabVal') {
            setTabVal(newVal);
        }
        if (val === 'subTabVal') {
            setSubTabVal(newVal);
        }
    }
    }

    return(
        <Paper >
        <Paper square >
          <Tabs
            value={tabVal}
            onChange={handleChange('tabVal')}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tooltip title="Identify your cluster" placement="top">
              <Tab
                icon={
                  <CloudIcon  />
                }
                label="Environment"
                data-cy="tabEnvironment"
              />
            </Tooltip>
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab
                label="Adapters"
                data-cy="tabServiceMeshes"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                icon={
                  <InsertChartIcon  />
                }
                label="Metrics"
                tab="tabMetrics"
              />
            </Tooltip>
          </Tabs>
          </Paper>
        {tabVal === 0 && (
          <MeshConfigComponent />
        )}
        {tabVal === 1 && (
           <MeshAdapterConfigComponent />
        )}
        {tabVal === 2 && (
            <>
              <AppBar position="static" color="default">
                <Tabs
                  value={subTabVal}
                  onChange={handleChange('subTabVal')}
                  indicatorColor="primary"
                  textColor="primary"
                  variant="fullWidth"
                >
                  <Tab label={(
                    <div>
                      Grafana
                      <CustomTabIcon src="/static/img/grafana_icon.svg"  />
                    </div>
                  )}
                  />
                  <Tab label={(
                    <div>
                      Prometheus
                      <CustomTabIcon src="/static/img/prometheus_logo_orange_circle.svg" />
                    </div>
                  )}
                  />
                </Tabs>
              </AppBar>
              {subTabVal === 0 && (
                <GrafanaComponent />   
              )}
              {subTabVal === 1 && (
                 <PrometheusComponent />  
              )}
              </>
          )}
        </Paper>
    )
}