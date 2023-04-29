import { createSlice } from "@reduxjs/toolkit";
import { HYDRATE } from "next-redux-wrapper";

export const themeSlice = createSlice({
    name: "theme",
    initialState: { darkTheme: false },
    reducers: {
        toggleTheme: (state) => {
            state.darkTheme = !state.darkTheme;
        },
    },
    extraReducers: {
        [HYDRATE]: (state, action) => {
            console.log('HYDRATE theme', action.payload);

            return {
                ...state,
                ...action.payload.theme,
            };
        },
    },
});

export const { toggleTheme } = themeSlice.actions;