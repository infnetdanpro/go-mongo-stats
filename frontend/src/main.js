import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

import "./assets/normalize.css";
import "./assets/skeleton.css";
{/* <link href="//fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css"></link> */}
const app = createApp(App);

app.use(createPinia());
app.use(router);

app.mount("#app");
