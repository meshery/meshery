import React from 'react'
import { Grid, IconButton, Menu, MenuItem, Tooltip } from '@material-ui/core';
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import useStyles from "./MesheryPatterns/Cards.styles";
import CloseIcon from '@material-ui/icons/Close';
import RJSFWrapper from './MesheryMeshInterface/PatternService/RJSF_wrapper';
import { ArrowDropDown } from '@material-ui/icons';
import { getSchema } from './MesheryMeshInterface/PatternService/helper';

const SchemaVersion = ({ schema_array, type, schemaChangeHandler }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <div>
      <Tooltip title="Schema_Changer">
        <IconButton component="span" onClick={e => setAnchorEl(e.currentTarget)}>
          <ArrowDropDown style={{ color: "#000" }} />
        </IconButton>
      </Tooltip>
      <Menu id="schema-menu" anchorEl={anchorEl} open={open} handleClose={handleClose}>
        {schema_array.map((version, index) => (
          <MenuItem
            id="schema-menu-item"
            key={index}
            selected={version === type}
            onClick={() => {
              schemaChangeHandler(version)
              handleClose();
            }}
          >
            {version}
          </MenuItem>
        ))}
      </Menu>
    </div>
  )
}


const RJSFWrapperComponent = (uiSchema) =>
    (
    /** @type {{ jsonSchema: any; children: React.DetailedReactHTMLElement<any, HTMLElement>; }} */ props
    ) => {

      // Clone the child to pass in additional props
      return React.cloneElement(props.children, {
        ...(props.children?.props || {}),
        uiSchema
      });
    };


function Modal({ open, title, handleClose, onChange, schema, formData, children, schema_array, type, schemaChangeHandler, onSubmit, uiSchema }) {
  const classes = useStyles();

  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>
        <DialogTitle>
          <div style={{
            display: "flex",
            justifyContent: "space-between",
          }}>
            <p style={{
              display: "flex",
              alignItems: "center",
            }}>
              <b id="simple-modal-title" style={{ textAlign: "center" }} > {title}</b>
              {schema_array?.length > 1 && (
                <SchemaVersion schema_array={schema_array} type={type} schemaChangeHandler={schemaChangeHandler} />
              )}
            </p>
            <IconButton aria-label="close" className={classes.closeButton} onClick={handleClose}>
              <CloseIcon />
            </IconButton>
          </div>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={24} alignItems="center">
            <RJSFWrapper
              key={type}
              formData={formData}
              jsonSchema={schema || getSchema(type)}
              onChange={onChange}
              hideTitle={true}
              uiSchema={uiSchema}
              RJSFWrapperComponent={RJSFWrapperComponent(uiSchema)}
            />
          </Grid>

        </DialogContent>
        <DialogActions>
          {children}
        </DialogActions>
      </Dialog>
    </>
  )
}

export default Modal;