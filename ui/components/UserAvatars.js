import React, { useState } from "react";
import Avatar from "@mui/material/Avatar";
import AvatarGroup from "@mui/material/AvatarGroup";
import Tooltip from "@mui/material/Tooltip";
import { withStyles } from '@material-ui/core/styles';

const styles = () => ({
  smallAvatar : {
    width : "30px !important",
    height : "30px !important",
    cursor : "pointer"
  },
  customTooltip : {
    display : "flex",
    top : "0",
    justifyContent : "center",
    backgroundColor : "#263238"
  },
  customArrow : {
    color : "#263238"
  }
});

const UserAvatars = ({ classes }) => {
  const [users,] = useState([{ name : "A", clientID : 1, avatar_url : "" },{ name : "B", clientID : 2, avatar_url : "" }])
  return (
    <div className={classes.avatars}>
      <AvatarGroup max={4}>
        {Object.entries(users).map((user) => (
          <Tooltip
            key={user.clientID}
            title={user.name}
            arrow
            classes={{
              tooltip : classes.customTooltip,
              arrow : classes.customArrow
            }}
          >
            <Avatar
              className={classes.smallAvatar}
              key={user.clientID}
              alt={user.name}
              style={{ border : `2px solid #00BBA6` }}
              src={user.avatar_url}
              imgProps={{ referrerPolicy : "no-referrer" }}
              onClick={() => { /* redirect */ }}
            />
          </Tooltip>
        ))}
      </AvatarGroup>
    </div>
  );
};

export default withStyles(styles)(UserAvatars);
