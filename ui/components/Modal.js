import React from "react";
import { IconButton, Menu, MenuItem, Tooltip, Button, Typography } from "@material-ui/core";
import { Dialog, DialogActions, makeStyles } from "@material-ui/core";
// import useStyles from "./MesheryPatterns/Cards.styles";
import CloseIcon from "@material-ui/icons/Close";
import PublicIcon from "@material-ui/icons/Public";
import InfoIcon from "@material-ui/icons/Info";
import RJSFWrapper from "./MesheryMeshInterface/PatternService/RJSF_wrapper";
import { ArrowDropDown } from "@material-ui/icons";
import { getSchema } from "./MesheryMeshInterface/PatternService/helper";
import Link from "next/link"

const useStyles = makeStyles((theme) => ({
  "@keyframes rotateCloseIcon" : {
    from : {
      transform : "rotate(0deg)",
    },
    to : {
      transform : "rotate(360deg)",
    },
  },
  infoIcon : {
    color : theme.palette.type === "dark" ? "#00B39F" : "#607d8b",
  },
  modalHeader : {
    display : "flex",
    justifyContent : "space-between",
    alignItems : "center",
    paddingBottom : 10,
    padding : "0 .5rem",
    paddingTop : 10,
    backgroundColor : theme.palette.secondary.mainBackground,
  },
  modelHeader : {
    fontSize : "1rem",
    color : "#fff",
  },
  iconStyle : {
    color : "#fff",
  },
  iconContainer : {
    transition : "all .3s",
    "&:hover" : {
      backgroundColor : "transparent !important",
      animation : "$rotateCloseIcon 1s",
    },
  },
  submitButton : {
    backgroundColor : theme.palette.type === "dark" ? "#00B39F" : "#607d8b",
    color : "#fff",
    width : "100%",
  },
  iconPatt : {
    marginRight : theme.spacing(1),
  },
  btnText : {
    textTransform : "none",
  },
  toolTip : {
    "&.MuiTooltip-popper" : {
      zIndex : "9999999 !important"
    }
  }
}));

const SchemaVersion = ({ schema_array, type, schemaChangeHandler }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <div>
      <Tooltip title="Schema_Changer">
        <IconButton component="span" onClick={(e) => setAnchorEl(e.currentTarget)}>
          <ArrowDropDown style={{ color : "#000" }} />
        </IconButton>
      </Tooltip>
      <Menu id="schema-menu" anchorEl={anchorEl} open={open} handleClose={handleClose}>
        {schema_array.map((version, index) => (
          <MenuItem
            id="schema-menu-item"
            key={index}
            selected={version === type}
            onClick={() => {
              schemaChangeHandler(version);
              handleClose();
            }}
          >
            {version}
          </MenuItem>
        ))}
      </Menu>
    </div>
  );
};
function Modal(props) {
  const {
    open,
    title,
    handleClose,
    onChange,
    schema,
    formData,
    schema_array,
    type,
    schemaChangeHandler,
    handleSubmit,
    payload,
    uiSchema = {},
  } = props;
  const classes = useStyles();

  return (
    <>
      <Dialog open={open} onClose={handleClose}>
        <div className={classes.modalHeader}>
          <Typography variant="h5"></Typography>
          <Typography className={classes.modelHeader} variant="h5">
            {title}
            {schema_array?.length < 1 && (
              <SchemaVersion schema_array={schema_array} type={type} schemaChangeHandler={schemaChangeHandler} />
            )}
          </Typography>
          <IconButton className={classes.iconContainer} onClick={handleClose}>
            <CloseIcon className={classes.iconStyle} />
          </IconButton>
        </div>
        {/* <Snackbar open={open} onClose={handleClose}>
        <Alert onClose={handleClose} severity="warning">
          Please use name which makes more sense!
        </Alert>
      </Snackbar> */}
        <RJSFWrapper
          key={type}
          formData={formData}
          jsonSchema={schema || getSchema(type)}
          uiSchema={uiSchema}
          onChange={onChange}
          hideTitle={true}
        />
        <DialogActions>
          <Button
            title="Submit"
            variant="contained"
            color="primary"
            className={classes.submitButton}
            onClick={() => {
              handleClose();
              handleSubmit(payload);
            }}
          >
            <PublicIcon className={classes.iconPatt} />
            <span className={classes.btnText}> Submit for Approval </span>
          </Button>
          <Link href="https://docs.meshery.io/concepts/catalog" passHref>
            <Tooltip title="View catalog approval flow" placement="top" className={classes.toolTip}>
              <a target="_blank" rel="noopener noreferrer">
                <IconButton color="primary">
                  <InfoIcon />
                </IconButton>
              </a>
            </Tooltip>
          </Link>
        </DialogActions>
      </Dialog>
    </>
  );
}

export default Modal;

