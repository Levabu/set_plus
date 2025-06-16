<script lang="ts">
  interface Props {
    roomLink: string
  }
  let { roomLink = $bindable() }: Props = $props();

  let copied = $state(false);
  // let value = $state<string>(roomLink)

  async function copyToClipboard() {
    await navigator.clipboard.writeText(roomLink);
    copied = true;
    setTimeout(() => (copied = false), 2000);
  }
</script>

<div class="copy-input-wrapper">
  <input class="copy-input" type="text" bind:value={roomLink} />
  <button class="copy-button" onclick={copyToClipboard}>
    {#if copied}
      ✅
    {:else}
      📋
    {/if}
  </button>
</div>

<style>
.copy-input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  background: #f9fafb;
  border: 1px solid #ccc;
  /* max-width: fit-content; */
}

.copy-input {
  border: none;
  background: transparent;
  color: #222;
  font-size: 1rem;
  font-weight: 600;
  outline: none;
  /* width: 10ch; */
  width: 100%;
}

.copy-button {
  border: 1px solid transparent;
  background: #e9ecef;
  padding: 0.3rem 0.6rem; 
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.copy-button:hover {
  background: #dfe4e8;
  border-color: #b3b3b3;
}
</style>