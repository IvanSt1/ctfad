interface UserModel {
  username: string;
  passwordHash: string;
  strategy: string;
  attack: number;
  defense: number;
};

interface DuelOutcome {
  attacker: string;
  defender: string;
  attacker_won: boolean;
};

interface RequestCtx {
  user: UserModel | null;
};
