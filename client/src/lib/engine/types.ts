export const FEATURES = {
  color: 'color',
  shape: 'shape',
  number: 'number',
  shading: 'shading',
  rotation: 'rotation',
} as const;

export type Feature = keyof typeof FEATURES;

export const COLORS = {
  red: '#FF0000',
  green: 'green',
  purple: '#800080',
  blue: 'blue',
} as const;
export const SHAPES = {
  diamond: 'diamond',
  squiggle: 'squiggle',
  oval: 'oval',
  arrow: 'arrow',
} as const;
export const SHADINGS = {
  solid: 'solid',
  striped: 'striped',
  empty: 'empty',
  dotted: 'dotted',
} as const;
export const ROTATIONS = [0, 90, 180, 270] as const;
export type VariationsNumber = 2 | 3 | 4;

export const featureValues = {
  color: Object.keys(COLORS) as Color[],
  shape: Object.keys(SHAPES) as Shape[],
  number: [1, 2, 3, 4],
  shading: Object.values(SHADINGS) as Shading[],
  rotation: ROTATIONS,
} as const;

export type FeatureValue = 
  | (keyof typeof COLORS)
  | (keyof typeof SHAPES)
  | (keyof typeof SHADINGS)
  | (typeof ROTATIONS[number])
  | VariationsNumber;

export type CardID = `${string}-${string}-${string}-${string}-${string}`;
export type Card = {
  id: CardID; 
  isVisible: boolean;
  isSelected: boolean;
  isDiscarded: boolean;
} & {
  // [key in Feature]?: keyof typeof featureValues[key];
  color: Color;
  shape: Shape;
  number: number;
  shading: Shading;
  rotation?: typeof ROTATIONS[number];
};

export type Color = keyof typeof COLORS;
export type Shape = keyof typeof SHAPES;
export type Shading = keyof typeof SHADINGS;