import route from "./route";
import server from "../../utils/server";

const app = server();
app.use("/replay", route);
