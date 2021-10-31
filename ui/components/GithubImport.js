import React from 'react'
import { Tooltip, IconButton, TextField, Button, Grid, Typography } from '@material-ui/core';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import GitHubIcon from '@material-ui/icons/GitHub';
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

}));
import GenericModal from "./GenericModal";


const GithubUploader = ({ onSubmit }) => {
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
    <>
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
                  <Grid
                    item xs={12}>
                    <TextField id="standard-basic" label="Paste URL here" fullWidth onChange={(e) => setInput(e.target.value)} />
                  </Grid>
                  <Grid
                    item xs={12} sm={6}>
                    <Button fullWidth variant="contained" color="primary" onClick={() => handleSubmit()}>Import</Button>
                  </Grid>
                  <Grid
                    item xs={12}>
                    <Typography>
                      Note: Private repository can not be used when using service as provider.
                    </Typography>
                  </Grid>
                </Grid>

              </div>
            </MuiThemeProvider>
          }
        />
      </label>
    </>
  )
}

export default GithubUploader
