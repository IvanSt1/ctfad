import { errorResponse, htmlResponse, jsonResponse, requireAuth, validateFormDataOrEror, validateQueryParamsOrError } from "./helpers";
import schema from "./schema";
import DuelsController from "./duels";

export default class UsersController {
  constructor(
    /** @type {DuelsController} */ duels,
    /** @type {import("bun:sqlite").Database} */ db,
    /** @type {import("nunjucks").Environment} */ templates
  ) {
    this.duels = duels;
    this.templates = templates;

    this.getLastUsersSql = db.prepare(`
      SELECT username FROM users
      ORDER BY rowid DESC
      LIMIT $limit OFFSET $offset;
    `);

    this.getUserSql = db.prepare(`
      SELECT username, attack, defense FROM users
      WHERE username == $username;
    `);

    this.getUserAsAnotherSql = db.prepare(`
      SELECT username, attack, defense, (
        CASE 
          WHEN
            $requestUser == $targetUser
          OR
            (SELECT COUNT(*) FROM duels WHERE (defender == $requestUser AND attacker == $targetUser AND attacker_won == 0)
                                           OR (attacker == $requestUser AND defender == $targetUser AND attacker_won == 1)
            ) > 0
        THEN strategy
        ELSE NULL
        END
      ) AS strategy
      FROM users
      WHERE username == $targetUser;
    `);
  }

  /** @returns {{username: string}p[]} */ getLastUsers(
    /** @type {{offset: number, limit: number}} */ queryParams = {}
  ) {
    const limit = queryParams.limit ?? 500;
    const offset = queryParams.offset ?? 0;
    return this.getLastUsersSql.all({
      $limit: limit,
      $offset: offset,
    });
  }

  /** @returns {{username: string, attack: number, defense: number}} */ getUser(
    /** @type {string} */ username,
  ) {
    return this.getUserSql.get({
      $username: username
    });
  }

  /** @returns {{username: string, strategy: string | null, attack: number, defense: number}} */ getUserAsAnother(
    /** @type {string} */ requestUser,
    /** @type {string} */ targetUser,
  ) {
    return this.getUserAsAnotherSql.get({
      $requestUser: requestUser,
      $targetUser: targetUser,
    });
  }

  /** @returns {Response} */ async handleIndexGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    const generateResponse = () => {
      return htmlResponse(new Response(this.templates.render("index.html", ctx)));
    }

    if (ctx.user === null) {
      return generateResponse();
    }

    ctx.challenges = {};

    try {
      ctx.challenges.incoming = this.duels.getIncomingChallenges(ctx.user.username);
    } catch (error) {
      return errorResponse(
        "Failed to list incoming challenges",
        error.toString(),
        500,
      );
    }

    try {
      ctx.challenges.outgoing = this.duels.getOutgoingChallenges(ctx.user.username);
    } catch (error) {
      return errorResponse(
        "Failed to list ougoing challenges",
        error.toString(),
        500,
      );
    }

    try {
      ctx.duels = this.duels.getUserDuels(ctx.user.username);
    } catch (error) {
      return errorResponse(
        "Failed to list user duels",
        error.toString(),
        500,
      );
    }

    return generateResponse();
  }

  /** @returns {Response} */ async handleUsersGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    ctx.users = this.getLastUsers()
    return htmlResponse(new Response(this.templates.render("users.html", ctx)));
  }

  /** @returns {Response} */ async handleUsersPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx,
  ) {
    const { data, errResp } = await validateFormDataOrEror(request, schema.getUsers);
    if (errResp !== null) {
      return errResp;
    }
    return jsonResponse(new Response(JSON.stringify(this.getLastUsers(data))));
  }

  /** @returns {Response} */ async handleProfileGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    {
      const errResp = requireAuth(ctx);
      if (errResp !== null) {
        return errResp;
      }
    }

    const { data, errResp } = await validateQueryParamsOrError(request, schema.getUser);
    if (errResp !== null) {
      return errResp;
    }

    const profile = this.getUserAsAnother(ctx.user.username, data.username);
    if (profile === null) {
      return errorResponse(
        `Failed to find user ${data.username}`,
        null,
        404
      );
    }
    ctx.profile = profile;
    ctx.challenge = {
      outgoing: this.duels.challengeId(ctx.user.username, ctx.profile.username),
      incoming: this.duels.challengeId(ctx.profile.username, ctx.user.username),
    };

    return htmlResponse(new Response(this.templates.render("profile.html", ctx)));
  }
}