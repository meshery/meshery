import React from 'react'
import {Container,InnerCard,FrontCard,BackCard} from "./FlipCard.style"



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