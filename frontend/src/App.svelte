<script lang="ts">
  import { onMount } from 'svelte';
  import Canvas from './lib/Canvas.svelte';
  import ProfileBar from './lib/ProfileBar.svelte';
  import Properties from './lib/Properties.svelte';
  import { config, niriOutputs, selectedProfileIndex, hasChanges } from './lib/stores';
  import type { Config, NiriOutput } from './lib/types';
  import { LoadConfig, DetectOutputs, SaveConfig, ReloadKanshi } from '../wailsjs/go/main/App';

  // Find the profile that best matches the currently connected outputs.
  // A profile matches if all its output criteria appear in the niri descriptions.
  function findMatchingProfile(cfg: Config, niri: NiriOutput[]): number {
    const niriDescs = new Set(niri.map(n => n.description));
    let bestIdx = 0;
    let bestCount = -1;

    for (let i = 0; i < cfg.profiles.length; i++) {
      const profile = cfg.profiles[i];
      const matchCount = profile.outputs.filter(o => niriDescs.has(o.criteria)).length;
      // Prefer profiles where ALL outputs match, then by match count
      if (matchCount === profile.outputs.length && matchCount > bestCount) {
        bestCount = matchCount;
        bestIdx = i;
      }
    }
    return bestIdx;
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
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault();
      if ($hasChanges) {
        save();
      }
    }
  }

  function handleBeforeUnload(e: BeforeUnloadEvent) {
    if ($hasChanges) {
      e.preventDefault();
      // Modern browsers show a generic message; returnValue is for legacy support
      e.returnValue = 'You have unsaved changes. Are you sure you want to leave?';
    }
  }

  onMount(async () => {
    let cfg: Config | null = null;
    let outputs: NiriOutput[] | null = null;

    try {
      cfg = await LoadConfig() as unknown as Config;
      if (cfg) {
        config.set(cfg);
      }
    } catch (e: any) {
      console.error('Failed to load config:', e);
    }

    try {
      outputs = await DetectOutputs() as unknown as NiriOutput[];
      if (outputs) {
        niriOutputs.set(outputs);
      }
    } catch (e: any) {
      console.error('Failed to detect outputs:', e);
    }

    // Auto-select the profile matching current outputs
    if (cfg && outputs && outputs.length > 0) {
      selectedProfileIndex.set(findMatchingProfile(cfg, outputs));
    }
  });
</script>

<svelte:window on:keydown={handleKeydown} on:beforeunload={handleBeforeUnload} />

<main>
  <ProfileBar />
  <Canvas />
  <Properties />
</main>

<style>
  :global(html) {
    color-scheme: dark;
  }

  :global(body) {
    margin: 0;
    padding: 0;
    background: #0f0f23;
    color: #eee;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    overflow: hidden;
  }

  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
  }
</style>
