import React from "react"
import Form from "@rjsf/material-ui";

//Rjsf  component
function Rjsf() {
  const schema = {
    "type": "object",
    "properties": {
      "rules": {
        "type": "array",
        "title": "Rules",
        "items": {
          "type": "object",
          "properties": {
            "route": {
              "type": "string",
              "title": "Route"
            },
            "ruleType":{
              "type":"string",
    
              "title": "Rule Type"
            },
            "parameters": {
              "type": "array",
              "title": "Parameters",
              "items":{
                "type":"object",
                "properties":{
                  "identifier":{
                    "type":"string",
                    "title":"Identifier"
                  },
                  "limit":{
                    "type":"string",
                    "title":"Limit of requests"
                  }
                }
              }
            }      
          }
        }
      }
    }
  }
  
  return (
    <div className="App">
      <Form schema = {schema} />      
    </div>
  );
}
  
export default Rjsf;
  