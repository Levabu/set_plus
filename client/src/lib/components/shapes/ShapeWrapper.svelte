<script lang="ts">
  import { COLORS, SHAPES, type Color, type Shading, type Shape as ShapeType } from "$lib/engine/types";
	import Arrow from "./Arrow.svelte";
	import Diamond from "./Diamond.svelte";
	import Oval from "./Oval.svelte";
	import Squiggle from "./Squiggle.svelte";

  interface Props {
    shape: ShapeType;
    color: Color;
    shading: Shading;
    width: number | string;
  }
  let { shape, color, shading, width }: Props = $props();

  const PATTERN_IDS: Record<Color, string> = {
    red: 'redHorizontalStripes',
    green: 'greenHorizontalStripes',
    blue: 'blueHorizontalStripes',
    purple: 'purpleHorizontalStripes',
  };

  function getShadingColor(shading: Shading, color: Color): string {
    switch (shading) {
      case 'solid':
        return color;
      case 'striped':
        return `url(#${PATTERN_IDS[color]})`;
      case 'empty':
      case 'dotted':
      default:
        return 'transparent';
    }
  }
  const shadingColor = getShadingColor(shading, color);
</script>

<div>
<svg width={width} height="90" viewBox="0 0 60 90" fill="none" xmlns="http://www.w3.org/2000/svg">
  <defs>
    {#each Object.keys(COLORS) as color}
      <pattern id={`${color}HorizontalStripes`} width="10" height="8" patternUnits="userSpaceOnUse">
        <rect width="10" height="4" fill={color}/>
        <rect y="5" width="10" height="4" fill="transparent"/>
      </pattern>

      <pattern id={`${color}Circles`} width="10" height="10" patternUnits="userSpaceOnUse">
        <circle cx="5" cy="5" r="2" fill={color}/>
        <circle cx="5" cy="5" r="1" fill="white"/>
      </pattern>
    {/each}
  </defs>

  {@render Shape({ baseColor: color, shadingColor: shadingColor })}
</svg>
</div>

{#snippet Shape({baseColor, shadingColor}: {baseColor: Color, shadingColor: string})}
  {#if shape === SHAPES.diamond}
    <Diamond baseColor={baseColor} shadingColor={shadingColor} />
  {:else if shape === SHAPES.squiggle}
    <Squiggle baseColor={baseColor} shadingColor={shadingColor} />
  {:else if shape === SHAPES.oval}
    <Oval baseColor={baseColor} shadingColor={shadingColor} />
  {:else if shape === SHAPES.arrow}
    <Arrow baseColor={baseColor} shadingColor={shadingColor} />
  {/if}

{/snippet}

<!-- <style>
  div {
    /* Static styles can remain here */
  }
</style> -->
