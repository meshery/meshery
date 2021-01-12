import React, {useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Collapse, Grid, Paper } from '@material-ui/core';

const useStyles = makeStyles(({ spacing }) => ({
  root: {
    margin: 'auto',
    borderRadius: spacing(2),
    boxShadow: '0 2px 7px 0 rgba(0,0,0,0.12)',
    position: 'relative',
    width: 370,
    overflow: 'initial',
    background: '#ffffff',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'column',
    padding: spacing(2),
    transition: '0.3s',
    '&:hover': {
      boxShadow: '0px 3.5px 15px rgba(34, 35, 58, 0.2)',
    },
  },
  media: props => ({
    width: '5em',
    height: "5em",
    backgroundSize: 'contain',
    borderRadius: spacing(2),
    backgroundColor: '#fff',
    position: 'relative',
    '&:after': {
      content: '" "',
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      backgroundImage: props.gradient,
      borderRadius: spacing(2),
      opacity: 0.5,
    },
  }),
  image: {
    width: '5em',
    height: "5em",
    backgroundSize: 'contain',
    borderRadius: spacing(2),
    position: 'relative',
    padding: '5px',
  },
  h4: {
    color: "rgba(105, 105, 105)",
  },
}));

export const MeshCardComponent = React.memo(function BlogCard(props) {
  const [isExpanded, setExpanded] = useState();
  const styles = useStyles(props.meshImg)

  return (
    <Grid item >
      <Paper className={styles.root}>
        <Grid container spacing={1} onClick={() => setExpanded(!isExpanded)}>
          <Grid item >
            <div className={styles.media}>
              {props.meshData!=null?<img className={styles.image} src={props.meshImg.img}></img>: <img className={styles.image} style={{filter: "grayscale(1)"}} src={props.meshImg.img}></img> }
              
            </div>
          </Grid >
          {/* <Grid item xs={3} spacing={1}>
            <FontAwesomeIcon icon={faBell} transform="grow-4" color="#23C552" fixedWidth />{"5"}
            <FontAwesomeIcon icon={faBell} transform="grow-4" color="#23C552" fixedWidth />{"5"}
          </Grid> */}
        </Grid>
        <Grid container spacing={2} onClick={() => setExpanded(!isExpanded)}>
          <Grid item xs={12}>
            {props.meshData!=null?<h3 >{props.meshData[0].namespace} </h3>: <h3 className={styles.h4}>{props.meshImg.name + " not found"} </h3> }
          </Grid >
        </Grid>
        <Collapse in={isExpanded} timeout="auto">
          {props.meshData!=null && props.meshData.map( data => 
            <Grid container spacing={1}>
              <Grid item xs={12}>
                <h3 className={styles.h4}>{data && data.component}</h3>
              </Grid >
              <Grid item sm={12} md={6}>
              Version: {data.version}
              </Grid>
              <Grid item sm={12} md={6}>
              Control Pod: {data.name.trimmed}
              </Grid>
            </Grid>
          )}
        </Collapse>
      </Paper>
    </Grid>
  );
});

export default MeshCardComponent
