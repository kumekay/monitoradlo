<script lang="ts">
  import { onMount } from 'svelte';
  import Canvas from './lib/Canvas.svelte';
  import ProfileBar from './lib/ProfileBar.svelte';
  import Properties from './lib/Properties.svelte';
  import { config, niriOutputs } from './lib/stores';

  onMount(async () => {
    try {
      // Load kanshi config
      // @ts-ignore
      const cfg = await window.go.main.App.LoadConfig();
      if (cfg) {
        config.set(cfg);
      }
    } catch (e: any) {
      console.error('Failed to load config:', e);
    }

    try {
      // Detect connected outputs
      // @ts-ignore
      const outputs = await window.go.main.App.DetectOutputs();
      if (outputs) {
        niriOutputs.set(outputs);
      }
    } catch (e: any) {
      console.error('Failed to detect outputs:', e);
    }
  });
</script>

<main>
  <ProfileBar />
  <Canvas />
  <Properties />
</main>

<style>
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
