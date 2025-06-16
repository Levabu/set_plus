<script lang="ts">
	import { GameVersions, type GameVersionKey } from "$lib/engine/types";
	import { MultiPlayerGameState } from "$lib/state/MultiPlayerGameState.svelte";
	import { onDestroy } from "svelte";
	import Card from "./Card.svelte";
	import Modal from "./Modal.svelte";
	import SelectGame from "./SelectGame.svelte";
	import { CONNECTION_STATUS } from "$lib/ws/ws.svelte";
	import { OUT_MESSAGES, type CreateRoomMessage, type JoinRoomMessage } from "$lib/ws/messages";
	import { Button } from "bits-ui";
	import DisplayRoomLink from "./DisplayRoomLink.svelte";

  // interface Props {
  //   gameVersion: GameVersionKey;
  // }

  // let { gameVersion }: Props = $props();

  let gameState = $state<MultiPlayerGameState | null>(null);
  let gameVersion = $state(GameVersions.classic.key) as GameVersionKey | null;
  // let isModalOpen = $state<boolean>(true)
  let isModalOpen = $derived<boolean>((() => {
    if (!gameState?.hasGameStarted) return true
    return false
  })())
  let joinRoomID = $state("")
  let modalButtonText = $derived((() => {
    if (joinRoomID && !gameState?.playerID) return "Join Room"
    if (joinRoomID && gameState?.playerID) return "Waiting for the game to start..."
    if (!gameState || !gameState.roomID) return "Create Room"
    return "Start Game"
  })())
  let roomLink = $derived((() => {
    const roomId = gameState?.roomID || joinRoomID
    if (!roomId) return ""
    return window.location.href + `/?roomID=${roomId}`
  })())

  $inspect({
    rommID: gameState?.roomID,
    playerID: gameState?.playerID,
    isRoomOwner: gameState?.isRoomOwner,
    started: gameState?.hasGameStarted,
    deck: gameState?.deck
  })
  $effect(() => {
    // if (!isModalOpen) return
    gameState = gameVersion !== null ? new MultiPlayerGameState(GameVersions[gameVersion as GameVersionKey]) : null;

    return () => {
      console.log("rerender")
      gameState?.ws.socket?.close()
    }
  });

  $effect(() => {
    if (!window) return
    const params = new URLSearchParams(window.location.search)
    const roomID = params.get("roomID")
    if (roomID) {
      joinRoomID = roomID
    }
  })

  function onClickModalButton() {
    if (!gameState || gameState.ws.connectionStatus !== CONNECTION_STATUS.CONNECTED) return
    
    console.log("clicked")
    if (gameState.roomID) {

    }
    if (joinRoomID) {
      gameState.ws.send({
        type: OUT_MESSAGES.JOIN_ROOM,
        roomID: joinRoomID
      } as JoinRoomMessage);
    } else if (!joinRoomID && !gameState.roomID) {
      gameState.ws.send({
        type: OUT_MESSAGES.CREATE_ROOM
      } as CreateRoomMessage);
    } else {
      gameState.handleStartGame()
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
<Modal open={isModalOpen}>
	<div class="modal-inner">
		{#if !joinRoomID}
			<div class="select-game">
				<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Choose Game Version</h3>
				<SelectGame bind:gameVersion />
			</div>
		{/if}

		<DisplayRoomLink {roomLink} />

		<div class="footer">
			<!-- {#if !(joinRoomID && gameState?.playerID)} -->
				<button onclick={onClickModalButton} class="button" disabled={Boolean(joinRoomID && gameState?.playerID)}>
					{modalButtonText}
				</button>
			<!-- {:else} -->
				<!-- <p class="waiting-message">Waiting for the game to start...</p> -->
			<!-- {/if} -->
		</div>
	</div>
</Modal>

{#if gameState?.hasGameStarted}
	<div class="game-info">
		<span>In Play Cards: {gameState.inPlayCards.length}</span>
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

	.modal-inner .waiting-message {
		margin-top: 1rem;
		font-size: 1.125rem; /* 18px */
		font-weight: 500;
		color: #666;
		background: #f8f9fa;
		padding: 0.75rem 1rem;
		border-radius: 6px;
		text-align: center;
		border: 1px solid #ddd;
		max-width: fit-content;
	}
</style>
