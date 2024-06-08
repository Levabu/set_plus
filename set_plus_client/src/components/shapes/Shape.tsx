import React from 'react'
import { SHAPES } from '../../lib/shapes'
import { COLORS } from '../../lib/colors'
import { Color } from '../../lib/features'

export default function ShapeWrapper({
  shape,
  color,
  shading,
  width
}: {
  shape: string,
  color: string,
  shading: string
  width: number,
}) {
  const Shape = SHAPES[shape as keyof typeof SHAPES]
  const baseColor = COLORS[color as keyof typeof COLORS]
  let shadingColor;
  switch (shading) {
    case 'solid':
      shadingColor = baseColor;
      break;
    case 'striped':
      shadingColor = `url(#${baseColor}HorizontalStripes)`;
      break;
    case 'empty':
      shadingColor = 'transparent';
      break;
    case 'dotted':
      shadingColor = 'transparent';
      break;
    default:
      shadingColor = 'transparent';
  }
  
  return (
    <svg width={`${width}%`} height="90" viewBox="0 0 60 90" fill="none" xmlns="http://www.w3.org/2000/svg">
      <defs>
        {Object.values(COLORS).map((color) => (
          <pattern key={color} id={`${color}HorizontalStripes`} width="10" height="8" patternUnits="userSpaceOnUse">
            <rect width="10" height="4" fill={color}/>
            <rect y="5" width="10" height="4" fill="transparent"/>
          </pattern>
        ))}
      </defs>
      <Shape baseColor={baseColor as Color} shadingColor={shadingColor} />
    </svg>
  )
}
