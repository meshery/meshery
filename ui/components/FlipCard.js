import React from 'react'
import { styled } from "@mui/material/styles";
import {
  Box,
  Paper
} from "@mui/material";

export const CardWrapper=styled(Paper)(({flipped,theme})=>({
      padding : theme.spacing(2),
      transformStyle:"preserve-3d",
      transition:"all 0.5s ease",
      boxShadow : "0 4px 8px 0 rgba(0,0,0,0.2)",
      cursor : "pointer",
      transform:`${flipped && 'rotateY(180deg)'}`
}))

export const FrontCard=styled(Box)(()=>({
  backfaceVisibility:"hidden",
}))

export const BackCard=styled(Box)(()=>({
  backfaceVisibility:"hidden",
  transform:"rotateY(180deg)"
}))


const getChild=(children,index)=>{
  if(children.length!==2) throw Error("Flipcard Must Have Two Child");
  return children[index];
}

const FlipCard = ({onClick,children}) => {
  const [flipped,setFlipped]=React.useState(false);
  const Front=getChild(children,0);
  const Back=getChild(children,1);

  return (
    
      <CardWrapper flipped={flipped}
        onClick={()=>{
         setFlipped(p=>!p);
         onClick()
        }}
      >
        <FrontCard>
          {!flipped && React.isValidElement(Front)
                ? Front
                : null}
        </FrontCard> 
        <BackCard>
          {flipped && React.isValidElement(Back)
              ? Back
              : null}
        </BackCard> 
      </CardWrapper>
  )
}

export default FlipCard