import React from 'react'
import {  Button, Grid, IconButton } from '@material-ui/core';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import validator from "@rjsf/validator-ajv6";
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import { Form } from '@rjsf/material-ui';

import PublishIcon from '@material-ui/icons/Publish';
import CloseIcon from '@material-ui/icons/Close';
const getMuiTheme = () => createTheme({
  palette : {
    primary : {
      main : "#607d8b"
    },
    secondary : {
      main : "#666666"
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


const styles = makeStyles(() => ({
  title : {
    display : 'flex',
    justifyContent : 'space-between',
    alignItems : 'center',
  },
  btn : {
    backgroundColor : '#607d8b',
  },
  iconPatt : {
    color : '#fff',
    marginRight : '0.5rem',
  },
  btnText : {
    color : '#fff',
  },

}));

function PublishModal(props) {
  const {  open, handleClose ,pattern,handlePublish } = props;
  const classes = styles();
  const schema={
    "type" : "object",
    "properties" : {
      "compatibility" : {
        "type" : "array",
        "items" : {
          "enum" : [
            "Istio",
            "Linkerd",
            "App Mesh",
            "OSM",
            "Nginx",
            "Kuma",
            "Consul",
            "NSM",
            "Traefik"
          ],
          "type" : "string"
        },
        "uniqueItems" : true,

      },
      "pattern_caveats" : {
        "type" : "string"
      },
      "pattern_info" : {
        "type" : "string"
      },
      "type" : {
        "type" : "string",
        "examples" : [
          "deployment",
          "observability",
          "resiliency",
          "scaling",
          "security",
          "traffic-management",
          "troubleshooting",
          "workloads"
        ]
      }
    }
  }
  const [data,setData]=React.useState(null)
  const [payload,setPayload]=React.useState({
    "id" : pattern?.id,
    "catalog_data" : pattern?.catalog_data
  })
  React.useEffect(() => {
    setData(pattern.catalog_data)
  },[pattern])
  React.useEffect(() => {
    setPayload({
      "id" : pattern?.id,
      "catalog_data" : data
    })
    console.log(payload)
  },[data])



  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>

        <MuiThemeProvider theme={getMuiTheme()}>
          <DialogTitle>
            <div className={classes.title}>

              <b id="simple-modal-title" style={{ textAlign : "center" }} > {pattern?.name}</b>
              <IconButton aria-label="close" className={classes.closeButton} onClick={handleClose}>
                <CloseIcon />
              </IconButton>
            </div>
          </DialogTitle>
          <DialogContent>
            <Grid container spacing={24} alignItems="center">
              <Form schema={schema} formData={data} validator={validator} onChange={(e) => setData(e.formData)} ><></></Form>
            </Grid>

          </DialogContent>
          <DialogActions>
            <Button
              title="Publish"
              variant="contained"
              className={classes.btn}
              onClick={() => {
                handleClose();
                handlePublish(payload)
              }}
            >
              <PublishIcon className={classes.iconPatt} />
              <span className={classes.btnText}> Publish </span>
            </Button>
          </DialogActions>
        </MuiThemeProvider>
      </Dialog>
    </>
  )
}

export default PublishModal;