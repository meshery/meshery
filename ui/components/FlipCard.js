import React from 'react'

const getChild=(children,index)=>{
  if(children.length!==2) throw Error("Flipcard Must Have Two Child");
  return children[index];
}

const FlipCard = ({onClick,children}) => {
  const [flipped,setFlipped]=React.useState(false);
  const Front=getChild(children,0);
  const Back=getChild(children,1);

  return (
    <div>
      <div
        onClick={()=>{
         setFlipped(p=>!p);
         onClick()
        }}
      >
        <div>
          {!flipped && React.isValidElement(Front)
                ? Front
                : null}
        </div> 
        <div>
          {flipped && React.isValidElement(Back)
              ? Back
              : null}
        </div> 
      </div>
    </div>
  )
}

export default FlipCard