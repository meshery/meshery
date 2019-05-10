import { CircularProgress } from "@material-ui/core";

const MesheryProgressBar = (props) => {
    const {showProgress, classes} = props;
    return (
        showProgress && 
        <CircularProgress className={classes.progress} />
    );
};

export default MesheryProgressBar;