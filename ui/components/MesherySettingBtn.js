import { Button, Typography } from "@material-ui/core"
import React, { useState } from "react"
import AddIconCircleBorder from "../assets/icons/AddIconCircleBorder"
import Modal from "./Modal"

export const MesherySettingBtn = ({
    type
}) => {
    const [modalState, setModalState] = useState({
        open: false
    })

    const handleClick = (ev) => {
        ev.stopPropagation();
        setModalState({
            open: true
        })
    }

    const handleModalClose = () => {
        setModalState({
            open: false
        })
    }

    return (
        <>
            <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={handleClick}
            style={{
            "padding" : "0.5rem",
            "borderRadius" : 5,
            "marginRight" : "2rem"
            }}
            data-cy="btnResetDatabase"
        >
            <AddIconCircleBorder style={{ width : "1.25rem" }} />
            <Typography
            style={{
                "paddingLeft" : "0.25rem" ,
                "marginRight" : "0.25rem"
            }}
            >
                Add {type}
            </Typography>
        
        </Button>
        <Modal open={modalState.open} handleClose={handleModalClose}  />
      </>
    )
}