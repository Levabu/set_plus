<script lang="ts">
	import Card from "$lib/components/Card.svelte";
	import type { Feature } from "$lib/engine/types";
	import { FEATURES } from "$lib/engine/types";

	import { SinglePlayerGameState } from "$lib/state/SinglePlayerGameState.svelte";
	import { onDestroy, onMount } from "svelte";

  const features: Feature[] = [
    FEATURES.shape,
    FEATURES.color,
    FEATURES.shading,
    FEATURES.number,
  ];

  const gameState = new SinglePlayerGameState(features, 3);
  
  $effect(() => {
    if (gameState.drawPile.length !== 0 || gameState.isSetAvailable) return;
    alert("You won!")
  })

  // for quick testing
  function onkeydown(event: KeyboardEvent) {
    if (event.code === "Space") {
      event.preventDefault();
      const set = gameState.findSet();
      if (set) {
        for (const card of set) {
          gameState.toggleSelectCard(card.id);
        }
      }
    }
  }

</script>

<svelte:window {onkeydown} />

<div class="page">
  <h1 class="text-xl font-semibold mb-4">Game Board</h1>
  <div class="board">
    {#each gameState.inPlayCards as card (card.id)}
      <Card
        card={card}
        variationsNumber={gameState.game.variationsNumber}
        onclick={() => {
          gameState.toggleSelectCard(card.id);
        }}
      />
    {/each}
  </div>
</div>

<style>
  .page {
    padding: 2rem;
    background-color: #f0f4f8;
    min-height: 100vh;
  }
  .board {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: wrap;
    gap: 1rem;
    padding: 1rem;
  }

  @media (max-width: 600px) {
    .board {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    }
  }
</style>