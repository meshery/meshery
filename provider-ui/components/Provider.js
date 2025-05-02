import React, { Fragment, useEffect, useState } from "react";
import PropTypes from "prop-types";
import dataFetch from "../lib/data-fetch";
import ProviderLayout from "./ProviderLayout";
import {
  CustomDiv,
  CustomDialog,
  MesheryLogo,
  MenuProviderDisabled,
  CustomDialogActions,
  LearnMore,
  CustomTypography,
  StyledPopover,
} from "./Provider.style";
import {
  Button,
  ButtonGroup,
  Divider,
  DialogContentText,
  DialogTitle,
  MenuList,
  MenuItem,
  Tooltip,
  Typography,
  IconButton,
  CircularProgress,
  styled,
  charcoal,
  accentGrey,
  CHINESE_SILVER,
  KEPPEL,
} from "@layer5/sistent";
import { CloseIcon, ClickAwayListener, DropDownIcon } from "@layer5/sistent";
function CustomDialogTitle(props) {
  const { children, onClose, ...other } = props;

  return (
    <DialogTitle sx={{ m: 0, p: 2 }} {...other}>
      {children}
      {onClose ? (
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position: "absolute",
            right: "1rem",
            top: "1rem",
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      ) : null}
    </DialogTitle>
  );
}

CustomDialogTitle.propTypes = {
  children: PropTypes.node,
  onClose: PropTypes.func.isRequired,
};
//Styled-components:
const StyledTypography = styled(Typography)(({ theme }) => ({
  fontWeight: 500,
  color: charcoal[100],
  marginBottom: theme.spacing(2), // Equivalent to `gutterBottom`
  "& a": {
    fontWeight: "normal",
  },
  "& :hover": {
    color: CHINESE_SILVER,
  },
}));

const StyledTooltip = styled(Tooltip)(({ theme }) => ({
  color: theme.palette.text.inverse,
  cursor: "pointer",
  fontWeight: "normal",
}));

const StyledCustomDialogTitle = styled(CustomDialogTitle)(({ theme }) => ({
  background: accentGrey[10],
  color: theme.palette.text.inverse,
}));

const StyledButton = styled(Button)(({ theme }) => ({
  backgroundColor: theme.palette.IconButton,
}));
const StyledButtonGroup = styled(ButtonGroup)(() => ({
  border: "none",
  "& .MuiButtonGroup-grouped": {
    border: "none !important",
  },
}));

const StyledDialogBox = styled(DialogContentText)(({ theme }) => ({
  color: theme.palette.text.inverse,
  backgroundColor: charcoal[40],
  padding: "1.2rem",
}));

export default function Provider() {
  const [anchorEl, setAnchorEl] = useState(null);
  const [availableProviders, setAvailableProviders] = useState({});
  const [selectedProvider, setSelectedProvider] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  /* eslint-disable no-unused-vars */
  const [openMenu, setOpenMenu] = useState(false);
  const [openModal, setModalOpen] = React.useState(false);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? "popover" : undefined;

  useEffect(() => {
    loadProvidersFromServer();
  }, []);

  function loadProvidersFromServer() {
    dataFetch(
      "/api/providers",
      {
        method: "GET",
        credentials: "include",
      },
      (result) => {
        if (typeof result !== "undefined") {
          Object.keys(result).forEach((key) => {
            if (result[key].ProviderType === "remote") {
              setSelectedProvider(selectedProvider);
            }
          });
          setAvailableProviders(result);
        }
      },
      (error) => {
        console.log(`there was an error fetching providers: ${error}`);
      }
    );
  }

  const handleMenuItemClick = (event, provider) => {
    event.preventDefault();
    setSelectedProvider(provider);
    setOpenMenu(false);
    setIsLoading(true);
    const existingQueryString = window.location.search;
    const existingQueryParams = new URLSearchParams(existingQueryString);
    existingQueryParams.append("provider", encodeURIComponent(provider));
    window.location.href = `/api/provider?${existingQueryParams.toString()}`;
  };

  const handleModalOpen = () => {
    setModalOpen(true);
  };

  const handleModalClose = () => {
    setModalOpen(false);
  };

  return (
    <ProviderLayout>
      <MesheryLogo
        src="/provider/static/img/meshery-logo/meshery-logo-dark-text-noBG.png"
        onError={(e) =>
          (e.target.src =
            "/static/img/meshery-logo/meshery-logo-dark-text-noBG.png")
        }
        alt="logo"
      />
      <CustomDiv>
        {availableProviders !== "" && (
          <Fragment>
            <StyledButtonGroup aria-label="split button">
              <Button
                size="large"
                variant="contained"
                aria-describedby={id}
                onClick={handleClick}
                aria-label="Select Provider"
                data-cy="select_provider"
                disableElevation
              >
                {isLoading && (
                  <CircularProgress
                    size={20}
                    sx={{ color: "white", marginRight: 8 }}
                  />
                )}
                {selectedProvider !== ""
                  ? selectedProvider
                  : "Select your provider"}

                <DropDownIcon />
              </Button>
            </StyledButtonGroup>

            <StyledPopover
              id={id}
              open={open}
              anchorEl={anchorEl}
              onClose={handleClose}
              anchorOrigin={{
                vertical: "bottom",
                horizontal: "center",
              }}
              transformOrigin={{
                vertical: "top",
                horizontal: "center",
              }}
            >
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList
                  sx={{
                    background: charcoal[20],
                    color: (theme) => theme.palette.text.inverse,
                  }}
                  id="split-button-menu"
                  autoFocusItem
                >
                  {Object.keys(availableProviders).map((key) => (
                    <MenuItem
                      key={key}
                      onClick={(e) => handleMenuItemClick(e, key)}
                      sx={{ "&:hover": { backgroundColor: accentGrey[20] } }}
                    >
                      {key}
                    </MenuItem>
                  ))}
                  <Divider
                    sx={{
                      my: 0.5,
                      backgroundColor: accentGrey[40],
                      width: "80%",
                      margin: "auto",
                      marginBottom: "0px",
                    }}
                  />
                  <MenuProviderDisabled
                    sx={{ marginTop: "0px" }}
                    disabled={true}
                    key="CNCF Labs"
                  >
                    CNCF Labs{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="Equinix US-DAL">
                    Equinix US-DAL{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="HPE Security">
                    HPE Security{"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="MIT">
                    Massachusetts Institute of Technology (MIT){"\u00A0"}
                    <span>Offline</span>
                  </MenuProviderDisabled>
                  <MenuProviderDisabled disabled={true} key="UT Austin">
                    The University of Texas at Austin{"\u00A0"}
                    <span> Offline</span>
                  </MenuProviderDisabled>
                  <Button
                    size="small"
                    variant="contained"
                    aria-describedby={id}
                    onClick={() =>
                      window.open(
                        "https://docs.meshery.io/extensibility/providers",
                        "_blank"
                      )
                    }
                    aria-label="Create Provider"
                    data-cy="create_provider"
                    disableElevation
                    sx={{
                      width: '100%',
                      marginTop: '0.4rem'
                    }}
                  >Create Your Own Provider&nbsp;
                    <img
                      src="/provider/static/img/external-link.svg"
                      onError={(e) =>
                        (e.target.src = "/static/img/external-link.svg")
                      }
                      width="16px"
                      alt="External link"
                      style={{
                        filter: 'brightness(0)'
                      }}
                    />
                  </Button>
                </MenuList>
              </ClickAwayListener>
            </StyledPopover>
          </Fragment>
        )}
      </CustomDiv>
      <LearnMore onClick={handleModalOpen}>
        <StyledTypography variant="h6" gutterBottom>
          <StyledTooltip
            title="Learn more about Meshery remote providers"
            placement="bottom"
            data-cy="providers-tooltip"
          >
            Learn more about providers
          </StyledTooltip>
        </StyledTypography>
      </LearnMore>
      <CustomDialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        open={openModal}
        disableScrollLock={true}
        data-cy="providers-modal"
      >
        <StyledCustomDialogTitle
          id="customized-dialog-title"
          onClose={handleModalClose}
        >
          <CustomTypography>Choosing a Provider</CustomTypography>
        </StyledCustomDialogTitle>

        <StyledDialogBox id="customized-dialog-content">
          Login to Meshery by choosing from the available providers. Providers
          extend Meshery by offering various plugins and services, including
          identity services, long-term persistence, advanced performance
          analysis, multi-player user collaboration, and so on.
          <h2>Available Providers</h2>
          {Object.keys(availableProviders).map((key) => {
            return (
              <React.Fragment key={availableProviders[key].provider_name}>
                <p style={{ fontWeight: 700 }}>
                  {availableProviders[key].provider_name}
                </p>
                <ul>
                  {availableProviders[key].provider_description?.map(
                    (desc, i) => (
                      <li key={`desc-${i}`}>{desc}</li>
                    )
                  )}
                </ul>
              </React.Fragment>
            );
          })}
          <p style={{ fontWeight: 700 }}>MIT</p>
          <ul>
            <li>Remote provider for performance testing</li>
            <li>Provides provenence of test results and their persistence</li>
            <li>Adaptive performance analysis - predictive optimization</li>
          </ul>
          <p style={{ fontWeight: 700 }}>The University of Texas at Austin</p>
          <ul>
            <li>Academic research and advanced studies by Ph.D. researchers</li>
            <li></li>
            <li>Used by school of Electrical and Computer Engineering (ECE)</li>
          </ul>
          <p style={{ fontWeight: 700 }}>
            Cloud Native Computing Foundation Infrastructure Lab
          </p>
          <ul>
            <li>
              Performance and compatibility-centric research and validation
            </li>
            <li></li>
            <li>Used by various cloud native projects</li>
          </ul>
          <p style={{ fontWeight: 700 }}>HPE Security</p>
          <ul>
            <li>Istio, SPIRE, and SPIFEE integration</li>
          </ul>
          <p style={{ fontWeight: 700 }}>Equinix</p>
          <ul>
            <li>Identity services</li>
            <li>Bare-metal Kubernetes configuration</li>
          </ul>
        </StyledDialogBox>

        <CustomDialogActions>
          <div className="learnmore">
            <a href="https://docs.meshery.io/extensibility/providers">
              Providers in Meshery Docs
              <img
                src="/provider/static/img/external-link.svg"
                onError={(e) =>
                  (e.target.src = "/static/img/external-link.svg")
                }
                width="16px"
              />
            </a>
          </div>

          <StyledButton
            onClick={handleModalClose}
            data-cy="providers-modal-button-ok"
            sx={{
              background: KEPPEL,
              marginRight: "1rem",
              color: (theme) => theme.palette.text.inverse,
            }}
          >
            {" "}
            OK
          </StyledButton>
        </CustomDialogActions>
      </CustomDialog>
    </ProviderLayout>
  );
}
