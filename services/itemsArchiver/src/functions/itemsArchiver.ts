import { app, InvocationContext, Timer } from "@azure/functions";
import { TableStorage } from "../storage/table";

export async function itemsArchiver(myTimer: Timer, context: InvocationContext): Promise<void> {
	const connectionString =
		"DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;";
	const table = new TableStorage(connectionString, "items");
	context.log("Getting all items...");
	const items = await table.getAllItems();
	context.log(`Found ${items.length} items.`);
	for (const item of items) {
		context.log(`Archiving item ${item.rowKey}...`);
	}
	context.log("Timer function processed request.");
}

app.timer("itemsArchiver", {
	schedule: "0 * * * * *",
	handler: itemsArchiver,
});
