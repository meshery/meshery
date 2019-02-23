import fetch from 'isomorphic-unfetch'
import IconButton from '@material-ui/core/IconButton';
import Avatar from '@material-ui/core/Avatar';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updateUser } from '../lib/store';

import Button from '@material-ui/core/Button';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';

class User extends React.Component {

  state = {
    user: null,
    userAnchorEl: null,
  }

  handleClick = event => {
    this.setState({ userAnchorEl: event.currentTarget });
  };

  handleClose = () => {
    this.setState({ userAnchorEl: null });
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
    const { userAnchorEl } = this.state;
    return (
      <div>
      <Button color={color} className={iconButtonClassName} 
      aria-owns={userAnchorEl ? 'user-menu' : undefined}
      aria-haspopup="true"
      onClick={this.handleClick}>
        <Avatar className={avatarClassName}  src={avatar_url} />
      </Button>
      {/* <Button
      aria-owns={anchorEl ? 'simple-menu' : undefined}
      aria-haspopup="true"
      onClick={this.handleClick}
    >
      Open Menu
    </Button> */}
    <Menu
        id="user-menu"
        anchorEl={userAnchorEl}
        open={Boolean(userAnchorEl)}
        onClose={this.handleClose}
      >
        <MenuItem onClick={this.handleClose}>Logout</MenuItem>
    </Menu>
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