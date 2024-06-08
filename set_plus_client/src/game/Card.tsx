import React from 'react'
import styled from 'styled-components'
import type { Card, Color, Shading, Shape } from '../lib/features';
import ShapeWrapper from '../components/shapes/Shape';



interface ICard {
  card: Card;
  optionsNumber: number;
}
const CARD_PADDING = 5;
const CADR_GAP = 5;

export default function Card({ card, optionsNumber }: ICard) {
  const shapeWidth = (100 - CARD_PADDING - CADR_GAP * Number(card.number) - 1) / optionsNumber;
  return (
    <Wrapper>
      {
        Array.from({ length: Number(card.number) }).map((_, index) => <ShapeWrapper key={index} {...{
          shape: card.shape as Shape,
          color: card.color as Color,
          shading: card.shading as Shading,
          width: shapeWidth
        }} />)
      }
    </Wrapper>
  )
}

const Wrapper = styled.ul`
  width: 15vw;
  border: solid black 1px;

  padding: ${CARD_PADDING}%;

  display: flex;
  justify-content: center;
  align-items: center;
  gap: ${CADR_GAP}%;
`;