import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr, Typography } from "@material-ui/core";
import FavoriteIcon from '@material-ui/icons/Favorite';
import MesheryLogo from "./MesheryLogo";
import Configuration from "./Config/Configuration";

const styles = (theme) => ({
  root: {
    textAlign: "center",
  },
  container: {
    width: "60%",
    marginLeft: "auto",
    marginRight: "auto",
    marginTop: theme.spacing(2),
  },
  footerText: {
    cursor: 'pointer',
    display: 'inline',
    verticalAlign: 'middle',
  },
  footerIcon: {
    display: 'inline',
    verticalAlign: 'top',
  },
});

class ExtensionComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
    };
  }

  render() {
    const { classes } = this.props;

    return (
      <NoSsr>
        <div className={classes.root}>
          <MesheryLogo />
          <Configuration />
          <Typography variant="body2" align="center" color="textSecondary" component="p">
                <span onClick={this.handleL5CommunityClick} className={classes.footerText}>
                  Built with <FavoriteIcon className={classes.footerIcon} /> by the Layer5 Community
                </span>
          </Typography>
        </div>
      </NoSsr>
    );
  }
}

ExtensionComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(ExtensionComponent);
