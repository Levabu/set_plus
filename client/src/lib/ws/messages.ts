import type Card from "$lib/components/Card.svelte";
import type { GameVersionKey } from "$lib/engine/types";

export const CLIENT_TO_SERVER_MESSAGES = {
  CREATE_GAME: 'CREATE_GAME',
} as const;

export interface CreateGameMessage {
  readonly type: typeof CLIENT_TO_SERVER_MESSAGES.CREATE_GAME;
  gameVersion: GameVersionKey;
}

export type ClientToServerMessage = 
  | CreateGameMessage;

export const SERVER_TO_CLIENT_MESSAGES = {
  GAME_CREATED: 'GAME_CREATED',
} as const;

export interface GameCreatedMessage {
  readonly type: typeof SERVER_TO_CLIENT_MESSAGES.GAME_CREATED;
  gameID: string;
  deck: Card[]
}
export type ServerToClientMessage =
  | GameCreatedMessage;
