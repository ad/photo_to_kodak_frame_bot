import { EmailMessage } from "cloudflare:email";
import { createMimeMessage } from "../utils/gas.js";

export function onRequest(context) {
  console.log("context.request", context.request);

  return sendEmail(context, []);
}

async function sendEmail(context, params) {
  const msg = createMimeMessage();
  msg.setSender({ name: "test", addr: "photo@workerdev.ru" });
  msg.setRecipient("photo@apatin.ru");
  msg.setSubject("An email generated in a worker");
  msg.addMessage({
      contentType: 'text/plain',
      data: `Congratulations, you just sent an email from a worker.`
  });

  var message = new EmailMessage(
    "photo@workerdev.ru",
    "photo@apatin.ru",
    msg.asRaw()
  );
  try {
    await env.SEB.send(message);
  } catch (e) {
    return new Response(e.message);
  }

  return new Response("Hello Send Email World!");
}
