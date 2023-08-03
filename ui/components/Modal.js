import React, { useState, useEffect } from "react";
import { IconButton, Menu, MenuItem, Tooltip, Button, Typography } from "@material-ui/core";
import { Dialog, DialogActions, makeStyles } from "@material-ui/core";
import { CustomTextTooltip } from "./MesheryMeshInterface/PatternService/CustomTextTooltip";
import CloseIcon from "@material-ui/icons/Close";
import PublicIcon from "@material-ui/icons/Public";
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import RJSFWrapper from "./MesheryMeshInterface/PatternService/RJSF_wrapper";
import { ArrowDropDown } from "@material-ui/icons";
import { getSchema } from "./MesheryMeshInterface/PatternService/helper";
import Link from "next/link";
import { Snackbar } from "@material-ui/core";
import { Alert } from "@material-ui/lab";

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
    position : "absolute",
    right : 10,
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
    backgroundColor : theme.palette.secondary.focused,
    color : "#fff",
  },
  iconPatt : {
    marginRight : theme.spacing(1),
  },
  btnText : {
    textTransform : "none",
  },
  toolTip : {
    textDecoration : "underline",
    color : theme.palette.secondary.link2,
  },
  dialogAction : {
    display : "flex",
    justifyContent : "center",
    padding : "0.5rem 1rem",
  },
  snackbar : {
    backgroundColor : theme.palette.secondary.elevatedComponents,
  },
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

/**
 * Renders common dialog component.
 *
 * @param {Object} props - Component props.
 * @param {boolean} props.open - Determines whether the modal is open or not.
 * @param {string} props.title - The title of the modal.
 * @param {Function} props.handleClose - Function to handle the close event of the modal.
 * @param {Function} props.onChange - Function to handle the change event of the form fields in the modal.
 * @param {Object} props.schema - The JSON schema for the form fields.
 * @param {Object} props.formData - The form data object.
 * @param {Array} props.schema_array - An array of schema versions.
 * @param {string} props.type - The selected schema version.
 * @param {Function} props.schemaChangeHandler - Function to handle the change of the schema version.
 * @param {Function} props.handleSubmit - Function to handle the submit event of the modal.
 * @param {Object} props.payload - The payload for the submit event.
 * @param {Object} props.showInfoIcon - Determines whether to show the info icon adjacent to the modal button.
 * @param {string} props.submitBtnText - The text for the submit button.
 * @param {Object} props.uiSchema - The UI schema for the form fields.
 */

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
    showInfoIcon,
    submitBtnText,
    uiSchema = {},
  } = props;
  const classes = useStyles();

  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const [snackbar, setSnackbar] = useState(false);
  const formRef = React.createRef();

  const renderTooltipContent = () => (
    <div>
      <span>{showInfoIcon.text}</span>
      {showInfoIcon.link && (
        <Link href={showInfoIcon.link} passHref onClick={(e) => e.stopPropagation()}>
          <a className={classes.toolTip} target="_blank" rel="noopener noreferrer">
            Learn more
          </a>
        </Link>
      )}
    </div>
  );

  const handleFormSubmit = () => {
    if (formRef.current && formRef.current.validateForm()) {
      handleClose();
      handleSubmit(payload);
    }
  };

  useEffect(() => {
    setCanNotSubmit(false);
    const handleDesignNameCheck = () => {
      const designName = title?.toLowerCase();
      const forbiddenWords = ["untitled design", "Untitled", "lfx"];

      for (const word of forbiddenWords) {
        if (designName?.includes(word)) {
          setSnackbar({
            severity : "warning",
            message : `Design name should not contain Untitled Design, Untitled, LFX`,
            open : true,
          });
          setCanNotSubmit(true);
          break;
        }
      }
    };
    handleDesignNameCheck();
  }, [title]);

  return (
    <>
      <Dialog style={{ zIndex : 9999 }} open={open} onClose={handleClose}>
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

        <RJSFWrapper
          key={type}
          formData={formData}
          jsonSchema={schema || getSchema(type)}
          uiSchema={uiSchema}
          onChange={onChange}
          liveValidate={false}
          formRef={formRef}
          hideTitle={true}
        />

        <DialogActions className={classes.dialogAction}>
          <Button
            title={submitBtnText ? submitBtnText : "Submit"}
            variant="contained"
            color="primary"
            className={classes.submitButton}
            disabled={canNotSubmit}
            onClick={handleFormSubmit}
          >
            {/* TODO: change below logic to support adding icon from prop */}
            {!submitBtnText && (
              <PublicIcon className={classes.iconPatt} />
            )}
            <span className={classes.btnText}>{submitBtnText ? submitBtnText : "Submit for Approval" }</span>
          </Button>
          {showInfoIcon && (
            <CustomTextTooltip
              backgroundColor="#3C494F"
              placement="top"
              interactive={true}
              title={renderTooltipContent()}
            >
              <IconButton className={classes.infoIcon} color="primary">
                <InfoOutlinedIcon />
              </IconButton>
            </CustomTextTooltip>
          )}
        </DialogActions>
        {snackbar && (
          <Snackbar
            anchorOrigin={{ vertical : "bottom", horizontal : "right" }}
            open={snackbar.open}
            autoHideDuration={6000}
            onClose={() => setSnackbar(null)}
          >
            <Alert className={classes.snackbar} onClose={() => setSnackbar(null)} severity={snackbar.severity}>
              {snackbar.message}
            </Alert>
          </Snackbar>
        )}
      </Dialog>
    </>
  );
}

export default Modal;

