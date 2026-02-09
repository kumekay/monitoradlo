import { writable, derived, get } from 'svelte/store';
import type { Config, Profile, Output, NiriOutput, MonitorRect } from './types';

// The full kanshi config
export const config = writable<Config>({ profiles: [] });

// Currently selected profile index
export const selectedProfileIndex = writable<number>(0);

// Currently selected output index within the profile
export const selectedOutputIndex = writable<number>(-1);

// Live niri outputs
export const niriOutputs = writable<NiriOutput[]>([]);

// Unsaved changes flag
export const hasChanges = writable<boolean>(false);

// Current profile (derived)
export const currentProfile = derived(
  [config, selectedProfileIndex],
  ([$config, $idx]) => {
    if ($config.profiles.length === 0) return null;
    return $config.profiles[Math.min($idx, $config.profiles.length - 1)] ?? null;
  }
);

// Selected output (derived)
export const selectedOutput = derived(
  [currentProfile, selectedOutputIndex],
  ([$profile, $idx]) => {
    if (!$profile || $idx < 0 || $idx >= $profile.outputs.length) return null;
    return $profile.outputs[$idx];
  }
);

// Monitor rectangles for the canvas (derived from current profile + niri data)
export const monitorRects = derived(
  [currentProfile, niriOutputs],
  ([$profile, $niri]): MonitorRect[] => {
    if (!$profile) return [];

    return $profile.outputs.map((output) => {
      // Find matching niri output by description
      const niriMatch = $niri.find(n => n.description === output.criteria);

      // Calculate logical size
      let width: number;
      let height: number;

      if (niriMatch?.logicalSize) {
        width = niriMatch.logicalSize.width;
        height = niriMatch.logicalSize.height;
      } else if (niriMatch?.currentMode) {
        const scale = output.scale ?? niriMatch?.scale ?? 1;
        width = Math.round(niriMatch.currentMode.width / scale);
        height = Math.round(niriMatch.currentMode.height / scale);
      } else if (output.mode) {
        const match = output.mode.match(/^(\d+)x(\d+)/);
        if (match) {
          const scale = output.scale ?? 1;
          width = Math.round(parseInt(match[1]) / scale);
          height = Math.round(parseInt(match[2]) / scale);
        } else {
          width = 1920;
          height = 1080;
        }
      } else {
        width = 1920;
        height = 1080;
      }

      return {
        output,
        connector: niriMatch?.connector ?? output.criteria,
        x: output.position?.x ?? 0,
        y: output.position?.y ?? 0,
        width,
        height,
        niriOutput: niriMatch,
      };
    });
  }
);

// Helper to update an output in the current profile
export function updateOutput(profileIdx: number, outputIdx: number, changes: Partial<Output>) {
  config.update(c => {
    const profile = c.profiles[profileIdx];
    if (profile && profile.outputs[outputIdx]) {
      profile.outputs[outputIdx] = { ...profile.outputs[outputIdx], ...changes };
    }
    return c;
  });
  hasChanges.set(true);
}

// Helper to update output position
export function updateOutputPosition(profileIdx: number, outputIdx: number, x: number, y: number) {
  config.update(c => {
    const profile = c.profiles[profileIdx];
    if (profile && profile.outputs[outputIdx]) {
      profile.outputs[outputIdx].position = { x, y };
    }
    return c;
  });
  hasChanges.set(true);
}
