import ProgressBarStyle from "./ProgressBarStyle";

const ProgressDisplay = ({stepNumber}) => {
    return (
      <ul className="timeline" id="timeline">

        <li className={stepNumber === 0 ? ("li active") : stepNumber > 0 ? ("li complete") : ("li")}>
          <div className="status">
          </div>
        </li>
        <li className={stepNumber === 1 ? ("li active") : stepNumber > 1 ? ("li complete") : ("li")}>
          <div className="status">
          </div>
        </li>
        <li className={stepNumber === 2 ? ("li active") : stepNumber > 2 ? ("li complete") : ("li")}>
          <div className="status">
          </div>
        </li>
        <li className={stepNumber === 3 ? ("li active") : stepNumber > 3 ? ("li complete") : ("li")}>
          <div className="status">
          </div>
        </li>
        <li className={stepNumber === 4 ? 
            ("li active") : stepNumber > 4 ? ("li complete") : ("li")}>
            <div className="status">
            </div>
        </li>
        
      </ul>
    );
  };
const ProgressBar = ({stepNumber}) =>{
    return (
      <ProgressBarStyle>
    <div>
      <ProgressDisplay stepNumber={stepNumber} />
    </div>
    </ProgressBarStyle>
    );
}


export default ProgressBar


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