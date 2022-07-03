import React, {useState} from "react";
import { styled } from "@mui/material/styles";
import {Autocomplete ,  Grid, Box,Button, TextField, Radio, RadioGroup, Tooltip,FormControl, FormControlLabel } from "@mui/material";
import SaveIcon from '@mui/icons-material/Save';


const loadGenerators = [
    'fortio',
    'wrk2',
    'nighthawk',
  ];

function MesherySettingsPerformanceComponent () {
   
    const [c, SetC] = useState(0);
    const [t, SetT] = useState([]);
    const [tValue, SetTValue] = useState(t);
    const [qps, SetQPS] = useState(0);
    const [gen, SetGen] = useState([]);


   const handleChange = (name) => (event) => {
        if (name === 'qps'){
            SetQPS(parseInt(event.target.value))
        }
        if (name === 'c'){
            SetC(parseInt(event.target.value))
        }
        if (name === 'gen'){
            SetGen((event.target.value))
        }
    }
   const handleDurationChange = (newValue) =>  {
    SetTValue(newValue)
    }
    return(
        <Box sx={{padding : "5rem"}} >
            <label><strong>Performance Load Test Defaults</strong></label>
            <Grid container spacing={3}>
              <Grid item xs={12} lg={4}>
                <TextField
                  required
                  id="c"
                  name="c"
                  label="Concurrent requests"
                  type="number"
                  fullWidth
                  value={c}
                  inputProps={{ min: '0', step: '1' }}
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('c')}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
              <Grid item xs={12} lg={4}>
                <TextField
                  required
                  id="qps"
                  name="qps"
                  label="Queries per second"
                  type="number"
                  fullWidth
                  value={qps}
                  inputProps={{ min: '0', step: '1' }}
                  margin="normal"
                  variant="outlined"
                  onChange={handleChange('qps')}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
              <Grid item xs={12} lg={4}>
                <Tooltip title={"Please use 'h', 'm' or 's' suffix for hour, minute or second respectively."}>
                  <Autocomplete
                    required
                    id="t"
                    name="t"
                    freeSolo
                    label="Duration*"
                    fullWidth
                    variant="outlined"
                    value={tValue}
                    inputValue={t}
                    onChange={handleDurationChange}
                    // onInputChange={this.handleInputDurationChange}
                    // options={durationOptions}
                    style={{ marginTop: '16px', marginBottom: '8px' }}
                    renderInput={(params) => <TextField {...params} label="Duration*" variant="outlined" />}
                  />
                </Tooltip>
              </Grid>
              <Grid item xs={12} lg={4}>
                <FormControl component="loadGenerator">
                  <label><strong>Default Load Generator</strong></label>
                  <RadioGroup aria-label="loadGenerator" name="loadGenerator" 
                  value={gen}
                   onChange={handleChange('gen')}
                    row>
                    {loadGenerators.map((lg) => (
                      <FormControlLabel value={lg} control={<Radio color="primary" />} label={lg} />
                    ))}
                  </RadioGroup>
                </FormControl>
              </Grid>
            </Grid>
              <Box sx={{display: 'flex', justifyContent: 'flex-end',}} >
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                //   onClick={this.handleSubmit}
                //   disabled={blockRunTest}
                >
                  <SaveIcon
                    style={{ marginRight: '3px' }}
                  />
                  Save
                  {/* {blockRunTest ? <CircularProgress size={30} /> : 'Save'} */}
                </Button>
              </Box>
          </Box>
    )

}

export default MesherySettingsPerformanceComponent;