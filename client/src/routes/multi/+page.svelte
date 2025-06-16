<script lang="ts">
	import Board from "$lib/components/Board.svelte";
	import SelectGame from "$lib/components/SelectGame.svelte";
	import type { GameVersionKey } from "$lib/engine/types";
	import { GameVersions } from "$lib/engine/types";



  // let gameVersion = $state(GameVersions.classic.key) as GameVersionKey | null;
</script>

<div class="page">
  <!-- <SelectGame bind:gameVersion={gameVersion} /> -->
  <!-- {#if gameVersion !== null} -->
    <!-- <Board gameVersion={gameVersion}/> -->
    <Board />
  <!-- {/if} -->
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