/* eslint-disable react/prop-types */
import React, { useEffect, useRef, useState } from "react";
import {
  IconButton,
  Avatar,
  Popper,
  MenuList,
  Grow,
  MenuItem,
  ListItemText,
  Paper,
  ClickAwayListener,
  NoSsr,
} from "@mui/material";
import ExtensionPointSchemaValidator from "../utils/extensionPointSchemaValidator";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import Link from "next/link";
import { useRouter, withRouter } from "next/router";
import dataFetch from "../lib/data-fetch";
import { updateUser } from "@/features/user";

function exportToJsonFile(jsonData, filename) {
  let dataStr = JSON.stringify(jsonData);
  let dataUri = "data:application/json;charset=utf-8," + encodeURIComponent(dataStr);

  let exportFileDefaultName = filename;

  let linkElement = document.createElement("a");
  linkElement.setAttribute("href", dataUri);
  linkElement.setAttribute("download", exportFileDefaultName);
  linkElement.click();
  linkElement.remove();
}

function ExtensionPointContent({ href, name }) {
  const content = <ListItemText>{name}</ListItemText>;

  if (href) return <Link href={href}>{content}</Link>;

  return content;
}

/**
 * @param {import("../utils/ExtensionPointSchemaValidator").AccountSchema[]} children
 */
function AccountExtension({ account }) {
  if (account && account.length > 0) {
    return account.map(({ href, title }) => {
      return (
        <MenuItem key={title + href}>
          <ExtensionPointContent href={href} name={title} />
        </MenuItem>
      );
    });
  } else {
    return null;
  }
}

const UserComponent = ({ updateUser, color }) => {
  const [user, setUser] = useState(null);
  const [open, setOpen] = useState(false);
  const [account, setAccount] = useState(ExtensionPointSchemaValidator("account")());
  // eslint-disable-next-line no-unused-vars
  const [providerType, setProviderType] = useState("");
  const anchorEl = useRef();
  const router = useRouter();

  useEffect(() => {
    dataFetch(
      "/api/user",
      {
        credentials: "same-origin",
      },
      (user) => {
        setUser(user);
        updateUser({ user });
      },
      (error) => ({
        error,
      })
    );

    dataFetch(
      "/api/provider/capabilities",
      {
        credentials: "same-origin",
        method: "GET",
      },
      (result) => {
        if (result) {
          setAccount(ExtensionPointSchemaValidator("account")(result?.extensions?.account));
          setProviderType(result?.provider_type);
        }
      },
      (err) => console.error(err)
    );
  }, []);

  const handleToggle = () => {
    setOpen(!open);
  };

  const handleClose = (event) => {
    if (anchorEl.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  const handleLogout = () => {
    window.location = "/user/logout";
  };

  const handlePreference = () => {
    router.push("/user/preferences");
  };

  const handleGetToken = () => {
    dataFetch(
      "/api/token",
      { credentials: "same-origin" },
      (data) => {
        exportToJsonFile(data, "auth.json");
      },
      (error) => ({ error })
    );
  };

  let avatar_url;
  if (user && user !== null) {
    avatar_url = user.avatar_url;
  }

  return (
    <div>
      <NoSsr>
        <div data-test="profile-button">
          <IconButton
            color={color}
            sx={{ padding: "4px" }}
            ref={anchorEl}
            aria-owns={open ? "menu-list-grow" : undefined}
            aria-haspopup="true"
            onClick={handleToggle}
          >
            <Avatar src={avatar_url} />
          </IconButton>
        </div>
        <Popper open={open} anchorEl={anchorEl.current} transition sx={{ zIndex: 10000 }} placement="top-end">
          {({ TransitionProps, placement }) => (
            <Grow
              {...TransitionProps}
              id="menu-list-grow"
              sx={{ transformOrigin: placement === "bottom" ? "left top" : "left bottom" }}
            >
              <Paper sx={{ color: "black" }}>
                <ClickAwayListener onClickAway={handleClose}>
                  <MenuList>
                    {account && account.length ? <AccountExtension account={account} /> : null}
                    <MenuItem onClick={handleGetToken}>Get Token</MenuItem>
                    <MenuItem onClick={handlePreference}>Preferences</MenuItem>
                    <MenuItem onClick={handleLogout}>Logout</MenuItem>
                  </MenuList>
                </ClickAwayListener>
              </Paper>
            </Grow>
          )}
        </Popper>
      </NoSsr>
    </div>
  );
};

const mapDispatchToProps = (dispatch) => ({ updateUser: bindActionCreators(updateUser, dispatch) });

export default connect(null, mapDispatchToProps)(withRouter(UserComponent));
