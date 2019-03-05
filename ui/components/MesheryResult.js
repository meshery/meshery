import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Typography from '@material-ui/core/Typography';
import { NoSsr } from '@material-ui/core';
import Moment from 'react-moment';

const styles = {
  card: {
    // maxWidth: 345,
    backgroundColor:  '#EBEFF1',  
  },
  media: {
    height: 140,
  },
};

function MesheryResult(props) {
  const { classes, data } = props; // data here maps to the MesheryResult model
//   const startTime = new Date(data.runner_results.StartTime).toString();
  return (
    // <NoSsr>
    <Card key={data.meshery_id} className={classes.card}>
      {/* <CardActionArea>
        <CardMedia
          className={classes.media}
          image="/static/images/cards/contemplative-reptile.jpg"
          title="Contemplative Reptile"
        /> */}
        <CardContent>
          <Typography gutterBottom variant="h5" component="h2">
            {data.mesh}
          </Typography>
          <Typography component="p">
            RunType: {data.runner_results.RunType}
           </Typography>
           <Typography component="p">
            Start time: <Moment format="LLLL">{data.runner_results.StartTime}</Moment>
            </Typography>
            <Typography component="p">
            Duration: {data.runner_results.RequestedDuration}
            </Typography>
            <Typography component="p">
            Actual QPS: {data.runner_results.ActualQPS.toFixed(2)}
            </Typography>
            <Typography component="p">
            Threads: {data.runner_results.NumThreads}
            </Typography>
        </CardContent>
      {/* </CardActionArea> */}
      {/* <CardActions>
        <Button size="small" color="primary">
          Share
        </Button>
        <Button size="small" color="primary">
          Learn More
        </Button>
      </CardActions> */}
    </Card>
    // </NoSsr>
  );
}

MesheryResult.propTypes = {
  classes: PropTypes.object.isRequired,
  data: PropTypes.object.isRequired,
};

export default withStyles(styles)(MesheryResult);