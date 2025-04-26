import vine from '@vinejs/vine';

const username = vine.string().minLength(3).maxLength(256).alphaNumeric({
  allowDashes: true,
  allowUnderscores: true,
  allowSpaces: false,
});

const password = vine.string().minLength(8).maxLength(256);

const register = vine.compile(
  vine.object({
    username: username,
    password: password,
    strategy: vine.string().maxLength(1024).optional(),
    attack: vine.number().min(0).optional(),
    defense: vine.number().min(0).optional(),
  }).use(
    vine.createRule(async (value, _, field) => {
      const maxStrength = 200;
      const attackStrength = 2 * value.attack;
      const defenseStrength = value.defense;
      const unusedStrength = maxStrength - attackStrength - defenseStrength;
      if (unusedStrength < 0) {
        field.report(
          `value of "2*<attack> + <defense>" cannot be more than ${maxStrength}; got ${maxStrength - unusedStrength}`,
          'totalStrength',
          field
        )
      }
    })()
  )
);

const login = vine.compile(
  vine.object({
    username: username,
    password: password,
  })
);

const selectLimitOffset = vine.compile(
  vine.object({
    limit: vine.number().min(0).max(1000).optional(),
    offset: vine.number().min(0).optional(),
  }).allowUnknownProperties()
);

const getUser = vine.compile(
  vine.object({
    username: username,
  })
);

const sendChallenge = vine.compile(
  vine.object({
    username: username,
  })
);

const challengeId = vine.number()

const acceptChallenge = vine.compile(
  vine.object({
    challengeId: challengeId,
  })
);

const denyChallenge = vine.compile(
  vine.object({
    challengeId: challengeId,
  })
);

export default {
  register,
  login,
  getUsers: selectLimitOffset,
  getUser,
  sendChallenge,
  acceptChallenge,
  denyChallenge,
  getDuels: selectLimitOffset,
};
