<script lang="ts">
  import { Label, RadioGroup } from 'bits-ui';
  import { GameVersions, type GameVersion, type GameVersionKey } from '$lib/engine/types';
  
  interface Props {
    gameVersion: GameVersionKey | null;
  }
  
  let { gameVersion = $bindable() }: Props = $props();
</script>

<div class="game-version-selector">
  <RadioGroup.Root bind:value={gameVersion} class="radio-group">
  {#each Object.entries(GameVersions) as [key, version]}
    <RadioGroup.Item value={key}>
      {#snippet children({ props, checked })}
        <div class="radio-item" {...props}>
          <div class="radio-indicator">
            {#if checked}
              <div class="indicator" />
            {/if}
          </div>
          <div class="option-content">
            <div class="option-header">
              <Label.Root
                class="title"
                for={key}
              >{version.title}</Label.Root>
              <!-- <h4 class="title">{version.title}</h4> -->
            </div>
          </div>
        </div>
      {/snippet}
    </RadioGroup.Item>
  {/each}
</RadioGroup.Root>
</div>

<style>
.game-version-selector {
  border-radius: 8px;
  max-width: fit-content;
}

.legend {
  margin-bottom: 1rem;
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

:global(.game-version-selector .title) {
  color: white;
}

:global(.game-version-selector .radio-group) {
  display: flex;
  flex-direction: row;
  gap: 0.75rem;
}

 .radio-item {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  transition: background 0.15s;
  border: 1px solid transparent;
}

.radio-item:hover,
.radio-item:focus-within {
  background: #e9ecef;
  border-color: #b3b3b3;
}

.radio-indicator {
  width: 22px;
  height: 22px;
  border: 2px solid #888;
  border-radius: 50%;
  margin-right: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  transition: border-color 0.15s;
}

.radio-item[aria-checked="true"] .radio-indicator {
  border-color: #0070f3;
}

.indicator {
  width: 12px;
  height: 12px;
  background: #0070f3;
  border-radius: 50%;
}

.option-content {
  /* flex: 1; */
}

:global([data-label-root]) {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #222;
}
</style>