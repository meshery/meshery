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
} from "@sistent/sistent";

import {
  CloseIcon,
  ClickAwayListener,
  DropDownIcon,
} from "@sistent/sistent";

import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

function CustomDialogTitle(props) {
  const { children, onClose, ...other } = props;

  return (
    <DialogTitle sx={{ m: 0, p: 2 }} {...other}>
      {children}

      {onClose && (
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position: "absolute",
            right: "16px",
            top: "16px",
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      )}
    </DialogTitle>
  );
}

CustomDialogTitle.propTypes = {
  children: PropTypes.node,
  onClose: PropTypes.func.isRequired,
};

// Styled Components

const StyledTypography = styled(Typography)(({ theme }) => ({
  fontWeight: 500,
  color: charcoal[100],
  marginBottom: theme.spacing(2),

  "& a": {
    fontWeight: "normal",
  },

  // FIXED: removed invalid spacing in hover selector
  "&:hover": {
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
  const [openModal, setModalOpen] = useState(false);

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
        if (typeof result !== "undefined" && result !== null) {
          setAvailableProviders(result);

          // Optional: auto-select first remote provider
          const remoteProviderKey = Object.keys(result).find(
            (key) => result[key].ProviderType === "remote"
          );

          if (remoteProviderKey) {
            setSelectedProvider(remoteProviderKey);
          }
        }
      },

      (error) => {
        console.error(
          `There was an error fetching providers: ${error}`
        );
      }
    );
  }

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleMenuItemClick = (event, provider) => {
    event.preventDefault();

    setSelectedProvider(provider);
    setIsLoading(true);
    handleClose();

    const existingQueryParams = new URLSearchParams(
      window.location.search
    );

    existingQueryParams.set(
      "provider",
      encodeURIComponent(provider)
    );

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
        onError={(e) => {
          e.target.src =
            "/static/img/meshery-logo/meshery-logo-dark-text-noBG.png";
        }}
        alt="Meshery logo"
      />

      {/* Provider Guidance Link */}
      <StyledTooltip
        title="Providers extend Meshery with identity services, persistence, performance analysis, and collaboration features."
        placement="bottom"
        data-cy="providers-tooltip"
        arrow
      >
        <LearnMore onClick={handleModalOpen}>
          <StyledTypography variant="h6" gutterBottom>
            Learn more about providers
          </StyledTypography>
        </LearnMore>
      </StyledTooltip>

      <CustomDiv>
        {availableProviders && (
          <Fragment>
            <StyledButtonGroup aria-label="provider selection">
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
                    sx={{
                      color: "white",
                      marginRight: 1,
                    }}
                  />
                )}

                {selectedProvider || "Select your provider"}

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
                  id="split-button-menu"
                  autoFocusItem
                  sx={{
                    background: charcoal[20],
                    color: (theme) => theme.palette.text.inverse,
                  }}
                >
                  {Object.keys(availableProviders).map((key) => (
                    <MenuItem
                      key={key}
                      onClick={(e) =>
                        handleMenuItemClick(e, key)
                      }
                      sx={{
                        "&:hover": {
                          backgroundColor: accentGrey[20],
                        },

                        display: "flex",
                        alignItems: "center",
                        justifyContent: "space-between",
                      }}
                    >
                      {key}

                      {/* Layer5 Info Affordance */}
                      {key.toLowerCase() === "layer5" && (
                        <Tooltip
                          title="Layer5 is Meshery's default remote provider. It enables persistent sessions, collaboration, authentication, and Layer5 Cloud integration."
                          placement="right"
                          arrow
                        >
                          <IconButton
                            size="small"
                            aria-label="More information about the Layer5 provider"
                            onClick={(e) => {
                              e.stopPropagation();
                            }}
                            sx={{
                              ml: 0.5,
                              p: 0.25,
                              color: accentGrey[60],

                              "&:hover": {
                                color: KEPPEL,
                              },
                            }}
                          >
                            <InfoOutlinedIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      )}
                    </MenuItem>
                  ))}

                  <Divider
                    sx={{
                      my: 0.5,
                      backgroundColor: accentGrey[40],
                      width: "80%",
                      margin: "auto",
                      marginBottom: 0,
                    }}
                  />

                  <MenuProviderDisabled disabled>
                    Exoscale Labs&nbsp;<span>Offline</span>
                  </MenuProviderDisabled>

                  <MenuProviderDisabled disabled>
                    Equinix US-DAL&nbsp;<span>Offline</span>
                  </MenuProviderDisabled>

                  <MenuProviderDisabled disabled>
                    HPE Security&nbsp;<span>Offline</span>
                  </MenuProviderDisabled>

                  <MenuProviderDisabled disabled>
                    F5 BIG IP iHealth&nbsp;<span>Offline</span>
                  </MenuProviderDisabled>

                  <MenuProviderDisabled disabled>
                    The University of Texas at Austin&nbsp;
                    <span>Offline</span>
                  </MenuProviderDisabled>

                  <Divider
                    sx={{
                      my: 0.5,
                      backgroundColor: accentGrey[40],
                      width: "80%",
                      margin: "auto",
                      marginBottom: 0,
                    }}
                  />

                  <Button
                    component="a"
                    href="https://docs.meshery.io/extensibility/providers"
                    target="_blank"
                    rel="noopener noreferrer"
                    aria-label="Create Provider"
                    data-cy="create_provider"
                    variant="text"
                    sx={{
                      display: "flex",
                      fontStyle: "italic",
                      margin: "auto 0.5rem",
                      textDecoration: "none",
                      cursor: "pointer",
                      width: "100%",
                      backgroundColor: "transparent",
                      justifySelf: "center",
                      textAlign: "center",
                      color: "#ccc",
                      alignItems: "center",

                      "&:hover": {
                        backgroundColor: accentGrey[20],
                        color: "#fff",
                      },
                    }}
                  >
                    Create Your Own Provider&nbsp;

                    <img
                      src="/provider/static/img/external-link.svg"
                      onError={(e) => {
                        e.target.src =
                          "/static/img/external-link.svg";
                      }}
                      width="16px"
                      alt="External link"
                      style={{
                        filter: "brightness(20)",
                      }}
                    />
                  </Button>
                </MenuList>
              </ClickAwayListener>
            </StyledPopover>
          </Fragment>
        )}
      </CustomDiv>

      {/* Provider Modal */}
      <CustomDialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        open={openModal}
        disableScrollLock
        data-cy="providers-modal"
      >
        <StyledCustomDialogTitle
          id="customized-dialog-title"
          onClose={handleModalClose}
        >
          <CustomTypography>
            Choosing a Provider
          </CustomTypography>
        </StyledCustomDialogTitle>

        <StyledDialogBox id="customized-dialog-content">
          Login to Meshery by choosing from the available
          providers. Providers extend Meshery by offering
          identity services, persistence, advanced performance
          analysis, collaboration, and more.

          <h2>Available Providers</h2>

          {Object.keys(availableProviders).map((key) => (
            <React.Fragment
              key={availableProviders[key].provider_name}
            >
              <p style={{ fontWeight: 700 }}>
                {availableProviders[key].provider_name}
              </p>

              <ul>
                {availableProviders[
                  key
                ].provider_description?.map((desc, i) => (
                  <li key={`desc-${i}`}>{desc}</li>
                ))}
              </ul>
            </React.Fragment>
          ))}
        </StyledDialogBox>

        <CustomDialogActions>
          <div className="learnmore">
            <Button
              component="a"
              href="https://docs.meshery.io/extensibility/providers"
              target="_blank"
              rel="noopener noreferrer"
              variant="text"
              sx={{
                color: (theme) =>
                  theme.palette.text.inverse,

                textTransform: "none",
                display: "flex",
                alignItems: "center",
                gap: "0.5rem",
                fontWeight: 400,
                fontSize: "1rem",

                "&:hover": {
                  textDecoration: "underline",
                  backgroundColor: "transparent",
                },
              }}
            >
              Providers in Meshery Docs

              <img
                src="/provider/static/img/external-link.svg"
                onError={(e) => {
                  e.target.src =
                    "/static/img/external-link.svg";
                }}
                width="16px"
                alt="External link"
              />
            </Button>
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
            OK
          </StyledButton>
        </CustomDialogActions>
      </CustomDialog>
    </ProviderLayout>
  );
}
