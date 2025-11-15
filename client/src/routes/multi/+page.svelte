<script lang="ts">
	import Board from "$lib/components/Board.svelte";
	import { remountKey } from "$lib/state/remountKey.svelte";
</script>

{#key remountKey.val}
<div class="page">
    <Board />
</div>
{/key}

<style>
  .page {
    display: flex;
    flex-direction: column;

    padding: 2rem;
    background-color: #f0f4f8;
    min-height: 100vh;
  }
</style>