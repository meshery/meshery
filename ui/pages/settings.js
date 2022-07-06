import React , {useState} from "react";
import { AppBar, Paper, Tooltip, IconButton, Button, Tabs, Tab } from "@mui/material";
import CloudIcon from '@mui/icons-material/Cloud';
import InsertChartIcon from '@mui/icons-material/InsertChart';
import MeshConfigComponent from "@/components/SettingsComponents/MeshConfigComponent"
import MeshAdapterConfigComponent from "@/components/SettingsComponents/MeshAdapterConfigComponent"

export default function Settings () {
    
    const [tabVal, setTabVal] = useState(0);
    const [subTabVal, setSubTabVal] = useState(0);


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
        <>
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
                // icon={
                //   <FontAwesomeIcon icon={faMendeley} transform={mainIconScale} />
                // }
                label="Service Meshes"
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
            <Tooltip title="Reset System" placement="top">
              <Tab
                // icon={
                //   <FontAwesomeIcon icon={faDatabase} transform={mainIconScale} fixedWidth />
                // }
                label="Reset"
                tab="systemReset"
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
                      <img src="/static/img/grafana_icon.svg"  />
                    </div>
                  )}
                  />
                  <Tab label={(
                    <div>
                      Prometheus
                      <img src="/static/img/prometheus_logo_orange_circle.svg" />
                    </div>
                  )}
                  />
                </Tabs>
              </AppBar>
              {subTabVal === 0 && (
                 <h1>  GrafanaComponent </h1>       
              )}
              {subTabVal === 1 && (
                 <h1>    PrometheusComponent </h1>   
              )}
              </>
          )}
        {/* {tabVal === 3 && (
          <TabContainer>
            <div className={classes.container}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={this.handleResetDatabase()}
                className={classes.DBBtn}
                data-cy="btnResetDatabase"

              >
                <Typography> System Reset </Typography>
              </Button>
            </div>
          </TabContainer>
        )} */}

        {/* {backToPlay} */}
        {/* <PromptComponent ref={this.systemResetRef} /> */}
        </>
    )
}