import { EmailMessage } from "cloudflare:email";
// import { createMimeMessage } from "../utils/gas.js";

export function onRequest(context) {
  console.log("context.request", context.request);

  try {
    sendEmail(context.env, "photo@workerdev.ru", "photo@apatin.ru", "test", "body");
    return new Response('Email sent successfully');
  } catch (e) {
    return new Response(`Internal Server Error: ${e.message}`, { status: 500 });
  }
}

async function sendEmail(env, from, to, subject, html) {
  const rawMessage = buildRawEmailMessage(from, to, subject, html);

  const message = new EmailMessage(
    from,
    to,
    rawMessage
  );

  await env.EMAIL.send(message);
}

function buildRawEmailMessage(from, to, subject, html) {
  const headers = [
    'MIME-Version: 1.0',
    `From: ${from}`,
    `To: ${to}`,
    `Subject: ${subject}`,
    'Content-Type: text/html; charset=UTF-8',
    '',
    html
  ];

  return headers.join('\r\n');
}