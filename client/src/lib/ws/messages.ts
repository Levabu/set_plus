import type Card from "$lib/components/Card.svelte";
import type { GameVersionKey } from "$lib/engine/types";

export const OUT_MESSAGES = {
  START_GAME: 'START_GAME',
  CREATE_ROOM: 'CREATE_ROOM',
  JOIN_ROOM: 'JOIN_ROOM',
  CHECK_SET: 'CHECK_SET',
} as const;

export interface StartGameMessage {
  readonly type: typeof OUT_MESSAGES.START_GAME;
  roomID: string;
  gameVersion: GameVersionKey;
}

export interface CreateRoomMessage {
  readonly type: typeof OUT_MESSAGES.CREATE_ROOM;
}

export interface JoinRoomMessage {
  readonly type: typeof OUT_MESSAGES.JOIN_ROOM;
  roomID: string;
}

export interface CheckSetMessage {
  readonly type: 'CHECK_SET';
  cardIDs: string[];
  playerID: string;
  roomID: string;
  gameID: string;
}

export type OutMessage =
  | StartGameMessage
  | CreateRoomMessage
  | JoinRoomMessage
  | CheckSetMessage;





export const IN_MESSAGES = {
  CREATED_ROOM: 'CREATED_ROOM',
  JOINED_ROOM: 'JOINED_ROOM',
  STARTED_GAME: 'STARTED_GAME',
  CHECK_SET_RESULT: 'CHECK_SET_RESULT',
  CHANGED_GAME_STATE: 'CHANGED_GAME_STATE',
  GAME_OVER: 'GAME_OVER',
} as const;

export type Player = { id: string; nickname: string; score: number }

export interface StartedGameMessage {
  readonly type: typeof IN_MESSAGES.STARTED_GAME;
  gameID: string;
  gameVersion: GameVersionKey;
  deck: Card[]
  players: Record<string, Player>;
}

export interface CreatedRoomMessage {
  readonly type: typeof IN_MESSAGES.CREATED_ROOM;
  roomID: string;
  playerID: string;
}

export interface JoinedRoomMessage {
  readonly type: typeof IN_MESSAGES.JOINED_ROOM;
  roomID: string;
  playerID: string;
  error: string | null;
}

export interface CheckSetResultMessage {
  readonly type: typeof IN_MESSAGES.CHECK_SET_RESULT;
  isSet: boolean;
}

export interface ChangedGameStateMessage {
  readonly type: typeof IN_MESSAGES.CHANGED_GAME_STATE;
  gameID: string;
  deck: Card[];
  playerID: string;
  players: Record<string, Player>;
}

export interface GameOverMessage {
  readonly type: typeof IN_MESSAGES.GAME_OVER;
  gameID: string;
  players: Record<string, Record<string, string | number>>;
}

export type InMessage =
  (
    | StartedGameMessage
    | CreatedRoomMessage
    | JoinedRoomMessage
    | CheckSetResultMessage
    | ChangedGameStateMessage
    | GameOverMessage
  ) & {
    isProcessed?: boolean;
    error?: string | null;
  };
