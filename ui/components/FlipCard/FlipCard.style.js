import { makeStyles } from "@material-ui/core";
const useStyle=makeStyles((theme)=>({
    container:{
      position:"relative",
      height:"100%",
      borderRadius : theme.spacing(1),
    },
    innerCard:{
      padding : theme.spacing(2),
      borderRadius : theme.spacing(1),
      position:"relative",
      height:"100%",
      transformStyle:"preserve-3d",
      transition:"all 0.5s ease",
      boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
      backgroundColor : "#fff",
      cursor : "pointer",
    },
    rotate:{
      transform:"rotateY(180deg)",
    },
    frontCard:{
      height:"100%",
      backfaceVisibility:"hidden",
    },
    backCard:{
      height:"100%",
      backfaceVisibility:"hidden",
      transform:"rotateY(180deg)"
    }

}))
export default useStyle;