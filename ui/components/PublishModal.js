import React from 'react';
import { useEffect } from 'react';
import { Button, Grid, IconButton } from '@material-ui/core';
// import { createTheme } from '@material-ui/core/styles';
import validator from "@rjsf/validator-ajv8";
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import { Form } from '@rjsf/material-ui';
import useStyles from "./MesheryPatterns/Cards.styles";
import PublicIcon from '@material-ui/icons/Public';
import CloseIcon from '@material-ui/icons/Close';
// const getMuiTheme = () => createTheme({
//   palette : {
//     primary : {
//       main : "#607d8b"
//     },
//     secondary : {
//       main : "#666666"
//     }
//   },
//   overrides : {
//     MuiGrid : {
//       input : {
//         color : '#607d8b'
//       }
//     },
//   }
// })

export default function PublishModal(props) {
  const { open, handleClose, pattern, handlePublish } = props;
  const classes = useStyles();
  const schema = {
    "type" : "object",
    "properties" : {
      "compatibility" : {
        "type" : "array",
        "items" : {
          "enum" : [            
            "Kubernetes",
            "Argo CD",
            "AWS App Mesh",
            "Consul",
            "Fluentd",
            "Istio",
            "Jaeger",
            "Kuma",
            "Linkerd",
            "Network Service Mesh",
            "NGINX Service Mesh",
            "Open Service Mesh",
            "Prometheus",
            "Traefik Mesh"
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
    },
    "required" : ["compatibility", "pattern_caveats", "pattern_info", "type"]
  }

  const [data, setData] = React.useState(null);

  useEffect(() => {
    if (pattern?.catalog_data) {
      setData(pattern.catalog_data)
    }
  }, []);

  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>
        <DialogTitle>
          <div className={classes.publishTitle}>

            <b id="simple-modal-title" style={{ textAlign : "center" }} > {pattern?.name}</b>
            <IconButton aria-label="close" className={classes.closeButton} onClick={handleClose}>
              <CloseIcon />
            </IconButton>
          </div>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={24} alignItems="center">
            <Form schema={schema} formData={data} validator={validator}
              onSubmit={(data) => {
                handleClose();
                handlePublish({
                  id : pattern.id,
                  catalog_data : data.formData
                })
              }}
            >
              <Button
                title="Publish"
                variant="contained"
                color="primary"
                type="submit"
                className={classes.testsButton}
              >
                <PublicIcon className={classes.iconPatt} />
                <span className={classes.btnText}> Submit for Approval </span>
              </Button>
            </Form>
          </Grid>

        </DialogContent>
        <DialogActions>

        </DialogActions>

      </Dialog>
    </>
  )
}