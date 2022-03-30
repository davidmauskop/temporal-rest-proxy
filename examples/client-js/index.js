import fetch from "node-fetch";
import { v4 as uuid } from "uuid";

// The base URL and auth token should be set as environment variables.
// The base URL might look something like "https://temporal-rest-proxy-abc1.onrender.com"
const baseUrl = process.env.BASE_URL;
const authToken = process.env.AUTH_TOKEN;
const encodedEncoding = Buffer.from("json/plain").toString("base64");
const namespace = "default";
const workflowName = "TransferMoney";
const workflowId = "transfer-money-workflow";
const taskQueueName = "TRANSFER_MONEY_TASK_QUEUE";

export async function triggerTransferMoneyWorkflow(amount, from, to) {
  // inputs to the TransferMoney workflow, which has already been registered with Temporal
  const workflowArgs = {
    Amount: amount,
    FromAccount: from,
    ToAccount: to,
    ReferenceID: uuid(),
  };
  const encodedInputData = Buffer.from(JSON.stringify(workflowArgs)).toString(
    "base64"
  );

  try {
    const response = await fetch(
      `${baseUrl}/api/v1/namespaces/${namespace}/workflows/${workflowId}/start`,
      {
        method: "POST",
        body: JSON.stringify({
          workflowType: { name: workflowName },
          taskQueue: { name: taskQueueName },
          requestId: uuid(),
          input: {
            payloads: [
              {
                data: encodedInputData,
                metadata: { encoding: encodedEncoding },
              },
            ],
          },
        }),
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      }
    );
    console.log(await response.json());
  } catch (e) {
    console.log(e);
  }
}

triggerTransferMoneyWorkflow(9.99, "Account 1", "Account 2");
