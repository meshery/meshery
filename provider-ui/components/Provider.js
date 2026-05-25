import React, { Fragment, useEffect, useRef, useState } from "react";
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
  CustomTooltip,
  IconButton,
  CircularProgress,
  InfoOutlinedIcon,
  Popover,
  styled,
  charcoal,
  accentGrey,
  KEPPEL,
} from "@sistent/sistent";
import { CloseIcon, ClickAwayListener, DropDownIcon } from "@sistent/sistent";
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
            right: "16px",
            top: "16px",
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

const PROVIDER_LOGO_SRC = "/provider/static/img/branding/meshery-logo-dark-text.png";
const PROVIDER_LOGO_FALLBACK_SRC = "/static/img/branding/meshery-logo-dark-text.png";
const EXTERNAL_LINK_ICON_SRC = "/provider/static/img/icons/external-link.svg";
const EXTERNAL_LINK_ICON_FALLBACK_SRC = "/static/img/icons/external-link.svg";

export default function Provider() {
  const [anchorEl, setAnchorEl] = useState(null);
  // Each entry: { ...ProviderProperties, _status: "checking"|"online"|"offline", _error?: string }
  // _-prefixed keys are client-only metadata stitched in from the SSE event
  // envelope; they never collide with server-emitted ProviderProperties fields.
  const [availableProviders, setAvailableProviders] = useState({});
  const [selectedProvider, setSelectedProvider] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  /* eslint-disable no-unused-vars */
  const [openMenu, setOpenMenu] = useState(false);
  const [openModal, setModalOpen] = React.useState(false);
  // Per-provider info popover state. The info button next to each remote
  // entry opens this with the provider's name and the providerDescription
  // list returned by that provider's /capabilities, so the chooser surfaces
  // exactly what each remote describes about itself rather than a generic
  // hardcoded blurb.
  const [providerInfo, setProviderInfo] = useState({
    anchor: null,
    provider: null,
  });
  const eventSourceRef = useRef(null);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? "popover" : undefined;

  const openProviderInfo = (event, provider) => {
    event.stopPropagation();
    setProviderInfo({ anchor: event.currentTarget, provider });
  };
  const closeProviderInfo = () => {
    setProviderInfo({ anchor: null, provider: null });
  };
  const providerInfoOpen = Boolean(providerInfo.anchor);

  useEffect(() => {
    // Prefer the SSE stream so the chooser renders the local provider
    // immediately and each remote updates independently as its probe
    // settles on the server. Fall back to the legacy /api/providers
    // polling endpoint if the browser does not implement EventSource or
    // if the stream errors out before delivering its first event - that
    // keeps older clients and proxies that strip text/event-stream usable.
    if (typeof EventSource === "undefined") {
      loadProvidersFromServer();
      return undefined;
    }

    const es = new EventSource("/api/providers/stream", { withCredentials: true });
    eventSourceRef.current = es;
    let receivedAny = false;
    let fellBack = false;

    const handleEvent = (e) => {
      try {
        const evt = JSON.parse(e.data);
        if (!evt || !evt.key) return;
        receivedAny = true;
        setAvailableProviders((prev) => ({
          ...prev,
          [evt.key]: {
            ...(evt.properties || {}),
            _status: evt.status,
            _error: evt.error || "",
          },
        }));
      } catch (err) {
        console.error("provider stream: failed to parse event", err);
      }
    };

    // The server names its events `provider`; default `message` listeners
    // would miss them. Subscribe to both for compatibility with any
    // intermediate proxy that drops the event-name line.
    es.addEventListener("provider", handleEvent);
    es.onmessage = handleEvent;
    es.onerror = (err) => {
      console.warn("provider stream: error", err);
      // The browser will auto-reconnect EventSource. If we still have
      // nothing to render after the first error, fall back so the
      // chooser is not stuck on a blank state.
      if (!receivedAny && !fellBack) {
        fellBack = true;
        loadProvidersFromServer();
      }
    };

    return () => {
      es.close();
      eventSourceRef.current = null;
    };
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
          // Polling fallback: stamp _status from capabilities presence so
          // the chooser renders without an SSE feed. Remotes with no
          // capabilities are treated as offline, matching how the SSE
          // stream would have classified them.
          const next = {};
          Object.keys(result).forEach((key) => {
            const p = result[key] || {};
            const isRemote = p.providerType === "remote";
            next[key] = {
              ...p,
              _status: isRemote
                ? p.capabilities?.length
                  ? "online"
                  : "offline"
                : "online",
            };
          });
          setAvailableProviders(next);
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
        src={PROVIDER_LOGO_SRC}
        onError={(e) => (e.target.src = PROVIDER_LOGO_FALLBACK_SRC)}
        alt="logo"
      />
      <CustomDiv>
        {availableProviders !== "" && (
          <Fragment>
            <CustomTooltip
              title="Learn more about Meshery remote providers"
              placement="bottom"
              data-cy="providers-tooltip"
              arrow
            >
              <LearnMore
                href="#provider-guidance-dialog"
                onClick={(e) => {
                  e.preventDefault();
                  handleModalOpen();
                }}
                aria-haspopup="dialog"
                aria-controls="provider-guidance-dialog"
              >
                Learn more about providers
              </LearnMore>
            </CustomTooltip>
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
                    minWidth: 260,
                  }}
                  id="split-button-menu"
                  autoFocusItem
                >
                  {Object.keys(availableProviders).map((key) => {
                    // /api/providers and /api/providers/stream both
                    // serialize ProviderProperties with camelCase JSON
                    // tags; the stream also stitches in _status / _error
                    // markers from the ProviderStatusEvent envelope.
                    const provider = availableProviders[key];
                    const isRemote = provider?.providerType === "remote";
                    const status = provider?._status || "online";
                    if (status === "offline") return null;
                    // Prefer the explicit providerName so the chooser shows
                    // the friendly name rather than the registration key
                    // (which can fall back to the URL host for remotes that
                    // were unreachable when the server booted).
                    const label = provider?.providerName || key;
                    const isChecking = status === "checking";
                    return (
                      <MenuItem
                        key={key}
                        onClick={(e) =>
                          isChecking ? undefined : handleMenuItemClick(e, key)
                        }
                        disabled={isChecking}
                        sx={{
                          "&:hover": { backgroundColor: accentGrey[20] },
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "space-between",
                          gap: 1,
                          opacity: isChecking ? 0.7 : 1,
                        }}
                      >
                        <span>{label}</span>
                        <span style={{ display: "flex", alignItems: "center", gap: 8 }}>
                          {isChecking && (
                            <CircularProgress
                              size={14}
                              sx={{ color: accentGrey[40] }}
                              aria-label="Verifying availability"
                            />
                          )}
                          {isRemote && (
                            <CustomTooltip
                              title="More information about this provider"
                              placement="right"
                              arrow
                            >
                              <IconButton
                                size="small"
                                aria-label="More information about the remote provider"
                                data-testid="provider-learn-more-button"
                                onClick={(e) => openProviderInfo(e, provider)}
                                sx={{
                                  ml: "auto",
                                  p: 0.25,
                                  flexShrink: 0,
                                  color: KEPPEL,
                                  opacity: 0.95,
                                  "&:hover": { color: KEPPEL, opacity: 1 },
                                }}
                              >
                                <InfoOutlinedIcon width={18} height={18} />
                              </IconButton>
                            </CustomTooltip>
                          )}
                        </span>
                      </MenuItem>
                    );
                  })}
                  <Divider
                    sx={{
                      my: 0.5,
                      backgroundColor: accentGrey[40],
                      width: "80%",
                      margin: "auto",
                      marginBottom: "0px",
                    }}
                  />
                  {Object.keys(availableProviders).map((key) => {
                    const provider = availableProviders[key];
                    if (provider?._status !== "offline") return null;
                    const label = provider?.providerName || key;
                    return (
                      <MenuProviderDisabled disabled={true} key={key}>
                        {label}{"\u00A0"}
                        <span>Offline</span>
                      </MenuProviderDisabled>
                    );
                  })}
                  {/*
                    The chooser is now strictly driven by the server's
                    /api/providers/stream feed (see useEffect above);
                    static "Offline" entries for providers the server
                    does not know about have been removed so the menu
                    cannot misrepresent Meshery's actual provider state.
                  */}
                  <Divider
                    sx={{
                      my: 0.5,
                      backgroundColor: accentGrey[40],
                      width: "80%",
                      margin: "auto",
                      marginBottom: "0px",
                    }}
                  />
                  {/* Use Sistent's Button as a link */}
                  <Button
                    component="a"
                    href="https://docs.meshery.io/extensibility/providers"
                    aria-describedby={id}
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
                      backgroundColor: "none",
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
                      src={EXTERNAL_LINK_ICON_SRC}
                      onError={(e) => (e.target.src = EXTERNAL_LINK_ICON_FALLBACK_SRC)}
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

            {/*
              Per-provider info popover. Anchored to the InfoOutlined icon
              next to each remote entry. Renders the descriptive blurb the
              remote returns from its own /capabilities (providerName +
              providerDescription) so the chooser explains what each
              specific provider offers instead of showing the legacy
              one-size-fits-all "Choosing a Provider" dialog.
            */}
            <Popover
              open={providerInfoOpen}
              anchorEl={providerInfo.anchor}
              onClose={closeProviderInfo}
              anchorOrigin={{ vertical: "center", horizontal: "right" }}
              transformOrigin={{ vertical: "center", horizontal: "left" }}
              PaperProps={{
                sx: {
                  background: charcoal[20],
                  color: (theme) => theme.palette.text.inverse,
                  maxWidth: 360,
                  p: 2,
                  borderLeft: `3px solid ${KEPPEL}`,
                },
              }}
            >
              {providerInfo.provider && (
                <div>
                  <div style={{ fontWeight: 600, fontSize: "1rem" }}>
                    {providerInfo.provider.providerName ||
                      "Remote Provider"}
                  </div>
                  {providerInfo.provider.providerUrl && (
                    <div
                      style={{
                        fontSize: "0.75rem",
                        opacity: 0.7,
                        marginTop: 2,
                        marginBottom: 8,
                        wordBreak: "break-all",
                      }}
                    >
                      {providerInfo.provider.providerUrl}
                    </div>
                  )}
                  {Array.isArray(providerInfo.provider.providerDescription) &&
                  providerInfo.provider.providerDescription.length > 0 ? (
                      <ul style={{ paddingLeft: "1.1rem", margin: 0 }}>
                        {providerInfo.provider.providerDescription.map(
                          (line, i) => (
                            <li
                              key={`${providerInfo.provider.providerName || "p"}-${i}`}
                              style={{ marginBottom: 4 }}
                            >
                              {line}
                            </li>
                          )
                        )}
                      </ul>
                    ) : (
                      <div style={{ fontSize: "0.85rem", opacity: 0.85 }}>
                        {providerInfo.provider._status === "checking"
                          ? "Verifying availability and loading capabilities..."
                          : "This provider has not published a description."}
                      </div>
                    )}
                </div>
              )}
            </Popover>
          </Fragment>
        )}
      </CustomDiv>
      <CustomDialog
        onClose={handleModalClose}
        aria-labelledby="customized-dialog-title"
        id="provider-guidance-dialog"
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
            const provider = availableProviders[key] || {};
            // /api/providers serializes ProviderProperties with camelCase
            // JSON tags (server/models/providers.go). The previous
            // snake_case lookup silently rendered nothing - producing an
            // empty list under "Available Providers" - because no field
            // by those names exists in the response.
            const name = provider.providerName || key;
            const description = Array.isArray(provider.providerDescription)
              ? provider.providerDescription
              : [];
            return (
              <React.Fragment key={key}>
                <p style={{ fontWeight: 700 }}>{name}</p>
                {description.length > 0 ? (
                  <ul>
                    {description.map((desc, i) => (
                      <li key={`${key}-desc-${i}`}>{desc}</li>
                    ))}
                  </ul>
                ) : (
                  <p style={{ fontStyle: "italic", opacity: 0.8 }}>
                    {provider._status === "checking"
                      ? "Loading description..."
                      : "No description available."}
                  </p>
                )}
              </React.Fragment>
            );
          })}
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
                color: (theme) => theme.palette.text.inverse,
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
                src={EXTERNAL_LINK_ICON_SRC}
                onError={(e) => (e.target.src = EXTERNAL_LINK_ICON_FALLBACK_SRC)}
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
            {" "}
            OK
          </StyledButton>
        </CustomDialogActions>
      </CustomDialog>
    </ProviderLayout>
  );
}
