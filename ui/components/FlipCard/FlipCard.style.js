import { styled } from "@mui/material/styles";

export const Container=styled('div')(({theme})=>({
  position:"relative",
  height:"100%",
  borderRadius : theme.spacing(1),
}))

export const InnerCard=styled('div')(({flipped,theme})=>({
      padding : theme.spacing(2),
      borderRadius : theme.spacing(1),
      position:"relative",
      height:"100%",
      transformStyle:"preserve-3d",
      transition:"all 0.5s ease",
      boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
      backgroundColor : "#fff",
      cursor : "pointer",
      transform:`${flipped && 'rotateY(180deg)'}`
}))

export const FrontCard=styled('div')(()=>({
  height:"100%",
  backfaceVisibility:"hidden",
}))

export const BackCard=styled('div')(()=>({
  height:"100%",
  backfaceVisibility:"hidden",
  transform:"rotateY(180deg)"
}))