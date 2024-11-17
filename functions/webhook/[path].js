import { EmailMessage } from "cloudflare:email";
import { createMimeMessage } from "mimetext";

export function onRequest(context) {
  const url = new URL(context.request.url);
  // console.log("context.params.path", context.params.path);

  return sendEmail(context, []);
}

async function sendEmail(context, params) {
  const msg = createMimeMessage();
  msg.setSender({ name: "test", addr: "sender@apatin.ru" });
  msg.setRecipient("recepient@apatin.ru");
  msg.setSubject("An email generated in a worker");
  msg.addMessage({
      contentType: 'text/plain',
      data: `Congratulations, you just sent an email from a worker.`
  });

  var message = new EmailMessage(
    "sender@apatin.ru",
    "recepient@apatin.ru",
    msg.asRaw()
  );
  try {
    await env.SEB.send(message);
  } catch (e) {
    return new Response(e.message);
  }

  return new Response("Hello Send Email World!");
}
