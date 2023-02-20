import config from "config";
import router from "./router";
import server from "../../utils/server";

const app = server(config.get<number>("microservices.replays.port"));
app.use("/replay", router);
