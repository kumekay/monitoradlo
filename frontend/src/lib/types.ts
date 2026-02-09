// Types matching Go structs

export interface Config {
  profiles: Profile[];
  preamble?: string;
}

export interface Profile {
  name: string;
  outputs: Output[];
  extraLines?: string[];
}

export interface Output {
  criteria: string;
  enabled?: boolean;
  mode?: string;
  scale?: number;
  position?: Position;
  transform?: string;
  adaptiveSync?: boolean;
}

export interface Position {
  x: number;
  y: number;
}

// Niri live output info
export interface NiriOutput {
  connector: string;
  make: string;
  model: string;
  serial: string;
  description: string;
  currentMode: NiriMode;
  availableModes: NiriMode[];
  logicalPosition?: { x: number; y: number };
  logicalSize?: { width: number; height: number };
  scale: number;
  transform: string;
  physicalSize?: { width: number; height: number };
}

export interface NiriMode {
  width: number;
  height: number;
  refreshRate: number;
  isCurrent: boolean;
  isPreferred: boolean;
}

// Canvas-specific types
export interface MonitorRect {
  output: Output;
  connector: string;
  x: number;
  y: number;
  width: number;
  height: number;
  niriOutput?: NiriOutput;
}

export interface SnapLine {
  orientation: 'horizontal' | 'vertical';
  position: number;
  start: number;
  end: number;
}
