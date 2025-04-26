import cookie from "cookie";
import nunjucks from "nunjucks";
import path from 'node:path';

import { pathsEqual } from "./util/path";
import { initDatabase } from "./db/init";
import { errorResponse } from "./controllers/helpers";

import AuthController from "./controllers/auth";
import DuelsController from "./controllers/duels";
import UsersController from "./controllers/users";

const HOST = process.env["HOST"] ?? "localhost";
const PORT = parseInt(process.env["PORT"] ?? "3000");
const DEVELOPMENT = process.env["NODE_ENV"] === "development";
const DB_PATH = process.env["DB"] ?? ":memory:";
const JWT_SECRET = process.env["JWT_SECRET"] ?? "i-am-so-secret";

const templates = nunjucks.configure('templates', {
  autoescape: true,
  dev: DEVELOPMENT,
  noCache: DEVELOPMENT,
})

const db = initDatabase(DB_PATH);

const auth = new AuthController(JWT_SECRET, db, templates);
const users = new UsersController(/* duels */ null, db, templates);
const duels = new DuelsController(users, db, templates);
users.duels = duels;

const pathMap = [
  { path: "/", method: "GET", handler: users.handleIndexGet.bind(users) },
  { path: "/users", method: "GET", handler: users.handleUsersGet.bind(users) },
  { path: "/users", method: "POST", handler: users.handleUsersPost.bind(users) },
  { path: "/profile", method: "GET", handler: users.handleProfileGet.bind(users) },

  { path: "/auth/login", method: "GET", handler: auth.handleLoginGet.bind(auth) },
  { path: "/auth/login", method: "POST", handler: auth.handleLoginPost.bind(auth) },
  { path: "/auth/register", method: "GET", handler: auth.handleRegisterGet.bind(auth) },
  { path: "/auth/register", method: "POST", handler: auth.handleRegisterPost.bind(auth) },
  { path: "/auth/logout", method: "POST", handler: auth.handleLogoutPost.bind(auth) },

  { path: "/duels", method: "GET", handler: duels.handleDuelsGet.bind(duels) },
  { path: "/duels", method: "POST", handler: duels.handleDuelsPost.bind(duels) },
  { path: "/challenge/send", method: "POST", handler: duels.handleChallengeSendPost.bind(duels) },
  { path: "/challenge/accept", method: "POST", handler: duels.handleChallengeAcceptPost.bind(duels) },
  { path: "/challenge/deny", method: "POST", handler: duels.handleChallengeDenyPost.bind(duels) },
]

let serverOptions = {
  hostname: HOST,
  port: PORT,
  development: DEVELOPMENT,
  fetch(request) {
    const /** @type {RequestCtx} */ ctx = Object.create(null);

    const rawCookies = request.headers.get("Cookie") ?? "";
    const cookies = cookie.parse(rawCookies);
    ctx.user = (
      cookies.session === undefined
        ? null
        : auth.authenticate(cookies.session)
    );

    const url = new URL(request.url);

    for (const { path, method, handler } of pathMap) {
      if (method === request.method && pathsEqual(path, url.pathname)) {
        return handler(request, ctx);
      }
    }

    return errorResponse("Not found", `Failed to match ${request.method} to ${url.pathname} to any handler`, 404);
  },
}

if (DEVELOPMENT) {
  const liveReload = require('bun-html-live-reload');
  serverOptions = liveReload.withHtmlLiveReload(
    serverOptions,
    {
      watchPath: path.resolve(import.meta.dir, "templates")
    }
  )
}

const server = Bun.serve(serverOptions);
console.log(`Server is listening on ${HOST}:${PORT}`);

process.on("SIGINT", () => {
  console.log("Stopping server...");
  server.stop(/* closeActiveConnections */ true);
  console.log("Closing database");
  db.close();
  process.exit(0);
});
