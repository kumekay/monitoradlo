<script lang="ts">
  import {
    selectedOutput,
    selectedOutputIndex,
    selectedProfileIndex,
    monitorRects,
    niriOutputs,
    updateOutput,
    updateOutputPosition,
  } from './stores';
  import type { NiriOutput } from './types';

  $: rect = $selectedOutputIndex >= 0 ? $monitorRects[$selectedOutputIndex] : null;
  $: output = $selectedOutput;
  $: niri = rect?.niriOutput ?? null;

  // Local form values (synced from store)
  $: enabled = output?.enabled !== false;
  $: mode = output?.mode ?? '';
  $: scale = output?.scale ?? 1;
  $: transform = output?.transform ?? '';
  $: posX = output?.position?.x ?? 0;
  $: posY = output?.position?.y ?? 0;

  const transforms = ['', 'normal', '90', '180', '270', 'flipped', 'flipped-90', 'flipped-180', 'flipped-270'];

  function setEnabled(val: boolean) {
    if ($selectedOutputIndex < 0) return;
    updateOutput($selectedProfileIndex, $selectedOutputIndex, { enabled: val });
  }

  function setMode(val: string) {
    if ($selectedOutputIndex < 0) return;
    updateOutput($selectedProfileIndex, $selectedOutputIndex, { mode: val || undefined });
  }

  function setScale(val: number) {
    if ($selectedOutputIndex < 0 || isNaN(val) || val <= 0) return;
    updateOutput($selectedProfileIndex, $selectedOutputIndex, { scale: val });
  }

  function setTransform(val: string) {
    if ($selectedOutputIndex < 0) return;
    updateOutput($selectedProfileIndex, $selectedOutputIndex, { transform: val || undefined });
  }

  function setPosition(x: number, y: number) {
    if ($selectedOutputIndex < 0) return;
    updateOutputPosition($selectedProfileIndex, $selectedOutputIndex, x, y);
  }

  function formatRefresh(rate: number): string {
    return rate.toFixed(3).replace(/0+$/, '').replace(/\.$/, '');
  }

  async function applyPreview() {
    if (!niri || !output) return;
    const props: Record<string, string> = {};

    if (output.enabled === false) {
      props['off'] = '';
    } else {
      props['on'] = '';
      if (output.mode) {
        props['mode'] = output.mode;
      }
      if (output.scale) {
        props['scale'] = String(output.scale);
      }
      if (output.position) {
        props['position'] = `${output.position.x} ${output.position.y}`;
      }
      if (output.transform) {
        props['transform'] = output.transform;
      }
    }

    try {
      // @ts-ignore
      await window.go.main.App.ApplyPreview(niri.connector, props);
    } catch (e: any) {
      alert('Preview failed: ' + e.message);
    }
  }
</script>

{#if output && rect}
  <div class="properties">
    <div class="header">
      <span class="connector">{rect.connector}</span>
      <span class="description">{output.criteria}</span>
    </div>

    <div class="fields">
      <label class="field">
        <input type="checkbox" checked={enabled} on:change={(e) => setEnabled(e.currentTarget.checked)} />
        Enabled
      </label>

      <label class="field">
        <span class="field-label">Mode</span>
        {#if niri && niri.availableModes.length > 0}
          <select value={mode} on:change={(e) => setMode(e.currentTarget.value)}>
            <option value="">Default</option>
            {#each niri.availableModes as m}
              <option value="{m.width}x{m.height}@{formatRefresh(m.refreshRate)}Hz">
                {m.width}x{m.height} @ {formatRefresh(m.refreshRate)}Hz
                {m.isPreferred ? '(preferred)' : ''}
                {m.isCurrent ? '(current)' : ''}
              </option>
            {/each}
          </select>
        {:else}
          <input type="text" value={mode} on:change={(e) => setMode(e.currentTarget.value)} placeholder="e.g. 1920x1080@60Hz" />
        {/if}
      </label>

      <label class="field">
        <span class="field-label">Scale</span>
        <input
          type="number"
          value={scale}
          min="0.25"
          step="0.25"
          on:change={(e) => setScale(parseFloat(e.currentTarget.value))}
        />
      </label>

      <label class="field">
        <span class="field-label">Transform</span>
        <select value={transform} on:change={(e) => setTransform(e.currentTarget.value)}>
          {#each transforms as t}
            <option value={t}>{t || 'None'}</option>
          {/each}
        </select>
      </label>

      <div class="field position-field">
        <span class="field-label">Position</span>
        <div class="pos-inputs">
          <label>
            x <input
              type="number"
              value={posX}
              on:change={(e) => setPosition(parseInt(e.currentTarget.value) || 0, posY)}
            />
          </label>
          <label>
            y <input
              type="number"
              value={posY}
              on:change={(e) => setPosition(posX, parseInt(e.currentTarget.value) || 0)}
            />
          </label>
        </div>
      </div>

      {#if niri}
        <div class="field">
          <button class="preview-btn" on:click={applyPreview}>
            Apply Preview
          </button>
        </div>
      {/if}
    </div>
  </div>
{:else}
  <div class="properties empty">
    <span class="hint">Click a monitor on the canvas to edit its properties</span>
  </div>
{/if}

<style>
  .properties {
    padding: 10px 14px;
    background: #16213e;
    border-top: 1px solid #333;
    flex-shrink: 0;
  }

  .properties.empty {
    padding: 16px;
    text-align: center;
  }

  .hint {
    color: #666;
    font-size: 13px;
  }

  .header {
    display: flex;
    align-items: baseline;
    gap: 10px;
    margin-bottom: 10px;
  }

  .connector {
    font-weight: bold;
    color: #eee;
    font-size: 16px;
  }

  .description {
    color: #888;
    font-size: 13px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .fields {
    display: flex;
    flex-wrap: wrap;
    gap: 10px 20px;
    align-items: center;
  }

  .field {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #ccc;
    font-size: 13px;
  }

  .field-label {
    color: #888;
    min-width: 60px;
  }

  .field select,
  .field input[type="text"],
  .field input[type="number"] {
    background: #1a1a2e;
    color: #eee;
    border: 1px solid #444;
    padding: 3px 6px;
    border-radius: 3px;
    font-size: 13px;
  }

  .field input[type="number"] {
    width: 70px;
  }

  .field input[type="checkbox"] {
    accent-color: #5b8def;
  }

  .position-field {
    display: flex;
    align-items: center;
  }

  .pos-inputs {
    display: flex;
    gap: 8px;
  }

  .pos-inputs label {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #888;
    font-size: 13px;
  }

  .preview-btn {
    background: #2a3a5a;
    color: #8ab4f8;
    border: 1px solid #4a6a8a;
    padding: 4px 12px;
    border-radius: 4px;
    font-size: 13px;
    cursor: pointer;
  }

  .preview-btn:hover {
    background: #3a4a6a;
  }
</style>
