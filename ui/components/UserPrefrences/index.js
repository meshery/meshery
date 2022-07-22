import React, {useState} from "react";
import { styled } from "@mui/material/styles";
import { Box, Paper, Tabs, Tab, Tooltip,FormControl, FormLabel, FormGroup, FormControlLabel, Switch } from "@mui/material";
import SettingsCellIcon from '@mui/icons-material/SettingsCell';
import SettingsRemoteIcon from '@mui/icons-material/SettingsRemote';
import PerformanceIcon from "@/components/Navbar/drawer-icons/performance_svg";
import MesherySettingsPerformanceComponent from "./MesherySettingsPerformanceComponent"
import ExtensionPointSchemaValidator from "@/utils/extensionPointSchemaValidator"

function UserPrefrences () {

  const CustomFormControl = styled(FormControl)(() => ({
    padding : 20,
    border : '1.5px solid #969696',
}))
 
const CustomFormContaainer = styled(Box)(() => ({
  display : 'flex',
  'flex-wrap' : 'wrap',
  'justify-content' : 'space-evenly',
  padding : 50
}))

  const [tabVal, SetTabVal] = useState(0);

 const userPrefs = ExtensionPointSchemaValidator("user_prefs")();

  const handleTabValChange  = (event, newVal) =>{
    SetTabVal(newVal);
  }
    return(
      <>
      <Paper square > 
        <Tabs
            value={tabVal}
            onChange={handleTabValChange}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tooltip title="General preferences" placement="top">
              <Tab
                icon={
                  <SettingsCellIcon />
                }
                
                label={<span>General</span>}
              />
            </Tooltip>
            <Tooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                icon={
                  <PerformanceIcon style={{width: "2.8rem"}} />
                }
                label={<span >Performance</span>}
              />
            </Tooltip>
            {/* NOTE: This tab's appearance is logical hence it must be put at last here! Otherwise added logic will need to be added for tab numbers!*/}
            {userPrefs && 
              <Tooltip title="Remote Provider preferences" placement="top">
                <Tab
                  icon={
                    <SettingsRemoteIcon />
                  }
                  label={<span>Remote Provider</span>}
                />
              </Tooltip>
             } 
          </Tabs>
      </Paper>
      <Paper>
      {tabVal == 0 &&
        <CustomFormContaainer >
          <CustomFormControl component="fieldset" >
            <FormLabel component="legend" sx={{fontSize : 20,}} >Analytics and Improvement Program</FormLabel>
            <FormGroup>
              <FormControlLabel
                key="UsageStatsPreference"
                control={(
                  <Switch
                    // checked={anonymousStats}
                    // onChange={this.handleToggle('anonymousUsageStats')}
                    color="primary"
                    data-cy="UsageStatsPreference"
                  />
                )}
                labelPlacement="end"
                label="Send Anonymous Usage Statistics"
              />
              <FormControlLabel
                key="PerfResultPreference"
                control={(
                  <Switch
                    // checked={perfResultStats}
                    // onChange={this.handleToggle('anonymousPerfResults')}
                    color="primary"
                    data-cy="PerfResultPreference"
                  />
                )}
                labelPlacement="end"
                label="Send Anonymous Performance Results"
              />
            </FormGroup>
          </CustomFormControl>
        </CustomFormContaainer>
      }
      {tabVal === 1 &&
        <MesherySettingsPerformanceComponent />
      }
      {tabVal == 2 && userPrefs &&
      <h1>HI</h1>
      }
    </Paper>
    </>
    )
}

export default UserPrefrences