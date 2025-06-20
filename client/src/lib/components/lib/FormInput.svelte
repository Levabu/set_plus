<script lang="ts">
  interface Props {
    value: string
    error?: string;
    label?: string;
    placeholder?: string;
    showClipboard?: boolean;
    readonly?: boolean;
    transformInput?: (val: string) => string;
  }
  let {
    value = $bindable(""),
    error = $bindable(""),
    label = "",
    placeholder = "",
    showClipboard = false,
    readonly = false,
    transformInput = (v) => v,
  }: Props = $props();

  let copied = $state(false);

  async function copyToClipboard() {
    await navigator.clipboard.writeText(value);
    copied = true;
    setTimeout(() => (copied = false), 2000);
  }
</script>

<div class="form-input-wrapper">
  {#if label}
    <label class="label" for="input-field">{label}</label>
  {/if}
  <div class="input-row">
    <input
      id="input-field"
      class="form-input"
      bind:value={
        () => value,
        (v) => {value = transformInput(v)}
      }
      readonly={readonly}
      placeholder={placeholder}
      oninput={(e) => {
        error = ""
        }}
    />
    {#if showClipboard}
      <button
        class="copy-button"
        type="button"
        onclick={copyToClipboard}
        aria-label="Copy to clipboard"
        title="Copy"
      >
        {#if copied}
          ✅
        {:else}
          📋
        {/if}
      </button>
    {/if}
  </div>
  <p class="error-text">{error}</p>
</div>

<style>
.form-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.95rem;
  font-weight: 600;
  color: #222;
}

.input-row {
  position: relative;
  display: flex;
  align-items: center;
}

.form-input {
  width: 100%;
  padding: 0.5rem 2.5rem 0.5rem 0.75rem; /* Add space on right for the icon */
  border-radius: 6px;
  border: 1px solid #ccc;
  background: #f9fafb;
  color: #222;
  font-size: 1rem;
  transition: background 0.15s, border-color 0.15s;
}

.form-input:focus {
  outline: none;
  background: #fff;
  border-color: #0070f3;
}

.copy-button {
  position: absolute;
  right: 0.5rem;
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
}

.error-text {
  min-height: 1.25rem;
  font-size: 0.85rem;
  color: #d33;
  margin: 0;
}
</style>