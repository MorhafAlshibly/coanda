import { app, InvocationContext, Timer } from "@azure/functions";

export async function itemsArchiver(myTimer: Timer, context: InvocationContext): Promise<void> {
    context.log('Timer function processed request.');
}

app.timer('itemsArchiver', {
    schedule: '0 0 0 * * *',
    handler: itemsArchiver
});
