import jwt from "jsonwebtoken";
import schema from "./schema";
import { errorResponse, validateFormDataOrEror } from "./helpers";
import cookie from "cookie";
import { htmlResponse } from "./helpers";


const SESSION_COOKIE = "session";


/** @returns {Response} */ function setSession(
  /** @type {Response} */ response,
  /** @type {string} */ sessionToken
) {
  response.headers.append(
    "Set-Cookie",
    cookie.serialize(SESSION_COOKIE, sessionToken, {
      path: '/',
      maxAge: 60 * 60 * 24, // 1 day
    })
  );
  return response;
}

/** @returns {Resposne} */ function resetSession(
  /** @type {Response} */ response
) {
  response.headers.append(
    "Set-Cookie",
    cookie.serialize(SESSION_COOKIE, "", {
      path: '/',
      maxAge: -1,
    })
  );
  return response;
}


export default class AuthController {
  constructor(
    /** @type {string} */ secret,
    /** @type {import("bun:sqlite").Database} */ db,
    /** @type {import("nunjucks").Environment} */ templates
  ) {
    this.secret = secret;
    this.templates = templates;

    this.createUserSql = db.prepare(`
      INSERT INTO users (username, password_hash, strategy, attack, defense)
      VALUES ($username, $password_hash, $strategy, $attack, $defense);
    `);

    this.getUserSql = db.prepare(`
      SELECT * FROM users WHERE username = $username;
    `);
  }

  /** @returns {UserModel} */ getUser(
    /** @type {string} */ username
  ) {
    return this.getUserSql.get({
      $username: username
    });
  }

  createUser(
    /** @type {UserModel} */ user
  ) {
    this.createUserSql.run({
      $username: user.username,
      $password_hash: user.passwordHash,
      $strategy: user.strategy,
      $attack: user.attack,
      $defense: user.defense,
    });
  }

  /** @returns {{token: string | null, error: string | null}} */ register(
    /** @type {string} */ username,
    /** @type {string} */ password,
    /** @type {string} */ strategy,
    /** @type {number} */ attack,
    /** @type {number} */ defense
  ) {
    const passwordHash = Bun.password.hashSync(password, {
      algorithm: "argon2id",
      memoryCost: 4,
      timeCost: 3,
    })
    try {
      this.createUser({ username, passwordHash, strategy, attack, defense });
      return {
        token: jwt.sign({ username: username }, this.secret),
        error: null,
      }
    } catch (error) {
      return {
        token: null,
        error: error.toString()
      }
    }
  }

  /** @returns {string | null} */ login(
    /** @type {string} */ username,
    /** @type {string} */password
  ) {
    const result = this.getUser(username);
    if (result === null || !Bun.password.verifySync(password, result.password_hash)) {
      return null;
    }

    return jwt.sign({ username: username }, this.secret);
  }

  /** @returns {string | null} */ authenticate(
    /** @type {string} */ jwtToken
  ) {
    try {
      const { username } = jwt.verify(jwtToken, this.secret);
      return this.getUser(username);
    } catch (error) {
      console.error(error);
    }
    return null;
  }

  /** @returns {Response} */ async handleLoginGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    return htmlResponse(new Response(this.templates.render('login.html', ctx)));
  }

  /** @returns {Response} */ async handleLoginPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    const { data, errResp } = await validateFormDataOrEror(request, schema.login);
    if (errResp !== null) {
      return errResp;
    }

    const token = this.login(data.username, data.password);
    if (token === null) {
      return errorResponse("Invalid credentials", null, 403);
    }

    let response = Response.redirect('/', 302);
    return setSession(response, token);
  }

  /** @returns {Response} */ async handleRegisterGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    return htmlResponse(new Response(this.templates.render('register.html', ctx)));
  }

  /** @returns {Response} */ async handleRegisterPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    const { data, errResp } = await validateFormDataOrEror(request, schema.register)
    if (errResp !== null) {
      return errResp;
    }

    const { token, error } = this.register(data.username, data.password, data.strategy, data.attack, data.defense);
    if (error != null) {
      return errorResponse(
        "Failed to register",
        error,
        400,
      )
    }

    let response = Response.redirect('/', 302)
    if (token !== null) {
      response = setSession(response, token)
    }
    return response
  }

  /** @returns {Response} */ async handleLogoutPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    return resetSession(Response.redirect('/', 302))
  }
}