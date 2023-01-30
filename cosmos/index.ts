import { CosmosClient } from "@azure/cosmos";
import dotenv from "dotenv";
dotenv.config();

import main from "./main/index.main";

// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const client = new CosmosClient({ endpoint: process.env.COSMOSENDPOINT!, key: process.env.COSMOSPRIMARYKEY! });

main(client);
