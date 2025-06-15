<script lang="ts">
	import { GameVersions, type GameVersionKey } from "$lib/engine/types";
	import { MultiPlayerGameState } from "$lib/state/MultiPlayerGameState.svelte";
	import { onDestroy } from "svelte";
	import Card from "./Card.svelte";

  interface Props {
    gameVersion: GameVersionKey;
  }

  let { gameVersion }: Props = $props();

  let gameState = $state<MultiPlayerGameState | null>(null);

  $effect(() => {
    gameState = gameVersion !== null ? new MultiPlayerGameState(GameVersions[gameVersion as GameVersionKey]) : null;

    return () => {
      console.log("rerender")
      gameState?.ws.socket?.close()
    }
  });

  onDestroy(() => {
    gameState?.ws.socket?.close()
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

  @media (max-width: 600px) {
    .board {
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    }
  }
</style>