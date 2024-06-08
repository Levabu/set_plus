export const features = {
  color: 'color',
  shape: 'shape',
  number: 'number',
  shading: 'shading',
  rotation: 'rotation',
} as const;

export type Feature = keyof typeof features;

const COLORS = ['red', 'green', 'purple', 'blue'] as const;
const SHAPES = ['diamond', 'squiggle', 'oval', 'arrow'] as const;
const SHADINGS = ['solid', 'striped', 'empty', 'dotted'] as const;
const ROTATIONS = [0, 90, 180, 270] as const;
// type Color = typeof COLORS[number];

export const featureValues = {
  color: COLORS,
  shape: SHAPES,
  number: [1, 2, 3, 4],
  shading: SHADINGS,
  rotation: ROTATIONS,
} as const;

export type FeatureValue = typeof featureValues[Feature][number];

export type OptionsNumber = 2 | 3 | 4;

export type Card = {
  [key in Feature]?: keyof typeof featureValues[key];
};

export type Color = typeof COLORS[number];
export type Shape = typeof SHAPES[number];
export type Shading = typeof SHADINGS[number];