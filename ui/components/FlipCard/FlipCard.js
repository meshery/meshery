import React from 'react'
import useStyle from "./FlipCard.style"




const getChild=(children,index)=>{
  if(children.length!==2) throw Error("Flipcard Must Have Two Child");
  return children[index];
}

const FlipCard = ({onClick,children}) => {
  const classes=useStyle()
  const [flipped,setFlipped]=React.useState(false);
  const Front=getChild(children,0);
  const Back=getChild(children,1);

  return (
    <div className={classes.container}>
      <div className={`${classes.innerCard} ${flipped && classes.rotate}`}
        onClick={()=>{
         setFlipped(p=>!p);
         onClick()
        }}
      >
        <div className={classes.frontCard}>
          {!flipped && React.isValidElement(Front)
                ? Front
                : null}
        </div> 
        <div className={classes.backCard}>
          {flipped && React.isValidElement(Back)
              ? Back
              : null}
        </div> 
      </div>
    </div>
  )
}

export default FlipCard