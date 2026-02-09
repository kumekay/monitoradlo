<script lang="ts">
  import { monitorRects, selectedOutputIndex, selectedProfileIndex, updateOutputPosition } from './stores';
  import type { MonitorRect, SnapLine } from './types';

  const SNAP_THRESHOLD = 12;
  const COLORS = ['#5b8def', '#e5574f', '#47b86b', '#f5a623', '#9b59b6', '#1abc9c', '#e67e22', '#3498db'];

  let svgEl: SVGSVGElement;
  let dragging = false;
  let dragIndex = -1;
  let dragOffsetX = 0;
  let dragOffsetY = 0;
  let activeSnapLines: SnapLine[] = [];
  let frozenViewBox: string | null = null;

  // Compute the view transform to fit all monitors, but freeze during drag
  $: computedViewBox = computeViewBox($monitorRects);
  $: viewBox = frozenViewBox ?? computedViewBox;

  // Scale factor for strokes, markers etc. based on total canvas extent
  $: canvasExtent = (() => {
    if ($monitorRects.length === 0) return 1920;
    let minX = Infinity, maxX = -Infinity, minY = Infinity, maxY = -Infinity;
    for (const r of $monitorRects) {
      minX = Math.min(minX, r.x);
      maxX = Math.max(maxX, r.x + r.width);
      minY = Math.min(minY, r.y);
      maxY = Math.max(maxY, r.y + r.height);
    }
    return Math.max(maxX - minX, maxY - minY);
  })();
  $: sw = canvasExtent * 0.002; // stroke width

  function computeViewBox(rects: MonitorRect[]): string {
    if (rects.length === 0) return '0 0 1920 1080';

    let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
    for (const r of rects) {
      minX = Math.min(minX, r.x);
      minY = Math.min(minY, r.y);
      maxX = Math.max(maxX, r.x + r.width);
      maxY = Math.max(maxY, r.y + r.height);
    }

    const w = maxX - minX;
    const h = maxY - minY;
    const padding = Math.max(w, h) * 0.05;
    return `${minX - padding} ${minY - padding} ${w + padding * 2} ${h + padding * 2}`;
  }

  function isEnabled(rect: MonitorRect): boolean {
    return rect.output.enabled !== false;
  }

  function onMouseDown(e: MouseEvent, index: number) {
    e.preventDefault();
    selectedOutputIndex.set(index);
    dragging = true;
    dragIndex = index;
    frozenViewBox = computedViewBox;

    const pt = svgPoint(e);
    const rect = $monitorRects[index];
    dragOffsetX = pt.x - rect.x;
    dragOffsetY = pt.y - rect.y;
  }

  function onMouseMove(e: MouseEvent) {
    if (!dragging || dragIndex < 0) return;
    e.preventDefault();

    const pt = svgPoint(e);
    let newX = Math.round(pt.x - dragOffsetX);
    let newY = Math.round(pt.y - dragOffsetY);

    // Snap
    const result = snap(newX, newY, dragIndex);
    newX = result.x;
    newY = result.y;
    activeSnapLines = result.snapLines;

    updateOutputPosition($selectedProfileIndex, dragIndex, newX, newY);
  }

  function onMouseUp() {
    dragging = false;
    dragIndex = -1;
    activeSnapLines = [];
    frozenViewBox = null;
  }

  function svgPoint(e: MouseEvent): { x: number; y: number } {
    const svg = svgEl;
    const ctm = svg.getScreenCTM();
    if (!ctm) return { x: 0, y: 0 };
    return {
      x: (e.clientX - ctm.e) / ctm.a,
      y: (e.clientY - ctm.f) / ctm.d,
    };
  }

  function snap(x: number, y: number, dragIdx: number): { x: number; y: number; snapLines: SnapLine[] } {
    const dragged = $monitorRects[dragIdx];
    const w = dragged.width;
    const h = dragged.height;
    const lines: SnapLine[] = [];

    // Edges and center of the dragged monitor
    const dragEdges = {
      left: x,
      right: x + w,
      top: y,
      bottom: y + h,
      centerX: x + w / 2,
      centerY: y + h / 2,
    };

    let snappedX = false;
    let snappedY = false;
    let bestDx = SNAP_THRESHOLD + 1;
    let bestDy = SNAP_THRESHOLD + 1;

    for (let i = 0; i < $monitorRects.length; i++) {
      if (i === dragIdx) continue;
      const other = $monitorRects[i];

      const otherEdges = {
        left: other.x,
        right: other.x + other.width,
        top: other.y,
        bottom: other.y + other.height,
        centerX: other.x + other.width / 2,
        centerY: other.y + other.height / 2,
      };

      // Horizontal snaps (X axis)
      const xSnaps = [
        { drag: dragEdges.left, target: otherEdges.left, adjust: 0 },
        { drag: dragEdges.left, target: otherEdges.right, adjust: 0 },
        { drag: dragEdges.right, target: otherEdges.left, adjust: -w },
        { drag: dragEdges.right, target: otherEdges.right, adjust: -w },
        { drag: dragEdges.centerX, target: otherEdges.centerX, adjust: -w / 2 },
      ];

      for (const s of xSnaps) {
        const d = Math.abs(s.drag - s.target);
        if (d < SNAP_THRESHOLD && d < bestDx) {
          bestDx = d;
          x = s.target + s.adjust;
          snappedX = true;
          // Record snap line
          const minY = Math.min(y, other.y);
          const maxY = Math.max(y + h, other.y + other.height);
          lines.push({
            orientation: 'vertical',
            position: s.target,
            start: minY - 10,
            end: maxY + 10,
          });
        }
      }

      // Vertical snaps (Y axis)
      const ySnaps = [
        { drag: dragEdges.top, target: otherEdges.top, adjust: 0 },
        { drag: dragEdges.top, target: otherEdges.bottom, adjust: 0 },
        { drag: dragEdges.bottom, target: otherEdges.top, adjust: -h },
        { drag: dragEdges.bottom, target: otherEdges.bottom, adjust: -h },
        { drag: dragEdges.centerY, target: otherEdges.centerY, adjust: -h / 2 },
      ];

      for (const s of ySnaps) {
        const d = Math.abs(s.drag - s.target);
        if (d < SNAP_THRESHOLD && d < bestDy) {
          bestDy = d;
          y = s.target + s.adjust;
          snappedY = true;
          const minX = Math.min(x, other.x);
          const maxX = Math.max(x + w, other.x + other.width);
          lines.push({
            orientation: 'horizontal',
            position: s.target,
            start: minX - 10,
            end: maxX + 10,
          });
        }
      }
    }

    return { x: Math.round(x), y: Math.round(y), snapLines: lines };
  }
</script>

<svelte:window on:mousemove={onMouseMove} on:mouseup={onMouseUp} />

<div class="canvas-container">
  <svg
    bind:this={svgEl}
    {viewBox}
    preserveAspectRatio="xMidYMid meet"
    class="canvas"
    class:dragging
  >
    <!-- Grid origin marker -->
    <line x1={-canvasExtent*0.01} y1="0" x2={canvasExtent*0.01} y2="0" stroke="#555" stroke-width={sw} stroke-dasharray={sw*4} />
    <line x1="0" y1={-canvasExtent*0.01} x2="0" y2={canvasExtent*0.01} stroke="#555" stroke-width={sw} stroke-dasharray={sw*4} />

    <!-- Monitor rectangles -->
    {#each $monitorRects as rect, i}
      <g
        class="monitor"
        class:selected={$selectedOutputIndex === i}
        class:disabled={!isEnabled(rect)}
        on:mousedown={(e) => onMouseDown(e, i)}
        role="button"
        tabindex="0"
      >
        <rect
          x={rect.x}
          y={rect.y}
          width={rect.width}
          height={rect.height}
          fill={isEnabled(rect) ? COLORS[i % COLORS.length] + '30' : '#33333330'}
          stroke={$selectedOutputIndex === i ? '#fff' : COLORS[i % COLORS.length]}
          stroke-width={$selectedOutputIndex === i ? sw * 2 : sw}
          rx={sw * 2}
        />
        {#if !isEnabled(rect)}
          <!-- Diagonal hatch for disabled -->
          <line
            x1={rect.x} y1={rect.y}
            x2={rect.x + rect.width} y2={rect.y + rect.height}
            stroke="#666" stroke-width={sw} stroke-dasharray={sw*8}
          />
          <line
            x1={rect.x + rect.width} y1={rect.y}
            x2={rect.x} y2={rect.y + rect.height}
            stroke="#666" stroke-width={sw} stroke-dasharray={sw*8}
          />
        {/if}
        <!-- Connector/name -->
        <text
          x={rect.x + rect.width / 2}
          y={rect.y + rect.height / 2 - Math.min(rect.width, rect.height) * 0.04}
          text-anchor="middle"
          dominant-baseline="middle"
          fill={isEnabled(rect) ? '#eee' : '#888'}
          font-size={Math.min(rect.width, rect.height) * 0.07}
          font-weight="bold"
          pointer-events="none"
        >
          {rect.connector}
        </text>
        <!-- Resolution -->
        <text
          x={rect.x + rect.width / 2}
          y={rect.y + rect.height / 2 + Math.min(rect.width, rect.height) * 0.05}
          text-anchor="middle"
          dominant-baseline="middle"
          fill={isEnabled(rect) ? '#ccc' : '#666'}
          font-size={Math.min(rect.width, rect.height) * 0.05}
          pointer-events="none"
        >
          {rect.width}x{rect.height}
        </text>
      </g>
    {/each}

    <!-- Snap guide lines -->
    {#each activeSnapLines as line}
      {#if line.orientation === 'vertical'}
        <line
          x1={line.position} y1={line.start}
          x2={line.position} y2={line.end}
          stroke="#f5a623" stroke-width={sw} stroke-dasharray={sw*4}
        />
      {:else}
        <line
          x1={line.start} y1={line.position}
          x2={line.end} y2={line.position}
          stroke="#f5a623" stroke-width={sw} stroke-dasharray={sw*4}
        />
      {/if}
    {/each}
  </svg>
</div>

<style>
  .canvas-container {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #1a1a2e;
    min-height: 200px;
    overflow: hidden;
  }

  .canvas {
    width: 100%;
    height: 100%;
    cursor: default;
  }

  .canvas.dragging {
    cursor: grabbing;
  }

  .monitor {
    cursor: grab;
  }

  .monitor:hover rect {
    filter: brightness(1.3);
  }

  .monitor.selected rect {
    filter: brightness(1.2);
  }

  .monitor.disabled {
    opacity: 0.5;
  }
</style>
