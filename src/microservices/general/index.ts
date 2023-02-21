import config from "config";
import router from "./router";
import server from "../../utils/server";

const app = server(config.get<number>("microservices.general.port"));
app.use("/", router);
