export const FEATURES = {
  color: 'color',
  shape: 'shape',
  number: 'number',
  shading: 'shading',
  rotation: 'rotation',
} as const;

export type Feature = keyof typeof FEATURES;

export const COLORS = {
  c1: '#FF0000',
  c2: 'green',
  c3: '#800080',
  c4: 'blue',
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
export const ROTATIONS = {
  vertical: 'vertical',
  horizontal: 'horizontal',
  diagonal: 'diagonal',
} as const;
export type VariationsNumber = 2 | 3 | 4;

export const featureValues = {
  color: Object.keys(COLORS) as ColorKey[],
  shape: Object.keys(SHAPES) as Shape[],
  number: [1, 2, 3, 4],
  shading: Object.values(SHADINGS) as Shading[],
  rotation: Object.keys(ROTATIONS) as Rotation[],
} as const;

export type FeatureValue = 
  | (keyof typeof COLORS)
  | (keyof typeof SHAPES)
  | (keyof typeof SHADINGS)
  | (keyof typeof ROTATIONS)
  | VariationsNumber;

export type CardID = `${string}-${string}-${string}-${string}-${string}`;
export type Card = {
  id: CardID; 
  isVisible: boolean;
  isSelected: boolean;
  isDiscarded: boolean;
} & {
  // [key in Feature]?: keyof typeof featureValues[key];
  color: ColorKey;
  shape: Shape;
  number: number;
  shading: Shading;
  rotation?: Rotation;
};

export type ColorKey = keyof typeof COLORS;
export type Color = (typeof COLORS)[ColorKey];
export type Shape = keyof typeof SHAPES;
export type Shading = keyof typeof SHADINGS;
export type Rotation = keyof typeof ROTATIONS;

export interface GameVersion {
  title: string;
  description: string;
  features: Feature[];
  variationsNumber: VariationsNumber;
}

export const GameVersions = {
  classic: {
    title: 'Classic',
    description: 'The original game with 4 features and 3 variations.',
    features: ['color', 'shape', 'number', 'shading'] as Feature[],
    variationsNumber: 3,
  } as GameVersion,
  v5x3: {
    title: '5x3',
    description: 'A variant with 5 features and 3 variations.',
    features: ['color', 'shape', 'number', 'shading', 'rotation'] as Feature[],
    variationsNumber: 3,
  } as GameVersion,
  v4x4: {
    title: '4x4',
    description: 'A variant with 4 features and 4 variations.',
    features: ['color', 'shape', 'number', 'shading'] as Feature[],
    variationsNumber: 4,
  } as GameVersion,
} as const;
