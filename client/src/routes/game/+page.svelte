<script>
	import SingleBoard from "$lib/components/SingleBoard.svelte"
	import { remountKey } from "$lib/state/remountKey.svelte";
</script>

{#key remountKey.val}
<div class="page">
  <SingleBoard />
</div>
{/key}