import Home from "../../pages/index";

export default {
    title: "Pages/Home",
    component: Home,
    parameters: {
        // More on how to position stories at: https://storybook.js.org/docs/react/configure/story-layout
        layout: 'fullscreen',
    },
};

export const HomePage = () => <Home />