import "./ProgressBarStyle.css";

const ProgressDisplay = ({stepNumber,content,totalStep}) => {
  let progress = [];
  for (let index = 0; index < totalStep; index++) {
    if(stepNumber===index){
      progress.push(
        <li className="li active">
          <div className="status">
            {content}
          </div>
        </li>
      )
    }
    else if(stepNumber>index){
      progress.push(
        <li className="li complete">
          <div className="status">
            {content}
          </div>
        </li>
      )
    }
    else{
      progress.push(
        <li className="li">
          <div className="status">
            {content}
          </div>
        </li>
      )
    } 
  }
    return (
      <ul className="timeline" id="timeline">
        {progress}
      </ul>
    );
  };
const ProgressBar = ({stepNumber}) =>{
    return (

    <div>
      <ProgressDisplay totalStep={5} stepNumber={stepNumber} content={""} />
    </div>

    );
}



// Example of main.js where we are calling the progress bar coponent

/* 
function App() {
  const [stepNumber, setStepNumber] = useState(0);

  const nextStep = () => {
    window.scrollTo(0, 0);
  };
  const laststep = () => {
    setStepNumber(stepNumber - 1);
  };

  const ToShowFirst = ({text})=>{
    const handleSubmit=()=>{
      
      console.log("button clicked")
      setStepNumber(1)
      nextStep()
      
    }
    return(
      <div>
      <form >
      <p>This is the {text} part of the form</p>
      <button type="submit" className="btn" onClick={handleSubmit}> nextStep </button>
      </form>
      
      </div>
    ) 
  }
  const ToShowSecond = ({text})=>{
    const handleSubmit=()=>{
      
      console.log("button clicked")
      setStepNumber(2)
      nextStep()
      
    }
    return(
      <div>
      <form >
      <p>This is the {text} part of the form</p>
      <button type="submit" className="btn" onClick={handleSubmit}> nextStep </button>
      </form>
      <button type="submit" onClick={laststep}>Previous Step</button>
      </div>
    ) 
  }
  const ToShowThird = ({text})=>{
    const handleSubmit=()=>{
      
      console.log("button clicked")
      setStepNumber(3)
      nextStep()
      
    }
    return(
      <div>
      <form >
      <p>This is the {text} part of the form</p>
      <button type="submit" className="btn" onClick={handleSubmit}> nextStep </button>
      </form>
      <button type="submit" onClick={laststep}>Previous Step</button>
      </div>
    ) 
  }
  const ToShowFourth = ({text})=>{
    const handleSubmit=()=>{
      
      console.log("button clicked")
      setStepNumber(4)
      nextStep()
      
    }
    return(
      <div>
      <form >
      <p>This is the {text} part of the form</p>
      <button type="submit" className="btn" onClick={handleSubmit}> nextStep </button>
      </form>
      <button type="submit" onClick={laststep}>Previous Step</button>
      </div>
    ) 
  }
  const ToShowFifth = ({text})=>{
    
    return(
      <div>
      <form >
      <p>This is the {text} part of the form</p>
      <br></br>
      <button type="submit" onClick={laststep}>Previous Step
      </button>
      </form>
      </div>
    ) 
  }
  return(
    <div>
  <ProgressBar 
  stepNumber={stepNumber} laststep={laststep} 
  setStepNumber={setStepNumber}/>
  {
        stepNumber === 0 &&
        <ToShowFirst text={"first"}/>
      }
      {
        stepNumber === 1 &&
        <ToShowSecond text={"second"}/>
      }
      {
        stepNumber === 2 &&
        <ToShowThird text={"third"}/>
      }
      {
        stepNumber === 3 &&
        <ToShowFourth text={"fourth"}/>
      }
      {
        stepNumber === 4 &&
        <ToShowFifth text={"fifth"}/>
      }
  </div>
  )
}

export default App;
*/