import React from 'react'
import { Tooltip, IconButton, TextField, Button, Grid, Typography } from '@material-ui/core';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import GitHubIcon from '@material-ui/icons/GitHub';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import { URLValidator } from '../utils/URLValidator';

const getMuiTheme = () => createTheme({
  palette : {
    primary : {
      main : "#607d8b"
    }
  },
  overrides : {
    MuiGrid : {
      input : {
        color : '#607d8b'
      }
    },
  }
})


const styles = makeStyles((theme) => ({
  paper : {
    position : 'absolute',
    width : 600,
    backgroundColor : theme.palette.background.paper,
    border : '0px solid #000',
    boxShadow : theme.shadows[5],
    padding : theme.spacing(2, 4, 3),
    top : '50%',
    left : '50%',
    transform : `translate(-50%, -50%)`,
    borderRadius : 10,
  },
  grid : {
    width : '100%'
  },
  GridTitle : {
    textAlign : 'center'
  },
  formControl : {
    margin : theme.spacing(1),
    minWidth : 120,
  },
  ghGroups : {
    padding : theme.spacing(0, 13, 0)
  },
  optionImport : {
    textAlign : 'center'
  }

}));
import GenericModal from "./GenericModal";


const GHImport = ({ onSubmit }) => {
  const classes = styles();
  const [open, setOpen] = React.useState(false);
  const [input, setInput] = React.useState('');
  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };
  const validURL = (str) => {
    return URLValidator(str);
  }
  const handleError = (input) => {
    console.log(input + ' is not valid a valid repository')
  }
  const handleSubmit = () => {
    validURL(input) ? onSubmit(input) : handleError(input);
    handleClose()
  }

  return (
    <React.Fragment>
      <label htmlFor="github-import-button">

        <Tooltip title="Github Import">
          <IconButton aria-label="Github-Import" component="span" onClick={handleOpen}>
            <GitHubIcon />
          </IconButton>
        </Tooltip>
        <GenericModal
          open={open}
          handleClose={handleClose}
          Content={
            <MuiThemeProvider theme={getMuiTheme()}>
              <div
                className={classes.paper}
              >
                <Grid
                  container spacing={2}
                >
                  <Grid
                    className={classes.GridTitle}
                    item xs={12}>
                    <h2 id="simple-modal-title">Import using Github</h2>
                  </Grid>
                  <div className={classes.ghGroups}>
                    <FormControl className={classes.formControl}>
                      <InputLabel htmlFor="grouped-native-select">Repository</InputLabel>
                      <Select defaultValue="" id="grouped-select">
                        <MenuItem value="">
                          <em>None</em>
                        </MenuItem>
                        <MenuItem value={1}>service-mesh-patterns</MenuItem>
                        <MenuItem value={2}>image-hub</MenuItem>
                      </Select>
                    </FormControl>
                    <FormControl className={classes.formControl}>
                      <InputLabel htmlFor="grouped-select">Branch</InputLabel>
                      <Select defaultValue="" id="grouped-select">
                        <MenuItem value="">
                          <em>None</em>
                        </MenuItem>
                        <MenuItem value={1}>master</MenuItem>
                        <MenuItem value={2}>main</MenuItem>
                      </Select>
                    </FormControl>
                  </div>
                  <Grid
                    item xs={12}
                    className={classes.optionImport}
                  >
                    <Typography align="center">OR</Typography>
                  </Grid>
                  <Grid
                    item xs={12}>
                    <TextField id="standard-basic" label="Paste URL here" fullWidth onChange={(e) => setInput(e.target.value)} />
                  </Grid>
                  <Grid
                    item xs={12}>
                    <Button fullWidth variant="contained" color="primary" onClick={() => handleSubmit()}>Import</Button>
                  </Grid>
                </Grid>
              </div>
            </MuiThemeProvider>
          }
        />
      </label>
    </React.Fragment>
  )
}

export default GHImport;