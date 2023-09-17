import { app, InvocationContext, Timer } from "@azure/functions";
import TableStorage from "../storage/table";

// This function is the entry point for the timer trigger.
export async function itemsArchiver(myTimer: Timer, context: InvocationContext): Promise<void> {
	context.log("Timer function triggered.", myTimer);
	const conn = process.env["AzureWebJobsStorage"] as string;
	const tableName = process.env["TableName"] as string;
	const table = new TableStorage(conn, tableName);
	context.log("Getting all items...");
	const items = await table.getAllItems();
	context.log(`Found ${items.length} items.`);
	for (const item of items) {
		// Check the expire timestamp
		if (item.Expire == null) {
			continue;
		}
		const expire = Date.parse(item.Expire as string);
		if (expire < Date.now()) {
			// Delete the item
			context.log(`Deleting item ${item.rowKey}...`);
			await table.deleteItem(item.rowKey as string, item.partitionKey as string);
			context.log(`Deleted item ${item.rowKey}.`);
		}
	}
	context.log("Timer function completed.");
}

app.timer("itemsArchiver", {
	schedule: "0 * * * * *",
	handler: itemsArchiver,
});
