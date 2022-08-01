import React from 'react'
import {Button, 
  Autocomplete ,
   Tooltip,
  MenuItem,
  IconButton,
  CircularProgress,
  FormControl,
  FormLabel,
  RadioGroup,
  Grid,
  FormControlLabel,
  Radio,
  Paper,
  Divider,
  Typography,
  TextField,
  Link,
  Accordion,
  AccordionSummary,
  AccordionDetails,} from "@mui/material"

  import CloseIcon from '@mui/icons-material/Close';  
  import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
  import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined';
  import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
  import PaperWithoutTitle from "@/components/Paper";
  import { useTheme } from "@mui/system";


  const loadGenerators = ["fortio", "wrk2", "nighthawk"];
  
function PerformanceModal() {
  const theme = useTheme();

  return (
    <PaperWithoutTitle sx={{ padding : theme.spacing(10), }}>
            <Grid container spacing={2} >
              <Grid item xs={12} md={6}>
                <Tooltip title="If a profile name is not provided, a random one will be generated for you.">
                  <TextField
                    id="profileName"
                    name="profileName"
                    label="Profile Name"
                    fullWidth
                    // value={profileName}
                    margin="normal"
                    variant="outlined"
                    // onChange={this.handleChange("profileName")}
                    inputProps={{ maxLength : 300 }}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  select
                  id="meshName"
                  name="meshName"
                  label="Service Mesh"
                  fullWidth
                  // value={meshName === "" && selectedMesh !== "" ? selectedMesh : meshName}
                  margin="normal"
                  variant="outlined"
                  // onChange={this.handleChange("meshName")}
                >
                    <MenuItem>
                  asdfe
                  </MenuItem>

                </TextField>
              </Grid>
              <Grid item xs={12}>
                <TextField
                  required
                  id="url"
                  name="url"
                  label="URL to test"
                  type="url"
                  fullWidth
                  // value={url}
                  // error={urlError}
                  // helperText={urlError ? "Please enter a valid URL along with protocol" : ""}
                  margin="normal"
                  variant="outlined"
                  // onChange={this.handleChange("url")}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="c"
                  name="c"
                  label="Concurrent requests"
                  type="number"
                  fullWidth
                  // value={c}
                  inputProps={{ min : "0", step : "1" }}
                  margin="normal"
                  variant="outlined"
                  // onChange={this.handleChange("c")}
                  InputLabelProps={{ shrink : true }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <TextField
                  required
                  id="qps"
                  name="qps"
                  label="Queries per second"
                  type="number"
                  fullWidth
                  // value={qps}
                  inputProps={{ min : "0", step : "1" }}
                  margin="normal"
                  variant="outlined"
                  // onChange={this.handleChange("qps")}
                  InputLabelProps={{ shrink : true }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                  <Autocomplete
                    required
                    id="t"
                    name="t"
                    freeSolo
                    label="Duration*"
                    fullWidth
                    variant="outlined"
                    // classes={{ root : tError }}
                    // value={tValue}
                    // inputValue={t}
                    // onChange={this.handleDurationChange}
                    // onInputChange={this.handleInputDurationChange}
                    // options={durationOptions}
                    style={{ marginTop : "16px", marginBottom : "8px" }}
                    renderInput={(params) => <TextField {...params} label="Duration*" variant="outlined" />}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} md={12}>
              <Accordion>
                <AccordionSummary  >
                <Typography align="center" color="textSecondary" variant="h6">
                      Advanced Options
                    </Typography>
                </AccordionSummary>
                <AccordionDetails>
                <Grid container spacing={1}>
                      <Grid item xs={12}>
                        <TextField
                          id="headers"
                          name="headers"
                          label='Request Headers e.g. {"host":"bookinfo.meshery.io"}'
                          fullWidth
                          // value={headers}
                          multiline
                          margin="normal"
                          variant="outlined"
                          // onChange={this.handleChange("headers")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label='Request Cookies e.g. {"yummy_cookie":"choco_chip"}'
                          fullWidth
                          // value={cookies}
                          multiline
                          margin="normal"
                          variant="outlined"
                          // onChange={this.handleChange("cookies")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12}>
                        <TextField
                          id="contentType"
                          name="contentType"
                          label="Content Type e.g. application/json"
                          fullWidth
                          // value={contentType}
                          multiline
                          margin="normal"
                          variant="outlined"
                          // onChange={this.handleChange("contentType")}
                        ></TextField>
                      </Grid>
                      <Grid item xs={12} md={12}>
                        <TextField
                          id="cookies"
                          name="cookies"
                          label='Request Body e.g. {"method":"post","url":"http://bookinfo.meshery.io/test"}'
                          fullWidth
                          // value={reqBody}
                          multiline
                          margin="normal"
                          variant="outlined"
                          // onChange={this.handleChange("reqBody")}
                        ></TextField>
                      </Grid>
                    </Grid>
                </AccordionDetails>
                </Accordion> 
              </Grid>  
              <Grid item xs={12} md={4}>
                <FormControl component="loadGenerator">
                  <FormLabel
                    component="loadGenerator"
                    style={{ display : "flex", alignItems : "center", flexWrap : "wrap" }}
                  >
                    Load generator
                    <Tooltip
                    //  title={infoloadGenerators}
                     interactive>
                      <HelpOutlineIcon />
                    </Tooltip>
                  </FormLabel>
                  <RadioGroup
                    aria-label="loadGenerator"
                    name="loadGenerator"
                    // value={loadGenerator}
                    // onChange={this.handleChange("loadGenerator")}
                    row
                  >
                    {loadGenerators.map((lg) => (
                      <FormControlLabel value={lg} control={<Radio color="primary" />} label={lg} />
                    ))}
                  </RadioGroup>
                </FormControl>
              </Grid>
              
            </Grid>
            <div >
              <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  // disabled={disableTest}
                  // onClick={() => this.handleAbort()}
                >
                  Clear
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  // onClick={() => this.submitProfile()}
                  // disabled={disableTest}
                  startIcon={<SaveOutlinedIcon />}
                >
                  Save Profile
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  // onClick={this.handleSubmit}
                  // disabled={blockRunTest || disableTest}
                >
                  Run Test
                </Button>
              </div>
    </PaperWithoutTitle>
  )
}

export default PerformanceModal