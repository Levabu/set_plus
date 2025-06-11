<script lang="ts">
  import { COLORS, SHAPES, type ColorKey, type Rotation, type Shading, type Shape as ShapeType, type Color } from "$lib/engine/types";
	import Arrow from "./Arrow.svelte";
	import Diamond from "./Diamond.svelte";
	import Oval from "./Oval.svelte";
	import Squiggle from "./Squiggle.svelte";

  interface Props {
    shape: ShapeType;
    color: ColorKey;
    shading: Shading;
    width?: number | string;
    rotation?: Rotation
  }
  let { shape, color, shading, rotation = "horizontal" }: Props = $props();
  
  // type Pattern = `${Color}_${Shading}`;
  // const PATTERN_IDS: Partial<Record<Pattern, string>> = {
  //   red_striped: 'redHorizontalStripes',
  //   green_striped: 'greenHorizontalStripes',
  //   blue_striped: 'blueHorizontalStripes',
  //   purple_striped: 'purpleHorizontalStripes',
  //   red_dotted: 'redCircles',
  //   green_dotted: 'greenCircles',
  //   blue_dotted: 'blueCircles',
  //   purple_dotted: 'purpleCircles',
  // };

  function getShadingColor(shading: Shading, color: ColorKey): string {
    switch (shading) {
      case 'solid':
        return color;
      case 'striped':
        return `url(#${color}_striped)`;
      case 'dotted':
        return `url(#${color}_dotted)`;
        // return `url(#${PATTERN_IDS[`${color}_${shading}`]})`;
      case 'empty':
      default:
        return 'transparent';
    }
  }
  function getRotationAngle(rotation: Rotation): string {
    switch (rotation) {
      case 'horizontal':
        return '0';
      case 'vertical':
        return '90';
      case 'diagonal':
        return '45';
      default:
        return '0';
    }
  }
  const shadingColor = getShadingColor(shading, color);
  const rotationAngle = getRotationAngle(rotation);
</script>

<div class="shape-wrapper" style={`
  transform: rotate(${rotationAngle}deg);`
}>
<svg viewBox="0 0 60 90" fill="none" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid meet"
>
  <defs>
    {#each Object.keys(COLORS) as color}
      <pattern id={`${color}_striped`} width="10" height="8" patternUnits="userSpaceOnUse">
        <rect width="10" height="4" fill={COLORS[color as ColorKey]}/>
        <rect y="5" width="10" height="4" fill="transparent"/>
      </pattern>

      <pattern id={`${color}_dotted`} width="10" height="10" patternUnits="userSpaceOnUse">
        <circle cx="5" cy="5" r="2" fill={COLORS[color as ColorKey]}/>
      </pattern>
    {/each}
  </defs>

  {@render Shape({ baseColor: COLORS[color], shadingColor: shadingColor })}
</svg>
</div>

{#snippet Shape({baseColor, shadingColor}: {baseColor: (typeof COLORS)[ColorKey], shadingColor: string})}
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

<style>
  .shape-wrapper {
    width: 16%;
  }
</style>
