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
  nickname: string;
}

export interface JoinRoomMessage {
  readonly type: typeof OUT_MESSAGES.JOIN_ROOM;
  roomID: string;
  nickname: string;
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
  LEFT_ROOM: 'LEFT_ROOM',
  RECONNECTED_TO_ROOM: 'RECONNECTED_TO_ROOM',
  SEND_STATE_TO_RECONNECTED: 'SEND_STATE_TO_RECONNECTED',
  STARTED_GAME: 'STARTED_GAME',
  CHECK_SET_RESULT: 'CHECK_SET_RESULT',
  CHANGED_GAME_STATE: 'CHANGED_GAME_STATE',
  GAME_OVER: 'GAME_OVER',
  ERROR: 'ERROR'
} as const;

export type RoomMember = { id: string; nickname: string; isConnected: boolean }
export type Player = RoomMember & { score: number }

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
  nickname: string;
}

export interface JoinedRoomMessage {
  readonly type: typeof IN_MESSAGES.JOINED_ROOM;
  roomID: string;
  playerID: string;
  nickname: string;
  players: Player[];
}

export interface LeftRoomMessage {
  readonly type: typeof IN_MESSAGES.LEFT_ROOM;
  playerID: string;
}

export interface SendStateToReconnectedMessage {
  readonly type: typeof IN_MESSAGES.SEND_STATE_TO_RECONNECTED;
  playerID: string;
  isOwner: boolean;
  roomID: string;
  gameID?: string;
  started: boolean;
  gameVersion?: GameVersionKey;
  deck?: Card[];
  players: Record<string, Player>;
}

export interface ReconnectedToRoomMessage {
  readonly type: typeof IN_MESSAGES.RECONNECTED_TO_ROOM;
  roomID: string;
  playerID: string;
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
  deck: Card[];
  players: Record<string, Player>;
}

export interface ErrorMessage {
  readonly type: typeof IN_MESSAGES.ERROR;
  refType: keyof typeof OUT_MESSAGES;
  field: string;
  reason: string;
}

export type InMessage =
  (
    | StartedGameMessage
    | CreatedRoomMessage
    | JoinedRoomMessage
    | LeftRoomMessage
    | SendStateToReconnectedMessage
    | ReconnectedToRoomMessage
    | CheckSetResultMessage
    | ChangedGameStateMessage
    | GameOverMessage
    | ErrorMessage
  ) & {
    isProcessed?: boolean;
    error?: string | null;
  };
