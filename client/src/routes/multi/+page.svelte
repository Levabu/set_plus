<script lang="ts">
	import Board from "$lib/components/Board.svelte";
</script>

<div class="page">
    <Board />
</div>

<style>
  .page {
    display: flex;
    flex-direction: column;

    padding: 2rem;
    background-color: #f0f4f8;
    min-height: 100vh;
  }
</style>