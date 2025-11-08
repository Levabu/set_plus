export class Session {
  roomID: string;

  constructor(roomID: string) {
    this.roomID = roomID;
  }

  private getSessionKey(): string {
    return `session:${this.roomID}`;
  }

  save(clientID: string, nickname: string, gameStarted: boolean = false): void {
    const sessionData = {
      clientID,
      nickname,
      gameStarted
    };
    localStorage.setItem(this.getSessionKey(), JSON.stringify(sessionData));
  }

  load(): { clientID: string; nickname: string; gameStarted?: boolean } | null {
    const sessionData = localStorage.getItem(this.getSessionKey());
    if (sessionData) {
      return JSON.parse(sessionData);
    }
    return null;
  }
}