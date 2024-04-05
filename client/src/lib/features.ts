export const features = {
  color: 'color',
  shape: 'shape',
  number: 'number',
  shading: 'shading',
  rotation: 'rotation',
} as const;

export type Feature = keyof typeof features;

export const featureValues = {
  color: ['red', 'green', 'purple', 'blue'],
  shape: ['diamond', 'squiggle', 'oval', 'arrow'],
  number: [1, 2, 3, 4],
  shading: ['solid', 'striped', 'empty', 'dotted'],
  rotation: [0, 90, 180, 270],
} as const;

export type FeatureValue = typeof featureValues[Feature][number];

export type OptionsNumber = 2 | 3 | 4;

export type Card = {
  [key in Feature]?: keyof typeof featureValues[key];
};