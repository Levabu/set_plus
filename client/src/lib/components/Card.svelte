<script lang="ts">
  import { type Card, type VariationsNumber } from "$lib/engine/types";
	import ShapeWrapper from "./shapes/ShapeWrapper.svelte";

  const CARD_PADDING = 5;
  const CARD_GAP = 5;

  interface Props {
    card: Card;
    variationsNumber: VariationsNumber;
    onclick: () => void;
  }

  let { card, variationsNumber, onclick }: Props = $props();
  let width = (100 - CARD_PADDING - CARD_GAP * Number(card.number) - 1) / variationsNumber;
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
        width={width}
      />
  {/each}
  <!-- <div class="mt-2 text-sm text-gray-600">
    {card.shape} - {card.color} - {card.shading} - {card.number}
  </div> -->
</button>

<style>
  .shapes {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 5px;

    padding: var(--card-padding, 5px);

    border: 1px solid #ccc;
    border-radius: 0.5rem;
    background-color: #f9f9f9;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 25vw;
    height: 100px;
    overflow: hidden;
    flex-wrap: wrap;
    box-sizing: border-box;
  }

  .shapes.selected {
    border-color: #4a90e2;
    box-shadow: 0 0 10px rgba(74, 144, 226, 0.5);
  }

  /* .shapes > div {
    height: 5rem;
  } */
</style>