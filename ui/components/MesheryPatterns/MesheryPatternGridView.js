//@ts-check
import { Grid } from "@material-ui/core";
import { Pagination } from "@material-ui/lab";
import React, { useState } from "react";
import MesheryPatternCard from "./MesheryPatternCard";
import { makeStyles } from "@material-ui/core/styles";
import PatternConfiguratorComponent from "../configuratorComponents/patternConfigurator"
import FILE_OPS from "../../utils/configurationFileHandlersEnum";


const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function PatternCardGridItem({ pattern, handleDeploy, handleUnDeploy, handleSubmit, setSelectedPatterns }) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(pattern.pattern_file);

  return (
    <Grid item {...gridProps}>
      <MesheryPatternCard
        // id={pattern.id}
        name={pattern.name}
        updated_at={pattern.updated_at}
        created_at={pattern.created_at}
        pattern_file={pattern.pattern_file}
        requestFullSize={() => setGridProps({ xl : 12, md : 12, xs : 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleDeploy={() => handleDeploy(pattern.pattern_file)}
        handleUnDeploy={() => handleUnDeploy(pattern.pattern_file)}
        deleteHandler={() => handleSubmit({ data : yaml, id : pattern.id, type : FILE_OPS.DELETE ,name : pattern.name })}
        updateHandler={() => handleSubmit({ data : yaml, id : pattern.id, type : FILE_OPS.UPDATE ,name : pattern.name })}
        setSelectedPatterns={() => setSelectedPatterns({ pattern : pattern, show : true })}
        setYaml={setYaml}
      />
    </Grid>
  );
}
const useStyles = makeStyles(() => ({
  pagination : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    marginTop : "2rem"
  }
}))

/**
 * MesheryPatternGrid is the react component for rendering grid
 * @param {{
 *  patterns:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  pattern_file: string,
 * }>,
 *  handleDeploy: (pattern_file: any) => void,
 *  handleUnDeploy: (pattern_file: any) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedPattern : ({show: boolean, pattern:any}) => void,
 *  selectedPattern: {show : boolean, pattern : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void
 * }} props props
 */

function MesheryPatternGrid({ patterns=[],handleDeploy, handleUnDeploy, handleSubmit, setSelectedPattern, selectedPattern, pages = 1,setPage, selectedPage }) {

  const classes = useStyles()
  return (
    <div>
      {selectedPattern.show &&
      <PatternConfiguratorComponent pattern={selectedPattern.pattern} show={setSelectedPattern}  onSubmit={handleSubmit} />
      }
      {!selectedPattern.show &&
      <Grid container spacing={3} style={{ padding : "1rem" }}>
        {patterns.map((pattern) => (
          <PatternCardGridItem
            key={pattern.id}
            pattern={pattern}
            handleDeploy={handleDeploy}
            handleUnDeploy={handleUnDeploy}
            handleSubmit={handleSubmit}
            setSelectedPatterns={setSelectedPattern}
          />
        ))}

      </Grid>
      }
      {patterns.length
        ? (
          <div className={classes.pagination} >
            <Pagination count={pages} page={selectedPage+1} onChange={(_, page) => setPage(page - 1)} />
          </div>
        )
        : null}
    </div>
  );
}

export default MesheryPatternGrid;