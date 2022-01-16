import makeStyles from '@mui/styles/makeStyles'

export const getStyles = () =>({
    timeline:{
        listStyleType: "none",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        
    },
    li: {
        transition: "all 150ms ease-in",
    },
    status: {
        '&::before':{
            content: '""',
            width: "25px",
            height: "25px",
            backgroundColor: "white",
            borderRadius: "25px",
            border: "5px solid #EFEFEF",
            position: "absolute",
            top: "-20px",
            left: "42%",
            transition: "all 200ms ease-in"
        },
        padding: "0px 40px",
        display: "flex",
        justifyContent:" center",
        borderTop: "5px solid #EFEFEF",
        position: "relative",
        transition: "all 100ms ease-in",
    },
    complete:{
        '& div':{
            '&::before':{
                backgroundColor: "#00B39F",
                border: "none",
                transition: "all 200ms ease-in",
                width:"35px",
                height:"35px",
            },
            borderTop:"5px solid #00B39F",
        },
    },
    active:{
        '& div':{
            '&::before':{
                borderColor: "#00B39F"
            },
            borderTop:"5px solid #00B39F",

        }

    }
      
      
      
});

export const useStyles = makeStyles(getStyles);