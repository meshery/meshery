import styled from "styled-components";

const ProgressBarStyle = styled.div`
  
  .timeline {
    list-style-type: none;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 50px 0;
  }
  
  .li {
    transition: all 200ms ease-in;
  }
  
  
  .status {
    padding: 0px 40px;
    display: flex;
    justify-content: center;
    border-top: 5px solid #EFEFEF;
    position: relative;
    transition: all 200ms ease-in;
  }
  .status:before {
    content: "";
    width: 25px;
    height: 25px;
    background-color: white;
    border-radius: 25px;
    border: 5px solid #EFEFEF;
    position: absolute;
    top: -20px;
    left: 42%;
    transition: all 200ms ease-in;
  }

  .li.active .status::before {
    border-color: #00B39F;
  }

  .li.complete .status,
  .li.active .status {
    border-top: 5px solid #00B39F;
  }
  .li.complete .status:before {
    background-color: #00B39F;
    border: none;
    transition: all 200ms ease-in;
    width:35px;
    height:35px;
  }


  @media (min-width: 320px) and (max-width: 700px) {
    .status {
      border-top: 3px solid #EFEFEF;
      padding: 0px 25px;
    }
    .status:before {
      content: "";
      width: 20px;
      height: 20px;
      border-radius: 20px;
      border: 3px solid #EFEFEF;
      top: -15px;
      left: -10%
    }
    .li.complete .status:before {
      width:25px;
      height:25px;
    }
  }
 
    

   
    }
`;

export default ProgressBarStyle;
