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
  ProviderInfoPopover,
  ProviderInfoSectionRule,
  ProviderInfoSectionHeading,
  ProviderInfoDescriptionList,
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
  Box,
  Chip,
  Link,
  Typography,
  ExternalLinkIcon,
  LockIcon,
  styled,
  charcoal,
  accentGrey,
  CARIBBEAN_GREEN,
  KEPPEL,
  SAFFRON,
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

// Order matters: the first matching entry wins, so more specific prefixes
// must come before broader ones (e.g. "persist-meshery-pattern" before
// "persist-"). Returning a stable category label per feature lets the info
// popover collapse 25+ raw capability names into a handful of meaningful
// chips, which is what users actually want to see at a glance.
const CAPABILITY_CATEGORY_RULES = [
  { label: "Identity", match: (f) => f.startsWith("users-") || f === "organizations" },
  { label: "Workspaces", match: (f) => f === "environments" || f === "workspaces" },
  { label: "Catalog", match: (f) => f.endsWith("-catalog") },
  {
    label: "Content",
    match: (f) =>
      f.startsWith("persist-meshery-pattern") ||
      f.startsWith("persist-meshery-filter") ||
      f.startsWith("clone-meshery-") ||
      f === "share-designs",
  },
  {
    label: "Performance",
    match: (f) =>
      f === "persist-results" ||
      f === "persist-result" ||
      f === "persist-metrics" ||
      f === "persist-schedules" ||
      f === "persist-performance-profiles" ||
      f === "persist-smp-test-profile",
  },
  { label: "Connections", match: (f) => f === "persist-connection" || f === "persist-credentials" },
  { label: "Events", match: (f) => f === "persist-events" },
  { label: "Sync", match: (f) => f === "sync-prefs" },
];

// Group a ProviderProperties.capabilities[] payload into category counts.
// Returns a sorted array of { label, count } so the popover renders chips in
// a deterministic order regardless of how the remote ordered its
// /capabilities response. Anything that does not match a known prefix is
// bucketed under "Other" so we surface the full count without losing it.
function summarizeCapabilities(capabilities) {
  if (!Array.isArray(capabilities) || capabilities.length === 0) return [];
  const counts = new Map();
  capabilities.forEach((c) => {
    const feature = c?.feature || "";
    if (!feature) return;
    const rule = CAPABILITY_CATEGORY_RULES.find((r) => r.match(feature));
    const label = rule ? rule.label : "Other";
    counts.set(label, (counts.get(label) || 0) + 1);
  });
  const order = CAPABILITY_CATEGORY_RULES.map((r) => r.label).concat("Other");
  return order
    .filter((label) => counts.has(label))
    .map((label) => ({ label, count: counts.get(label) }));
}

// Returns the label to render for a provider. Prefers the friendly
// providerName the remote published in /capabilities; falls back to the
// registration key lowercased - that key is the URL host when the remote
// was unreachable at boot (so it is a URL fragment and should not carry
// arbitrary uppercase the user typed into PROVIDER_BASE_URLS).
function formatProviderLabel(key, provider) {
  const name = provider?.providerName;
  if (name) return name;
  return (key || "").toLowerCase();
}

// Count navigator-extension components advertised by the provider. The
// payload nests them as { extensions: { collaborator: [...], ... } }; a
// total across known buckets gives the popover a single Extensions chip
// without needing to enumerate every component path.
function countExtensionComponents(extensions) {
  if (!extensions || typeof extensions !== "object") return 0;
  return Object.values(extensions).reduce((total, list) => {
    return total + (Array.isArray(list) ? list.length : 0);
  }, 0);
}

// Two-tone "shield" chip used for the per-category capability counts.
// The left half carries the category label on a KEPPEL-tinted background;
// the right half carries the count on a darker KEPPEL fill so the number
// reads as a distinct value, not part of the label text. Built from raw
// Boxes (not MUI Chip) because Chip renders a single contiguous fill.
function CapabilityChip({ label, count }) {
  return (
    <Box
      role="img"
      aria-label={`${count} ${label} capabilities`}
      sx={{
        // flexShrink:0 + whiteSpace:nowrap keep each chip at its natural
        // content width. Without them flexbox shrinks the chip below the
        // label's intrinsic width and the inner Box clips ("Identit"
        // instead of "Identity"). The parent Stack uses flexWrap:"wrap",
        // so non-shrinking chips push to the next row when there's no
        // horizontal room - which is the desired behavior.
        display: "inline-flex",
        flexShrink: 0,
        maxWidth: "100%",
        minWidth: 0,
        alignItems: "stretch",
        borderRadius: "12px",
        overflow: "hidden",
        border: "1px solid rgba(0,179,159,0.35)",
        fontSize: "0.7rem",
        fontWeight: 600,
        lineHeight: 1.4,
        whiteSpace: "nowrap",
      }}
    >
      <Box
        sx={{
          background: "rgba(0,179,159,0.12)",
          color: KEPPEL,
          padding: "1px 8px",
        }}
      >
        {label}
      </Box>
      <Box
        sx={{
          background: "rgba(0,179,159,0.32)",
          color: "#fff",
          padding: "1px 8px",
          borderLeft: "1px solid rgba(0,179,159,0.35)",
        }}
      >
        {count}
      </Box>
    </Box>
  );
}

CapabilityChip.propTypes = {
  label: PropTypes.string.isRequired,
  count: PropTypes.number.isRequired,
};

// Renders the body of the per-provider info popover. Defined at module
// scope (not as an IIFE inside the parent JSX) so the layout, derived
// values, and per-section visibility rules are testable in isolation and
// readable apart from the chooser's state plumbing. Returns null when
// `provider` is absent so the parent can mount the component
// unconditionally inside an empty Popover.
function ProviderInfoContent({ provider }) {
  if (!provider) return null;

  const status = provider._status || "online";
  const totalCaps = Array.isArray(provider.capabilities)
    ? provider.capabilities.length
    : 0;
  const categories = summarizeCapabilities(provider.capabilities);
  const extensionCount = countExtensionComponents(provider.extensions);
  const restricted = Boolean(provider.restrictedAccess?.isMesheryUIRestricted);
  const statusColor =
    status === "online"
      ? CARIBBEAN_GREEN
      : status === "checking"
        ? SAFFRON
        : accentGrey[40];
  const statusLabel =
    status === "checking"
      ? "Checking"
      : status === "offline"
        ? "Offline"
        : "Online";
  const hasDescription =
    Array.isArray(provider.providerDescription) &&
    provider.providerDescription.length > 0;
  const hasMetaChips = Boolean(provider.packageVersion) || totalCaps > 0 || restricted;

  return (
    <Box sx={{ width: "100%" }}>
      {/* Use an explicit three-column grid so the name, URL, and status chip
          stay on a single visual row. Flexbox let the centered URL fall onto
          a second line in constrained widths, which is what the review
          screenshot called out. */}
      <Box
        sx={{
          width: "100%",
          display: "grid",
          gridTemplateColumns: "auto minmax(0, 1fr) auto",
          alignItems: "center",
          columnGap: 1,
        }}
      >
        <Typography
          component="div"
          sx={{
            fontWeight: 700,
            fontSize: "1.05rem",
            lineHeight: 1.3,
            flexShrink: 0,
            minWidth: 0,
          }}
        >
          {provider.providerName || "Remote Provider"}
        </Typography>
        {provider.providerUrl ? (
          <Link
            href={provider.providerUrl}
            target="_blank"
            rel="noopener noreferrer"
            underline="hover"
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 0.5,
              minWidth: 0,
              justifyContent: "center",
              justifySelf: "center",
              width: "100%",
              maxWidth: "100%",
              fontSize: "0.75rem",
              color: "inherit",
              opacity: 0.75,
              whiteSpace: "nowrap",
              overflow: "hidden",
              textOverflow: "ellipsis",
              "&:hover": { opacity: 1, color: KEPPEL },
            }}
          >
            <Box
              component="span"
              sx={{
                display: "block",
                minWidth: 0,
                overflow: "hidden",
                textOverflow: "ellipsis",
                whiteSpace: "nowrap",
              }}
            >
              {provider.providerUrl.toLowerCase()}
            </Box>
            <ExternalLinkIcon
              width={12}
              height={12}
              aria-hidden="true"
              style={{ flexShrink: 0 }}
            />
          </Link>
        ) : (
          // Keep the grid skeleton consistent so the status chip stays
          // anchored on the right whether or not a URL is present.
          <Box aria-hidden="true" sx={{ minWidth: 0 }} />
        )}
        <Chip
          size="small"
          label={statusLabel}
          aria-label={`Provider status: ${statusLabel}`}
          icon={
            <span
              aria-hidden="true"
              style={{
                width: 8,
                height: 8,
                borderRadius: "50%",
                background: statusColor,
                marginLeft: 8,
              }}
            />
          }
          sx={{
            flexShrink: 0,
            justifySelf: "end",
            background: "rgba(255,255,255,0.08)",
            color: "inherit",
            fontWeight: 600,
            letterSpacing: "0.02em",
            ".MuiChip-icon": { ml: 1 },
          }}
        />
      </Box>

      {hasMetaChips && (
        <>
          <ProviderInfoSectionRule />
          <Box
            sx={{
              display: "flex",
              flexWrap: "wrap",
              gap: 0.75,
              alignItems: "center",
              width: "100%",
              minWidth: 0,
            }}
          >
            {provider.packageVersion && (
              <Chip
                size="small"
                label={provider.packageVersion}
                aria-label={`Package version ${provider.packageVersion}`}
                sx={{
                  background: "rgba(255,255,255,0.08)",
                  color: "inherit",
                  fontFamily:
                    "ui-monospace, SFMono-Regular, Menlo, monospace",
                  fontSize: "0.7rem",
                }}
              />
            )}
            {totalCaps > 0 && (
              <Chip
                size="small"
                label={`${totalCaps} capabilities`}
                aria-label={`${totalCaps} capabilities`}
                sx={{
                  background: "rgba(255,255,255,0.08)",
                  color: "inherit",
                  fontSize: "0.7rem",
                }}
              />
            )}
            {restricted && (
              <CustomTooltip
                title="This provider restricts which Meshery UI components are available"
                placement="top"
                arrow
              >
                <Chip
                  size="small"
                  icon={
                    <LockIcon width={12} height={12} aria-hidden="true" />
                  }
                  label="Restricted UI"
                  aria-label="Restricted UI"
                  sx={{
                    background: "rgba(255,184,0,0.12)",
                    color: SAFFRON,
                    fontSize: "0.7rem",
                    ".MuiChip-icon": { color: SAFFRON },
                  }}
                />
              </CustomTooltip>
            )}
          </Box>
        </>
      )}

      {categories.length > 0 && (
        <>
          <ProviderInfoSectionRule />
          <ProviderInfoSectionHeading component="h4">
            Capabilities
          </ProviderInfoSectionHeading>
          {/* width:100% so the wrap boundary is the popover Paper, not
              the chips' intrinsic content width - otherwise a long row
              of chips overflows past the Paper's right edge instead of
              wrapping. */}
          <Box
            sx={{
              display: "flex",
              flexWrap: "wrap",
              gap: 0.5,
              width: "100%",
              minWidth: 0,
            }}
          >
            {categories.map(({ label, count }) => (
              <CapabilityChip key={label} label={label} count={count} />
            ))}
          </Box>
        </>
      )}

      <ProviderInfoSectionRule />
      <ProviderInfoSectionHeading component="h4">
        What this provider offers
      </ProviderInfoSectionHeading>
      {hasDescription ? (
        <ProviderInfoDescriptionList>
          {provider.providerDescription.map((line, i) => (
            <li key={`${provider.providerName || "p"}-${i}`}>{line}</li>
          ))}
        </ProviderInfoDescriptionList>
      ) : (
        <Typography
          component="div"
          sx={{ fontSize: "0.85rem", opacity: 0.8, fontStyle: "italic" }}
        >
          {status === "checking"
            ? "Verifying availability and loading capabilities..."
            : "This provider has not published a description."}
        </Typography>
      )}

      {extensionCount > 0 && (
        <>
          <ProviderInfoSectionRule />
          <ProviderInfoSectionHeading component="h4">
            Extensions
          </ProviderInfoSectionHeading>
          {/* Chips sit on their own row below the section heading, matching
              the Capabilities section pattern (heading on its own line,
              then chips beneath). */}
          <Box
            sx={{
              display: "flex",
              flexWrap: "wrap",
              gap: 0.75,
              width: "100%",
              minWidth: 0,
              mt: 0.25,
            }}
          >
            {Object.keys(provider.extensions || {}).map((slot) => {
              const list = provider.extensions[slot];
              const n = Array.isArray(list) ? list.length : 0;
              if (n === 0) return null;
              return (
                <Chip
                  key={slot}
                  size="small"
                  label={n > 1 ? `${slot} · ${n}` : slot}
                  aria-label={`${slot} extension`}
                  sx={{
                    background: "rgba(255,255,255,0.08)",
                    color: "inherit",
                    fontSize: "0.7rem",
                  }}
                />
              );
            })}
          </Box>
        </>
      )}
    </Box>
  );
}

ProviderInfoContent.propTypes = {
  provider: PropTypes.object,
};

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
        {Object.keys(availableProviders).length > 0 && (
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
                sx={{
                  // Push spinner to the left edge and the DropDownIcon to
                  // the right edge with the label centered between them.
                  // 16px (px:2 in MUI's 8px spacing) padding lives on the
                  // button so the icons sit just inside the rounded ends.
                  px: 2,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "space-between",
                  textTransform: "none",
                  minWidth: 240,
                }}
              >
                {/* Leading slot: spinner when loading, otherwise an
                    invisible spacer the same width so the centered
                    label does not shift as loading toggles. */}
                {isLoading ? (
                  <CircularProgress
                    size={20}
                    sx={{ color: "white" }}
                    aria-label="Loading"
                  />
                ) : (
                  <span
                    aria-hidden="true"
                    style={{ width: 20, height: 20, display: "inline-block" }}
                  />
                )}
                <span
                  style={{
                    flex: 1,
                    textAlign: "center",
                    padding: "0 8px",
                    whiteSpace: "nowrap",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                  }}
                >
                  {selectedProvider !== ""
                    ? formatProviderLabel(
                      selectedProvider,
                      availableProviders[selectedProvider]
                    )
                    : "Select your provider"}
                </span>
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
                    // Background is provided by StyledPopover's Paper so
                    // the rounded corners clip a single continuous surface
                    // (avoids the square artifact at each corner where a
                    // mismatched MenuList background bleeds past the
                    // Paper's border-radius).
                    color: (theme) => theme.palette.text.inverse,
                    minWidth: 260,
                  }}
                  id="split-button-menu"
                  autoFocusItem
                >
                  {(() => {
                    // Local is pinned to the top of the chooser with a
                    // divider immediately after, so the built-in option is
                    // visually distinct from the remote providers below.
                    // Match Local by an explicit `providerType === "local"`
                    // (positive) rather than `!== "remote"` (negative).
                    // The negative form pulled in providers whose
                    // capabilities had not yet loaded (providerType absent),
                    // putting unreachable remotes in the Local bucket.
                    // There are only two categories of providers: Local and
                    // Remote.
                    const entries = Object.keys(availableProviders).map(
                      (key) => ({ key, provider: availableProviders[key] })
                    );
                    const localEntries = entries.filter(
                      (e) => e.provider?.providerType === "local"
                    );
                    const remoteEntries = entries.filter(
                      (e) => e.provider?.providerType !== "local"
                    );
                    const renderRow = ({ key, provider }) => {
                      const status = provider?._status || "online";
                      if (status === "offline") return null;
                      const label = formatProviderLabel(key, provider);
                      const isChecking = status === "checking";
                      return (
                        <MenuItem
                          key={key}
                          onClick={(e) =>
                            isChecking
                              ? undefined
                              : handleMenuItemClick(e, key)
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
                          <span
                            style={{
                              display: "flex",
                              alignItems: "center",
                              gap: 8,
                            }}
                          >
                            {isChecking && (
                              <CircularProgress
                                size={14}
                                sx={{ color: accentGrey[40] }}
                                aria-label="Verifying availability"
                              />
                            )}
                            <IconButton
                              size="small"
                              aria-label={`More information about ${label}`}
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
                          </span>
                        </MenuItem>
                      );
                    };
                    const renderSectionLabel = (label) => (
                      <Typography
                        component="div"
                        sx={{
                          px: 2,
                          pt: 1,
                          pb: 0.5,
                          fontSize: "0.68rem",
                          fontWeight: 700,
                          letterSpacing: "0.06em",
                          textTransform: "uppercase",
                          color: "inherit",
                          opacity: 0.65,
                        }}
                      >
                        {label}
                      </Typography>
                    );
                    return (
                      <>
                        {localEntries.length > 0 && renderSectionLabel("Local")}
                        {localEntries.map(renderRow)}
                        {localEntries.length > 0 && remoteEntries.length > 0 && (
                          <Divider
                            sx={{
                              mt: 0.75,
                              mb: 0,
                              mx: 0,
                              backgroundColor: accentGrey[40],
                              opacity: 0.45,
                              width: "100%",
                            }}
                          />
                        )}
                        {remoteEntries.length > 0 && renderSectionLabel("Remote")}
                        {remoteEntries.map(renderRow)}
                      </>
                    );
                  })()}
                  {Object.values(availableProviders).some(
                    (provider) => provider?._status === "offline"
                  ) && (
                    <>
                      <Divider
                        sx={{
                          mt: 0.75,
                          mb: 0,
                          mx: 0,
                          backgroundColor: accentGrey[40],
                          opacity: 0.45,
                          width: "100%",
                        }}
                      />
                      <Typography
                        component="div"
                        sx={{
                          px: 2,
                          pt: 1,
                          pb: 0.5,
                          fontSize: "0.68rem",
                          fontWeight: 700,
                          letterSpacing: "0.06em",
                          textTransform: "uppercase",
                          color: (theme) => theme.palette.text.inverse,
                          opacity: 0.65,
                        }}
                      >
                        Offline
                      </Typography>
                      {Object.keys(availableProviders).map((key) => {
                        const provider = availableProviders[key];
                        if (provider?._status !== "offline") return null;
                        const label = formatProviderLabel(key, provider);
                        return (
                          <MenuProviderDisabled disabled={true} key={key}>
                            {label}{"\u00A0"}
                            <span>Offline</span>
                          </MenuProviderDisabled>
                        );
                      })}
                    </>
                  )}
                  {/*
                    The chooser is now strictly driven by the server's
                    /api/providers/stream feed (see useEffect above);
                    static "Offline" entries for providers the server
                    does not know about have been removed so the menu
                    cannot misrepresent Meshery's actual provider state.
                  */}
                  <Divider
                    sx={{
                      mt: 0.5,
                      mb: 0,
                      mx: "auto",
                      backgroundColor: accentGrey[40],
                      width: "80%",
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
              next to each remote entry. Renders a structured summary of
              what the remote published from its own /capabilities -
              identity (name, URL, status), package version, restricted-UI
              badge, the categorized capability count, the human-readable
              providerDescription bullets, and any navigator extensions -
              so the chooser surfaces every worthwhile field the remote
              advertises instead of just the name/URL/blurb triplet.
            */}
            <ProviderInfoPopover
              open={providerInfoOpen}
              anchorEl={providerInfo.anchor}
              onClose={closeProviderInfo}
              anchorOrigin={{ vertical: "center", horizontal: "right" }}
              transformOrigin={{ vertical: "center", horizontal: "left" }}
              data-testid="provider-info-popover"
            >
              <ProviderInfoContent provider={providerInfo.provider} />
            </ProviderInfoPopover>
          </Fragment>
        )}
      </CustomDiv>
      {Object.keys(availableProviders).length > 0 && (
        // LearnMore is a sibling of (not a child of) CustomDiv so the
        // logo -> link -> dropdown layout stacks as three blocks with
        // balanced vertical spacing; LearnMore's own marginTop/marginBottom
        // centers it within the gap between the logo and the dropdown.
        <LearnMore
          href="#provider-guidance-dialog"
          onClick={(e) => {
            e.preventDefault();
            handleModalOpen();
          }}
          aria-haspopup="dialog"
          aria-controls="provider-guidance-dialog"
          data-cy="providers-learn-more"
        >
          Learn more about providers
        </LearnMore>
      )}
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
          <p style={{ marginTop: 0 }}>
            A <strong>Provider</strong> extends Meshery with services beyond
            what the local install offers - identity and multi-tenancy,
            long-term persistence, design and performance history, sharing,
            and team collaboration. Meshery ships with a built-in{" "}
            <strong>Local</strong> provider; remote providers (such as Meshery
            Cloud) layer additional services on top.
          </p>
          <p>
            Per-provider details are available from the info{" "}
            <span aria-hidden="true">&#9432;</span> icon next to each entry in
            the chooser, populated live from that provider&apos;s own{" "}
            <code>/capabilities</code>.
          </p>
          <p>
            <strong>Pre-selecting a provider.</strong> To skip this chooser on
            launch, set the <code>PROVIDER</code> environment variable to the
            provider&apos;s name (for example <code>PROVIDER=Local</code> or{" "}
            <code>PROVIDER=Meshery</code>) before starting the server. An
            invalid value is ignored and the chooser is shown as usual.
          </p>
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
