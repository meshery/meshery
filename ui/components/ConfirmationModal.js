import {
  Button, Checkbox, Chip, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, Tab, Tabs, TextField,
  Tooltip, Typography
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { styled } from "@mui/material/styles";
import { Search } from "@mui/icons-material";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { setK8sContexts, updateProgress } from "../lib/store";
import { closeButtonForSnackbarAction, errorHandlerGenerator, successHandlerGenerator } from "./ConnectionWizard/helpers/common";
import { pingKubernetes } from "./ConnectionWizard/helpers/kubernetesHelpers";
import { getK8sConfigIdsFromK8sConfig } from "../utils/multi-ctx";
import { bindActionCreators } from "redux";
import { useEffect, useState } from "react";
import UndeployIcon from "../public/static/img/UndeployIcon";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import AddIcon from '@mui/icons-material/Add';
import DoneIcon from "@mui/icons-material/Done";
import Link from 'next/link';
import Operator from "../assets/img/Operator";
import { ACTIONS } from "../utils/Enum";
import { iconMedium, iconSmall } from "../css/icons.styles";

const DialogDialogBox = styled(Dialog)(() => ({}));

const DivChip = styled("div")(({ theme }) => ({
  height : "50px",
  fontSize : "15px",
  position : "relative",
  top : theme.spacing(0.5),
  [theme.breakpoints.down("md")] : { fontSize : "12px" },
}));

const ChipCtxChip = styled(Chip)(({ theme }) => ({
  backgroundColor : "white",
  cursor : "pointer",
  marginRight : theme.spacing(1),
  marginLeft : theme.spacing(1),
  marginBottom : theme.spacing(1),
  height : "100%",
  padding : theme.spacing(0.5),
}));

const ImgCtxIcon = styled("img")(({ theme }) => ({
  display : "inline",
  verticalAlign : "text-top",
  width : theme.spacing(2.5),
  marginLeft : theme.spacing(0.5),
}));

const DialogTitleTitle = styled(DialogTitle)(({ theme }) => ({
  textAlign : "center",

  // minWidth : 300,
  padding : theme.spacing(1),

  color : "#fff",
  backgroundColor : "rgb(57, 102, 121)",
  fontSize : "1rem",
}));

const DialogContentTextSubtitle = styled(DialogContentText)(() => ({
  minWidth : 400,
  overflowWrap : "anywhere",
  textAlign : "center",
  padding : "5px",
}));

const StyledButton = styled(Button)(({ theme }) => ({
  margin : theme.spacing(0.5),
  padding : theme.spacing(1),
  borderRadius : 5,
  minWidth : 100,
  color : "#fff",

  "&:hover" : {
    boxShadow :
      "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)",
  },
}));

const DialogActionsActions = styled(DialogActions)(() => ({
  display : "flex",
  justifyContent : "space-evenly",
}));

const DivAll = styled("div")(() => ({
  display : "table",
}));

const DivContexts = styled("div")(() => ({
  display : "flex",
  flexWrap : "wrap",
}));

const StyledTabs = styled(Tabs)(() => ({
  marginLeft : 0,

  "& .MuiTab-root.Mui-selected" : {
    backgroundColor : "#D9D9D9",
  },
}));

const SpanTabLabel = styled("span")(({ theme }) => ({
  tabLabel : {
    [theme.breakpoints.up("sm")] : {
      fontSize : "1em",
    },
    [theme.breakpoints.between("xs", "sm")] : {
      fontSize : "0.8em",
    },
  },
}));

const StyledAddIcon = styled(AddIcon)(({ theme }) => ({
  width : theme.spacing(2.5),
  paddingRight : theme.spacing(0.5),
}));

const DivTextContent = styled("div")(() => ({
  display : "flex",
  flexDirection : "column",
  alignItems : "center",
  justifyContent : "center",
  marginTop : "1rem",
  backgroundColor : "rgb(234, 235, 236)",
  padding : "10px",
  borderRadius : "10px",
}));

function ConfirmationMsg(props) {
  const {
    classes,
    open,
    handleClose,
    submit,
    selectedK8sContexts,
    k8scontext,
    title,
    validationBody,
    setK8sContexts,
    enqueueSnackbar,
    closeSnackbar,
    componentCount,
    tab,
  } = props;

  const [tabVal, setTabVal] = useState(tab);
  const [disabled, setDisabled] = useState(true);
  const [context, setContexts] = useState([]);
  let isDisabled =
    typeof selectedK8sContexts.length === "undefined" ||
    selectedK8sContexts.length === 0;

  useEffect(() => {
    setTabVal(tab);
    setContexts(k8scontext);
  }, [open]);

  useEffect(() => {
    setDisabled(isDisabled);
  }, [selectedK8sContexts]);

  const handleTabValChange = (event, newVal) => {
    setTabVal(newVal);
  };

  const handleKubernetesClick = (ctxID) => {
    updateProgress({ showProgress : true });
    pingKubernetes(
      successHandlerGenerator(
        enqueueSnackbar,
        closeButtonForSnackbarAction(closeSnackbar),
        "Kubernetes pinged",
        () => updateProgress({ showProgress : false })
      ),
      errorHandlerGenerator(
        enqueueSnackbar,
        closeButtonForSnackbarAction(closeSnackbar),
        "Kubernetes not pinged",
        () => updateProgress({ showProgress : false })
      ),
      ctxID
    );
  };

  const handleSubmit = () => {
    if (selectedK8sContexts.length === 0) {
      enqueueSnackbar(
        "Please select Kubernetes context(s) before proceeding with the operation",
        {
          variant : "info",
          preventDuplicate : true,
          action : (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => closeSnackbar(key)}
            >
              <CloseIcon style={iconMedium} />
            </IconButton>
          ),
          autoHideDuration : 3000,
        }
      );
    }

    if (tabVal === 2) {
      submit.deploy();
    } else if (tabVal === 1) {
      submit.unDeploy();
    }
    handleClose();
  };

  const searchContexts = (search) => {
    if (search === "") {
      setContexts(k8scontext);
      return;
    }
    let matchedCtx = [];
    k8scontext.forEach((ctx) => {
      if (ctx.name.includes(search)) {
        matchedCtx.push(ctx);
      }
    });
    setContexts(matchedCtx);
  };

  const setContextViewer = (id) => {
    if (id === "all") {
      if (selectedK8sContexts?.includes("all")) {
        // updateProgress({ showProgress : true })
        setK8sContexts({ selectedK8sContexts : [] });
      } else {
        setK8sContexts({ selectedK8sContexts : ["all"] });
      }
      return;
    }

    if (selectedK8sContexts?.includes(id)) {
      const filteredContexts = selectedK8sContexts.filter((cid) => cid !== id);
      setK8sContexts({ selectedK8sContexts : filteredContexts });
    } else if (selectedK8sContexts[0] === "all") {
      const allContextIds = getK8sConfigIdsFromK8sConfig(k8scontext);
      setK8sContexts({
        selectedK8sContexts : allContextIds.filter((cid) => cid !== id),
      });
    } else {
      if (selectedK8sContexts.length === k8scontext.length - 1) {
        setK8sContexts({ selectedK8sContexts : ["all"] });
        return;
      }
      setK8sContexts({ selectedK8sContexts : [...selectedK8sContexts, id] });
    }
  };
  return (
    <div className={classes.root}>
      <DialogDialogBox
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <>
          <DialogTitleTitle id="alert-dialog-title">{title}</DialogTitleTitle>
          {/* <Paper square className={classes.paperRoot}> */}
          <StyledTabs
            value={validationBody ? tabVal : tabVal === 2 ? 1 : 0}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            {!!validationBody && (
              <Tab
                data-cy="validate-btn-modal"
                className={classes.tab}
                onClick={(event) => handleTabValChange(event, 0)}
                label={
                  <div style={{ display : "flex" }}>
                    {" "}
                    <DoneIcon
                      style={{ margin : "2px", ...iconSmall }}
                      fontSize="small"
                    />
                    <SpanTabLabel>Validate</SpanTabLabel>{" "}
                  </div>
                }
              />
            )}
            <Tab
              data-cy="Undeploy-btn-modal"
              className={classes.tab}
              onClick={(event) => handleTabValChange(event, 1)}
              label={
                <div style={{ display : "flex" }}>
                  {" "}
                  <div style={{ margin : "2px" }}>
                    {" "}
                    <UndeployIcon
                      style={iconSmall}
                      fill="rgba(0, 0, 0, 0.54)"
                      width="20"
                      height="20"
                    />{" "}
                  </div>{" "}
                  <SpanTabLabel>Undeploy</SpanTabLabel>{" "}
                </div>
              }
            />
            <Tab
              data-cy="deploy-btn-modal"
              className={classes.tab}
              onClick={(event) => handleTabValChange(event, 2)}
              label={
                <div style={{ display : "flex" }}>
                  {" "}
                  <DoneAllIcon
                    style={{ margin : "2px", ...iconSmall }}
                    fontSize="small"
                  />{" "}
                  <SpanTabLabel>Deploy</SpanTabLabel>{" "}
                </div>
              }
            />
          </StyledTabs>

          {(tabVal === ACTIONS.DEPLOY || tabVal === ACTIONS.UNDEPLOY) && (
            <DialogContent>
              <DialogContentTextSubtitle id="alert-dialog-description">
                <Typography
                  variant="subtitle1"
                  style={{ marginBottom : "0.8rem" }}
                >
                  {" "}
                  {componentCount !== undefined ? (
                    <>
                      {" "}
                      {componentCount} component{componentCount > 1 ? "s" : ""}{" "}
                    </>
                  ) : (
                    ""
                  )}
                </Typography>
                {k8scontext.length > 0 ? (
                  <Typography variant="body1">
                    <TextField
                      id="search-ctx"
                      label="Search"
                      size="small"
                      variant="outlined"
                      onChange={(event) => searchContexts(event.target.value)}
                      style={{
                        width : "100%",
                        backgroundColor : "rgba(102, 102, 102, 0.12)",
                        margin : "1px 1px 8px ",
                      }}
                      InputProps={{
                        endAdornment : <Search style={iconMedium} />,
                      }}
                      // margin="none"
                    />
                    {context.length > 0 ? (
                      <DivAll>
                        <Checkbox
                          checked={selectedK8sContexts?.includes("all")}
                          onChange={() => setContextViewer("all")}
                          color="primary"
                        />
                        <span style={{ fontWeight : "bolder" }}>select all</span>
                      </DivAll>
                    ) : (
                      <Typography variant="subtitle1">
                        No Context found
                      </Typography>
                    )}

                    <DivContexts>
                      {context.map((ctx) => (
                        <DivChip id={ctx.id} key={ctx.id}>
                          <Tooltip title={`Server: ${ctx.server}`}>
                            <div
                              style={{
                                display : "flex",
                                justifyContent : "flex-wrap",
                                alignItems : "center",
                              }}
                            >
                              <Checkbox
                                checked={
                                  selectedK8sContexts?.includes(ctx.id) ||
                                  (selectedK8sContexts?.length > 0 &&
                                    selectedK8sContexts[0] === "all")
                                }
                                onChange={() => setContextViewer(ctx.id)}
                                color="primary"
                              />
                              <ChipCtxChip
                                label={ctx.name}
                                onClick={() => handleKubernetesClick(ctx.id)}
                                icon={
                                  <ImgCtxIcon src="/static/img/kubernetes.svg" />
                                }
                                variant="outlined"
                                data-cy="chipContextName"
                              />
                            </div>
                          </Tooltip>
                        </DivChip>
                      ))}
                    </DivContexts>
                  </Typography>
                ) : (
                  <DivTextContent>
                    <Operator />
                    <Typography variant="h5">
                      No cluster connected yet
                    </Typography>

                    <Link href="/settings">
                      <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        style={{
                          margin : "0.6rem 0.6rem",
                          whiteSpace : "nowrap",
                        }}
                      >
                        <StyledAddIcon />
                        Connect Clusters
                      </Button>
                    </Link>
                  </DivTextContent>
                )}
              </DialogContentTextSubtitle>
            </DialogContent>
          )}
          {tabVal === ACTIONS.VERIFY && ( // Validate
            <DialogContent>
              <DialogContentText>{validationBody}</DialogContentText>
            </DialogContent>
          )}
          {/* </Paper> */}

          <DialogActionsActions>
            {tabVal == ACTIONS.DEPLOY || tabVal === ACTIONS.UNDEPLOY ? (
              <>
                <Button onClick={handleClose} type="submit" variant="contained">
                  <Typography variant body2>
                    {" "}
                    CANCEL{" "}
                  </Typography>
                </Button>
                <Button
                  disabled
                  className={
                    tabVal === ACTIONS.UNDEPLOY ? classes.disabledBtnDel : ""
                  }
                  type="submit"
                  variant="contained"
                  color="primary"
                >
                  <Typography variant body2>
                    {" "}
                    {tabVal === ACTIONS.UNDEPLOY
                      ? "UNDEPLOY LATER"
                      : "DEPLOY LATER"}{" "}
                  </Typography>
                  {/* colorchange  */}
                </Button>
                <Button
                  onClick={handleSubmit}
                  className={
                    isDisabled
                      ? tabVal === ACTIONS.UNDEPLOY
                        ? classes.disabledBtnDel
                        : classes.button
                      : tabVal === ACTIONS.UNDEPLOY
                        ? classes.undeployBtn
                        : classes.button
                  }
                  autoFocus
                  type="submit"
                  variant="contained"
                  color="primary"
                  data-cy="deploy-btn-confirm"
                  disabled={disabled}
                >
                  <Typography variant body2>
                    {" "}
                    {tabVal === ACTIONS.UNDEPLOY ? "UNDEPLOY" : "DEPLOY"}{" "}
                  </Typography>
                </Button>
              </>
            ) : (
              <StyledButton
                onClick={handleClose}
                autoFocus
                type="submit"
                variant="contained"
                color="primary"
              >
                <Typography variant body2>
                  {" "}
                  OK{" "}
                </Typography>
              </StyledButton>
            )}
          </DialogActionsActions>
        </>
      </DialogDialogBox>
    </div>
  );
}

const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get("selectedK8sContexts");
  const k8scontext = state.get("k8sConfig");
  return { selectedK8sContexts : selectedK8sContexts, k8scontext : k8scontext };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  setK8sContexts : bindActionCreators(setK8sContexts, dispatch),
});

export default (
  connect(mapStateToProps, mapDispatchToProps)(withSnackbar(ConfirmationMsg))
);
