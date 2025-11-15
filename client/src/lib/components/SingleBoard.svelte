<script lang="ts">
	import Card from "$lib/components/Card.svelte";
	import Modal from "$lib/components/lib/Modal.svelte";
	import SelectGame from "$lib/components/SelectGame.svelte";
	import type { GameVersionKey } from "$lib/engine/types";
	import { GameVersions } from "$lib/engine/types";

	import { SinglePlayerGameState } from "$lib/state/SinglePlayerGameState.svelte";

  let gameVersion = $state(GameVersions.classic.key) as GameVersionKey | null;
  let gameState = $state<SinglePlayerGameState | null>(null);
  let clickedStart = $state(false);

  $effect(() => {
    if (!clickedStart) return;
    gameState = gameVersion !== null ? new SinglePlayerGameState(GameVersions[gameVersion as GameVersionKey]) : null;
  });

  $effect(() => {
    if (gameState?.drawPile.length !== 0 || gameState.isSetAvailable) return;
    alert("You won!")
  })

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

  function onClickStartButton() {
    clickedStart = true;
  }

</script>

<svelte:window {onkeydown} />

<div class="page">
  {#if gameState === null}
    <Modal open={true}>
      <div class="modal-inner">
        <div class="select-game">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Choose Game Version</h3>
          <SelectGame bind:gameVersion />
        </div>

        <div class="footer">
          <button onclick={onClickStartButton} class="button">
            Start Game
          </button>
        </div>
      </div>
    </Modal>
  {/if}
  {#if gameState !== null}
  <div class="board">
    {#each gameState.inPlayCards as card (card.id)}
      <Card
        card={card}
        onclick={() => {
          gameState?.toggleSelectCard(card.id);
        }}
      />
    {/each}
  </div>
  {/if}
</div>

<style>
  .page {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    padding: 2rem;
    background-color: #f0f4f8;
    min-height: 100vh;
  }
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

  @media (max-width: 600px) {
    .board {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    }
  }
</style>