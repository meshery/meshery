import React, { useState, useEffect } from "react";
import { IconButton, Menu, MenuItem, Tooltip, Button, Typography } from "@material-ui/core";
import { Dialog, DialogActions, makeStyles } from "@material-ui/core";
import { CustomTextTooltip } from "./MesheryMeshInterface/PatternService/CustomTextTooltip";
import CloseIcon from "@material-ui/icons/Close";
import PublicIcon from "@material-ui/icons/Public";
import InfoIcon from "@material-ui/icons/Info";
import RJSFWrapper from "./MesheryMeshInterface/PatternService/RJSF_wrapper";
import { ArrowDropDown } from "@material-ui/icons";
import { getSchema } from "./MesheryMeshInterface/PatternService/helper";
import Link from "next/link";
import { useSnackbar } from "notistack";
import { Cancel } from "@material-ui/icons";

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
    textDecoration : "underline",
    color : theme.palette.secondary.link2,
  },
  dialogAction : {
    padding : "0.5rem 1rem",
  },
  snackbar : {
    zIndex : 9999,
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

  const [canNotSubmit, setCanNotSubmit] = useState(false);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const renderTooltipContent = () => (
    <div>
      <span>Upon submitting your catalog item, an approval flow will be initiated. </span>
      <Link href="https://docs.meshery.io/concepts/catalog" passHref onClick={(e) => e.stopPropagation()}>
        <a className={classes.toolTip} target="_blank" rel="noopener noreferrer">
          Learn more
        </a>
      </Link>
    </div>
  );

  useEffect(() => {
    setCanNotSubmit(false);
    const handleDesignNameCheck = () => {
      const designName = title?.toLowerCase();
      const forbiddenWords = ["untitled design", "untitle", "lfx"];

      for (const word of forbiddenWords) {
        if (designName?.includes(word)) {
          enqueueSnackbar("Please provide a valid name", {
            variant : "warning",
            autoHideDuration : 4000,
            preventDuplicate : true,
            contentProps : { className : classes.snackbar },
            action : (key) => (
              <IconButton onClick={() => closeSnackbar(key)} color="secondary">
                <Cancel />
              </IconButton>
            ),
          });
          setCanNotSubmit(true);
          break;
        }
      }
    };
    handleDesignNameCheck();
  }, [title, enqueueSnackbar, closeSnackbar]);

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
        <RJSFWrapper
          key={type}
          formData={formData}
          jsonSchema={schema || getSchema(type)}
          uiSchema={uiSchema}
          onChange={onChange}
          hideTitle={true}
        />
        <DialogActions className={classes.dialogAction}>
          <Button
            title="Submit"
            variant="contained"
            color="primary"
            className={classes.submitButton}
            disabled={canNotSubmit}
            onClick={() => {
              handleClose();
              handleSubmit(payload);
            }}
          >
            <PublicIcon className={classes.iconPatt} />
            <span className={classes.btnText}> Submit for Approval </span>
          </Button>
          <CustomTextTooltip
            backgroundColor="#3C494F"
            placement="top"
            interactive={true}
            title={renderTooltipContent()}
          >
            <IconButton color="primary">
              <InfoIcon />
            </IconButton>
          </CustomTextTooltip>
        </DialogActions>
      </Dialog>
    </>
  );
}

export default Modal;

