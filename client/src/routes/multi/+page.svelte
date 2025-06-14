<script lang="ts">
	import Card from "$lib/components/Card.svelte";
	import SelectGame from "$lib/components/SelectGame.svelte";
	import type { GameVersionKey } from "$lib/engine/types";
	import { GameVersions } from "$lib/engine/types";
	import { MultiPlayerGameState } from "$lib/state/MultiPlayerGameState.svelte";

	import { SinglePlayerGameState } from "$lib/state/SinglePlayerGameState.svelte";
	import type { WS } from "$lib/ws/ws.svelte";
	import { onDestroy } from "svelte";

  let gameVersion = $state(null) as GameVersionKey | null;
  let gameState = $state<MultiPlayerGameState | null>(null);

  $effect(() => {
    gameState = gameVersion !== null ? new MultiPlayerGameState(GameVersions[gameVersion as GameVersionKey]) : null;

    // if (gameState === null) return
    // gameState.ws.socket.op
  });

  // $effect(() => {

  // })

  onDestroy(() => {
    if (gameState == null) return
    gameState.ws.socket?.close()
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

</script>

<svelte:window {onkeydown} />

<div class="page">
  <SelectGame bind:gameVersion={gameVersion} />
  <!-- <h1 class="text-xl font-semibold mb-4">Game Board</h1> -->
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

  @media (max-width: 600px) {
    .board {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    }
  }
</style>