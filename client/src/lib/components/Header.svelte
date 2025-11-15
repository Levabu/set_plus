<script lang="ts">
  import { page } from '$app/state';
	import { remountKey } from '$lib/state/remountKey.svelte';

  function triggerRemount(pathname: string) {
    if (page?.url?.pathname !== pathname) return;
    remountKey.val += 1;
  };
</script>

<nav class="bg-gray-800 text-white p-4 flex gap-4">
  <a class="hover:underline" href="/">Home</a>
  <a class="hover:underline" href="/game" onclick={() => {
    triggerRemount("/game")
  }}>Game</a>
  <a class="hover:underline" href="/multi" onclick={() => {
    triggerRemount("/multi")
  }}>Mutli</a>
</nav>