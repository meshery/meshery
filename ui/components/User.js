import fetch from 'isomorphic-unfetch'
import IconButton from '@material-ui/core/IconButton';
import Avatar from '@material-ui/core/Avatar';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updateUser } from '../lib/store';

import Button from '@material-ui/core/Button';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
import MenuItem from '@material-ui/core/MenuItem';
import Popper from '@material-ui/core/Popper';
import Paper from '@material-ui/core/Paper';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';

class User extends React.Component {

  state = {
    user: null,
    open: false,
  }

  handleToggle = () => {
    this.setState(state => ({ open: !state.open }));
  };

  handleClose = event => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ open: false });
  };

  handleLogout = () => {
    window.location = "/logout";
  };

  componentDidMount() {
    console.log("fetching user data");
    fetch('/api/user', { credentials: 'same-origin' })
    .then(res => {
      if (res.ok) {
        // console.log(`res ok: ${res.ok}`);
        return res.json();
      } else {
      }
    }).then(user => {
      this.setState({user})
      this.props.updateUser({user})
    })
    .catch(error => {
      return {
        error
      };
    });
  }

  render() {
    console.log("user render called. . .")
    const {color, iconButtonClassName, avatarClassName, ...other} = this.props;
    let avatar_url, user_id;
    if (this.state.user !== null){
      avatar_url = this.state.user.avatar_url;
      user_id = this.state.user.user_id;
    }
    const { open } = this.state;
    return (
      <div>
      <IconButton color={color} className={iconButtonClassName} 
      buttonRef={node => {
        this.anchorEl = node;
      }}
      aria-owns={open ? 'menu-list-grow' : undefined}
      aria-haspopup="true"
      onClick={this.handleToggle}>
        <Avatar className={avatarClassName}  src={avatar_url} />
      </IconButton>
      {/* <Menu
        id="user-menu"
        anchorEl={userAnchorEl}
        open={Boolean(userAnchorEl)}
        onClose={this.handleClose}
      >
        <MenuItem onClick={this.handleLogout}>Logout</MenuItem>
    </Menu> */}
    <Popper open={open} anchorEl={this.anchorEl} transition disablePortal placement='top-end'>
            {({ TransitionProps, placement }) => (
              <Grow
                {...TransitionProps}
                id="menu-list-grow"
                style={{ transformOrigin: placement === 'bottom' ? 'left top' : 'left bottom' }}
              >
                <Paper>
                  <ClickAwayListener onClickAway={this.handleClose}>
                    <MenuList>
                      <MenuItem onClick={this.handleLogout}>Logout</MenuItem>
                    </MenuList>
                  </ClickAwayListener>
                </Paper>
              </Grow>
            )}
          </Popper>
    </div>
    )
  }
}

const mapDispatchToProps = dispatch => {
  return {
    updateUser: bindActionCreators(updateUser, dispatch)
  }
}

export default connect(
  null,
  mapDispatchToProps
)(User);