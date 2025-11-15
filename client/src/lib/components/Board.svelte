<script lang="ts">
	import { GameVersions, type GameVersionKey } from "$lib/engine/types";
	import { MultiPlayerGameState } from "$lib/state/MultiPlayerGameState.svelte";
	import Card from "./Card.svelte";
	import Modal from "./lib/Modal.svelte";
	import SelectGame from "./SelectGame.svelte";
	import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
	import { OUT_MESSAGES, type CreateRoomMessage, type JoinRoomMessage } from "$lib/ws/messages";
	import FormInput from "./lib/FormInput.svelte";
	import { generateNickname } from "$lib/utils/nicknames";
	import { browser } from "$app/environment";
	import WaitingList from "./WaitingList.svelte";
	import { page } from "$app/state";
	import { Session } from "$lib/utils/sessions";

  let ws = $state<WS | null>(null);
  let gameState = $derived<MultiPlayerGameState | null>(ws?.game || null);
  let gameVersion = $state(GameVersions.classic.key) as GameVersionKey | null;
	let cardsLeft = $derived((() => {
    if (!gameState) return 0
    const total = gameState.variationsNumber ** gameState.features.length
    const discarded = Object.values(gameState.players).map(p => p.score).reduce((cur, acc) => acc + cur) * gameState.variationsNumber
    return total - discarded - gameState.inPlayCards.length
  })())
	
  let isModalOpen = $state<boolean>(false)
  let isReconnectingToOngoingGame = $state<boolean>(false)
  let joinRoomID = $state(page?.url.searchParams.get("roomID") || "")
	let roomLink = $state((() => {
		if (joinRoomID) return page.url.href + `/?roomID=${joinRoomID}`
		return ""
		})())
	let roomLinkError = $state("")
  let modalButtonText = $derived((() => {
		if (joinRoomID && !roomLink && !ws?.roomID) return "Create Room"
    if (joinRoomID && !ws?.playerID) return "Join Room"
		if (ws?.playerID && ws.isRoomOwner && !ws?.started) return "Start Game"
    if (ws?.playerID && !ws.isRoomOwner) return "Waiting for the game to start..."
		if (!joinRoomID && !ws?.playerID && roomLink) return "Join Room"
    if (!ws?.roomID) return "Create Room"
    return "Start Game"
  })())
	let nickname = $state(generateNickname())
	let nicknameError = $state("")
	let isButtonDisabled = $derived(!ws?.isRoomOwner && Boolean(joinRoomID && ws?.playerID) || !!nicknameError)

	$effect(() => {
		if (ws?.started) {
			isModalOpen = false
			isReconnectingToOngoingGame = false
		}

		if (!ws?.started && !isReconnectingToOngoingGame) {
			isModalOpen = true
		}
	})

  $effect(() => {
		let clientID: string | null = null
		if (joinRoomID) {
			const session = new Session(joinRoomID).load()
			if (session !== null) {
				clientID = session.clientID
				isReconnectingToOngoingGame = session.gameStarted === true
			}
		}
		ws = new WS("ws://localhost:8080/ws", clientID)

    return () => {
      ws?.socket?.close(1000)
			ws = null
    }
  })

  $effect(() => {
    if (ws?.roomID && page != null) {
      roomLink = page.url.href + `/?roomID=${ws.roomID}`
    }
  })

	$effect(() => {
		if (browser) {
			const session = new Session(joinRoomID).load()
			
			if (session !== null) nickname = session.nickname
		}
	})

	$effect(() => {
		if (!ws) return
		if (ws.errors.nickname) {
			nicknameError = ws.errors.nickname
		}
		if (ws.errors.roomLink) {
			roomLinkError = ws.errors.roomLink
		}
	})

	function isValidNickname() {
		return !(!nickname || nickname.length > 20)
	}

  function onClickModalButton() {
    if (!ws || ws.connectionStatus !== CONNECTION_STATUS.CONNECTED || !gameVersion) return

		if (!isValidNickname()) {
			nicknameError = "Nickname should be 1 to 20 characters long"
			return
		}

		if (ws?.isRoomOwner) {
      ws.handleStartGame(gameVersion)
		}
    
    if (joinRoomID || (roomLink && !ws.playerID)) {
      const params = new URLSearchParams(roomLink.split("/").at(-1))
      const roomID = params.get("roomID") || joinRoomID

      ws.send({
        type: OUT_MESSAGES.JOIN_ROOM,
        roomID: roomID,
				nickname
      } as JoinRoomMessage);
    } else if (!joinRoomID && !ws.roomID) {
      ws.send({
        type: OUT_MESSAGES.CREATE_ROOM,
				nickname
      } as CreateRoomMessage);
    }
  }

  // for quick testing
  function onkeydown(event: KeyboardEvent) {
    if (event.code === "Space") {
      event.preventDefault();
      const set = gameState?.findSet();
      if (set) {
        for (const card of set) {
          gameState?.toggleSelectCard(card.id);
        }
      }
    }
  }
</script>

<svelte:window {onkeydown} />
{#if (!ws?.started && !isReconnectingToOngoingGame)}
<Modal bind:open={isModalOpen}>
	<div class="modal-inner">
		{#if !joinRoomID}
			<div class="select-game">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Choose Game Version</h3>
				<SelectGame bind:gameVersion />
			</div>
		{/if}

		<FormInput
			bind:value={nickname}
			bind:error={nicknameError}
			placeholder="Nickname (1 to 20 characters)"
			transformInput={(v) => v.trim().slice(0, 20)}
		/>
    
    {#if !(ws?.roomID && !ws.isRoomOwner) || roomLinkError}
		<FormInput
			bind:value={roomLink}
			bind:error={roomLinkError}
			placeholder="Paste a room url"
			readonly={ws?.isRoomOwner || false}
			showClipboard={true}
		/>
    {/if}

		{#if (ws?.roomMembers.length)}
			<WaitingList playerID={ws.playerID} roomMembers={ws.roomMembers} />
		{/if}

		<div class="footer">
				<button onclick={onClickModalButton} class="button" disabled={isButtonDisabled}>
					{modalButtonText}
				</button>
		</div>
	</div>
</Modal>
{/if}

{#if ws?.started && gameState}
	<div class="game-info">
		<span>Cards In Play: {gameState.inPlayCards.length}</span>
    <span>Cards Left: {cardsLeft}</span>
		<span>Sets found: {gameState.score}</span>
	</div>
	<div class="board">
		{#each gameState.inPlayCards as card (card.id)}
			<Card
				{card}
				onclick={() => {
					gameState?.toggleSelectCard(card.id);
				}}
			/>
		{/each}
	</div>
{/if}

<style>
	.board {
		border: 1px solid #ccc;
		border-radius: 0.5rem;
		background-color: #fff;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

		display: grid;
		grid-template-columns: repeat(4, minmax(10rem, 1fr));
		grid-gap: 1rem;
		padding: 1rem;
		box-sizing: border-box;
	}

	.game-info {
		display: flex;
		justify-content: space-between;
		margin-bottom: 1rem;
		font-size: 1.2rem;
		color: #333;
	}

	@media (max-width: 600px) {
		.board {
			grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
		}
	}

	.modal-inner {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.select-game h3 {
    margin-bottom: 1rem;
  }

	.modal-inner .footer {
		display: flex;
		justify-content: center;
	}

	.modal-inner .button {
		padding: 0.5rem 0.75rem;
		border-radius: 6px;
		border: 1px solid transparent;
		background: #f9fafb;
		color: #222;
		font-weight: 600;
		text-align: center;
		transition:
			background 0.15s,
			border-color 0.15s;
	}

	.modal-inner .button:hover,
	.modal-inner .button:focus-visible {
		cursor: pointer;
		background: #e9ecef;
		border-color: #b3b3b3;
	}

  .modal-inner .button:disabled {
	background: #f0f0f0;
	border-color: #dcdcdc;
	color: #999;
	cursor: not-allowed;
	pointer-events: none;
}
</style>
