<script lang="ts">
	import type { ColorKey, COLORS } from "$lib/engine/types";

  interface Props {
    baseColor: (typeof COLORS)[ColorKey];
    shadingColor: string;
  }
  const { shadingColor, baseColor }: Props = $props()
</script>

<mask id="path-1-inside-1_201_50" fill="white">
  <path fill-rule="evenodd" clip-rule="evenodd" d="M30 0L60 45L30 90L0 45L30 0Z"/>
</mask>
<path fill-rule="evenodd" clip-rule="evenodd" d="M30 0L60 45L30 90L0 45L30 0Z" fill={shadingColor}/>
<path d="M30 0L32.4962 -1.6641L30 -5.40833L27.5038 -1.6641L30 0ZM60 45L62.4962 46.6641L63.6056 45L62.4962 43.3359L60 45ZM30 90L27.5038 91.6641L30 95.4083L32.4962 91.6641L30 90ZM0 45L-2.49615 43.3359L-3.60555 45L-2.49615 46.6641L0 45ZM27.5038 1.6641L57.5038 46.6641L62.4962 43.3359L32.4962 -1.6641L27.5038 1.6641ZM57.5038 43.3359L27.5038 88.3359L32.4962 91.6641L62.4962 46.6641L57.5038 43.3359ZM32.4962 88.3359L2.49615 43.3359L-2.49615 46.6641L27.5038 91.6641L32.4962 88.3359ZM2.49615 46.6641L32.4962 1.6641L27.5038 -1.6641L-2.49615 43.3359L2.49615 46.6641Z" fill={baseColor} mask="url(#path-1-inside-1_201_50)"/>
