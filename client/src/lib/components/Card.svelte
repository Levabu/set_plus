<script lang="ts">
  import { type Card, type VariationsNumber } from "$lib/engine/types";
	import ShapeWrapper from "./shapes/ShapeWrapper.svelte";

  interface Props {
    card: Card;
    onclick: () => void;
  }

  let { card, onclick }: Props = $props();
</script>

<button
  class={`shapes ${card.isSelected ? 'selected' : ''}`}
  onclick={onclick}
>
  {#each { length: card.number }}
      <ShapeWrapper
        shape={card.shape}
        color={card.color}
        shading={card.shading}
        rotation={card.rotation || "vertical"}
      />
  {/each}
</button>

<style>
  .shapes {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10%;
    overflow: hidden;
    flex-wrap: wrap;
    box-sizing: border-box;

    padding: 0.5rem;

    border: 1px solid #ccc;
    border-radius: 0.5rem;
    background-color: #f9f9f9;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    
    min-height: 9rem;

    cursor: pointer;
  }

  .shapes.selected {
    border-color: #4a90e2;
    box-shadow: 0 0 10px rgba(74, 144, 226, 0.5);
  }
</style>