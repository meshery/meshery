import React from 'react'
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


const getChild=(children,index)=>{
  if(children.length!==2) throw Error("Flipcard Must Have Two Child");
  return children[index];
}

const FlipCard = ({onClick,children}) => {
  const [flipped,setFlipped]=React.useState(false);
  const Front=getChild(children,0);
  const Back=getChild(children,1);

  return (
    <Container>
      <InnerCard flipped={flipped}
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
      </InnerCard>
    </Container>
  )
}

export default FlipCard