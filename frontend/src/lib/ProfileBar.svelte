<script lang="ts">
  import { config, selectedProfileIndex, hasChanges, niriOutputs } from './stores';
  import type { Profile } from './types';
  import { SaveConfig, ReloadKanshi } from '../../wailsjs/go/main/App';

  let renaming = false;
  let renameValue = '';

  $: profiles = $config.profiles;
  $: currentName = profiles[$selectedProfileIndex]?.name ?? '';

  function selectProfile(idx: number) {
    selectedProfileIndex.set(idx);
  }

  function addProfile() {
    const name = `Profile ${profiles.length + 1}`;
    config.update(c => {
      // Try to pre-populate with detected outputs
      const outputs = $niriOutputs.map(n => ({
        criteria: n.description,
        enabled: true,
        scale: n.scale || undefined,
        position: n.logicalPosition ? { x: n.logicalPosition.x, y: n.logicalPosition.y } : { x: 0, y: 0 },
      }));

      c.profiles.push({
        name,
        outputs: outputs.length > 0 ? outputs : [],
      });
      return c;
    });
    selectedProfileIndex.set(profiles.length - 1);
    hasChanges.set(true);
  }

  function startRename() {
    renaming = true;
    renameValue = currentName;
  }

  function finishRename() {
    if (renameValue.trim()) {
      config.update(c => {
        c.profiles[$selectedProfileIndex].name = renameValue.trim();
        return c;
      });
      hasChanges.set(true);
    }
    renaming = false;
  }

  function deleteProfile() {
    if (profiles.length <= 1) return;
    if (!confirm(`Delete profile "${currentName}"?`)) return;

    config.update(c => {
      c.profiles.splice($selectedProfileIndex, 1);
      return c;
    });
    if ($selectedProfileIndex >= profiles.length) {
      selectedProfileIndex.set(Math.max(0, profiles.length - 1));
    }
    hasChanges.set(true);
  }

  async function save() {
    try {
      await SaveConfig($config as any);
      await ReloadKanshi();
      hasChanges.set(false);
    } catch (e: any) {
      alert('Save failed: ' + (e?.message ?? e));
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (renaming && e.key === 'Enter') {
      finishRename();
    } else if (renaming && e.key === 'Escape') {
      renaming = false;
    }
  }
</script>

<div class="profile-bar">
  <span class="label">Profile:</span>

  {#if renaming}
    <input
      class="rename-input"
      bind:value={renameValue}
      on:blur={finishRename}
      on:keydown={handleKeydown}
      autofocus
    />
  {:else}
    <select
      class="profile-select"
      value={$selectedProfileIndex}
      on:change={(e) => selectProfile(parseInt(e.currentTarget.value))}
    >
      {#each profiles as profile, i}
        <option value={i}>{profile.name || `Profile ${i + 1}`}</option>
      {/each}
    </select>
  {/if}

  <div class="actions">
    <button on:click={addProfile} title="New profile">+ New</button>
    <button on:click={startRename} title="Rename profile">Rename</button>
    <button
      on:click={deleteProfile}
      title="Delete profile"
      disabled={profiles.length <= 1}
    >Delete</button>
    <button
      class="save-btn"
      class:has-changes={$hasChanges}
      on:click={save}
      title="Save to kanshi config"
    >Save</button>
  </div>
</div>

<style>
  .profile-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background: #16213e;
    border-bottom: 1px solid #333;
    flex-shrink: 0;
  }

  .label {
    color: #aaa;
    font-size: 14px;
    white-space: nowrap;
  }

  .profile-select {
    background: #1a1a2e;
    color: #eee;
    border: 1px solid #444;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 14px;
    min-width: 150px;
  }

  .profile-select option {
    background: #1a1a2e;
    color: #eee;
  }

  .rename-input {
    background: #1a1a2e;
    color: #eee;
    border: 1px solid #5b8def;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 14px;
    min-width: 150px;
    outline: none;
  }

  .actions {
    display: flex;
    gap: 4px;
    margin-left: auto;
  }

  button {
    background: #2a2a4a;
    color: #ccc;
    border: 1px solid #444;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 13px;
    cursor: pointer;
  }

  button:hover {
    background: #3a3a5a;
    color: #fff;
  }

  button:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .save-btn {
    background: #2a4a2a;
    border-color: #4a6a4a;
  }

  .save-btn:hover {
    background: #3a6a3a;
  }

  .save-btn.has-changes {
    background: #4a6a2a;
    border-color: #6a8a4a;
    color: #fff;
  }
</style>
