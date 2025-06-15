import type Card from "$lib/components/Card.svelte";
import type { GameVersionKey } from "$lib/engine/types";

export const OUT_MESSAGES = {
  START_GAME: 'START_GAME',
  CREATE_ROOM: 'CREATE_ROOM',
  JOIN_ROOM: 'JOIN_ROOM',
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

export type OutMessage = 
  | StartGameMessage
  | CreateRoomMessage
  | JoinRoomMessage;





export const IN_MESSAGES = {
  CREATED_ROOM: 'CREATED_ROOM',
  JOINED_ROOM: 'JOINED_ROOM',
  STARTED_GAME: 'STARTED_GAME',
} as const;

export interface StartedGameMessage {
  readonly type: typeof IN_MESSAGES.STARTED_GAME;
  gameID: string;
  deck: Card[]
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

export type InMessage =
  | StartedGameMessage
  | CreatedRoomMessage
  | JoinedRoomMessage;
