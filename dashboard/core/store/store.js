import { createWrapper } from "next-redux-wrapper";

const { combineReducers, configureStore } = require("@reduxjs/toolkit");
const { themeSlice } = require("./theme/themeSlice");

const reducers = {
    [themeSlice.name]: themeSlice.reducer
};

const reducer = combineReducers(reducers);

const makeStore = () =>
    configureStore({
        reducer,
        devTools: true,
    });

export const wrapper = createWrapper(makeStore);